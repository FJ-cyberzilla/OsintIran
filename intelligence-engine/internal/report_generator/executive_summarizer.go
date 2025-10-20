package report_generator

type ExecutiveSummarizer struct {
    templateEngine *TemplateEngine
    riskCalculator *RiskCalculator
}

type ExecutiveSummary struct {
    Overview       string         `json:"overview"`
    KeyFindings    []KeyFinding   `json:"key_findings"`
    RiskAssessment *RiskAssessment `json:"risk_assessment"`
    Recommendations []string       `json:"recommendations"`
    Conclusion     string         `json:"conclusion"`
}

type KeyFinding struct {
    Category      string  `json:"category"`
    Title         string  `json:"title"`
    Description   string  `json:"description"`
    Confidence    float64 `json:"confidence"`
    Impact        string  `json:"impact"` // HIGH, MEDIUM, LOW
    Evidence      []string `json:"evidence"`
}

func (es *ExecutiveSummarizer) GenerateExecutiveSummary(report *IntelligenceReport) (*ExecutiveSummary, error) {
    summary := &ExecutiveSummary{
        KeyFindings:    make([]KeyFinding, 0),
        Recommendations: make([]string, 0),
    }
    
    // Generate overview
    summary.Overview = es.generateOverview(report)
    
    // Extract key findings
    summary.KeyFindings = es.extractKeyFindings(report)
    
    // Assess risks
    summary.RiskAssessment = es.assessRisks(report)
    
    // Generate recommendations
    summary.Recommendations = es.generateRecommendations(summary.KeyFindings, summary.RiskAssessment)
    
    // Write conclusion
    summary.Conclusion = es.generateConclusion(summary)
    
    return summary, nil
}

func (es *ExecutiveSummarizer) extractKeyFindings(report *IntelligenceReport) []KeyFinding {
    var findings []KeyFinding
    
    // Identity findings
    if report.DetailedAnalysis.IdentityAnalysis != nil {
        findings = append(findings, es.analyzeIdentityFindings(report.DetailedAnalysis.IdentityAnalysis)...)
    }
    
    // Social network findings
    if report.DetailedAnalysis.SocialNetwork != nil {
        findings = append(findings, es.analyzeSocialFindings(report.DetailedAnalysis.SocialNetwork)...)
    }
    
    // Behavioral findings
    if report.DetailedAnalysis.BehavioralAnalysis != nil {
        findings = append(findings, es.analyzeBehavioralFindings(report.DetailedAnalysis.BehavioralAnalysis)...)
    }
    
    // Threat findings
    if report.DetailedAnalysis.ThreatAssessment != nil {
        findings = append(findings, es.analyzeThreatFindings(report.DetailedAnalysis.ThreatAssessment)...)
    }
    
    return findings
}
