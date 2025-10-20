// api-gateway/internal/services/rate_limit_service.go
package services

import (
    "context"
    "fmt"
    "strconv"
    "strings"
    "time"

    "github.com/go-redis/redis/v8"
)

type RateLimitService struct {
    redisClient *redis.Client
    tenantService *TenantService
}

type RateLimitConfig struct {
    RequestsPerMinute int           `json:"requests_per_minute"`
    Burst            int           `json:"burst"`
    Window           time.Duration `json:"window"`
    Scope            string        `json:"scope"` // "user", "tenant", "ip", "endpoint"
}

type RateLimitResult struct {
    Allowed    bool  `json:"allowed"`
    Remaining  int64 `json:"remaining"`
    ResetTime  int64 `json:"reset_time"`
    Limit      int64 `json:"limit"`
    RetryAfter int64 `json:"retry_after,omitempty"`
}

func NewRateLimitService(redisClient *redis.Client, tenantService *TenantService) *RateLimitService {
    return &RateLimitService{
        redisClient:   redisClient,
        tenantService: tenantService,
    }
}

// CheckRateLimit performs advanced rate limiting with multiple dimensions
func (rls *RateLimitService) CheckRateLimit(ctx context.Context, key string, config *RateLimitConfig) (*RateLimitResult, error) {
    now := time.Now().Unix()
    window := int64(config.Window.Seconds())
    
    // Redis key for this rate limit window
    redisKey := fmt.Sprintf("rate_limit:%s:%d", key, now/window)
    
    // Use Redis pipeline for atomic operations
    pipe := rls.redisClient.Pipeline()
    
    // Increment counter
    incr := pipe.Incr(ctx, redisKey)
    
    // Set expiry if this is a new key
    pipe.Expire(ctx, redisKey, config.Window)
    
    // Get current count
    current := pipe.Get(ctx, redisKey)
    
    // Execute pipeline
    _, err := pipe.Exec(ctx)
    if err != nil && err != redis.Nil {
        return nil, fmt.Errorf("redis error: %w", err)
    }
    
    currentCount, _ := strconv.ParseInt(current.Val(), 10, 64)
    
    // Check if limit exceeded
    allowed := currentCount <= int64(config.RequestsPerMinute)
    
    result := &RateLimitResult{
        Allowed:   allowed,
        Remaining: max(0, int64(config.RequestsPerMinute)-currentCount),
        ResetTime: ((now/window)+1)*window,
        Limit:     int64(config.RequestsPerMinute),
    }
    
    if !allowed {
        result.RetryAfter = result.ResetTime - now
    }
    
    return result, nil
}

// MultiDimensionalRateLimit checks rate limits across multiple dimensions
func (rls *RateLimitService) MultiDimensionalRateLimit(
    ctx context.Context,
    tenantID string,
    userID string,
    clientIP string,
    endpoint string,
) (*RateLimitResult, error) {
    
    // Get tenant rate limit configuration
    tenantConfig, err := rls.getTenantRateLimitConfig(tenantID)
    if err != nil {
        return nil, err
    }
    
    var results []*RateLimitResult
    
    // Check tenant-level rate limit
    tenantResult, err := rls.CheckRateLimit(ctx, 
        fmt.Sprintf("tenant:%s", tenantID),
        &RateLimitConfig{
            RequestsPerMinute: tenantConfig.RequestsPerMinute,
            Window: time.Minute,
        })
    if err != nil {
        return nil, err
    }
    results = append(results, tenantResult)
    
    // Check user-level rate limit (if user is authenticated)
    if userID != "" {
        userResult, err := rls.CheckRateLimit(ctx,
            fmt.Sprintf("user:%s:%s", tenantID, userID),
            &RateLimitConfig{
                RequestsPerMinute: tenantConfig.UserLimit,
                Window: time.Minute,
            })
        if err != nil {
            return nil, err
        }
        results = append(results, userResult)
    }
    
    // Check IP-level rate limit
    ipResult, err := rls.CheckRateLimit(ctx,
        fmt.Sprintf("ip:%s", clientIP),
        &RateLimitConfig{
            RequestsPerMinute: tenantConfig.IPLimit,
            Window: time.Minute,
        })
    if err != nil {
        return nil, err
    }
    results = append(results, ipResult)
    
    // Check endpoint-specific rate limit
    endpointLimit := tenantConfig.EndpointLimits[endpoint]
    if endpointLimit == 0 {
        endpointLimit = tenantConfig.EndpointLimits["default"]
    }
    
    endpointResult, err := rls.CheckRateLimit(ctx,
        fmt.Sprintf("endpoint:%s:%s", tenantID, endpoint),
        &RateLimitConfig{
            RequestsPerMinute: endpointLimit,
            Window: time.Minute,
        })
    if err != nil {
        return nil, err
    }
    results = append(results, endpointResult)
    
    // Return the most restrictive result
    return rls.getMostRestrictiveResult(results), nil
}

// AdaptiveRateLimit adjusts limits based on client behavior
func (rls *RateLimitService) AdaptiveRateLimit(
    ctx context.Context,
    tenantID string,
    clientIP string,
    endpoint string,
    requestSize int,
) (*RateLimitResult, error) {
    
    // Calculate request cost (larger requests cost more)
    requestCost := rls.calculateRequestCost(endpoint, requestSize)
    
    // Get client behavior score
    behaviorScore, err := rls.getClientBehaviorScore(ctx, tenantID, clientIP)
    if err != nil {
        return nil, err
    }
    
    // Adjust limits based on behavior score
    baseLimit := rls.getBaseLimit(tenantID, endpoint)
    adjustedLimit := int(float64(baseLimit) * behaviorScore)
    
    // Use token bucket algorithm for more flexible rate limiting
    bucketKey := fmt.Sprintf("token_bucket:%s:%s:%s", tenantID, clientIP, endpoint)
    
    return rls.tokenBucketRateLimit(ctx, bucketKey, adjustedLimit, requestCost)
}

// Token bucket rate limiting algorithm
func (rls *RateLimitService) tokenBucketRateLimit(
    ctx context.Context,
    key string,
    limit int,
    cost int,
) (*RateLimitResult, error) {
    
    now := time.Now().Unix()
    bucketKey := fmt.Sprintf("%s:bucket", key)
    lastUpdateKey := fmt.Sprintf("%s:last_update", key)
    
    pipe := rls.redisClient.Pipeline()
    
    // Get current token count and last update time
    tokenCmd := pipe.Get(ctx, bucketKey)
    lastUpdateCmd := pipe.Get(ctx, lastUpdateKey)
    
    if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
        return nil, err
    }
    
    tokens, _ := strconv.ParseFloat(tokenCmd.Val(), 64)
    lastUpdate, _ := strconv.ParseInt(lastUpdateCmd.Val(), 10, 64)
    
    // Calculate new token count based on time elapsed
    if lastUpdate > 0 {
        elapsed := float64(now - lastUpdate)
        refillRate := float64(limit) / 60.0 // tokens per second
        tokens = min(float64(limit), tokens + elapsed * refillRate)
    } else {
        tokens = float64(limit) // Start with full bucket
    }
    
    // Check if enough tokens available
    allowed := tokens >= float64(cost)
    
    if allowed {
        tokens -= float64(cost)
    }
    
    // Update bucket state
    pipe.Set(ctx, bucketKey, tokens, time.Minute)
    pipe.Set(ctx, lastUpdateKey, now, time.Minute)
    pipe.Exec(ctx)
    
    return &RateLimitResult{
        Allowed:   allowed,
        Remaining: int64(tokens),
        Limit:     int64(limit),
        ResetTime: now + 60, // Reset in 60 seconds
    }, nil
}

func (rls *RateLimitService) calculateRequestCost(endpoint string, size int) int {
    baseCosts := map[string]int{
        "phone_lookup":     1,
        "email_discovery":  3,
        "bulk_operation":   10,
        "report_generation": 5,
    }
    
    baseCost := baseCosts[endpoint]
    if baseCost == 0 {
        baseCost = 1
    }
    
    // Adjust cost based on request size
    sizeMultiplier := max(1, size/1024) // 1 cost per KB
    
    return baseCost * sizeMultiplier
}

func (rls *RateLimitService) getClientBehaviorScore(ctx context.Context, tenantID, clientIP string) (float64, error) {
    // Analyze client behavior: success rate, error rate, etc.
    // Return score between 0.5 (suspicious) and 1.5 (excellent)
    
    behaviorKey := fmt.Sprintf("behavior:%s:%s", tenantID, clientIP)
    
    // This would analyze historical request patterns
    // For now, return a default score
    return 1.0, nil
}

func (rls *RateLimitService) getMostRestrictiveResult(results []*RateLimitResult) *RateLimitResult {
    if len(results) == 0 {
        return &RateLimitResult{Allowed: true}
    }
    
    mostRestrictive := results[0]
    for _, result := range results[1:] {
        if !result.Allowed || result.Remaining < mostRestrictive.Remaining {
            mostRestrictive = result
        }
    }
    
    return mostRestrictive
}
