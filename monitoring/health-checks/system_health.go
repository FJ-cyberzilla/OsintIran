// monitoring/health-checks/system_health.go
package health

type SystemHealthMonitor struct {
    metricsCollector *MetricsCollector
    alertManager     *AlertManager
    autoHealing      *AutoHealing
}

func (shm *SystemHealthMonitor) StartMonitoring() {
    go shm.monitorProxyHealth()
    go shm.monitorAIAgents()
    go shm.monitorBrowserWorkers()
    go shm.monitorSystemResources()
}

func (shm *SystemHealthMonitor) monitorProxyHealth() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            healthyCount := shm.checkProxyHealth()
            if healthyCount < 5 {
                shm.alertManager.SendAlert("LOW_PROXY_COUNT", 
                    fmt.Sprintf("Only %d healthy proxies remaining", healthyCount))
                
                // Trigger auto-healing
                shm.autoHealing.ReplenishProxies()
            }
        }
    }
}

func (shm *SystemHealthMonitor) monitorAIAgents() {
    ticker := time.NewTicker(60 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            for _, agent := range shm.listAIAgents() {
                if !agent.IsHealthy() {
                    shm.alertManager.SendAlert("AI_AGENT_DOWN", 
                        fmt.Sprintf("AI agent %s is down", agent.Name))
                    
                    shm.autoHealing.RestartAgent(agent)
                }
            }
        }
    }
}
