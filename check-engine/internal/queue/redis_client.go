// check-engine/internal/queue/redis_client.go
package queue

type RedisClient struct {
    client *redis.Client
    pubsub *redis.PubSub
}

func (rc *RedisClient) PublishTask(task *Task) error {
    taskJSON, err := json.Marshal(task)
    if err != nil {
        return err
    }

    return rc.client.Publish(context.Background(), "tasks", taskJSON).Err()
}

func (rc *RedisClient) SubscribeToResults(taskID string) <-chan *TaskResult {
    ch := make(chan *TaskResult)
    
    go func() {
        pubsub := rc.client.Subscribe(context.Background(), fmt.Sprintf("task_result:%s", taskID))
        defer pubsub.Close()
        
        for msg := range pubsub.Channel() {
            var result TaskResult
            if err := json.Unmarshal([]byte(msg.Payload), &result); err == nil {
                ch <- &result
            }
        }
    }()
    
    return ch
}
