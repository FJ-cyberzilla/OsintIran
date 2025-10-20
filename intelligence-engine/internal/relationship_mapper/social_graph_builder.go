// intelligence-engine/internal/relationship_mapper/social_graph_builder.go
package relationship_mapper

import (
    "encoding/json"
    "fmt"
    "math"
    "sort"
    "time"

    "github.com/graphql-go/graphql"
)

type SocialGraphBuilder struct {
    graphDB        *GraphDatabase
    similarityEngine *SimilarityEngine
    networkAnalyzer *NetworkAnalyzer
}

type SocialGraph struct {
    Nodes          map[string]*PersonNode     `json:"nodes"`
    Edges          map[string]*RelationshipEdge `json:"edges"`
    Communities    []*Community               `json:"communities"`
    CentralNodes   []*CentralityScore         `json:"central_nodes"`
    Metrics        *GraphMetrics              `json:"metrics"`
    GeneratedAt    time.Time                  `json:"generated_at"`
}

type PersonNode struct {
    ID            string                 `json:"id"`
    Type          NodeType               `json:"type"` // phone, email, social_profile, organization
    Value         string                 `json:"value"`
    Profiles      map[string]SocialProfile `json:"profiles"`
    Attributes    map[string]interface{} `json:"attributes"`
    Metadata      *NodeMetadata          `json:"metadata"`
    Centrality    *CentralityMetrics     `json:"centrality"`
    RiskScore     float64                `json:"risk_score"`
    DiscoveredAt  time.Time              `json:"discovered_at"`
    LastSeen      time.Time              `json:"last_seen"`
}

type RelationshipEdge struct {
    ID            string               `json:"id"`
    Source        string               `json:"source"`
    Target        string               `json:"target"`
    Type          RelationshipType     `json:"type"`
    Strength      float64              `json:"strength"` // 0.0 - 1.0
    Confidence    float64              `json:"confidence"`
    Evidence      []Evidence           `json:"evidence"`
    FirstSeen     time.Time            `json:"first_seen"`
    LastSeen      time.Time            `json:"last_seen"`
    Weight        float64              `json:"weight"`
    Direction     DirectionType        `json:"direction"` // undirected, directed
}

type Community struct {
    ID           string        `json:"id"`
    Name         string        `json:"name"`
    Nodes        []string      `json:"nodes"`
    Size         int           `json:"size"`
    Density      float64       `json:"density"`
    Modularity   float64       `json:"modularity"`
    CentralNode  string        `json:"central_node"`
    RiskLevel    string        `json:"risk_level"`
}

func NewSocialGraphBuilder(graphDB *GraphDatabase) *SocialGraphBuilder {
    return &SocialGraphBuilder{
        graphDB:        graphDB,
        similarityEngine: NewSimilarityEngine(),
        networkAnalyzer: NewNetworkAnalyzer(),
    }
}

// BuildComprehensiveSocialGraph creates advanced relationship mapping
func (sgb *SocialGraphBuilder) BuildComprehensiveSocialGraph(phoneNumber string, intelligenceData *IntelligenceData) (*SocialGraph, error) {
    graph := &SocialGraph{
        Nodes:       make(map[string]*PersonNode),
        Edges:       make(map[string]*RelationshipEdge),
        Communities: make([]*Community, 0),
        GeneratedAt: time.Now(),
    }

    // Step 1: Build core graph from intelligence data
    if err := sgb.buildCoreGraph(graph, phoneNumber, intelligenceData); err != nil {
        return nil, fmt.Errorf("failed to build core graph: %w", err)
    }

    // Step 2: Analyze connections and relationships
    sgb.analyzeConnections(graph)

    // Step 3: Detect communities and clusters
    sgb.detectCommunities(graph)

    // Step 4: Calculate network metrics and centrality
    sgb.calculateNetworkMetrics(graph)

    // Step 5: Identify influence and risk patterns
    sgb.identifyInfluencePatterns(graph)

    // Step 6: Score relationships and risks
    sgb.scoreRelationshipsAndRisks(graph)

    return graph, nil
}

func (sgb *SocialGraphBuilder) buildCoreGraph(graph *SocialGraph, phoneNumber string, intelligenceData *IntelligenceData) error {
    // Create root node for the phone number
    rootNode := &PersonNode{
        ID:           sgb.generateNodeID("phone", phoneNumber),
        Type:         NodeTypePhone,
        Value:        phoneNumber,
        Profiles:     make(map[string]SocialProfile),
        Attributes:   make(map[string]interface{}),
        Metadata:     &NodeMetadata{},
        DiscoveredAt: time.Now(),
        LastSeen:     time.Now(),
    }
    graph.Nodes[rootNode.ID] = rootNode

    // Add discovered emails as nodes and create relationships
    for _, email := range intelligenceData.DiscoveredEmails {
        emailNode := &PersonNode{
            ID:           sgb.generateNodeID("email", email.Email),
            Type:         NodeTypeEmail,
            Value:        email.Email,
            Profiles:     make(map[string]SocialProfile),
            Attributes:   map[string]interface{}{"confidence": email.Confidence, "source": email.Source},
            Metadata:     &NodeMetadata{},
            DiscoveredAt: time.Now(),
            LastSeen:     time.Now(),
        }
        graph.Nodes[emailNode.ID] = emailNode

        // Create relationship between phone and email
        edgeID := sgb.generateEdgeID(rootNode.ID, emailNode.ID)
        graph.Edges[edgeID] = &RelationshipEdge{
            ID:         edgeID,
            Source:     rootNode.ID,
            Target:     emailNode.ID,
            Type:       RelationshipTypeSamePerson,
            Strength:   email.Confidence,
            Confidence: email.Confidence,
            Evidence:   []Evidence{{Type: "email_discovery", Confidence: email.Confidence}},
            FirstSeen:  time.Now(),
            LastSeen:   time.Now(),
            Weight:     email.Confidence,
            Direction:  DirectionTypeUndirected,
        }
    }

    // Add social profiles and their connections
    for platform, profiles := range intelligenceData.SocialProfiles {
        for _, profile := range profiles {
            profileNode := &PersonNode{
                ID:    sgb.generateNodeID("social", profile.Username+"@"+platform),
                Type:  NodeTypeSocialProfile,
                Value: profile.Username,
                Profiles: map[string]SocialProfile{
                    platform: profile,
                },
                Attributes: map[string]interface{}{
                    "platform":    platform,
                    "verified":    profile.Verified,
                    "followers":   profile.FollowerCount,
                    "last_active": profile.LastActive,
                },
                Metadata:     &NodeMetadata{},
                DiscoveredAt: time.Now(),
                LastSeen:     time.Now(),
            }
            graph.Nodes[profileNode.ID] = profileNode

            // Connect profiles to likely email nodes
            sgb.connectProfilesToEmails(graph, profileNode, intelligenceData.DiscoveredEmails)

            // Connect profiles to phone number through behavioral patterns
            sgb.connectProfilesToPhone(graph, profileNode, rootNode, intelligenceData.BehavioralPatterns)
        }
    }

    // Add organizational connections if available
    sgb.addOrganizationalConnections(graph, intelligenceData.OrganizationalData)

    return nil
}

func (sgb *SocialGraphBuilder) analyzeConnections(graph *SocialGraph) {
    nodes := sgb.getNodeList(graph.Nodes)
    
    // Analyze all possible connections between nodes
    for i, node1 := range nodes {
        for j, node2 := range nodes {
            if i >= j {
                continue // Avoid duplicate comparisons
            }

            // Calculate multiple similarity measures
            similarityScore := sgb.calculateComprehensiveSimilarity(node1, node2)
            
            if similarityScore > 0.3 { // Threshold for connection
                edgeID := sgb.generateEdgeID(node1.ID, node2.ID)
                
                // Determine relationship type based on similarity patterns
                relationshipType := sgb.determineRelationshipType(node1, node2, similarityScore)
                
                // Calculate confidence based on evidence
                confidence := sgb.calculateRelationshipConfidence(node1, node2, similarityScore)
                
                graph.Edges[edgeID] = &RelationshipEdge{
                    ID:         edgeID,
                    Source:     node1.ID,
                    Target:     node2.ID,
                    Type:       relationshipType,
                    Strength:   similarityScore,
                    Confidence: confidence,
                    Evidence:   sgb.collectRelationshipEvidence(node1, node2),
                    FirstSeen:  time.Now(),
                    LastSeen:   time.Now(),
                    Weight:     similarityScore * confidence,
                    Direction:  DirectionTypeUndirected,
                }
            }
        }
    }
}

func (sgb *SocialGraphBuilder) calculateComprehensiveSimilarity(node1, node2 *PersonNode) float64 {
    var totalScore float64
    var weightSum float64

    // 1. Profile-based similarity
    if profileScore := sgb.calculateProfileSimilarity(node1, node2); profileScore > 0 {
        totalScore += profileScore * 0.4
        weightSum += 0.4
    }

    // 2. Behavioral similarity
    if behaviorScore := sgb.calculateBehavioralSimilarity(node1, node2); behaviorScore > 0 {
        totalScore += behaviorScore * 0.3
        weightSum += 0.3
    }

    // 3. Temporal similarity (activity patterns)
    if temporalScore := sgb.calculateTemporalSimilarity(node1, node2); temporalScore > 0 {
        totalScore += temporalScore * 0.2
        weightSum += 0.2
    }

    // 4. Content similarity
    if contentScore := sgb.calculateContentSimilarity(node1, node2); contentScore > 0 {
        totalScore += contentScore * 0.1
        weightSum += 0.1
    }

    if weightSum == 0 {
        return 0
    }

    return totalScore / weightSum
}

func (sgb *SocialGraphBuilder) detectCommunities(graph *SocialGraph) {
    // Use Louvain method for community detection
    communities := sgb.networkAnalyzer.DetectCommunitiesLouvain(graph)
    
    for i, community := range communities {
        communityID := fmt.Sprintf("community-%d", i+1)
        
        // Calculate community metrics
        density := sgb.calculateCommunityDensity(community, graph)
        modularity := sgb.calculateModularity(community, graph)
        centralNode := sgb.findCommunityCentralNode(community, graph)
        riskLevel := sgb.assessCommunityRisk(community, graph)

        graph.Communities = append(graph.Communities, &Community{
            ID:          communityID,
            Name:        fmt.Sprintf("Community %d", i+1),
            Nodes:       community,
            Size:        len(community),
            Density:     density,
            Modularity:  modularity,
            CentralNode: centralNode,
            RiskLevel:   riskLevel,
        })
    }

    // Sort communities by size
    sort.Slice(graph.Communities, func(i, j int) bool {
        return graph.Communities[i].Size > graph.Communities[j].Size
    })
}

func (sgb *SocialGraphBuilder) calculateNetworkMetrics(graph *SocialGraph) {
    metrics := &GraphMetrics{}
    
    // Basic metrics
    metrics.NodeCount = len(graph.Nodes)
    metrics.EdgeCount = len(graph.Edges)
    metrics.Density = sgb.calculateGraphDensity(graph)
    
    // Centrality metrics
    metrics.DegreeCentrality = sgb.calculateDegreeCentrality(graph)
    metrics.BetweennessCentrality = sgb.calculateBetweennessCentrality(graph)
    metrics.ClosenessCentrality = sgb.calculateClosenessCentrality(graph)
    
    // Clustering metrics
    metrics.AverageClustering = sgb.calculateAverageClustering(graph)
    metrics.ConnectedComponents = sgb.findConnectedComponents(graph)
    
    // Path metrics
    metrics.AveragePathLength = sgb.calculateAveragePathLength(graph)
    metrics.Diameter = sgb.calculateDiameter(graph)
    
    graph.Metrics = metrics
    
    // Update node centrality scores
    sgb.updateNodeCentralityScores(graph)
}

func (sgb *SocialGraphBuilder) identifyInfluencePatterns(graph *SocialGraph) {
    // Identify key influencers in the network
    influencers := sgb.findInfluencers(graph)
    
    // Detect bridge nodes (connect different communities)
    bridges := sgb.findBridgeNodes(graph)
    
    // Identify isolated clusters
    isolated := sgb.findIsolatedClusters(graph)
    
    // Detect anomalous connection patterns
    anomalies := sgb.detectAnomalousPatterns(graph)
    
    // Store influence patterns in graph metadata
    graph.Metrics.Influencers = influencers
    graph.Metrics.BridgeNodes = bridges
    graph.Metrics.IsolatedClusters = isolated
    graph.Metrics.Anomalies = anomalies
}

func (sgb *SocialGraphBuilder) scoreRelationshipsAndRisks(graph *SocialGraph) {
    // Score each relationship based on multiple factors
    for _, edge := range graph.Edges {
        edge.RiskScore = sgb.calculateRelationshipRisk(edge, graph)
    }
    
    // Score each node based on network position and relationships
    for _, node := range graph.Nodes {
        node.RiskScore = sgb.calculateNodeRisk(node, graph)
    }
    
    // Identify high-risk relationships and nodes
    sgb.identifyHighRiskElements(graph)
}
