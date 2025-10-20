// tests/integration/proxies/proxy_rotation.test.go
package integration

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "secure-iran-intel/proxy-pool/internal/proxy"
)

type ProxyRotationTestSuite struct {
    suite.Suite
    rotationEngine *proxy.ProxyRotationEngine
    testProxies    []*proxy.Proxy
    ctx            context.Context
}

func TestProxyRotationSuite(t *testing.T) {
    suite.Run(t, new(ProxyRotationTestSuite))
}

func (suite *ProxyRotationTestSuite) SetupTest() {
    suite.ctx = context.Background()
    suite.rotationEngine = proxy.NewProxyRotationEngine(proxy.Random)
    
    // Create test proxies
    suite.testProxies = []*proxy.Proxy{
        {
            ID:          "proxy-1",
            IP:          "192.168.1.1",
            Port:        8080,
            Type:        proxy.HTTP,
            Location:    "Tehran",
            ISP:         "MCI",
            Speed:       100,
            SuccessRate: 0.95,
            IsHealthy:   true,
        },
        {
            ID:          "proxy-2", 
            IP:          "192.168.1.2",
            Port:        8080,
            Type:        proxy.HTTPS,
            Location:    "Mashhad",
            ISP:         "MTN",
            Speed:       150,
            SuccessRate: 0.98,
            IsHealthy:   true,
        },
        {
            ID:          "proxy-3",
            IP:          "192.168.1.3",
            Port:        8080,
            Type:        proxy.SOCKS5,
            Location:    "Isfahan",
            ISP:         "Rightel",
            Speed:       200,
            SuccessRate: 0.92,
            IsHealthy:   false, // Mark as unhealthy for testing
        },
    }
    
    err := suite.rotationEngine.LoadProxies(suite.testProxies)
    suite.NoError(err)
}

func (suite *ProxyRotationTestSuite) TestProxyRotation() {
    // Test round-robin rotation
    suite.rotationEngine.SetRotationStrategy(proxy.RoundRobin)
    
    proxies := make([]*proxy.Proxy, 0)
    for i := 0; i < 5; i++ {
        proxy, err := suite.rotationEngine.GetNextProxy("https://facebook.com")
        suite.NoError(err)
        proxies = append(proxies, proxy)
    }
    
    // Should only return healthy proxies
    for _, p := range proxies {
        suite.True(p.IsHealthy, "Should only return healthy proxies")
    }
    
    // Should rotate through available proxies
    suite.Equal(2, len(uniqueProxyIDs(proxies)), "Should use both healthy proxies")
}

func (suite *ProxyRotationTestSuite) TestHealthChecking() {
    // Mark a proxy as unhealthy
    unhealthyProxy := suite.testProxies[0]
    suite.rotationEngine.ReportProxyFailure(unhealthyProxy.ID)
    
    // Health check should mark it as unhealthy
    suite.rotationEngine.CheckProxyHealth(unhealthyProxy.ID)
    
    updatedProxy, err := suite.rotationEngine.GetProxy(unhealthyProxy.ID)
    suite.NoError(err)
    suite.False(updatedProxy.IsHealthy, "Failed proxy should be marked unhealthy")
}

func (suite *ProxyRotationTestSuite) TestGeographicRotation() {
    suite.rotationEngine.SetRotationStrategy(proxy.Geographic)
    
    // Test with Iranian target
    proxy, err := suite.rotationEngine.GetNextProxy("https://facebook.com")
    suite.NoError(err)
    suite.Equal("Iran", proxy.Location, "Should prefer Iranian proxies for Iranian targets")
    
    // Test with US target  
    proxy, err = suite.rotationEngine.GetNextProxy("https://linkedin.com")
    suite.NoError(err)
    // Should still return available proxy even if geographic match not perfect
    suite.NotNil(proxy)
}

func (suite *ProxyRotationTestSuite) TestConcurrentProxyAccess() {
    const numGoroutines = 10
    const accessesPerGoroutine = 5
    
    results := make(chan *proxy.Proxy, numGoroutines*accessesPerGoroutine)
    errors := make(chan error, numGoroutines*accessesPerGoroutine)
    
    for i := 0; i < numGoroutines; i++ {
        go func(id int) {
            for j := 0; j < accessesPerGoroutine; j++ {
                proxy, err := suite.rotationEngine.GetNextProxy("https://test.com")
                if err != nil {
                    errors <- err
                } else {
                    results <- proxy
                }
                time.Sleep(time.Millisecond * 10)
            }
        }(i)
    }
    
    // Wait for all goroutines to complete
    time.Sleep(time.Second * 2)
    close(results)
    close(errors)
    
    // Collect results
    var proxies []*proxy.Proxy
    for p := range results {
        proxies = append(proxies, p)
    }
    
    var errs []error
    for e := range errors {
        errs = append(errs, e)
    }
    
    suite.Empty(errs, "Should not have any errors in concurrent access")
    suite.Len(proxies, numGoroutines*accessesPerGoroutine, "Should handle concurrent access")
    
    // Verify all proxies are healthy
    for _, p := range proxies {
        suite.True(p.IsHealthy, "All returned proxies should be healthy")
    }
}

func uniqueProxyIDs(proxies []*proxy.Proxy) []string {
    seen := make(map[string]bool)
    var unique []string
    for _, p := range proxies {
        if !seen[p.ID] {
            seen[p.ID] = true
            unique = append(unique, p.ID)
        }
    }
    return unique
}
