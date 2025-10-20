// admin-backend/internal/handlers/bulk_jobs.go
package handlers

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"
)

type BulkJobHandler struct {
    jobService     *service.JobService
    fileProcessor  *FileProcessor
}

func NewBulkJobHandler() *BulkJobHandler {
    return &BulkJobHandler{
        jobService:    service.NewJobService(),
        fileProcessor: NewFileProcessor(),
    }
}

// ProcessBulkUpload handles large file uploads for bulk jobs
func (h *BulkJobHandler) ProcessBulkUpload(c *gin.Context) {
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File upload required"})
        return
    }
    defer file.Close()

    // Parse upload options
    platforms := strings.Split(c.PostForm("platforms"), ",")
    priority := c.PostForm("priority")
    jobName := c.PostForm("job_name")

    // Process file based on type
    var phoneNumbers []string
    switch {
    case strings.HasSuffix(header.Filename, ".csv"):
        phoneNumbers, err = h.processCSV(file)
    case strings.HasSuffix(header.Filename, ".txt"):
        phoneNumbers, err = h.processTXT(file)
    case strings.HasSuffix(header.Filename, ".xlsx"):
        phoneNumbers, err = h.processExcel(file)
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file format"})
        return
    }

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File processing failed: %v", err)})
        return
    }

    // Create bulk job
    jobID, err := h.jobService.CreateBulkJob(service.BulkJobRequest{
        Name:         jobName,
        PhoneNumbers: phoneNumbers,
        Platforms:    platforms,
        Priority:     priority,
        CreatedBy:    c.GetString("user_id"),
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bulk job"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "job_id": jobID,
        "total_numbers": len(phoneNumbers),
        "status": "processing",
    })
}

func (h *BulkJobHandler) processCSV(file io.Reader) ([]string, error) {
    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    var phoneNumbers []string
    for i, record := range records {
        if i == 0 && h.isHeaderRow(record) {
            continue // Skip header row
        }
        
        if len(record) > 0 {
            phoneNumbers = append(phoneNumbers, strings.TrimSpace(record[0]))
        }
    }

    return phoneNumbers, nil
}

func (h *BulkJobHandler) processTXT(file io.Reader) ([]string, error) {
    content, err := io.ReadAll(file)
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(content), "\n")
    var phoneNumbers []string
    
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line != "" {
            phoneNumbers = append(phoneNumbers, line)
        }
    }

    return phoneNumbers, nil
}

// GetBulkJobStatus returns status of a bulk job
func (h *BulkJobHandler) GetBulkJobStatus(c *gin.Context) {
    jobID := c.Param("id")
    
    status, err := h.jobService.GetBulkJobStatus(jobID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
        return
    }

    c.JSON(http.StatusOK, status)
}

// ListBulkJobs returns paginated list of bulk jobs
func (h *BulkJobHandler) ListBulkJobs(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
    status := c.Query("status")

    jobs, total, err := h.jobService.ListBulkJobs(page, limit, status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "jobs": jobs,
        "pagination": gin.H{
            "page":  page,
            "limit": limit,
            "total": total,
            "pages": (total + limit - 1) / limit,
        },
    })
}
