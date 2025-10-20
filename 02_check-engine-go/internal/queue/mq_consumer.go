// 02_check-engine-go/internal/queue/mq_consumer.go
package queue

import (
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/streadway/amqp"
    "secure-iran-intel/02_check-engine-go/internal/checker"
)

type MQConsumer struct {
    checker    *checker.SimpleChecker
    connection *amqp.Connection
    channel    *amqp.Channel
    running    bool
}

func NewMQConsumer(checker *checker.SimpleChecker) *MQConsumer {
    return &MQConsumer{
        checker: checker,
        running: false,
    }
}

func (mc *MQConsumer) Start() error {
    var err error
    
    // Connect to RabbitMQ
    mc.connection, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }
    
    mc.channel, err = mc.connection.Channel()
    if err != nil {
        return fmt.Errorf("failed to open channel: %w", err)
    }
    
    // Declare queue
    _, err = mc.channel.QueueDeclare(
        "intelligence_tasks", // name
        true,                 // durable
        false,                // delete when unused
        false,                // exclusive
        false,                // no-wait
        nil,                  // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare queue: %w", err)
    }
    
    // Set prefetch count to limit concurrent tasks
    err = mc.channel.Qos(
        10,     // prefetch count
        0,      // prefetch size
        false,  // global
    )
    if err != nil {
        return fmt.Errorf("failed to set QoS: %w", err)
    }
    
    // Start consuming
    messages, err := mc.channel.Consume(
        "intelligence_tasks", // queue
        "",                   // consumer
        false,                // auto-ack
        false,                // exclusive
        false,                // no-local
        false,                // no-wait
        nil,                  // args
    )
    if err != nil {
        return fmt.Errorf("failed to register consumer: %w", err)
    }
    
    mc.running = true
    
    // Process messages
    go mc.processMessages(messages)
    
    log.Println("✅ MQ Consumer started successfully")
    return nil
}

func (mc *MQConsumer) processMessages(messages <-chan amqp.Delivery) {
    for message := range messages {
        if !mc.running {
            break
        }
        
        var task Task
        if err := json.Unmarshal(message.Body, &task); err != nil {
            log.Printf("❌ Failed to unmarshal task: %v", err)
            message.Nack(false, false) // Don't requeue malformed messages
            continue
        }
        
        log.Printf("🔍 Processing task: %s for %s on %s", task.ID, task.PhoneNumber, task.Platform)
        
        // Process the task
        result, err := mc.checker.ProcessTask(&task)
        if err != nil {
            log.Printf("❌ Task failed: %s - %v", task.ID, err)
            message.Nack(false, true) // Requeue failed tasks
            continue
        }
        
        // Publish result
        if err := mc.publishResult(result); err != nil {
            log.Printf("❌ Failed to publish result: %v", err)
        }
        
        // Acknowledge message
        message.Ack(false)
        
        log.Printf("✅ Task completed: %s - Success: %t", task.ID, result.Success)
        
        // Rate limiting - don't overwhelm targets
        time.Sleep(100 * time.Millisecond)
    }
}

func (mc *MQConsumer) publishResult(result *checker.TaskResult) error {
    resultJSON, err := json.Marshal(result)
    if err != nil {
        return err
    }
    
    return mc.channel.Publish(
        "",                    // exchange
        "task_results",        // routing key
        false,                 // mandatory
        false,                 // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        resultJSON,
            Timestamp:   time.Now(),
        },
    )
}

func (mc *MQConsumer) Stop() {
    mc.running = false
    if mc.channel != nil {
        mc.channel.Close()
    }
    if mc.connection != nil {
        mc.connection.Close()
    }
}
