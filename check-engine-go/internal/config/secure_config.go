// check-engine-go/internal/config/secure_config.go
package config

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/json"
    "fmt"
    "io"
    "os"
)

type SecureConfig struct {
    DatabaseURL    string `json:"database_url" encrypted:"true"`
    RabbitMQURL    string `json:"rabbitmq_url" encrypted:"true"`
    EncryptionKey  []byte `json:"-"` // Never serialize to config
    APISecrets     map[string]string `json:"api_secrets" encrypted:"true"`
}

func LoadSecureConfig() (*SecureConfig, error) {
    // Load encryption key from secure location
    encryptionKey, err := loadEncryptionKey()
    if err != nil {
        return nil, fmt.Errorf("failed to load encryption key: %w", err)
    }

    // Read encrypted config file
    encryptedData, err := os.ReadFile("config/secure_config.enc")
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    // Decrypt config
    decryptedData, err := decryptConfig(encryptedData, encryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt config: %w", err)
    }

    var config SecureConfig
    if err := json.Unmarshal(decryptedData, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }

    config.EncryptionKey = encryptionKey
    return &config, nil
}

func decryptConfig(encryptedData []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    if len(encryptedData) < nonceSize {
        return nil, fmt.Errorf("ciphertext too short")
    }

    nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}

func loadEncryptionKey() ([]byte, error) {
    // Load from secure environment or HSM
    key := os.Getenv("CONFIG_ENCRYPTION_KEY")
    if key == "" {
        return nil, fmt.Errorf("encryption key not found in environment")
    }
    return []byte(key), nil
}
