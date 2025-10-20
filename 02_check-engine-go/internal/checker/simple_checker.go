// 02_check-engine-go/internal/checker/simple_checker.go
package checker

import (
    "fmt"
    "io"
    "net/http"
    "time"

    "secure-iran-intel/02_check-engine-go/internal/proxy"
)

type SimpleChecker struct {
    rotationEngine *proxy.ProxyRotationEngine
    httpClient     *http.Client
}

func NewSimpleChecker(rotationEngine *proxy.ProxyRotationEngine) *SimpleChecker {
    return &SimpleChecker{
        rotationEngine: rotationEngine,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (sc *SimpleChecker) ProcessTask(task *queue.Task) (*TaskResult, error) {
    var result *TaskResult
    
    // Each request uses a different proxy from the 23,000 pool
    proxy, err := sc.rotationEngine.GetNextProxy(task.Platform)
    if err != nil {
        return nil, fmt.Errorf("failed to get proxy: %w", err)
    }
    
    defer sc.rotationEngine.ReleaseProxy(proxy)

    // Execute check with the selected proxy
    switch task.Platform {
    case "facebook":
        result = sc.checkFacebook(task, proxy)
    case "linkedin":
        result = sc.checkLinkedIn(task, proxy)
    case "whatsapp":
        result = sc.checkWhatsApp(task, proxy)
    case "twitter":
        result = sc.checkTwitter(task, proxy)
    default:
        result = sc.genericCheck(task, proxy)
    }
    
    return result, nil
}

func (sc *SimpleChecker) checkFacebook(task *queue.Task, proxy *proxy.Proxy) *TaskResult {
    // Create request with proxy
    req, err := sc.createRequestWithProxy("https://facebook.com", proxy)
    if err != nil {
        return &TaskResult{
            Success: false,
            Error:   err.Error(),
        }
    }
    
    // Add Facebook-specific headers
    req.Header.Add("User-Agent", getRandomUserAgent())
    req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    
    return sc.executeCheck(req, task, "facebook")
}

func (sc *SimpleChecker) checkLinkedIn(task *queue.Task, proxy *proxy.Proxy) *TaskResult {
    req, err := sc.createRequestWithProxy("https://linkedin.com", proxy)
    if err != nil {
        return &TaskResult{
            Success: false,
            Error:   err.Error(),
        }
    }
    
    // LinkedIn-specific headers
    req.Header.Add("User-Agent", getProfessionalUserAgent())
    req.Header.Add("Accept-Language", "en-US,en;q=0.9")
    
    return sc.executeCheck(req, task, "linkedin")
}

func (sc *SimpleChecker) checkWhatsApp(task *queue.Task, proxy *proxy.Proxy) *TaskResult {
    req, err := sc.createRequestWithProxy("https://web.whatsapp.com", proxy)
    if err != nil {
        return &TaskResult{
            Success: false,
            Error:   err.Error(),
        }
    }
    
    // WhatsApp-specific headers
    req.Header.Add("User-Agent", getMobileUserAgent())
    
    return sc.executeCheck(req, task, "whatsapp")
}

func (sc *SimpleChecker) createRequestWithProxy(url string, proxy *proxy.Proxy) (*http.Request, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    // Set proxy for HTTP client
    proxyURL := fmt.Sprintf("http://%s:%d", proxy.IP, proxy.Port)
    sc.httpClient.Transport = &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    }
    
    return req, nil
}

func (sc *SimpleChecker) executeCheck(req *http.Request, task *queue.Task, platform string) *TaskResult {
    startTime := time.Now()
    
    resp, err := sc.httpClient.Do(req)
    if err != nil {
        return &TaskResult{
            Success:    false,
            Error:      err.Error(),
            Platform:   platform,
            PhoneNumber: task.PhoneNumber,
            Duration:   time.Since(startTime),
        }
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return &TaskResult{
            Success:    false,
            Error:      err.Error(),
            Platform:   platform,
            PhoneNumber: task.PhoneNumber,
            Duration:   time.Since(startTime),
        }
    }
    
    // Analyze response for phone number presence
    found := sc.analyzeResponseForPhone(body, task.Normalized)
    
    return &TaskResult{
        Success:     true,
        Found:       found,
        Platform:    platform,
        PhoneNumber: task.PhoneNumber,
        Duration:    time.Since(startTime),
        StatusCode:  resp.StatusCode,
    }
}

func (sc *SimpleChecker) analyzeResponseForPhone(body []byte, phoneNumber string) bool {
    // Implement phone number detection logic in response
    // This would vary by platform
    return false // Simplified for example
}
