// auth-service/internal/services/tenant_service.go
package services

import (
    "context"
    "fmt"
    "time"

    "gorm.io/gorm"
)

type TenantService struct {
    db *gorm.DB
}

type Tenant struct {
    ID                  string    `json:"id" gorm:"type:uuid;primary_key"`
    Name                string    `json:"name" gorm:"not null"`
    Slug                string    `json:"slug" gorm:"uniqueIndex;not null"`
    PlanType            string    `json:"plan_type" gorm:"default:'starter'"`
    MaxUsers            int       `json:"max_users" gorm:"default:5"`
    MaxRequestsPerMonth int       `json:"max_requests_per_month" gorm:"default:10000"`
    BillingEmail        string    `json:"billing_email"`
    Status              string    `json:"status" gorm:"default:'active'"`
    CreatedAt           time.Time `json:"created_at"`
    UpdatedAt           time.Time `json:"updated_at"`
    TrialEndsAt         *time.Time `json:"trial_ends_at"`
    Settings            JSONMap   `json:"settings" gorm:"type:jsonb"`
}

type CreateTenantRequest struct {
    Name         string `json:"name" binding:"required"`
    Slug         string `json:"slug" binding:"required"`
    PlanType     string `json:"plan_type"`
    BillingEmail string `json:"billing_email" binding:"required,email"`
}

func NewTenantService(db *gorm.DB) *TenantService {
    return &TenantService{db: db}
}

// CreateTenant creates a new tenant with default settings
func (ts *TenantService) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*Tenant, error) {
    // Check if slug is available
    var existingTenant Tenant
    if err := ts.db.WithContext(ctx).Where("slug = ?", req.Slug).First(&existingTenant).Error; err == nil {
        return nil, fmt.Errorf("tenant slug already exists")
    }

    tenant := &Tenant{
        Name:         req.Name,
        Slug:         req.Slug,
        PlanType:     req.PlanType,
        BillingEmail: req.BillingEmail,
        Settings: JSONMap{
            "allowed_platforms": getAllowedPlatformsForPlan(req.PlanType),
            "max_concurrent_jobs": getMaxConcurrentJobsForPlan(req.PlanType),
            "data_retention_days": getDataRetentionForPlan(req.PlanType),
        },
    }

    if req.PlanType == "trial" {
        trialEnd := time.Now().Add(30 * 24 * time.Hour) // 30-day trial
        tenant.TrialEndsAt = &trialEnd
    }

    if err := ts.db.WithContext(ctx).Create(tenant).Error; err != nil {
        return nil, fmt.Errorf("failed to create tenant: %w", err)
    }

    // Create default roles for the tenant
    if err := ts.createDefaultRoles(ctx, tenant.ID); err != nil {
        return nil, fmt.Errorf("failed to create default roles: %w", err)
    }

    return tenant, nil
}

func (ts *TenantService) createDefaultRoles(ctx context.Context, tenantID string) error {
    defaultRoles := []Role{
        {
            TenantID:      tenantID,
            Name:          "admin",
            Permissions:   JSONMap{"*": true},
            IsSystemRole:  true,
        },
        {
            TenantID:     tenantID,
            Name:         "user",
            Permissions:  JSONMap{
                "phone_lookup:execute": true,
                "reports:read":         true,
                "exports:read":         true,
            },
            IsSystemRole: true,
        },
        {
            TenantID:     tenantID,
            Name:         "viewer", 
            Permissions:  JSONMap{
                "reports:read": true,
                "exports:read": true,
            },
            IsSystemRole: true,
        },
    }

    return ts.db.WithContext(ctx).Create(&defaultRoles).Error
}

// GetTenantBySlug retrieves tenant by slug
func (ts *TenantService) GetTenantBySlug(ctx context.Context, slug string) (*Tenant, error) {
    var tenant Tenant
    if err := ts.db.WithContext(ctx).Where("slug = ? AND status = 'active'", slug).First(&tenant).Error; err != nil {
        return nil, fmt.Errorf("tenant not found: %w", err)
    }
    return &tenant, nil
}

// GetTenantUsage returns current usage for a tenant
func (ts *TenantService) GetTenantUsage(ctx context.Context, tenantID string) (*TenantUsage, error) {
    currentMonth := time.Now().UTC().Format("2006-01-01")
    
    var usage TenantUsage
    err := ts.db.WithContext(ctx).
        Where("tenant_id = ? AND month = ?", tenantID, currentMonth).
        First(&usage).Error

    if err != nil {
        // Create usage record if it doesn't exist
        usage = TenantUsage{
            TenantID:   tenantID,
            Month:      currentMonth,
        }
        if createErr := ts.db.WithContext(ctx).Create(&usage).Error; createErr != nil {
            return nil, createErr
        }
    }

    return &usage, nil
}

// CheckRateLimit checks if tenant has exceeded rate limits
func (ts *TenantService) CheckRateLimit(ctx context.Context, tenantID string, endpoint string) (bool, error) {
    var tenant Tenant
    if err := ts.db.WithContext(ctx).First(&tenant, "id = ?", tenantID).Error; err != nil {
        return false, err
    }

    // Get current usage
    usage, err := ts.GetTenantUsage(ctx, tenantID)
    if err != nil {
        return false, err
    }

    // Check monthly request limit
    if usage.RequestsCount >= tenant.MaxRequestsPerMonth {
        return false, fmt.Errorf("monthly request limit exceeded")
    }

    // Check plan-specific rate limits
    planLimits := getRateLimitsForPlan(tenant.PlanType)
    endpointLimit, exists := planLimits[endpoint]
    if !exists {
        endpointLimit = planLimits["default"]
    }

    // Here you would check per-minute/second limits using Redis
    // This is a simplified version
    return true, nil
}

// IncrementUsage increments usage counters for a tenant
func (ts *TenantService) IncrementUsage(ctx context.Context, tenantID string, usageType string, count int64) error {
    currentMonth := time.Now().UTC().Format("2006-01-01")
    
    return ts.db.WithContext(ctx).Model(&TenantUsage{}).
        Where("tenant_id = ? AND month = ?", tenantID, currentMonth).
        UpdateColumn(usageType, gorm.Expr(fmt.Sprintf("%s + ?", usageType), count)).
        Error
}

func getAllowedPlatformsForPlan(planType string) []string {
    switch planType {
    case "enterprise":
        return []string{"all"}
    case "professional":
        return []string{"facebook", "instagram", "linkedin", "twitter", "whatsapp"}
    case "starter":
        return []string{"facebook", "instagram"}
    default:
        return []string{"facebook"}
    }
}

func getMaxConcurrentJobsForPlan(planType string) int {
    switch planType {
    case "enterprise":
        return 100
    case "professional":
        return 20
    case "starter":
        return 5
    default:
        return 1
    }
}

func getRateLimitsForPlan(planType string) map[string]int {
    switch planType {
    case "enterprise":
        return map[string]int{
            "default":          1000, // requests per minute
            "phone_lookup":     500,
            "email_discovery":  200,
            "bulk_operations":  100,
        }
    case "professional":
        return map[string]int{
            "default":          500,
            "phone_lookup":     200,
            "email_discovery":  100,
            "bulk_operations":  50,
        }
    case "starter":
        return map[string]int{
            "default":          100,
            "phone_lookup":     50,
            "email_discovery":  20,
            "bulk_operations":  10,
        }
    default:
        return map[string]int{"default": 10}
    }
}
