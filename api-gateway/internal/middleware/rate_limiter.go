// api-gateway/internal/middleware/rate_limiter.go
package middleware

import (
    "fmt"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "secure-iran-intel/api-gateway/internal/services"
)

type RateLimitMiddleware struct {
    rateLimitService *services.RateLimitService
    tenantService    *services.TenantService
}

func NewRateLimitMiddleware(rateLimitService *services.RateLimitService, tenantService *services.TenantService) *RateLimitMiddleware {
    return &RateLimitMiddleware{
        rateLimitService: rateLimitService,
        tenantService:    tenantService,
    }
}

// RateLimitByTenant applies rate limiting based on tenant
func (rlm *RateLimitMiddleware) RateLimitByTenant() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := c.Request.Context()
        
        // Extract tenant ID from context
        tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
            c.Abort()
            return
        }

        // Extract user ID if available
        userID, _ := ctx.Value(middleware.UserIDKey).(string)
        
        // Get client IP
        clientIP := c.ClientIP()
        
        // Get endpoint
        endpoint := c.FullPath()
        
        // Check multi-dimensional rate limits
        result, err := rlm.rateLimitService.MultiDimensionalRateLimit(
            ctx, tenantID, userID, clientIP, endpoint)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit service error"})
            c.Abort()
            return
        }

        if !result.Allowed {
            c.Header("X-RateLimit-Limit", strconv.FormatInt(result.Limit, 10))
            c.Header("X-RateLimit-Remaining", strconv.FormatInt(result.Remaining, 10))
            c.Header("X-RateLimit-Reset", strconv.FormatInt(result.ResetTime, 10))
            c.Header("Retry-After", strconv.FormatInt(result.RetryAfter, 10))
            
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": result.RetryAfter,
                "limit": result.Limit,
                "remaining": result.Remaining,
            })
            c.Abort()
            return
        }

        // Set rate limit headers
        c.Header("X-RateLimit-Limit", strconv.FormatInt(result.Limit, 10))
        c.Header("X-RateLimit-Remaining", strconv.FormatInt(result.Remaining, 10))
        c.Header("X-RateLimit-Reset", strconv.FormatInt(result.ResetTime, 10))

        c.Next()
    }
}

// AdaptiveRateLimit applies intelligent rate limiting
func (rlm *RateLimitMiddleware) AdaptiveRateLimit() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := c.Request.Context()
        
        tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
            c.Abort()
            return
        }

        clientIP := c.ClientIP()
        endpoint := c.FullPath()
        
        // Estimate request size
        requestSize := rlm.estimateRequestSize(c)
        
        // Apply adaptive rate limiting
        result, err := rlm.rateLimitService.AdaptiveRateLimit(
            ctx, tenantID, clientIP, endpoint, requestSize)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit service error"})
            c.Abort()
            return
        }

        if !result.Allowed {
            c.Header("X-RateLimit-Limit", strconv.FormatInt(result.Limit, 10))
            c.Header("X-RateLimit-Remaining", strconv.FormatInt(result.Remaining, 10))
            c.Header("X-RateLimit-Reset", strconv.FormatInt(result.ResetTime, 10))
            
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": result.RetryAfter,
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// QuotaMiddleware checks monthly usage quotas
func (rlm *RateLimitMiddleware) QuotaMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := c.Request.Context()
        
        tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
            c.Abort()
            return
        }

        // Check if tenant has exceeded monthly quota
        allowed, err := rlm.tenantService.CheckRateLimit(ctx, tenantID, c.FullPath())
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Quota service error"})
            c.Abort()
            return
        }

        if !allowed {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Monthly request quota exceeded",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

func (rlm *RateLimitMiddleware) estimateRequestSize(c *gin.Context) int {
    size := 0
    
    // Estimate based on headers
    for key, values := range c.Request.Header {
        size += len(key)
        for _, value := range values {
            size += len(value)
        }
    }
    
    // Estimate based on query parameters
    size += len(c.Request.URL.RawQuery)
    
    // Estimate based on route parameters
    for _, param := range c.Params {
        size += len(param.Value)
    }
    
    return size
}
