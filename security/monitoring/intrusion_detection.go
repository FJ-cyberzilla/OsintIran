// security/monitoring/intrusion_detection.go
package monitoring

import (
    "fmt"
    "log"
    "time"
)

type IntrusionDetectionSystem struct {
    suspiciousActivities []SuspiciousActivity
    alertThreshold       int
    monitoringEnabled    bool
}

type SuspiciousActivity struct {
    Type        string
    Source      string
    Timestamp   time.Time
    Severity    string
    Details     map[string]interface{}
}

func NewIDS() *IntrusionDetectionSystem {
    ids := &IntrusionDetectionSystem{
        alertThreshold:    5,
        monitoringEnabled: true,
    }
    
    go ids.monitorActivities()
    return ids
}

func (ids *IntrusionDetectionSystem) LogSuspiciousActivity(activityType, source string, severity string, details map[string]interface{}) {
    activity := SuspiciousActivity{
        Type:      activityType,
        Source:    source,
        Timestamp: time.Now(),
        Severity:  severity,
        Details:   details,
    }
    
    ids.suspiciousActivities = append(ids.suspiciousActivities, activity)
    
    // Check if alert threshold is reached
    if ids.shouldAlert(activityType, source) {
        ids.triggerAlert(activity)
    }
}

func (ids *IntrusionDetectionSystem) shouldAlert(activityType, source string) bool {
    recentActivities := 0
    cutoffTime := time.Now().Add(-5 * time.Minute)
    
    for _, activity := range ids.suspiciousActivities {
        if activity.Timestamp.After(cutoffTime) && 
           activity.Type == activityType && 
           activity.Source == source {
            recentActivities++
        }
    }
    
    return recentActivities >= ids.alertThreshold
}

func (ids *IntrusionDetectionSystem) triggerAlert(activity SuspiciousActivity) {
    alertMessage := fmt.Sprintf(
        "INTRUSION ALERT: %s from %s - Severity: %s",
        activity.Type,
        activity.Source,
        activity.Severity,
    )
    
    log.Printf("🚨 %s", alertMessage)
    
    // Additional alert actions (email, SMS, etc.)
    ids.notifySecurityTeam(alertMessage, activity)
}

func (ids *IntrusionDetectionSystem) monitorActivities() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        ids.analyzePatterns()
        ids.cleanupOldActivities()
    }
}
