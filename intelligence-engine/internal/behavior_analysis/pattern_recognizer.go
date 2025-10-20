// intelligence-engine/internal/behavior_analysis/pattern_recognizer.go
package behavior_analysis

import (
    "math"
    "sort"
    "time"
)

type PatternRecognizer struct {
    behaviorModels map[string]*BehaviorModel
    anomalyDB      *AnomalyDatabase
    mlModel        *MLPredictor
}

type BehaviorModel struct {
    UserID         string                    `json:"user_id"`
    ActivityPatterns map[string]ActivityPattern `json:"activity_patterns"`
    CommunicationStyle *CommunicationStyle     `json:"communication_style"`
    RiskIndicators []RiskIndicator           `json:"risk_indicators"`
    Confidence     float64                   `json:"confidence"`
    LastUpdated    time.Time                 `json:"last_updated"`
}

type ActivityPattern struct {
    Platform      string    `json:"platform"`
    TypicalHours  []int     `json:"typical_hours"` // 0-23
    PostFrequency float64   `json:"post_frequency"` // posts per day
    SessionLength time.Duration `json:"session_length"`
    ContentTypes  []string  `json:"content_types"`
}

type RiskScorer struct {
    patternRecognizer *PatternRecognizer
    ruleEngine        *RuleEngine
    threatIntelligence *ThreatIntelligence
}

func NewRiskScorer() *RiskScorer {
    return &RiskScorer{
        patternRecognizer: NewPatternRecognizer(),
        ruleEngine:        NewRuleEngine(),
        threatIntelligence: NewThreatIntelligence(),
    }
}

// Comprehensive risk assessment
func (rs *RiskScorer) CalculateComprehensiveRisk(phoneIntelligence *PhoneIntelligence) *RiskAssessment {
    assessment := &RiskAssessment{
        PhoneNumber: phoneIntelligence.PhoneNumber,
        Timestamp:   time.Now(),
        Factors:     make([]RiskFactor, 0),
    }

    // 1. Behavioral Anomaly Risk
    behavioralRisk := rs.assessBehavioralRisk(phoneIntelligence)
    assessment.Factors = append(assessment.Factors, behavioralRisk...)

    // 2. Social Graph Risk
    socialRisk := rs.assessSocialGraphRisk(phoneIntelligence)
    assessment.Factors = append(assessment.Factors, socialRisk...)

    // 3. Digital Footprint Risk
    footprintRisk := rs.assessDigitalFootprintRisk(phoneIntelligence)
    assessment.Factors = append(assessment.Factors, footprintRisk...)

    // 4. Threat Intelligence Risk
    threatRisk := rs.assessThreatIntelligenceRisk(phoneIntelligence)
    assessment.Factors = append(assessment.Factors, threatRisk...)

    // Calculate overall score
    assessment.OverallScore = rs.calculateOverallScore(assessment.Factors)
    assessment.RiskLevel = rs.determineRiskLevel(assessment.OverallScore)
    assessment.Confidence = rs.calculateConfidence(assessment.Factors)

    return assessment
}

func (rs *RiskScorer) assessBehavioralRisk(intel *PhoneIntelligence) []RiskFactor {
    var factors []RiskFactor
    
    // Activity pattern analysis
    for platform, activities := range intel.ActivityPatterns {
        pattern := rs.patternRecognizer.AnalyzeActivityPattern(activities)
        
        // Check for anomalies
        if anomalyScore := rs.detectActivityAnomalies(pattern); anomalyScore > 0.7 {
            factors = append(factors, RiskFactor{
                Type:        "behavioral_anomaly",
                Description: fmt.Sprintf("Unusual activity pattern on %s", platform),
                Score:       anomalyScore,
                Evidence:    []string{"activity_pattern_analysis"},
            })
        }
        
        // Check for bot-like behavior
        if botScore := rs.detectAutomatedBehavior(pattern); botScore > 0.8 {
            factors = append(factors, RiskFactor{
                Type:        "automated_behavior",
                Description: fmt.Sprintf("Potential automated activity on %s", platform),
                Score:       botScore,
                Evidence:    []string{"behavior_pattern_analysis"},
            })
        }
    }
    
    return factors
}

func (rs *RiskScorer) assessSocialGraphRisk(intel *PhoneIntelligence) []RiskFactor {
    var factors []RiskFactor
    
    graph := rs.patternRecognizer.BuildSocialGraph(intel)
    
    // Network centrality risk
    if centrality := graph.CalculateCentrality(); centrality > 0.8 {
        factors = append(factors, RiskFactor{
            Type:        "high_centrality",
            Description: "Highly central position in social network",
            Score:       centrality,
            Evidence:    []string{"network_analysis"},
        })
    }
    
    // Suspicious connection patterns
    if suspiciousConnections := rs.detectSuspiciousConnections(graph); len(suspiciousConnections) > 0 {
        factors = append(factors, RiskFactor{
            Type:        "suspicious_connections",
            Description: fmt.Sprintf("Found %d suspicious network connections", len(suspiciousConnections)),
            Score:       math.Min(float64(len(suspiciousConnections))/10.0, 1.0),
            Evidence:    []string{"connection_pattern_analysis"},
        })
    }
    
    return factors
}

func (rs *RiskScorer) detectActivityAnomalies(pattern *ActivityPattern) float64 {
    // Implement anomaly detection algorithm
    // This would use statistical methods to detect deviations from normal patterns
    
    var anomalyScore float64
    
    // Check for unusual posting times
    if rs.isUnusualPostingTime(pattern) {
        anomalyScore += 0.3
    }
    
    // Check for inconsistent activity patterns
    if rs.hasInconsistentPattern(pattern) {
        anomalyScore += 0.4
    }
    
    // Check for robotic timing
    if rs.hasRoboticTiming(pattern) {
        anomalyScore += 0.3
    }
    
    return anomalyScore
}

// ML-based risk prediction
func (rs *RiskScorer) predictRiskWithML(features *RiskFeatures) float64 {
    // Extract features for ML model
    featureVector := rs.extractFeatures(features)
    
    // Use trained model to predict risk
    prediction := rs.mlModel.Predict(featureVector)
    
    return prediction.RiskScore
}
