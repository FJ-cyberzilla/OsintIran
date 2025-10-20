// check-engine-go/internal/checker/secure_http.go
package checker

import (
    "database/sql"
    "fmt"
    "net/http"
    "time"

    "secure-iran-intel/pkg/security/sql_injection"
)

type SecureHTTPChecker struct {
    sqlProtector *sql_injection.SQLInjectionProtector
    client       *http.Client
    db           *sql_injection.SecureDB
}

func NewSecureHTTPChecker(sqlProtector *sql_injection.SQLInjectionProtector, cfg *config.Config) *SecureHTTPChecker {
    return &SecureHTTPChecker{
        sqlProtector: sqlProtector,
        client: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 10,
                IdleConnTimeout:     90 * time.Second,
            },
        },
        db: sql_injection.NewSecureDB(cfg.DB),
    }
}

func (shc *SecureHTTPChecker) ProcessSecureTask(task interface{}) error {
    // Validate task structure
    if err := shc.validateTask(task); err != nil {
        return fmt.Errorf("task validation failed: %w", err)
    }

    // Secure database operations
    result, err := shc.db.SecureQuery(
        "SELECT * FROM intelligence_data WHERE phone_number = ? AND platform = ?",
        task.PhoneNumber, task.Platform,
    )
    if err != nil {
        return fmt.Errorf("secure query failed: %w", err)
    }
    defer result.Close()

    // Secure HTTP request
    resp, err := shc.secureHTTPRequest(task.URL, task.Headers)
    if err != nil {
        return fmt.Errorf("secure HTTP request failed: %w", err)
    }
    defer resp.Body.Close()

    // Process response with security checks
    return shc.processSecureResponse(resp, task)
}

func (shc *SecureHTTPChecker) secureHTTPRequest(url string, headers map[string]string) (*http.Response, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    // Add security headers
    req.Header.Add("User-Agent", "SecureIntelBot/1.0")
    req.Header.Add("X-Security-Token", generateSecurityToken())
    
    for key, value := range headers {
        // Validate header values for injection attempts
        if shc.sqlProtector.ContainsSQLInjection(value) {
            return nil, fmt.Errorf("potential SQL injection in header: %s", key)
        }
        req.Header.Add(key, value)
    }

    return shc.client.Do(req)
}
