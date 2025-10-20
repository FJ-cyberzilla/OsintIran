// backend/internal/middleware/license_check.go
package middleware

func LicenseCheck() gin.HandlerFunc {
    return func(c *gin.Context) {
        licenseHeader := c.GetHeader("X-License-Key")
        if licenseHeader == "" {
            c.JSON(401, gin.H{"error": "License required"})
            c.Abort()
            return
        }

        validator := security.NewLicenseValidator()
        license, err := validator.ValidateLicense([]byte(licenseHeader))
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid license"})
            c.Abort()
            return
        }

        // Add license info to context
        c.Set("license", license)
        c.Next()
    }
}

func AntiTamper() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check request integrity
        if !verifyRequestIntegrity(c.Request) {
            c.JSON(400, gin.H{"error": "Request tampered"})
            c.Abort()
            return
        }

        // Validate timestamp to prevent replay attacks
        if !validateTimestamp(c) {
            c.JSON(400, gin.H{"error": "Invalid timestamp"})
            c.Abort()
            return
        }

        c.Next()
    }
}
