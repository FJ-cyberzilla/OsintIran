// 04_proxy-pool/internal/proxy/rotation_engine.go
package proxy

import (
    "crypto/rand"
    "fmt"
    "math/big"
    "sync"
    "time"
)

type RotationStrategy int

const (
    RoundRobin RotationStrategy = iota
    Random
    Weighted
    StickySession
    Geographic
)

type ProxyRotationEngine struct {
    proxies          []*Proxy
    healthyProxies   []*Proxy
    rotationStrategy RotationStrategy
    currentIndex     int
    mu               sync.RWMutex
    stats            *RotationStats
    geoLocator       *GeoLocator
}

type Proxy struct {
    ID            string    `json:"id"`
    IP            string    `json:"ip"`
    Port          int       `json:"port"`
    Type          ProxyType `json:"type"` // HTTP, HTTPS, SOCKS5
    Provider      string    `json:"provider"`
    Location      string    `json:"location"` // Country/City
    ISP           string    `json:"isp"`
    Speed         int       `json:"speed"`     // ms response time
    SuccessRate   float64   `json:"success_rate"`
    LastUsed      time.Time `json:"last_used"`
    LastChecked   time.Time `json:"last_checked"`
    IsHealthy     bool      `json:"is_healthy"`
    FailureCount  int       `json:"failure_count"`
    ConcurrentUse int       `json:"concurrent_use"`
    MaxConcurrent int       `json:"max_concurrent"`
    Weight        int       `json:"weight"` // For weighted rotation
}

type RotationStats struct {
    TotalRotations   int64         `json:"total_rotations"`
    SuccessfulProxies int64        `json:"successful_proxies"`
    FailedProxies    int64         `json:"failed_proxies"`
    AvgResponseTime  time.Duration `json:"avg_response_time"`
}

func NewProxyRotationEngine(strategy RotationStrategy) *ProxyRotationEngine {
    return &ProxyRotationEngine{
        rotationStrategy: strategy,
        stats:            &RotationStats{},
        geoLocator:       NewGeoLocator(),
        proxies:          make([]*Proxy, 0),
        healthyProxies:   make([]*Proxy, 0),
    }
}

// Load 23,000+ proxies into the rotation engine
func (pre *ProxyRotationEngine) LoadProxies(proxies []*Proxy) error {
    pre.mu.Lock()
    defer pre.mu.Unlock()
    
    pre.proxies = proxies
    pre.healthyProxies = make([]*Proxy, 0)
    
    // Initial health check
    for _, proxy := range pre.proxies {
        if pre.healthCheck(proxy) {
            pre.healthyProxies = append(pre.healthyProxies, proxy)
        }
    }
    
    fmt.Printf("Loaded %d proxies, %d healthy\n", len(pre.proxies), len(pre.healthyProxies))
    return nil
}

// GetNextProxy - Main rotation method for 23,000 different IPs
func (pre *ProxyRotationEngine) GetNextProxy(targetURL string) (*Proxy, error) {
    pre.mu.Lock()
    defer pre.mu.Unlock()
    
    if len(pre.healthyProxies) == 0 {
        return nil, fmt.Errorf("no healthy proxies available")
    }
    
    var selectedProxy *Proxy
    
    switch pre.rotationStrategy {
    case RoundRobin:
        selectedProxy = pre.roundRobinSelection()
    case Random:
        selectedProxy = pre.randomSelection()
    case Weighted:
        selectedProxy = pre.weightedSelection()
    case StickySession:
        selectedProxy = pre.stickySessionSelection(targetURL)
    case Geographic:
        selectedProxy = pre.geographicSelection(targetURL)
    default:
        selectedProxy = pre.randomSelection()
    }
    
    if selectedProxy == nil {
        return nil, fmt.Errorf("failed to select proxy")
    }
    
    // Update usage stats
    selectedProxy.LastUsed = time.Now()
    selectedProxy.ConcurrentUse++
    pre.stats.TotalRotations++
    
    return selectedProxy, nil
}

// Round-robin selection - cycles through all 23,000 proxies
func (pre *ProxyRotationEngine) roundRobinSelection() *Proxy {
    if pre.currentIndex >= len(pre.healthyProxies) {
        pre.currentIndex = 0
    }
    
    proxy := pre.healthyProxies[pre.currentIndex]
    pre.currentIndex++
    
    return proxy
}

// Random selection - truly random from 23,000 proxies
func (pre *ProxyRotationEngine) randomSelection() *Proxy {
    n, err := rand.Int(rand.Reader, big.NewInt(int64(len(pre.healthyProxies))))
    if err != nil {
        return pre.healthyProxies[0]
    }
    return pre.healthyProxies[n.Int64()]
}

// Weighted selection - based on success rate and speed
func (pre *ProxyRotationEngine) weightedSelection() *Proxy {
    totalWeight := 0
    for _, proxy := range pre.healthyProxies {
        weight := pre.calculateProxyWeight(proxy)
        totalWeight += weight
    }
    
    if totalWeight == 0 {
        return pre.randomSelection()
    }
    
    r, _ := rand.Int(rand.Reader, big.NewInt(int64(totalWeight)))
    randomWeight := int(r.Int64())
    
    currentWeight := 0
    for _, proxy := range pre.healthyProxies {
        currentWeight += pre.calculateProxyWeight(proxy)
        if currentWeight >= randomWeight {
            return proxy
        }
    }
    
    return pre.healthyProxies[0]
}

// Sticky session - same proxy for same target
func (pre *ProxyRotationEngine) stickySessionSelection(targetURL string) *Proxy {
    // Hash the target URL to consistently select same proxy
    hash := pre.hashString(targetURL)
    index := hash % uint32(len(pre.healthyProxies))
    return pre.healthyProxies[index]
}

// Geographic selection - match proxy location to target
func (pre *ProxyRotationEngine) geographicSelection(targetURL string) *Proxy {
    targetCountry := pre.geoLocator.DetectCountryFromURL(targetURL)
    
    // Try to find proxy in same country
    var countryProxies []*Proxy
    for _, proxy := range pre.healthyProxies {
        if pre.geoLocator.GetCountry(proxy.Location) == targetCountry {
            countryProxies = append(countryProxies, proxy)
        }
    }
    
    if len(countryProxies) > 0 {
        n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(countryProxies))))
        return countryProxies[n.Int64()]
    }
    
    // Fallback to random selection
    return pre.randomSelection()
}

func (pre *ProxyRotationEngine) calculateProxyWeight(proxy *Proxy) int {
    baseWeight := 100
    
    // Adjust based on success rate
    successWeight := int(proxy.SuccessRate * 100)
    
    // Adjust based on speed (faster = higher weight)
    speedWeight := 0
    if proxy.Speed < 100 {
        speedWeight = 50
    } else if proxy.Speed < 500 {
        speedWeight = 25
    }
    
    // Adjust based on recent usage
    usageWeight := 0
    if time.Since(proxy.LastUsed) > 5*time.Minute {
        usageWeight = 30 // Prefer unused proxies
    }
    
    // Adjust based on concurrent usage
    concurrentWeight := 0
    if proxy.ConcurrentUse < proxy.MaxConcurrent/2 {
        concurrentWeight = 20
    }
    
    return baseWeight + successWeight + speedWeight + usageWeight + concurrentWeight
}

// Health checking for 23,000 proxies
func (pre *ProxyRotationEngine) StartHealthMonitoring() {
    ticker := time.NewTicker(30 * time.Second)
    go func() {
        for range ticker.C {
            pre.checkAllProxiesHealth()
        }
    }()
}

func (pre *ProxyRotationEngine) checkAllProxiesHealth() {
    pre.mu.Lock()
    defer pre.mu.Unlock()
    
    var wg sync.WaitGroup
    healthyCount := 0
    
    // Check proxies in batches to avoid overwhelming
    batchSize := 100
    for i := 0; i < len(pre.proxies); i += batchSize {
        end := i + batchSize
        if end > len(pre.proxies) {
            end = len(pre.proxies)
        }
        
        batch := pre.proxies[i:end]
        
        for _, proxy := range batch {
            wg.Add(1)
            go func(p *Proxy) {
                defer wg.Done()
                if pre.healthCheck(p) {
                    healthyCount++
                }
            }(proxy)
        }
        
        wg.Wait()
        time.Sleep(1 * time.Second) // Rate limiting between batches
    }
    
    // Update healthy proxies list
    pre.healthyProxies = pre.healthyProxies[:0]
    for _, proxy := range pre.proxies {
        if proxy.IsHealthy {
            pre.healthyProxies = append(pre.healthyProxies, proxy)
        }
    }
    
    fmt.Printf("Health check completed: %d/%d proxies healthy\n", healthyCount, len(pre.proxies))
}

func (pre *ProxyRotationEngine) healthCheck(proxy *Proxy) bool {
    // Implement actual health check logic
    // Test proxy connectivity, speed, etc.
    
    // Simulate health check
    proxy.LastChecked = time.Now()
    
    // For demo - 95% success rate
    if proxy.FailureCount > 10 {
        proxy.IsHealthy = false
        return false
    }
    
    proxy.IsHealthy = true
    return true
}

func (pre *ProxyRotationEngine) hashString(s string) uint32 {
    var h uint32
    for i := 0; i < len(s); i++ {
        h = 31*h + uint32(s[i])
    }
    return h
}
