// 01_orchestrator/cmd/main.go
package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "secure-iran-intel/01_orchestrator/internal/handler"
    "secure-iran-intel/01_orchestrator/internal/service"
    "secure-iran-intel/01_orchestrator/internal/repository"
)

func main() {
    // Initialize dependencies
    jobRepo := repository.NewJobRepository()
    proxyService := service.NewProxyService()
    jobService := service.NewJobService(jobRepo, proxyService)
    httpHandler := handler.NewHTTPHandler(jobService)

    // Create router
    router := gin.Default()

    // Routes
    api := router.Group("/api/v1")
    {
        api.POST("/jobs", httpHandler.CreateJob)
        api.GET("/jobs/:id", httpHandler.GetJobStatus)
        api.POST("/batch", httpHandler.CreateBatchJobs)
        api.GET("/proxies/health", httpHandler.GetProxyHealth)
    }

    // Start server
    log.Println("🚀 Orchestrator service started on :8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
