// api-gateway/internal/middleware/normalization.go
package middleware

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"

    "github.com/gin-gonic/gin"
    "secure-iran-intel/pkg/normalizer"
)

type NormalizationMiddleware struct {
    normalizer *normalizer.PhoneNormalizer
}

func NewNormalizationMiddleware() *NormalizationMiddleware {
    return &NormalizationMiddleware{
        normalizer: normalizer.NewPhoneNormalizer("IR"),
    }
}

// AutoNormalizePhoneNumbers - Automatically normalize phone numbers in requests
func (nm *NormalizationMiddleware) AutoNormalizePhoneNumbers() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Only process JSON requests
        if c.ContentType() != "application/json" {
            c.Next()
            return
        }
        
        // Read and parse the request body
        var body map[string]interface{}
        if err := c.ShouldBindJSON(&body); err != nil {
            c.Next()
            return
        }
        
        // Normalize phone numbers in the request body
        nm.normalizeRequestBody(body)
        
        // Set the normalized body back to the context
        normalizedBody, err := json.Marshal(body)
        if err != nil {
            c.Next()
            return
        }
        
        c.Request.Body = io.NopCloser(bytes.NewBuffer(normalizedBody))
        c.Set("normalized_body", body)
        
        c.Next()
    }
}

func (nm *NormalizationMiddleware) normalizeRequestBody(body map[string]interface{}) {
    for key, value := range body {
        switch v := value.(type) {
        case string:
            if nm.looksLikePhoneNumber(v) {
                normalized, err := nm.normalizer.NormalizePhone(v, "")
                if err == nil {
                    body[key] = normalized.Normalized
                }
            }
        case []interface{}:
            nm.normalizeArray(v)
        case map[string]interface{}:
            nm.normalizeRequestBody(v)
        }
    }
}

func (nm *NormalizationMiddleware) normalizeArray(arr []interface{}) {
    for i, item := range arr {
        switch v := item.(type) {
        case string:
            if nm.looksLikePhoneNumber(v) {
                normalized, err := nm.normalizer.NormalizePhone(v, "")
                if err == nil {
                    arr[i] = normalized.Normalized
                }
            }
        case map[string]interface{}:
            nm.normalizeRequestBody(v)
        case []interface{}:
            nm.normalizeArray(v)
        }
    }
}

func (nm *NormalizationMiddleware) looksLikePhoneNumber(input string) bool {
    // Simple heuristic to identify phone number fields
    phoneKeywords := []string{"phone", "mobile", "tel", "number", "contact"}
    inputLower := strings.ToLower(input)
    
    for _, keyword := range phoneKeywords {
        if strings.Contains(inputLower, keyword) {
            return true
        }
    }
    
    // Check if the string contains mostly digits and common phone characters
    digitCount := 0
    for _, char := range input {
        if char >= '0' && char <= '9' {
            digitCount++
        }
    }
    
    return digitCount >= 7 && digitCount <= 15
}
