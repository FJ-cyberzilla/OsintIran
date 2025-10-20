// internal/pkg/resilience/circuit_breaker.go
package resilience

type IranCircuitBreaker struct {
    failureThreshold int
    successThreshold int
    timeout          time.Duration
    state            CircuitState
    lastFailure      time.Time
    iranSpecific     IranConfig
}

type IranConfig struct {
    VPNFallback      bool
    ProxyRotation    bool
    MobileDataFallback bool
    TimeoutMultiplier time.Duration // Longer timeouts for Iran
}

func (cb *IranCircuitBreaker) Execute(ctx context.Context, operation func() error) error {
    if cb.state == StateOpen {
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = StateHalfOpen
        } else {
            return ErrCircuitOpen
        }
    }
    
    // Iran-specific: Try primary, then fallbacks
    err := cb.executeWithIranFallbacks(ctx, operation)
    if err != nil {
        cb.recordFailure()
        return err
    }
    
    cb.recordSuccess()
    return nil
}

func (cb *IranCircuitBreaker) executeWithIranFallbacks(ctx context.Context, op func() error) error {
    strategies := []func() error{
        op, // Primary operation
        cb.withVPNFallback,
        cb.withProxyRotation,
        cb.withMobileData,
    }
    
    for _, strategy := range strategies {
        if err := strategy(); err == nil {
            return nil
        }
        time.Sleep(cb.iranSpecific.TimeoutMultiplier)
    }
    return ErrAllFallbacksFailed
}
