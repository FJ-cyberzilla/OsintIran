// 02_check-engine-go/cmd/worker.go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "secure-iran-intel/02_check-engine-go/internal/checker"
    "secure-iran-intel/02_check-engine-go/internal/proxy"
    "secure-iran-intel/02_check-engine-go/internal/queue"
)

func main() {
    log.Println("🚀 Starting Check Engine Worker...")

    // Initialize proxy rotation engine
    rotationEngine := proxy.NewProxyRotationEngine(proxy.Random)
    
    // Load 23,000 proxies
    if err := rotationEngine.LoadProxies(loadProxiesFromDatabase()); err != nil {
        log.Fatalf("Failed to load proxies: %v", err)
    }

    // Start health monitoring
    rotationEngine.StartHealthMonitoring()

    // Initialize checker
    simpleChecker := checker.NewSimpleChecker(rotationEngine)

    // Initialize message queue consumer
    mqConsumer := queue.NewMQConsumer(simpleChecker)

    // Start consuming tasks
    if err := mqConsumer.Start(); err != nil {
        log.Fatalf("Failed to start MQ consumer: %v", err)
    }

    log.Println("✅ Check Engine Worker started successfully")

    // Wait for shutdown signal
    waitForShutdown(mqConsumer)
}

func loadProxiesFromDatabase() []*proxy.Proxy {
    // Load 23,000 proxies from database or file
    var proxies []*proxy.Proxy
    
    // This would typically come from a database
    // For demo, create sample proxies
    for i := 1; i <= 23000; i++ {
        proxies = append(proxies, &proxy.Proxy{
            ID:      fmt.Sprintf("proxy-%d", i),
            IP:      fmt.Sprintf("192.168.%d.%d", (i/255)+1, i%255),
            Port:    8080,
            Type:    proxy.HTTP,
            Location: "Various",
            ISP:     "Various ISPs",
            Speed:   100 + (i % 500), // Varying speeds
            SuccessRate: 0.95,
            IsHealthy:   true,
            MaxConcurrent: 5,
        })
    }
    
    return proxies
}

func waitForShutdown(consumer *queue.MQConsumer) {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    <-sigChan
    log.Println("🛑 Shutdown signal received")
    
    consumer.Stop()
    log.Println("✅ Check Engine Worker stopped gracefully")
}
