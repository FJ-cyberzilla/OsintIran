// auth-service/internal/middleware/auth_middleware.go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
    jwtSecret        string
    tenantService    *services.TenantService
    userService      *services.UserService
}

type ContextKey string

const (
    TenantIDKey ContextKey = "tenant_id"
    UserIDKey   ContextKey = "user_id"
    RoleKey     ContextKey = "role"
)

func NewAuthMiddleware(jwtSecret string, tenantService *services.TenantService, userService *services.UserService) *AuthMiddleware {
    return &AuthMiddleware{
        jwtSecret:     jwtSecret,
        tenantService: tenantService,
        userService:   userService,
    }
}

// JWTAuthMiddleware validates JWT tokens and sets context
func (am *AuthMiddleware) JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Extract token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }

        tokenString := parts[1]

        // Parse and validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(am.jwtSecret), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        // Extract claims
        tenantID, ok := claims["tenant_id"].(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid tenant ID in token"})
            c.Abort()
            return
        }

        userID, ok := claims["user_id"].(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
            c.Abort()
            return
        }

        role, _ := claims["role"].(string)

        // Verify tenant is active
        tenant, err := am.tenantService.GetTenantByID(c.Request.Context(), tenantID)
        if err != nil || tenant.Status != "active" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not active"})
            c.Abort()
            return
        }

        // Set context values
        ctx := context.WithValue(c.Request.Context(), TenantIDKey, tenantID)
        ctx = context.WithValue(ctx, UserIDKey, userID)
        ctx = context.WithValue(ctx, RoleKey, role)
        
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}

// APIKeyAuthMiddleware validates API keys
func (am *AuthMiddleware) APIKeyAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        if apiKey == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
            c.Abort()
            return
        }

        // Hash the provided API key
        hashedKey := am.hashAPIKey(apiKey)

        // Look up API key in database
        apiKeyRecord, err := am.userService.GetAPIKeyByHash(c.Request.Context(), hashedKey)
        if err != nil || !apiKeyRecord.IsActive {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
            c.Abort()
            return
        }

        // Check if API key has expired
        if apiKeyRecord.ExpiresAt != nil && apiKeyRecord.ExpiresAt.Before(time.Now()) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "API key expired"})
            c.Abort()
            return
        }

        // Update last used timestamp
        go am.userService.UpdateAPIKeyLastUsed(c.Request.Context(), apiKeyRecord.ID)

        // Set context
        ctx := context.WithValue(c.Request.Context(), TenantIDKey, apiKeyRecord.TenantID)
        ctx = context.WithValue(ctx, UserIDKey, "") // No user for API keys
        ctx = context.WithValue(ctx, RoleKey, "api_key")
        
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}

// PermissionMiddleware checks if user has required permissions
func (am *AuthMiddleware) PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := c.Request.Context()
        role, ok := ctx.Value(RoleKey).(string)
        if !ok {
            c.JSON(http.StatusForbidden, gin.H{"error": "Role not found in context"})
            c.Abort()
            return
        }

        // Get role permissions
        permissions, err := am.userService.GetRolePermissions(ctx, role)
        if err != nil {
            c.JSON(http.StatusForbidden, gin.H{"error": "Failed to get permissions"})
            c.Abort()
            return
        }

        // Check if role has required permission
        if !am.hasPermission(permissions, requiredPermission) {
            c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
}

func (am *AuthMiddleware) hasPermission(permissions map[string]bool, required string) bool {
    // Check for wildcard permission
    if permissions["*"] {
        return true
    }

    // Check specific permission
    return permissions[required]
}

func (am *AuthMiddleware) hashAPIKey(apiKey string) string {
    hash := sha256.Sum256([]byte(apiKey))
    return hex.EncodeToString(hash[:])
}
