// intelligence-engine/internal/report_generator/report_builder.go
package report_generator

import (
    "encoding/json"
    "fmt"
    "html/template"
    "os"
    "path/filepath"
    "time"
)

type ReportGenerator struct {
    templates    map[string]*template.Template
    exportFormats []ExportFormat
    reportDB     *ReportDatabase
}

type IntelligenceReport struct {
    ReportID      string                 `json:"report_id"`
    PhoneNumber   string                 `json:"phone_number"`
    GeneratedAt   time.Time              `json:"generated_at"`
    ReportType    string                 `json:"report_type"`
    ExecutiveSummary *ExecutiveSummary     `json:"executive_summary"`
    EmailDiscovery *EmailDiscoveryReport  `json:"email_discovery"`
    SocialGraph   *SocialGraphReport      `json:"social_graph"`
    RiskAssessment *RiskAssessment        `json:"risk_assessment"`
    BehavioralAnalysis *BehavioralAnalysis `json:"behavioral_analysis"`
    Recommendations []Recommendation       `json:"recommendations"`
    RawData       map[string]interface{} `json:"raw_data"`
    Confidence    float64                `json:"confidence"`
}

type ExecutiveSummary struct {
    OverallRisk     string   `json:"overall_risk"`
    KeyFindings     []string `json:"key_findings"`
    CriticalAlerts  []string `json:"critical_alerts"`
    Recommendations []string `json:"recommendations"`
}

func NewReportGenerator() *ReportGenerator {
    rg := &ReportGenerator{
        templates:    make(map[string]*template.Template),
        exportFormats: []ExportFormat{PDF, HTML, JSON, CSV},
        reportDB:     NewReportDatabase(),
    }
    rg.loadTemplates()
    return rg
}

// GenerateComprehensiveReport creates a complete intelligence report
func (rg *ReportGenerator) GenerateComprehensiveReport(phoneNumber string, intelligence *PhoneIntelligence) (*IntelligenceReport, error) {
    report := &IntelligenceReport{
        ReportID:    rg.generateReportID(),
        PhoneNumber: phoneNumber,
        GeneratedAt: time.Now(),
        ReportType:  "comprehensive",
        RawData:     make(map[string]interface{}),
    }

    // Generate all report sections
    var wg sync.WaitGroup
    errors := make(chan error, 5)

    // Email discovery section
    wg.Add(1)
    go func() {
        defer wg.Done()
        emailReport, err := rg.generateEmailDiscoverySection(intelligence)
        if err != nil {
            errors <- fmt.Errorf("email discovery: %w", err)
            return
        }
        report.EmailDiscovery = emailReport
    }()

    // Social graph section
    wg.Add(1)
    go func() {
        defer wg.Done()
        socialReport, err := rg.generateSocialGraphSection(intelligence)
        if err != nil {
            errors <- fmt.Errorf("social graph: %w", err)
            return
        }
        report.SocialGraph = socialReport
    }()

    // Risk assessment section
    wg.Add(1)
    go func() {
        defer wg.Done()
        riskReport, err := rg.generateRiskAssessmentSection(intelligence)
        if err != nil {
            errors <- fmt.Errorf("risk assessment: %w", err)
            return
        }
        report.RiskAssessment = riskReport
    }()

    // Behavioral analysis section
    wg.Add(1)
    go func() {
        defer wg.Done()
        behaviorReport, err := rg.generateBehavioralAnalysisSection(intelligence)
        if err != nil {
            errors <- fmt.Errorf("behavioral analysis: %w", err)
            return
        }
        report.BehavioralAnalysis = behaviorReport
    }()

    // Executive summary (depends on other sections)
    wg.Add(1)
    go func() {
        defer wg.Done()
        // Wait for other sections to complete
        wg.Wait()
        close(errors)
        
        execSummary := rg.generateExecutiveSummary(report)
        report.ExecutiveSummary = execSummary
        report.Recommendations = rg.generateRecommendations(report)
        report.Confidence = rg.calculateOverallConfidence(report)
    }()

    // Collect any errors
    var errorList []string
    for err := range errors {
        errorList = append(errorList, err.Error())
    }

    if len(errorList) > 0 {
        return report, fmt.Errorf("report generation completed with errors: %v", errorList)
    }

    // Store report in database
    if err := rg.reportDB.StoreReport(report); err != nil {
        return report, fmt.Errorf("failed to store report: %w", err)
    }

    return report, nil
}

func (rg *ReportGenerator) generateExecutiveSummary(report *IntelligenceReport) *ExecutiveSummary {
    summary := &ExecutiveSummary{
        KeyFindings:    make([]string, 0),
        CriticalAlerts: make([]string, 0),
        Recommendations: make([]string, 0),
    }

    // Determine overall risk level
    if report.RiskAssessment.OverallScore >= 0.8 {
        summary.OverallRisk = "HIGH"
    } else if report.RiskAssessment.OverallScore >= 0.5 {
        summary.OverallRisk = "MEDIUM" 
    } else {
        summary.OverallRisk = "LOW"
    }

    // Extract key findings
    if report.EmailDiscovery != nil && report.EmailDiscovery.TotalFound > 0 {
        summary.KeyFindings = append(summary.KeyFindings, 
            fmt.Sprintf("Discovered %d email addresses", report.EmailDiscovery.TotalFound))
    }

    if report.SocialGraph != nil {
        summary.KeyFindings = append(summary.KeyFindings,
            fmt.Sprintf("Mapped %d social connections", report.SocialGraph.TotalConnections))
    }

    // Add critical alerts from risk assessment
    for _, factor := range report.RiskAssessment.Factors {
        if factor.Score >= 0.8 {
            summary.CriticalAlerts = append(summary.CriticalAlerts,
                fmt.Sprintf("High risk: %s", factor.Description))
        }
    }

    // Generate recommendations
    summary.Recommendations = rg.generateSummaryRecommendations(report)

    return summary
}

// Export report in multiple formats
func (rg *ReportGenerator) ExportReport(report *IntelligenceReport, format ExportFormat) ([]byte, error) {
    switch format {
    case JSON:
        return json.MarshalIndent(report, "", "  ")
    case HTML:
        return rg.exportHTML(report)
    case PDF:
        return rg.exportPDF(report)
    case CSV:
        return rg.exportCSV(report)
    default:
        return nil, fmt.Errorf("unsupported export format: %s", format)
    }
}

func (rg *ReportGenerator) exportHTML(report *IntelligenceReport) ([]byte, error) {
    template, exists := rg.templates["comprehensive"]
    if !exists {
        return nil, fmt.Errorf("template not found")
    }

    var buf bytes.Buffer
    if err := template.Execute(&buf, report); err != nil {
        return nil, fmt.Errorf("template execution failed: %w", err)
    }

    return buf.Bytes(), nil
}
// ... continuing from the previous code ...

            RiskScore:     rg.calculatePlatformRisk(platform, profiles),
        }
        analysis.PlatformPresence = append(analysis.PlatformPresence, presence)
    }

    // Analyze data breaches
    if intelligenceData.DataBreaches != nil {
        analysis.DataBreaches = &BreachAnalysis{
            TotalBreaches:      len(intelligenceData.DataBreaches.Breaches),
            CompromisedData:    rg.aggregateCompromisedData(intelligenceData.DataBreaches),
            FirstBreach:        rg.findFirstBreach(intelligenceData.DataBreaches),
            LatestBreach:       rg.findLatestBreach(intelligenceData.DataBreaches),
            HighRiskBreaches:   rg.countHighRiskBreaches(intelligenceData.DataBreaches),
        }
    }

    // Calculate identity confidence metrics
    analysis.IdentityConfidence = rg.calculateIdentityConfidence(analysis)

    return analysis, nil
}

func (rg *ReportGenerator) generateSocialNetworkAnalysis(intelligenceData *IntelligenceData) (*SocialNetworkAnalysis, error) {
    analysis := &SocialNetworkAnalysis{
        NetworkMap:      make(map[string][]SocialConnection),
        InfluenceMetrics: &InfluenceMetrics{},
        CommunityStructure: &CommunityStructure{},
    }

    // Build network connections
    if intelligenceData.SocialConnections != nil {
        for platform, connections := range intelligenceData.SocialConnections.Connections {
            analysis.NetworkMap[platform] = rg.transformConnections(connections)
        }
    }

    // Calculate influence metrics
    analysis.InfluenceMetrics = rg.calculateInfluenceMetrics(intelligenceData)

    // Analyze community structure
    analysis.CommunityStructure = rg.analyzeCommunityStructure(analysis.NetworkMap)

    return analysis, nil
}

func (rg *ReportGenerator) generateBehavioralAnalysis(intelligenceData *IntelligenceData) (*BehavioralAnalysis, error) {
    analysis := &BehavioralAnalysis{
        ActivityPatterns: &ActivityPatterns{},
        ContentAnalysis:  &ContentAnalysis{},
        RiskIndicators:   make([]RiskIndicator, 0),
    }

    // Analyze activity patterns
    analysis.ActivityPatterns = rg.analyzeActivityPatterns(intelligenceData)

    // Perform content analysis
    analysis.ContentAnalysis = rg.analyzeContentBehavior(intelligenceData)

    // Identify risk indicators
    analysis.RiskIndicators = rg.identifyBehavioralRisks(analysis)

    return analysis, nil
}

func (rg *ReportGenerator) generateThreatAssessment(intelligenceData *IntelligenceData) (*ThreatAssessment, error) {
    assessment := &ThreatAssessment{
        ThreatLevel:      "LOW",
        RiskFactors:      make([]RiskFactor, 0),
        MitigationStrategies: make([]string, 0),
        MonitoringRecommendations: make([]string, 0),
    }

    // Assess various threat vectors
    riskFactors := rg.assembleRiskFactors(intelligenceData)
    assessment.RiskFactors = riskFactors

    // Calculate overall threat level
    assessment.ThreatLevel = rg.calculateThreatLevel(riskFactors)

    // Generate mitigation strategies
    assessment.MitigationStrategies = rg.generateMitigationStrategies(riskFactors)

    // Create monitoring recommendations
    assessment.MonitoringRecommendations = rg.generateMonitoringRecommendations(assessment)

    return assessment, nil
}

func (rg *ReportGenerator) generateRecommendations(report *IntelligenceReport) []Recommendation {
    recommendations := make([]Recommendation, 0)

    // Generate recommendations based on risk assessment
    if report.ExecutiveSummary.RiskAssessment.OverallRisk >= 0.7 {
        recommendations = append(recommendations, Recommendation{
            Type:        "IMMEDIATE_ACTION",
            Title:       "High Risk Identified - Immediate Review Required",
            Description: "Subject demonstrates multiple high-risk indicators requiring immediate attention",
            Priority:    "CRITICAL",
            Actions:     []string{"Escalate to security team", "Initiate enhanced monitoring", "Conduct deeper investigation"},
        })
    }

    // Identity protection recommendations
    if report.DetailedAnalysis.IdentityAnalysis.DataBreaches != nil &&
        report.DetailedAnalysis.IdentityAnalysis.DataBreaches.TotalBreaches > 0 {
        recommendations = append(recommendations, Recommendation{
            Type:        "IDENTITY_PROTECTION",
            Title:       "Data Breach Exposure Detected",
            Description: fmt.Sprintf("Subject appears in %d known data breaches", report.DetailedAnalysis.IdentityAnalysis.DataBreaches.TotalBreaches),
            Priority:    "HIGH",
            Actions:     []string{"Recommend password changes", "Enable two-factor authentication", "Monitor for identity theft"},
        })
    }

    // Social network recommendations
    if len(report.DetailedAnalysis.SocialNetwork.RiskConnections) > 0 {
        recommendations = append(recommendations, Recommendation{
            Type:        "NETWORK_ANALYSIS",
            Title:       "Risky Social Connections Identified",
            Description: fmt.Sprintf("Found %d potentially risky social connections", len(report.DetailedAnalysis.SocialNetwork.RiskConnections)),
            Priority:    "MEDIUM",
            Actions:     []string{"Review connection patterns", "Monitor for suspicious activity", "Assess relationship risks"},
        })
    }

    // Behavioral risk recommendations
    behavioralRisks := len(report.DetailedAnalysis.BehavioralAnalysis.RiskIndicators)
    if behavioralRisks > 0 {
        recommendations = append(recommendations, Recommendation{
            Type:        "BEHAVIORAL_MONITORING",
            Title:       "Behavioral Risk Indicators Present",
            Description: fmt.Sprintf("Identified %d behavioral risk indicators", behavioralRisks),
            Priority:    "MEDIUM",
            Actions:     []string{"Establish behavioral baseline", "Monitor for changes", "Document concerning patterns"},
        })
    }

    return recommendations
}

func (rg *ReportGenerator) calculateOverallConfidence(report *IntelligenceReport) float64 {
    var totalConfidence float64
    var weightCount int

    // Identity confidence (30% weight)
    if report.DetailedAnalysis.IdentityAnalysis != nil &&
        report.DetailedAnalysis.IdentityAnalysis.IdentityConfidence != nil {
        totalConfidence += report.DetailedAnalysis.IdentityAnalysis.IdentityConfidence.OverallScore * 0.3
        weightCount++
    }

    // Social network confidence (25% weight)
    if report.DetailedAnalysis.SocialNetwork != nil &&
        report.DetailedAnalysis.SocialNetwork.ConfidenceScore > 0 {
        totalConfidence += report.DetailedAnalysis.SocialNetwork.ConfidenceScore * 0.25
        weightCount++
    }

    // Behavioral analysis confidence (20% weight)
    if report.DetailedAnalysis.BehavioralAnalysis != nil &&
        report.DetailedAnalysis.BehavioralAnalysis.Confidence > 0 {
        totalConfidence += report.DetailedAnalysis.BehavioralAnalysis.Confidence * 0.2
        weightCount++
    }

    // Data completeness (25% weight)
    completenessScore := rg.calculateDataCompleteness(report)
    totalConfidence += completenessScore * 0.25
    weightCount++

    // Normalize if we have fewer components
    if weightCount < 4 {
        totalConfidence = totalConfidence / float64(weightCount) * 4
    }

    return math.Min(totalConfidence, 1.0)
}

// Helper methods for analysis calculations
func (rg *ReportGenerator) analyzeBreachInvolvement(discovery *EmailDiscovery) *BreachInvolvement {
    involvement := &BreachInvolvement{
        BreachedEmails: make([]string, 0),
        TotalBreaches:  0,
    }

    // Implementation would check breach databases
    // This is a simplified version
    return involvement
}

func (rg *ReportGenerator) countVerifiedProfiles(profiles []SocialProfile) int {
    count := 0
    for _, profile := range profiles {
        if profile.Verified {
            count++
        }
    }
    return count
}

func (rg *ReportGenerator) calculateActivityLevel(profiles []SocialProfile) string {
    if len(profiles) == 0 {
        return "LOW"
    }

    // Simplified activity calculation
    totalPosts := 0
    for _, profile := range profiles {
        totalPosts += profile.PostCount
    }

    avgPosts := totalPosts / len(profiles)
    switch {
    case avgPosts > 1000:
        return "VERY_HIGH"
    case avgPosts > 100:
        return "HIGH"
    case avgPosts > 10:
        return "MEDIUM"
    default:
        return "LOW"
    }
}

func (rg *ReportGenerator) calculatePlatformRisk(platform string, profiles []SocialProfile) float64 {
    // Platform-specific risk calculations
    baseRisk := map[string]float64{
        "facebook":  0.3,
        "twitter":   0.4,
        "instagram": 0.3,
        "linkedin":  0.2,
        "reddit":    0.5,
    }

    risk := baseRisk[platform]
    
    // Adjust based on activity and verification
    activityLevel := rg.calculateActivityLevel(profiles)
    switch activityLevel {
    case "VERY_HIGH":
        risk *= 1.5
    case "HIGH":
        risk *= 1.2
    case "LOW":
        risk *= 0.8
    }

    return math.Min(risk, 1.0)
}

func (rg *ReportGenerator) generateReportID() string {
    return fmt.Sprintf("INTEL-%s-%d", time.Now().Format("20060102"), rand.Intn(1000))
}

// Export methods
func (rg *ReportGenerator) ExportToPDF(report *IntelligenceReport, filePath string) error {
    return rg.exporter.ExportToPDF(report, filePath)
}

func (rg *ReportGenerator) ExportToHTML(report *IntelligenceReport, filePath string) error {
    return rg.exporter.ExportToHTML(report, filePath)
}

func (rg *ReportGenerator) ExportToJSON(report *IntelligenceReport, filePath string) error {
    return rg.exporter.ExportToJSON(report, filePath)
}

// Additional helper methods would be implemented here...
func (rg *ReportGenerator) aggregateCompromisedData(breaches *DataBreachData) []string {
    // Implementation for aggregating compromised data types
    return []string{}
}

func (rg *ReportGenerator) findFirstBreach(breaches *DataBreachData) time.Time {
    // Implementation for finding first breach
    return time.Time{}
}

func (rg *ReportGenerator) findLatestBreach(breaches *DataBreachData) time.Time {
    // Implementation for finding latest breach
    return time.Now()
}

func (rg *ReportGenerator) countHighRiskBreaches(breaches *DataBreachData) int {
    // Implementation for counting high-risk breaches
    return 0
}

func (rg *ReportGenerator) calculateIdentityConfidence(analysis *IdentityAnalysis) *ConfidenceMetrics {
    // Implementation for calculating identity confidence
    return &ConfidenceMetrics{
        OverallScore: 0.8,
        DataSources:  5,
        Verification: 0.7,
    }
}

func (rg *ReportGenerator) transformConnections(connections []SocialConnection) []SocialConnection {
    // Implementation for transforming connections
    return connections
}

func (rg *ReportGenerator) calculateInfluenceMetrics(data *IntelligenceData) *InfluenceMetrics {
    // Implementation for calculating influence metrics
    return &InfluenceMetrics{
        FollowerCount:    0,
        EngagementRate:   0.0,
        NetworkDensity:  0.0,
        InfluenceScore:  0.0,
    }
}

func (rg *ReportGenerator) analyzeCommunityStructure(networkMap map[string][]SocialConnection) *CommunityStructure {
    // Implementation for community structure analysis
    return &CommunityStructure{
        Communities:    0,
        LargestSize:   0,
        Modularity:    0.0,
    }
}

func (rg *ReportGenerator) analyzeActivityPatterns(data *IntelligenceData) *ActivityPatterns {
    // Implementation for activity pattern analysis
    return &ActivityPatterns{
        PeakHours:      []int{},
        PostFrequency:  0.0,
        ActivityScore:  0.0,
    }
}

func (rg *ReportGenerator) analyzeContentBehavior(data *IntelligenceData) *ContentAnalysis {
    // Implementation for content analysis
    return &ContentAnalysis{
        Sentiment:      0.0,
        Topics:         []string{},
        RiskKeywords:   []string{},
    }
}

func (rg *ReportGenerator) identifyBehavioralRisks(analysis *BehavioralAnalysis) []RiskIndicator {
    // Implementation for identifying behavioral risks
    return []RiskIndicator{}
}

func (rg *ReportGenerator) assembleRiskFactors(data *IntelligenceData) []RiskFactor {
    // Implementation for assembling risk factors
    return []RiskFactor{}
}

func (rg *ReportGenerator) calculateThreatLevel(riskFactors []RiskFactor) string {
    // Implementation for calculating threat level
    return "LOW"
}

func (rg *ReportGenerator) generateMitigationStrategies(riskFactors []RiskFactor) []string {
    // Implementation for generating mitigation strategies
    return []string{}
}

func (rg *ReportGenerator) generateMonitoringRecommendations(assessment *ThreatAssessment) []string {
    // Implementation for generating monitoring recommendations
    return []string{}
}
            return {
                success: false,
                platform: 'instagram',
                error: 'USER_NOT_FOUND',
                timestamp: new Date()
            };
            
        } catch (error) {
            return {
                success: false,
                platform: 'instagram',
                error: error.message,
                timestamp: new Date()
            };
        }
    }
    
    private async fetchProfile(page: Page, username: string): Promise<ProfileData> {
        // Navigate to profile page
        await page.goto(`${this.baseUrl}/${username}/`);
        
        // Extract profile information
        const profileData = await page.evaluate(() => {
            const nameElement = document.querySelector('h1');
            const bioElement = document.querySelector('.bio');
            const postsElement = document.querySelector('li span');
            const followersElement = document.querySelector('a[href*="followers"] span');
            const followingElement = document.querySelector('a[href*="following"] span');
            
            return {
                fullName: nameElement?.textContent?.trim(),
                biography: bioElement?.textContent?.trim(),
                postsCount: postsElement?.textContent ? parseInt(postsElement.textContent) : 0,
                followersCount: followersElement?.textContent ? parseInt(followersElement.textContent) : 0,
                followingCount: followingElement?.textContent ? parseInt(followingElement.textContent) : 0,
                isPrivate: document.querySelector('.profile-private') !== null,
                isVerified: document.querySelector('.verified-badge') !== null
            };
        });
        
        return profileData;
    }
}

func (rg *ReportGenerator) calculateDataCompleteness(report *IntelligenceReport) float64 {
    // Implementation for calculating data completeness
    return 0.8
}
