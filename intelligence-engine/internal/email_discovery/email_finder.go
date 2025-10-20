// intelligence-engine/internal/email_discovery/email_finder.go
package email_discovery

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "regexp"
    "strings"
    "sync"

    "github.com/nyaruka/phonenumbers"
)

type EmailDiscoveryEngine struct {
    patternGenerator *PatternGenerator
    breachChecker    *BreachChecker
    socialCrawler    *SocialCrawler
    similarityScorer *SimilarityScorer
}

type EmailDiscoveryResult struct {
    Email           string  `json:"email"`
    Confidence      float64 `json:"confidence"`
    Source          string  `json:"source"`
    FoundInBreach   bool    `json:"found_in_breach"`
    SocialProfiles  []SocialProfile `json:"social_profiles"`
    PatternType     string  `json:"pattern_type"`
    IsVerified      bool    `json:"is_verified"`
    LastSeen        string  `json:"last_seen,omitempty"`
}

type SocialProfile struct {
    Platform    string `json:"platform"`
    Username    string `json:"username"`
    ProfileURL  string `json:"profile_url"`
    Verified    bool   `json:"verified"`
}

func NewEmailDiscoveryEngine() *EmailDiscoveryEngine {
    return &EmailDiscoveryEngine{
        patternGenerator: NewPatternGenerator(),
        breachChecker:    NewBreachChecker(),
        socialCrawler:    NewSocialCrawler(),
        similarityScorer: NewSimilarityScorer(),
    }
}

// DiscoverEmailsFromPhone - Main email discovery function
func (ede *EmailDiscoveryEngine) DiscoverEmailsFromPhone(phoneNumber string) (*EmailDiscoveryReport, error) {
    // Normalize phone number first
    normalized, err := ede.normalizePhone(phoneNumber)
    if err != nil {
        return nil, fmt.Errorf("phone normalization failed: %w", err)
    }

    var wg sync.WaitGroup
    results := make(chan *EmailDiscoveryResult, 50)
    var discoveredEmails []*EmailDiscoveryResult

    // Method 1: Pattern-based generation
    wg.Add(1)
    go func() {
        defer wg.Done()
        patternEmails := ede.patternGenerator.GenerateEmailPatterns(normalized)
        for _, email := range patternEmails {
            results <- email
        }
    }()

    // Method 2: Breach data lookup
    wg.Add(1)
    go func() {
        defer wg.Done()
        breachEmails := ede.breachChecker.FindEmailsInBreaches(normalized)
        for _, email := range breachEmails {
            results <- email
        }
    }()

    // Method 3: Social media crawling
    wg.Add(1)
    go func() {
        defer wg.Done()
        socialEmails := ede.socialCrawler.FindEmailsFromSocialProfiles(normalized)
        for _, email := range socialEmails {
            results <- email
        }
    }()

    // Method 4: Cross-platform correlation
    wg.Add(1)
    go func() {
        defer wg.Done()
        correlatedEmails := ede.correlateAcrossPlatforms(normalized)
        for _, email := range correlatedEmails {
            results <- email
        }
    }()

    // Close results channel when all goroutines complete
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    for result := range results {
        discoveredEmails = append(discoveredEmails, result)
    }

    // Deduplicate and rank results
    finalEmails := ede.deduplicateAndRank(discoveredEmails)

    report := &EmailDiscoveryReport{
        PhoneNumber:    phoneNumber,
        Normalized:     normalized,
        DiscoveredEmails: finalEmails,
        TotalFound:     len(finalEmails),
        Confidence:     ede.calculateOverallConfidence(finalEmails),
        GeneratedAt:    time.Now(),
    }

    return report, nil
}

// PatternGenerator - Generates email patterns based on phone number
type PatternGenerator struct {
    commonDomains   []string
    iranianDomains  []string
    nameDatabases   map[string][]string // Common first/last names by country
}

func NewPatternGenerator() *PatternGenerator {
    return &PatternGenerator{
        commonDomains: []string{
            "gmail.com", "yahoo.com", "outlook.com", "hotmail.com",
            "protonmail.com", "aol.com", "icloud.com",
        },
        iranianDomains: []string{
            "chmail.ir", "yahoo.com", "gmail.com", "protonmail.com",
            "mailfa.com", "iran.ir",
        },
        nameDatabases: loadNameDatabases(),
    }
}

func (pg *PatternGenerator) GenerateEmailPatterns(phone *NormalizedPhone) []*EmailDiscoveryResult {
    var results []*EmailDiscoveryResult
    nationalNum := phonenumbers.GetNationalSignificantNumber(phone.ParsedNumber)
    
    // Extract components from phone number
    last4 := nationalNum[len(nationalNum)-4:]
    last6 := nationalNum[len(nationalNum)-6:]
    fullNational := strings.TrimLeft(nationalNum, "0")

    // Iranian-specific patterns
    if phone.CountryCode == 98 {
        results = append(results, pg.generateIranianPatterns(fullNational, last4, last6)...)
    }

    // Global patterns
    results = append(results, pg.generateGlobalPatterns(fullNational, last4, last6, phone.CountryCode)...)

    // Name-based patterns (if we can infer names)
    results = append(results, pg.generateNameBasedPatterns(phone)...)

    return results
}

func (pg *PatternGenerator) generateIranianPatterns(fullNational, last4, last6 string) []*EmailDiscoveryResult {
    var results []*EmailDiscoveryResult
    basePatterns := []string{
        "98%s", "0%s", "ir%s", "09%s", "9%s", "%s", "98ir%s",
    }

    for _, domain := range pg.iranianDomains {
        for _, pattern := range basePatterns {
            email := fmt.Sprintf(pattern, fullNational) + "@" + domain
            results = append(results, &EmailDiscoveryResult{
                Email:       email,
                Confidence:  0.6,
                Source:      "pattern_generation",
                PatternType: "iranian_numeric",
            })

            // With last digits only
            emailLast4 := fmt.Sprintf(pattern, last4) + "@" + domain
            results = append(results, &EmailDiscoveryResult{
                Email:       emailLast4,
                Confidence:  0.4,
                Source:      "pattern_generation", 
                PatternType: "iranian_short",
            })
        }
    }

    return results
}

func (pg *PatternGenerator) generateNameBasedPatterns(phone *NormalizedPhone) []*EmailDiscoveryResult {
    var results []*EmailDiscoveryResult
    
    // Get common names for the detected region
    regionNames := pg.nameDatabases[phone.Region]
    if len(regionNames) == 0 {
        return results
    }

    // Try combinations with phone number components
    nationalNum := phonenumbers.GetNationalSignificantNumber(phone.ParsedNumber)
    last4 := nationalNum[len(nationalNum)-4:]

    for _, name := range regionNames[:10] { // Limit to top 10 names
        namePatterns := []string{
            name.First + "." + name.Last + last4,
            name.First[0] + name.Last + last4, 
            name.First + last4,
            name.Last + last4,
            name.First + "." + last4,
        }

        for _, pattern := range namePatterns {
            for _, domain := range pg.commonDomains {
                email := strings.ToLower(pattern) + "@" + domain
                results = append(results, &EmailDiscoveryResult{
                    Email:       email,
                    Confidence:  0.3, // Lower confidence for name-based
                    Source:      "pattern_generation",
                    PatternType: "name_based",
                })
            }
        }
    }

    return results
}
