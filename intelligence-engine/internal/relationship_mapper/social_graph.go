// intelligence-engine/internal/relationship_mapper/social_graph.go
package relationship_mapper

import (
    "encoding/json"
    "fmt"
    "math"
    "sort"
)

type SocialGraph struct {
    nodes map[string]*PersonNode
    edges map[string]*RelationshipEdge
}

type PersonNode struct {
    ID           string                 `json:"id"`
    Type         string                 `json:"type"` // "phone", "email", "social_profile"
    Value        string                 `json:"value"`
    Profiles     map[string]SocialProfile `json:"profiles"`
    Metadata     map[string]interface{} `json:"metadata"`
    Centrality   float64                `json:"centrality"`
    RiskScore    float64                `json:"risk_score"`
    DiscoveredAt time.Time              `json:"discovered_at"`
}

type RelationshipEdge struct {
    Source      string    `json:"source"`
    Target      string    `json:"target"`
    Type        string    `json:"type"` // "same_person", "family", "friend", "professional"
    Strength    float64   `json:"strength"` // 0.0 - 1.0
    Evidence    []string  `json:"evidence"`
    FirstSeen   time.Time `json:"first_seen"`
    LastSeen    time.Time `json:"last_seen"`
}

type ConnectionAnalyzer struct {
    graph         *SocialGraph
    similarity    *SimilarityScorer
    patternDB     *PatternDatabase
}

func NewConnectionAnalyzer() *ConnectionAnalyzer {
    return &ConnectionAnalyzer{
        graph:      NewSocialGraph(),
        similarity: NewSimilarityScorer(),
        patternDB:  NewPatternDatabase(),
    }
}

// BuildSocialGraphFromIntelligence - Creates relationship graph from intelligence data
func (ca *ConnectionAnalyzer) BuildSocialGraphFromIntelligence(phoneNumber string, intelligence *PhoneIntelligence) *SocialGraph {
    // Start with the phone number as root node
    rootNode := &PersonNode{
        ID:    ca.generateNodeID("phone", phoneNumber),
        Type:  "phone",
        Value: phoneNumber,
        Profiles: make(map[string]SocialProfile),
    }
    ca.graph.AddNode(rootNode)

    // Add all discovered emails as nodes
    for _, email := range intelligence.DiscoveredEmails {
        emailNode := &PersonNode{
            ID:    ca.generateNodeID("email", email.Email),
            Type:  "email", 
            Value: email.Email,
        }
        ca.graph.AddNode(emailNode)
        
        // Create edge between phone and email
        ca.graph.AddEdge(&RelationshipEdge{
            Source:   rootNode.ID,
            Target:   emailNode.ID,
            Type:     "same_person",
            Strength: email.Confidence,
            Evidence: []string{"email_discovery"},
        })
    }

    // Add social profiles and their connections
    for platform, profiles := range intelligence.SocialProfiles {
        for _, profile := range profiles {
            profileNode := &PersonNode{
                ID:   ca.generateNodeID("social", profile.Username+"@"+platform),
                Type: "social_profile",
                Value: profile.Username,
                Profiles: map[string]SocialProfile{
                    platform: profile,
                },
            }
            ca.graph.AddNode(profileNode)
            
            // Connect to likely email nodes
            for _, email := range intelligence.DiscoveredEmails {
                if ca.isLikelyMatch(email.Email, profile.Username, platform) {
                    ca.graph.AddEdge(&RelationshipEdge{
                        Source:   emailNode.ID,
                        Target:   profileNode.ID,
                        Type:     "same_person",
                        Strength: 0.8,
                        Evidence: []string{"username_email_correlation"},
                    })
                }
            }
        }
    }

    // Analyze connections between different profiles
    ca.analyzeCrossPlatformConnections()
    
    // Calculate network metrics
    ca.calculateCentrality()
    ca.calculateRiskScores()

    return ca.graph
}

// AnalyzeCrossPlatformConnections finds relationships across different platforms
func (ca *ConnectionAnalyzer) analyzeCrossPlatformConnections() {
    nodes := ca.graph.GetAllNodes()
    
    for i, node1 := range nodes {
        for j, node2 := range nodes {
            if i >= j { // Avoid duplicate comparisons
                continue
            }
            
            similarity := ca.calculateNodeSimilarity(node1, node2)
            if similarity > 0.7 { // Threshold for connection
                ca.graph.AddEdge(&RelationshipEdge{
                    Source:   node1.ID,
                    Target:   node2.ID,
                    Type:     ca.determineRelationshipType(node1, node2),
                    Strength: similarity,
                    Evidence: []string{"profile_similarity"},
                })
            }
        }
    }
}

func (ca *ConnectionAnalyzer) calculateNodeSimilarity(node1, node2 *PersonNode) float64 {
    var similarity float64
    
    // Username similarity
    if node1.Type == "social_profile" && node2.Type == "social_profile" {
        similarity += ca.similarity.CompareUsernames(node1.Value, node2.Value) * 0.3
    }
    
    // Profile information similarity
    similarity += ca.compareProfileMetadata(node1, node2) * 0.4
    
    // Connection pattern similarity  
    similarity += ca.compareConnectionPatterns(node1, node2) * 0.3
    
    return similarity
}

func (ca *ConnectionAnalyzer) determineRelationshipType(node1, node2 *PersonNode) string {
    similarity := ca.calculateNodeSimilarity(node1, node2)
    
    switch {
    case similarity > 0.9:
        return "same_person"
    case similarity > 0.7:
        return "close_connection" 
    case similarity > 0.5:
        return "known_connection"
    default:
        return "weak_connection"
    }
}

// Calculate network centrality metrics
func (ca *ConnectionAnalyzer) calculateCentrality() {
    nodes := ca.graph.GetAllNodes()
    
    for _, node := range nodes {
        // Degree centrality (number of connections)
        degree := len(ca.graph.GetEdgesForNode(node.ID))
        maxPossibleDegree := len(nodes) - 1
        
        if maxPossibleDegree > 0 {
            node.Centrality = float64(degree) / float64(maxPossibleDegree)
        }
    }
}
