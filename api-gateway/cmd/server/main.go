// api-gateway/cmd/server/main.go
package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"

    "secure-iran-intel/api-gateway/internal/middleware"
    "secure-iran-intel/api-gateway/internal/services"
    "secure-iran-intel/auth-service/internal/services"
)

func main() {
    // Initialize dependencies
    db := initDatabase()
    redisClient := initRedis()
    
    // Initialize services
    tenantService := auth_services.NewTenantService(db)
    rateLimitService := services.NewRateLimitService(redisClient, tenantService)
    quotaService := services.NewQuotaService(db, tenantService)
    
    // Initialize middleware
    authMiddleware := middleware.NewAuthMiddleware(
        os.Getenv("JWT_SECRET"),
        tenantService,
        auth_services.NewUserService(db),
    )
    rateLimitMiddleware := middleware.NewRateLimitMiddleware(rateLimitService, tenantService)

    // Create router
    router := gin.Default()

    // Public routes (authentication)
    public := router.Group("/api/v1/auth")
    {
        public.POST("/login", authHandler.Login)
        public.POST("/register", authHandler.Register)
        public.POST("/refresh", authHandler.RefreshToken)
    }

    // Tenant-aware API routes
    api := router.Group("/api/v1")
    {
        // Apply authentication to all API routes
        api.Use(authMiddleware.JWTAuthMiddleware())
        
        // Apply rate limiting based on tenant plan
        api.Use(rateLimitMiddleware.RateLimitByTenant())
        api.Use(rateLimitMiddleware.QuotaMiddleware())

        // Phone intelligence endpoints
        intelligence := api.Group("/intelligence")
        {
            intelligence.POST("/phone-lookup", phoneHandler.LookupPhone)
            intelligence.POST("/email-discovery", emailHandler.DiscoverEmails)
            intelligence.POST("/bulk-operations", bulkHandler.ProcessBulk)
            intelligence.GET("/reports/:id", reportHandler.GetReport)
        }

        // Admin endpoints (require admin permissions)
        admin := api.Group("/admin")
        admin.Use(authMiddleware.PermissionMiddleware("admin"))
        {
            admin.GET("/users", adminHandler.GetUsers)
            admin.POST("/users", adminHandler.CreateUser)
            admin.GET("/tenants", adminHandler.GetTenants)
            admin.POST("/tenants", adminHandler.CreateTenant)
            admin.GET("/usage", adminHandler.GetUsage)
        }
    }

    // API key routes (different rate limits)
    apiKeyRoutes := router.Group("/api/v1")
    {
        apiKeyRoutes.Use(authMiddleware.APIKeyAuthMiddleware())
        apiKeyRoutes.Use(rateLimitMiddleware.AdaptiveRateLimit())
        
        apiKeyRoutes.POST("/lookup", phoneHandler.LookupPhone)
        apiKeyRoutes.GET("/status/:id", jobHandler.GetJobStatus)
    }

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("🚀 Enhanced API Gateway started on :%s", port)
    log.Printf("📊 Multi-tenant support: ENABLED")
    log.Printf("🚦 Advanced rate limiting: ENABLED")
    
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func initDatabase() *gorm.DB {
    // Database initialization
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"), 
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    
    return db
}

func initRedis() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_URL"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })
}
