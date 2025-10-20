// proxy-pool/internal/proxy/iran_proxy_pool.go
package proxy

type IranProxyPool struct {
    proxies         []*Proxy
    healthyProxies  []*Proxy
    residentialIPs  []*ResidentialProxy
    mobileProxies   []*MobileProxy
    healthChecker   *HealthChecker
    rotationStrategy RotationStrategy
    mu              sync.RWMutex
}

type Proxy struct {
    ID          string
    IP          string
    Port        int
    Type        ProxyType // HTTP, HTTPS, SOCKS5
    Provider    string
    Location    string   // Iranian cities
    Speed       time.Duration
    SuccessRate float64
    LastUsed    time.Time
    IsHealthy   bool
}

func NewIranProxyPool() *IranProxyPool {
    pool := &IranProxyPool{
        healthChecker: NewHealthChecker(),
        rotationStrategy: &SmartRotation{},
    }
    
    // Load Iranian proxies
    pool.loadIranianProxies()
    pool.loadResidentialIPs()
    pool.loadMobileProxies()
    
    // Start health checking
    go pool.continuousHealthCheck()
    
    return pool
}

func (pool *IranProxyPool) GetNextProxy() (*Proxy, error) {
    pool.mu.RLock()
    defer pool.mu.RUnlock()
    
    if len(pool.healthyProxies) == 0 {
        return nil, fmt.Errorf("no healthy proxies available")
    }
    
    return pool.rotationStrategy.SelectProxy(pool.healthyProxies), nil
}

func (pool *IranProxyPool) continuousHealthCheck() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            pool.checkAllProxies()
        }
    }
}

func (pool *IranProxyPool) loadIranianProxies() {
    // Load Iran-specific proxies
    iranProxies := []*Proxy{
        {
            ID:   "ir-proxy-1",
            IP:   "185.123.456.789", // Example Iranian IP
            Port: 8080,
            Type: HTTP,
            Location: "Tehran",
            Provider: "IranCell",
        },
        // Add more Iranian proxies
    }
    
    pool.proxies = append(pool.proxies, iranProxies...)
}
