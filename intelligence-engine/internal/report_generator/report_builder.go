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
