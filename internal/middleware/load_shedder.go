// internal/middleware/load_shedder.go
package middleware

type LoadShedder struct {
    maxConcurrent int64
    current       int64
    rejectProbability float64
    metrics       *MetricsCollector
}

func NewLoadShedder(maxConcurrent int64) *LoadShedder {
    return &LoadShedder{
        maxConcurrent: maxConcurrent,
        rejectProbability: 0.0,
        metrics: NewMetricsCollector(),
    }
}

func (ls *LoadShedder) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Iran-specific: Check connectivity health first
        if ls.shouldRejectRequest() {
            ls.metrics.IncrementRejections()
            w.Header().Set("X-Load-Shedding", "true")
            http.Error(w, "Service temporarily unavailable", http.StatusServiceUnavailable)
            return
        }
        
        atomic.AddInt64(&ls.current, 1)
        defer atomic.AddInt64(&ls.current, -1)
        
        // Adaptive rejection based on Iran network conditions
        ls.adaptRejectionProbability()
        
        next.ServeHTTP(w, r)
    })
}

func (ls *LoadShedder) adaptRejectionProbability() {
    health := ls.metrics.GetNetworkHealth()
    
    // Increase rejection probability during poor Iran connectivity
    if health < 0.7 {
        ls.rejectProbability = 0.3
    } else if health < 0.9 {
        ls.rejectProbability = 0.1
    } else {
        ls.rejectProbability = 0.0
    }
}
