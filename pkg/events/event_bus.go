// pkg/events/event_bus.go
type EventBus interface {
    Publish(ctx context.Context, event Event) error
    Subscribe(ctx context.Context, handler EventHandler)
}

// Kafka implementation
type KafkaEventBus struct {
    producer sarama.SyncProducer
    consumer sarama.Consumer
}
