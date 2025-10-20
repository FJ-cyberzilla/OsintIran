// 01_orchestrator/internal/handler/http_handler.go
package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "secure-iran-intel/01_orchestrator/internal/service"
)

type HTTPHandler struct {
    jobService *service.JobService
}

func NewHTTPHandler(jobService *service.JobService) *HTTPHandler {
    return &HTTPHandler{
        jobService: jobService,
    }
}

type CreateJobRequest struct {
    PhoneNumbers []string          `json:"phone_numbers" binding:"required"`
    Platforms    []string          `json:"platforms" binding:"required"`
    Priority     string            `json:"priority,omitempty"`
    Options      map[string]interface{} `json:"options,omitempty"`
}

func (h *HTTPHandler) CreateJob(c *gin.Context) {
    var req CreateJobRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    jobID, err := h.jobService.CreateIntelligenceJob(req.PhoneNumbers, req.Platforms, req.Priority, req.Options)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "job_id": jobID,
        "status": "created",
        "message": "Job queued for processing",
    })
}

func (h *HTTPHandler) CreateBatchJobs(c *gin.Context) {
    var requests []CreateJobRequest
    if err := c.ShouldBindJSON(&requests); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var jobIDs []string
    for _, req := range requests {
        jobID, err := h.jobService.CreateIntelligenceJob(req.PhoneNumbers, req.Platforms, req.Priority, req.Options)
        if err != nil {
            // Continue with other jobs even if some fail
            continue
        }
        jobIDs = append(jobIDs, jobID)
    }

    c.JSON(http.StatusOK, gin.H{
        "job_ids": jobIDs,
        "total_created": len(jobIDs),
        "total_failed": len(requests) - len(jobIDs),
    })
}

func (h *HTTPHandler) GetJobStatus(c *gin.Context) {
    jobID := c.Param("id")
    
    status, err := h.jobService.GetJobStatus(jobID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
        return
    }

    c.JSON(http.StatusOK, status)
}

func (h *HTTPHandler) GetProxyHealth(c *gin.Context) {
    health := h.jobService.GetProxyHealthStats()
    
    c.JSON(http.StatusOK, gin.H{
        "total_proxies": health.TotalProxies,
        "healthy_proxies": health.HealthyProxies,
        "health_percentage": health.HealthPercentage,
        "avg_response_time": health.AvgResponseTime,
    })
}
