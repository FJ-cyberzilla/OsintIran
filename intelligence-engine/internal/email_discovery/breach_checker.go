// intelligence-engine/internal/email_discovery/breach_checker.go
package email_discovery

import (
    "crypto/sha1"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
)

type BreachChecker struct {
    apiKeys         map[string]string
    localBreachDB   *LocalBreachDatabase
    httpClient      *http.Client
}

type BreachResult struct {
    Email         string   `json:"email"`
    Breaches      []string `json:"breaches"`
    Found         bool     `json:"found"`
    FirstBreached string   `json:"first_breached,omitempty"`
    LastBreached  string   `json:"last_breached,omitempty"`
}

func NewBreachChecker() *BreachChecker {
    return &BreachChecker{
        apiKeys: map[string]string{
            "haveibeenpwned": "your-api-key",
            "dehashed":       "your-api-key", 
        },
        localBreachDB: NewLocalBreachDatabase(),
        httpClient:    &http.Client{Timeout: 30 * time.Second},
    }
}

func (bc *BreachChecker) FindEmailsInBreaches(phone *NormalizedPhone) []*EmailDiscoveryResult {
    var results []*EmailDiscoveryResult
    
    // Generate potential emails from phone number
    potentialEmails := bc.generatePotentialEmails(phone)
    
    // Check each potential email against breach databases
    for _, email := range potentialEmails {
        // Check Have I Been Pwned API
        if breachResult := bc.checkHIBP(email); breachResult.Found {
            results = append(results, &EmailDiscoveryResult{
                Email:         email,
                Confidence:    0.9, // High confidence if found in breaches
                Source:        "breach_database",
                FoundInBreach: true,
                PatternType:   "breach_verified",
            })
        }
        
        // Check local breach database
        if bc.localBreachDB.CheckEmail(email) {
            results = append(results, &EmailDiscoveryResult{
                Email:         email,
                Confidence:    0.8,
                Source:        "local_breach_db", 
                FoundInBreach: true,
                PatternType:   "breach_verified",
            })
        }
        
        // Check Dehashed API
        if dehashedResult := bc.checkDehashed(email); dehashedResult.Found {
            results = append(results, &EmailDiscoveryResult{
                Email:         email,
                Confidence:    0.85,
                Source:        "dehashed",
                FoundInBreach: true,
                PatternType:   "breach_verified",
            })
        }
    }
    
    return results
}

func (bc *BreachChecker) checkHIBP(email string) *BreachResult {
    // Use k-anonymity model for privacy
    hash := sha1.Sum([]byte(strings.ToLower(email)))
    hashPrefix := hex.EncodeToString(hash[:5])
    hashSuffix := hex.EncodeToString(hash[5:])
    
    url := fmt.Sprintf("https://api.pwnedpasswords.com/range/%s", hashPrefix)
    
    resp, err := bc.httpClient.Get(url)
    if err != nil {
        return &BreachResult{Email: email, Found: false}
    }
    defer resp.Body.Close()
    
    // Parse response to check if our hash suffix exists
    // Implementation would parse the response body
    
    return &BreachResult{
        Email:    email,
        Found:    true, // Simplified for example
        Breaches: []string{"Adobe2013", "LinkedIn2012"},
    }
}

func (bc *BreachChecker) checkDehashed(email string) *BreachResult {
    // Implementation for Dehashed API
    // This would make authenticated requests to Dehashed
    return &BreachResult{
        Email: email,
        Found: false, // Placeholder
    }
}

// LocalBreachDatabase for faster lookups
type LocalBreachDatabase struct {
    breachData map[string][]BreachRecord
    loaded     bool
}

type BreachRecord struct {
    Email     string    `json:"email"`
    Breach    string    `json:"breach"`
    Date      time.Time `json:"date"`
    Source    string    `json:"source"`
}

func NewLocalBreachDatabase() *LocalBreachDatabase {
    db := &LocalBreachDatabase{
        breachData: make(map[string][]BreachRecord),
    }
    // Load breach data on initialization
    db.loadBreachData()
    return db
}

func (lbd *LocalBreachDatabase) CheckEmail(email string) bool {
    records, exists := lbd.breachData[strings.ToLower(email)]
    return exists && len(records) > 0
}
