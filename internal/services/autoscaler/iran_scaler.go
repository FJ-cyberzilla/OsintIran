// internal/services/autoscaler/iran_scaler.go
package autoscaler

type IranScaler struct {
    baseReplicas    int
    maxReplicas     int
    networkMetrics  *NetworkMetrics
    proxyHealth     *ProxyHealthChecker
}

func (is *IranScaler) ShouldScale() (int, error) {
    base := is.baseReplicas
    
    // Iran-specific: Scale based on network conditions
    if is.networkMetrics.GetLatency() > 500*time.Millisecond {
        base += 2 // More instances to handle slow networks
    }
    
    if is.proxyHealth.HealthyProxies() < 5 {
        base += 1 // Scale up if proxy pool is depleted
    }
    
    // Scale based on success rate (Iran networks are less reliable)
    successRate := is.networkMetrics.GetSuccessRate()
    if successRate < 0.8 {
        base = int(float64(base) * (1.0 / successRate))
    }
    
    return min(base, is.maxReplicas), nil
}
