// 01_orchestrator/internal/service/job_service.go
package service

import (
    "encoding/json"
    "fmt"
    "time"

    "secure-iran-intel/01_orchestrator/internal/repository"
    "secure-iran-intel/pkg/normalizer"
)

type JobService struct {
    jobRepo      *repository.JobRepository
    proxyService *ProxyService
    normalizer   *normalizer.PhoneNormalizer
    mqProducer   *MQProducer
}

func NewJobService(jobRepo *repository.JobRepository, proxyService *ProxyService) *JobService {
    return &JobService{
        jobRepo:      jobRepo,
        proxyService: proxyService,
        normalizer:   normalizer.NewPhoneNormalizer("IR"),
        mqProducer:   NewMQProducer(),
    }
}

func (js *JobService) CreateIntelligenceJob(phoneNumbers []string, platforms []string, priority string, options map[string]interface{}) (string, error) {
    // Step 1: Normalize phone numbers
    normalizedNumbers, errors := js.normalizer.NormalizeBatch(phoneNumbers, "")
    if len(normalizedNumbers) == 0 {
        return "", fmt.Errorf("no valid phone numbers after normalization")
    }

    // Step 2: Create job record
    jobID := generateJobID()
    job := &repository.Job{
        ID:           jobID,
        PhoneNumbers: phoneNumbers,
        Platforms:    platforms,
        Priority:     priority,
        Status:       "queued",
        CreatedAt:    time.Now(),
        Options:      options,
    }

    if err := js.jobRepo.Create(job); err != nil {
        return "", fmt.Errorf("failed to create job record: %w", err)
    }

    // Step 3: Create individual tasks for each phone-platform combination
    tasks := js.createTasks(jobID, normalizedNumbers, platforms, priority)

    // Step 4: Send tasks to message queue
    if err := js.mqProducer.SendTasks(tasks); err != nil {
        return "", fmt.Errorf("failed to queue tasks: %w", err)
    }

    // Step 5: Update job status
    job.Status = "processing"
    js.jobRepo.Update(job)

    return jobID, nil
}

func (js *JobService) createTasks(jobID string, numbers map[string]*normalizer.NormalizedPhone, platforms []string, priority string) []*Task {
    var tasks []*Task
    
    for original, normalized := range numbers {
        for _, platform := range platforms {
            task := &Task{
                ID:           generateTaskID(),
                JobID:        jobID,
                PhoneNumber:  original,
                Normalized:   normalized.Normalized,
                Platform:     platform,
                Priority:     priority,
                Status:       "pending",
                CreatedAt:    time.Now(),
                // Each task gets a different proxy
                ProxyConfig:  js.proxyService.GetProxyConfigForTask(platform, normalized.Normalized),
            }
            tasks = append(tasks, task)
        }
    }
    
    return tasks
}

func (js *JobService) GetJobStatus(jobID string) (*repository.JobStatus, error) {
    return js.jobRepo.GetStatus(jobID)
}

func (js *JobService) GetProxyHealthStats() *ProxyHealthStats {
    return js.proxyService.GetHealthStats()
}

type Task struct {
    ID          string                      `json:"id"`
    JobID       string                      `json:"job_id"`
    PhoneNumber string                      `json:"phone_number"`
    Normalized  string                      `json:"normalized"`
    Platform    string                      `json:"platform"`
    Priority    string                      `json:"priority"`
    Status      string                      `json:"status"`
    CreatedAt   time.Time                   `json:"created_at"`
    ProxyConfig *ProxyConfig                `json:"proxy_config"`
}

type ProxyConfig struct {
    ProxyID    string `json:"proxy_id"`
    IP         string `json:"ip"`
    Port       int    `json:"port"`
    Type       string `json:"type"`
    Username   string `json:"username,omitempty"`
    Password   string `json:"password,omitempty"`
}
