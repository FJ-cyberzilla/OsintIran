// check-engine-go/cmd/worker.go
package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "secure-iran-intel/check-engine-go/internal/config"
    "secure-iran-intel/check-engine-go/internal/checker"
    "secure-iran-intel/check-engine-go/internal/queue"
    "secure-iran-intel/check-engine-go/internal/security"
    "secure-iran-intel/pkg/security/sql_injection"
)

func main() {
    // Initialize code protection
    codeProtector := security.NewCodeProtector()
    defer codeProtector.Cleanup()

    // Load secure configuration
    cfg, err := config.LoadSecureConfig()
    if err != nil {
        log.Fatalf("Failed to load secure config: %v", err)
    }

    // Initialize SQL injection protector
    sqlProtector := sql_injection.NewSQLInjectionProtector()

    // Initialize secure message queue
    secureQueue, err := queue.NewSecureRabbitMQ(cfg.QueueConfig, cfg.EncryptionKey)
    if err != nil {
        log.Fatalf("Failed to initialize secure queue: %v", err)
    }
    defer secureQueue.Close()

    // Initialize secure HTTP checker
    secureChecker := checker.NewSecureHTTPChecker(sqlProtector, cfg)

    // Start secure message consumer
    err = secureQueue.ConsumeSecureMessages(func(message interface{}) error {
        return secureChecker.ProcessSecureTask(message)
    })
    if err != nil {
        log.Fatalf("Failed to start secure consumer: %v", err)
    }

    log.Println("Secure check engine started successfully")

    // Wait for termination signal
    waitForShutdown()
}

func waitForShutdown() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    log.Println("Shutdown signal received")
}
