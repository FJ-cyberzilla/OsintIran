// intelligence-engine/internal/correlation_engine/platform_correlator.go
package correlation_engine

import (
    "fmt"
    "regexp"
    "strings"
    "time"

    "github.com/sajari/fuzzy"
)

type PlatformCorrelator struct {
    identityResolver   *IdentityResolver
    patternMatcher     *PatternMatcher
    confidenceCalculator *ConfidenceCalculator
    timelineAnalyzer   *TimelineAnalyzer
}

type CrossPlatformIdentity struct {
    PrimaryID      string                     `json:"primary_id"`
    Identities     map[string]PlatformIdentity `json:"identities"` // platform -> identity
    Confidence     float64                    `json:"confidence"`
    Evidence       []CorrelationEvidence      `json:"evidence"`
    Timeline       *UnifiedTimeline           `json:"timeline"`
    RiskIndicators []RiskIndicator            `json:"risk_indicators"`
}

type PlatformIdentity struct {
    Platform      string                 `json:"platform"`
    Username      string                 `json:"username"`
    ProfileURL    string                 `json:"profile_url"`
    Verified      bool                   `json:"verified"`
    Attributes    map[string]interface{} `json:"attributes"`
    FirstSeen     time.Time              `json:"first_seen"`
    LastSeen      time.Time              `json:"last_seen"`
    ActivityLevel string                 `json:"activity_level"`
}

type CorrelationEvidence struct {
    Type        string  `json:"type"`
    PlatformA   string  `json:"platform_a"`
    PlatformB   string  `json:"platform_b"`
    Confidence  float64 `json:"confidence"`
    Description string  `json:"description"`
    Timestamp   time.Time `json:"timestamp"`
}

func NewPlatformCorrelator() *PlatformCorrelator {
    return &PlatformCorrelator{
        identityResolver:   NewIdentityResolver(),
        patternMatcher:     NewPatternMatcher(),
        confidenceCalculator: NewConfidenceCalculator(),
        timelineAnalyzer:   NewTimelineAnalyzer(),
    }
}

// CorrelateAcrossPlatforms performs comprehensive cross-platform identity resolution
func (pc *PlatformCorrelator) CorrelateAcrossPlatforms(phoneNumber string, platformData map[string]PlatformData) (*CrossPlatformIdentity, error) {
    correlation := &CrossPlatformIdentity{
        PrimaryID:  phoneNumber,
        Identities: make(map[string]PlatformIdentity),
        Evidence:   make([]CorrelationEvidence, 0),
        Timeline:   &UnifiedTimeline{Events: make([]TimelineEvent, 0)},
    }

    // Step 1: Extract and normalize identities from each platform
    normalizedIdentities := pc.extractNormalizedIdentities(platformData)

    // Step 2: Perform pairwise platform correlation
    correlationResults := pc.correlatePlatformPairs(normalizedIdentities)

    // Step 3: Resolve identities across all platforms
    resolvedIdentity := pc.resolveGlobalIdentity(correlationResults, normalizedIdentities)

    // Step 4: Build unified timeline
    unifiedTimeline := pc.buildUnifiedTimeline(resolvedIdentity, platformData)

    // Step 5: Calculate overall confidence
    overallConfidence := pc.calculateOverallConfidence(correlationResults)

    // Step 6: Identify risk indicators
    riskIndicators := pc.identifyRiskIndicators(resolvedIdentity, unifiedTimeline)

    correlation.Identities = resolvedIdentity
    correlation.Confidence = overallConfidence
    correlation.Timeline = unifiedTimeline
    correlation.RiskIndicators = riskIndicators

    return correlation, nil
}

func (pc *PlatformCorrelator) extractNormalizedIdentities(platformData map[string]PlatformData) map[string]NormalizedIdentity {
    normalized := make(map[string]NormalizedIdentity)
    
    for platform, data := range platformData {
        normalized[platform] = pc.normalizePlatformIdentity(platform, data)
    }
    
    return normalized
}

func (pc *PlatformCorrelator) normalizePlatformIdentity(platform string, data PlatformData) NormalizedIdentity {
    normalized := NormalizedIdentity{
        Platform: platform,
        Username: data.Username,
        Profiles: data.Profiles,
    }

    // Extract and normalize name information
    if fullName, exists := data.Attributes["full_name"]; exists {
        normalized.FullName = pc.normalizeName(fullName.(string))
    }

    // Extract and normalize location information
    if location, exists := data.Attributes["location"]; exists {
        normalized.Location = pc.normalizeLocation(location.(string))
    }

    // Extract contact information
    normalized.Emails = pc.extractEmails(data)
    normalized.PhoneNumbers = pc.extractPhoneNumbers(data)

    // Extract behavioral patterns
    normalized.BehavioralPatterns = pc.extractBehavioralPatterns(data)

    // Extract temporal patterns
    normalized.TemporalPatterns = pc.extractTemporalPatterns(data)

    return normalized
}

func (pc *PlatformCorrelator) correlatePlatformPairs(identities map[string]NormalizedIdentity) map[string]CorrelationResult {
    results := make(map[string]CorrelationResult)
    platforms := pc.getPlatformKeys(identities)

    // Perform pairwise correlation between all platforms
    for i := 0; i < len(platforms); i++ {
        for j := i + 1; j < len(platforms); j++ {
            platformA := platforms[i]
            platformB := platforms[j]
            
            identityA := identities[platformA]
            identityB := identities[platformB]

            // Multiple correlation methods
            correlationScore := pc.calculateMultiFactorCorrelation(identityA, identityB)
            
            // Only include significant correlations
            if correlationScore.OverallConfidence > 0.3 {
                key := fmt.Sprintf("%s-%s", platformA, platformB)
                results[key] = CorrelationResult{
                    PlatformA:        platformA,
                    PlatformB:        platformB,
                    CorrelationScore: correlationScore,
                    Evidence:         pc.collectCorrelationEvidence(identityA, identityB),
                }
            }
        }
    }

    return results
}

func (pc *PlatformCorrelator) calculateMultiFactorCorrelation(identityA, identityB NormalizedIdentity) CorrelationScore {
    var score CorrelationScore
    
    // 1. Username similarity (25% weight)
    if usernameScore := pc.calculateUsernameSimilarity(identityA.Username, identityB.Username); usernameScore > 0 {
        score.UsernameSimilarity = usernameScore
        score.OverallConfidence += usernameScore * 0.25
    }

    // 2. Name similarity (20% weight)
    if nameScore := pc.calculateNameSimilarity(identityA.FullName, identityB.FullName); nameScore > 0 {
        score.NameSimilarity = nameScore
        score.OverallConfidence += nameScore * 0.20
    }

    // 3. Contact information match (25% weight)
    if contactScore := pc.calculateContactSimilarity(identityA, identityB); contactScore > 0 {
        score.ContactSimilarity = contactScore
        score.OverallConfidence += contactScore * 0.25
    }

    // 4. Behavioral pattern similarity (15% weight)
    if behaviorScore := pc.calculateBehavioralSimilarity(identityA.BehavioralPatterns, identityB.BehavioralPatterns); behaviorScore > 0 {
        score.BehavioralSimilarity = behaviorScore
        score.OverallConfidence += behaviorScore * 0.15
    }

    // 5. Temporal pattern similarity (15% weight)
    if temporalScore := pc.calculateTemporalSimilarity(identityA.TemporalPatterns, identityB.TemporalPatterns); temporalScore > 0 {
        score.TemporalSimilarity = temporalScore
        score.OverallConfidence += temporalScore * 0.15
    }

    return score
}

func (pc *PlatformCorrelator) calculateUsernameSimilarity(usernameA, usernameB string) float64 {
    if usernameA == "" || usernameB == "" {
        return 0
    }

    var similarity float64

    // Exact match
    if strings.ToLower(usernameA) == strings.ToLower(usernameB) {
        return 1.0
    }

    // Levenshtein distance for similar usernames
    distance := fuzzy.Levenshtein(&usernameA, &usernameB)
    maxLen := max(len(usernameA), len(usernameB))
    if maxLen > 0 {
        similarity = 1.0 - float64(distance)/float64(maxLen)
    }

    // Pattern-based similarity (common username variations)
    if pc.isCommonVariation(usernameA, usernameB) {
        similarity = math.Max(similarity, 0.8)
    }

    return similarity
}

func (pc *PlatformCorrelator) calculateNameSimilarity(nameA, nameB string) float64 {
    if nameA == "" || nameB == "" {
        return 0
    }

    // Normalize names
    normalizedA := pc.normalizeName(nameA)
    normalizedB := pc.normalizeName(nameB)

    // Exact match after normalization
    if normalizedA == normalizedB {
        return 1.0
    }

    // Token-based similarity
    tokensA := strings.Fields(normalizedA)
    tokensB := strings.Fields(normalizedB)

    commonTokens := 0
    for _, tokenA := range tokensA {
        for _, tokenB := range tokensB {
            if strings.EqualFold(tokenA, tokenB) {
                commonTokens++
                break
            }
        }
    }

    totalTokens := len(tokensA) + len(tokensB)
    if totalTokens > 0 {
        return 2.0 * float64(commonTokens) / float64(totalTokens)
    }

    return 0
}

func (pc *PlatformCorrelator) resolveGlobalIdentity(correlationResults map[string]CorrelationResult, identities map[string]NormalizedIdentity) map[string]PlatformIdentity {
    resolved := make(map[string]PlatformIdentity)
    
    // Create graph of platform connections
    connectionGraph := pc.buildConnectionGraph(correlationResults)
    
    // Find connected components (clusters of related platforms)
    components := pc.findConnectedComponents(connectionGraph)
    
    // For each component, resolve to a single identity
    for _, component := range components {
        if len(component) == 0 {
            continue
        }

        // Use the platform with highest confidence as anchor
        anchorPlatform := pc.findAnchorPlatform(component, identities)
        anchorIdentity := identities[anchorPlatform]

        // Resolve other platforms in component to anchor
        for _, platform := range component {
            if platform == anchorPlatform {
                resolved[platform] = pc.createPlatformIdentity(anchorIdentity)
                continue
            }

            // Merge identity information
            mergedIdentity := pc.mergeIdentities(anchorIdentity, identities[platform])
            resolved[platform] = pc.createPlatformIdentity(mergedIdentity)
        }
    }

    // Handle platforms not in any component
    for platform, identity := range identities {
        if _, exists := resolved[platform]; !exists {
            resolved[platform] = pc.createPlatformIdentity(identity)
        }
    }

    return resolved
}

func (pc *PlatformCorrelator) buildUnifiedTimeline(resolvedIdentity map[string]PlatformIdentity, platformData map[string]PlatformData) *UnifiedTimeline {
    timeline := &UnifiedTimeline{
        Events: make([]TimelineEvent, 0),
    }

    // Collect all events from all platforms
    allEvents := make([]TimelineEvent, 0)
    
    for platform, identity := range resolvedIdentity {
        if platformData, exists := platformData[platform]; exists {
            events := pc.extractPlatformEvents(platform, identity, platformData)
            allEvents = append(allEvents, events...)
        }
    }

    // Sort events by timestamp
    sort.Slice(allEvents, func(i, j int) bool {
        return allEvents[i].Timestamp.Before(allEvents[j].Timestamp)
    })

    // Merge duplicate events and resolve conflicts
    timeline.Events = pc.mergeTimelineEvents(allEvents)

    // Calculate timeline metrics
    timeline.Metrics = pc.calculateTimelineMetrics(timeline.Events)

    return timeline
}

func (pc *PlatformCorrelator) identifyRiskIndicators(resolvedIdentity map[string]PlatformIdentity, timeline *UnifiedTimeline) []RiskIndicator {
    indicators := make([]RiskIndicator, 0)

    // 1. Identity fragmentation risk
    if fragmentationScore := pc.calculateIdentityFragmentation(resolvedIdentity); fragmentationScore > 0.7 {
        indicators = append(indicators, RiskIndicator{
            Type:        "identity_fragmentation",
            Score:       fragmentationScore,
            Description: "High identity fragmentation across platforms",
            Confidence:  0.8,
        })
    }

    // 2. Behavioral anomaly risk
    if anomalyScore := pc.detectBehavioralAnomalies(timeline); anomalyScore > 0.6 {
        indicators = append(indicators, RiskIndicator{
            Type:        "behavioral_anomaly",
            Score:       anomalyScore,
            Description: "Unusual behavioral patterns detected",
            Confidence:  0.7,
        })
    }

    // 3. Temporal inconsistency risk
    if inconsistencyScore := pc.detectTemporalInconsistencies(timeline); inconsistencyScore > 0.5 {
        indicators = append(indicators, RiskIndicator{
            Type:        "temporal_inconsistency",
            Score:       inconsistencyScore,
            Description: "Inconsistent activity patterns across platforms",
            Confidence:  0.6,
        })
    }

    // 4. Platform correlation risk
    if correlationRisk := pc.assessCorrelationRisk(resolvedIdentity); correlationRisk > 0.4 {
        indicators = append(indicators, RiskIndicator{
            Type:        "low_correlation",
            Score:       correlationRisk,
            Description: "Low correlation between platform identities",
            Confidence:  0.5,
        })
    }

    return indicators
}
