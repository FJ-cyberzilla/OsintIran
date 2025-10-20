// api-gateway/internal/handlers/job_creator.go
package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "secure-iran-intel/pkg/normalizer"
)

type JobCreationHandler struct {
    normalizer *normalizer.PhoneNormalizer
    queue      JobQueue
}

func NewJobCreationHandler(queue JobQueue) *JobCreationHandler {
    return &JobCreationHandler{
        normalizer: normalizer.NewPhoneNormalizer("IR"),
        queue:      queue,
    }
}

type IntelligenceJobRequest struct {
    PhoneNumbers []string `json:"phone_numbers" binding:"required"`
    CountryHint  string   `json:"country_hint,omitempty"`
    Platforms    []string `json:"platforms" binding:"required"`
    Priority     string   `json:"priority,omitempty"`
    Options      JobOptions `json:"options,omitempty"`
}

type NormalizedJob struct {
    JobID        string                     `json:"job_id"`
    Original     []string                   `json:"original_numbers"`
    Normalized   map[string]normalizer.NormalizedPhone `json:"normalized_numbers"`
    Platforms    []string                   `json:"platforms"`
    Priority     string                     `json:"priority"`
    CreatedAt    string                     `json:"created_at"`
}

func (jch *JobCreationHandler) CreateIntelligenceJob(c *gin.Context) {
    var req IntelligenceJobRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }
    
    // Normalize all phone numbers
    normalizedNumbers, errors := jch.normalizer.NormalizeBatch(req.PhoneNumbers, req.CountryHint)
    
    // Filter out invalid numbers if strict mode
    validNumbers := make(map[string]normalizer.NormalizedPhone)
    for original, normalized := range normalizedNumbers {
        if normalized.IsValid || normalized.IsPossible {
            validNumbers[original] = *normalized
        }
    }
    
    if len(validNumbers) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "No valid phone numbers found after normalization",
            "details": errors,
        })
        return
    }
    
    // Create normalized job
    job := NormalizedJob{
        JobID:      generateJobID(),
        Original:   req.PhoneNumbers,
        Normalized: validNumbers,
        Platforms:  req.Platforms,
        Priority:   req.Priority,
        CreatedAt:  time.Now().Format(time.RFC3339),
    }
    
    // Submit job to queue
    if err := jch.queue.SubmitJob(job); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to submit job to queue",
        })
        return
    }
    
    response := gin.H{
        "success": true,
        "job_id":  job.JobID,
        "data": gin.H{
            "total_submitted":   len(req.PhoneNumbers),
            "valid_numbers":     len(validNumbers),
            "invalid_numbers":   len(req.PhoneNumbers) - len(validNumbers),
            "normalized_format": "E.164",
        },
    }
    
    if len(errors) > 0 {
        response["normalization_errors"] = errors
    }
    
    c.JSON(http.StatusOK, response)
}
