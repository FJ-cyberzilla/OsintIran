// api-gateway/internal/handlers/phone_normalize.go
package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/gin-gonic/gin"
    "secure-iran-intel/pkg/normalizer"
)

type NormalizationHandler struct {
    normalizer *normalizer.PhoneNormalizer
}

func NewNormalizationHandler() *NormalizationHandler {
    return &NormalizationHandler{
        normalizer: normalizer.NewPhoneNormalizer("IR"), // Default to Iran
    }
}

// NormalizeSingle - Normalize a single phone number
func (nh *NormalizationHandler) NormalizeSingle(c *gin.Context) {
    var req struct {
        PhoneNumber string `json:"phone_number" binding:"required"`
        Country     string `json:"country,omitempty"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }
    
    normalized, err := nh.normalizer.NormalizePhone(req.PhoneNumber, req.Country)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to normalize phone number",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    normalized,
    })
}

// NormalizeBatch - Normalize multiple phone numbers
func (nh *NormalizationHandler) NormalizeBatch(c *gin.Context) {
    var req struct {
        PhoneNumbers []string `json:"phone_numbers" binding:"required"`
        Country      string   `json:"country,omitempty"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }
    
    results, errors := nh.normalizer.NormalizeBatch(req.PhoneNumbers, req.Country)
    
    response := gin.H{
        "success": true,
        "data":    results,
    }
    
    if len(errors) > 0 {
        response["errors"] = errors
    }
    
    c.JSON(http.StatusOK, response)
}

// ValidatePhone - Validate phone number format and existence
func (nh *NormalizationHandler) ValidatePhone(c *gin.Context) {
    var req struct {
        PhoneNumber string `json:"phone_number" binding:"required"`
        Country     string `json:"country,omitempty"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }
    
    normalized, err := nh.normalizer.NormalizePhone(req.PhoneNumber, req.Country)
    if err != nil {
        c.JSON(http.StatusOK, gin.H{
            "valid":   false,
            "message": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "valid":    normalized.IsValid,
        "possible": normalized.IsPossible,
        "type":     normalized.Type,
        "carrier":  normalized.Carrier,
        "region":   normalized.Region,
        "formats": gin.H{
            "e164":         normalized.Normalized,
            "international": normalized.International,
            "national":      normalized.National,
        },
    })
}
