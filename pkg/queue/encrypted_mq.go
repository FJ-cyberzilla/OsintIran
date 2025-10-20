// pkg/queue/encrypted_mq.go
package queue

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/json"
    "fmt"
    "io"

    "github.com/streadway/amqp"
)

type SecureMessageQueue struct {
    conn         *amqp.Connection
    channel      *amqp.Channel
    encryptionKey []byte
    queueName    string
}

type SecureMessage struct {
    EncryptedData []byte `json:"encrypted_data"`
    IV           []byte `json:"iv"`
    HMAC         []byte `json:"hmac"`
    Timestamp    int64  `json:"timestamp"`
}

func NewSecureMessageQueue(amqpURL, queueName string, encryptionKey []byte) (*SecureMessageQueue, error) {
    conn, err := amqp.Dial(amqpURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    channel, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to open channel: %w", err)
    }

    // Declare secure queue
    _, err = channel.QueueDeclare(
        queueName, // name
        true,      // durable
        false,     // delete when unused
        false,     // exclusive
        false,     // no-wait
        nil,       // arguments
    )
    if err != nil {
        return nil, fmt.Errorf("failed to declare queue: %w", err)
    }

    return &SecureMessageQueue{
        conn:         conn,
        channel:      channel,
        encryptionKey: encryptionKey,
        queueName:    queueName,
    }, nil
}

func (smq *SecureMessageQueue) PublishSecureMessage(message interface{}) error {
    // Encrypt message
    encryptedMsg, err := smq.encryptMessage(message)
    if err != nil {
        return fmt.Errorf("failed to encrypt message: %w", err)
    }

    // Publish to queue
    err = smq.channel.Publish(
        "",              // exchange
        smq.queueName,   // routing key
        false,           // mandatory
        false,           // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        encryptedMsg,
            Timestamp:   time.Now(),
        })
    if err != nil {
        return fmt.Errorf("failed to publish message: %w", err)
    }

    return nil
}

func (smq *SecureMessageQueue) ConsumeSecureMessages(handler func(interface{}) error) error {
    msgs, err := smq.channel.Consume(
        smq.queueName, // queue
        "",            // consumer
        true,          // auto-ack
        false,         // exclusive
        false,         // no-local
        false,         // no-wait
        nil,           // args
    )
    if err != nil {
        return fmt.Errorf("failed to register consumer: %w", err)
    }

    go func() {
        for delivery := range msgs {
            var secureMsg SecureMessage
            if err := json.Unmarshal(delivery.Body, &secureMsg); err != nil {
                fmt.Printf("Failed to unmarshal secure message: %v\n", err)
                continue
            }

            // Decrypt message
            decryptedData, err := smq.decryptMessage(secureMsg)
            if err != nil {
                fmt.Printf("Failed to decrypt message: %v\n", err)
                continue
            }

            // Process message
            if err := handler(decryptedData); err != nil {
                fmt.Printf("Failed to process message: %v\n", err)
            }
        }
    }()

    return nil
}

func (smq *SecureMessageQueue) encryptMessage(data interface{}) ([]byte, error) {
    // Serialize data
    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    // Generate IV
    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }

    // Encrypt data
    block, err := aes.NewCipher(smq.encryptionKey)
    if err != nil {
        return nil, err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    encrypted := make([]byte, len(jsonData))
    stream.XORKeyStream(encrypted, jsonData)

    // Create secure message
    secureMsg := SecureMessage{
        EncryptedData: encrypted,
        IV:           iv,
        Timestamp:    time.Now().Unix(),
    }

    return json.Marshal(secureMsg)
}

func (smq *SecureMessageQueue) decryptMessage(secureMsg SecureMessage) (interface{}, error) {
    block, err := aes.NewCipher(smq.encryptionKey)
    if err != nil {
        return nil, err
    }

    stream := cipher.NewCFBDecrypter(block, secureMsg.IV)
    decrypted := make([]byte, len(secureMsg.EncryptedData))
    stream.XORKeyStream(decrypted, secureMsg.EncryptedData)

    var data interface{}
    if err := json.Unmarshal(decrypted, &data); err != nil {
        return nil, err
    }

    return data, nil
}
