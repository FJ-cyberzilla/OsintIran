package correlation_engine

type IdentityResolver struct {
    correlationEngine *PlatformCorrelator
    confidenceThreshold float64
}

type ResolvedIdentity struct {
    CanonicalID    string            `json:"canonical_id"`
    Entities       []*SocialEntity   `json:"entities"`
    PrimaryEntity  *SocialEntity     `json:"primary_entity"`
    Confidence     float64           `json:"confidence"`
    Evidence       []Evidence        `json:"evidence"`
    Platforms      []string          `json:"platforms"`
    FirstSeen      time.Time         `json:"first_seen"`
    LastSeen       time.Time         `json:"last_seen"`
}

func (ir *IdentityResolver) ResolveIdentities(entities []*SocialEntity) []ResolvedIdentity {
    // Create initial clusters
    clusters := ir.initialClustering(entities)
    
    // Merge clusters based on correlation
    mergedClusters := ir.mergeClusters(clusters)
    
    // Build resolved identities
    var resolvedIdentities []ResolvedIdentity
    for _, cluster := range mergedClusters {
        if ir.validateCluster(cluster) {
            resolvedIdentity := ir.buildResolvedIdentity(cluster)
            if resolvedIdentity.Confidence >= ir.confidenceThreshold {
                resolvedIdentities = append(resolvedIdentities, resolvedIdentity)
            }
        }
    }
    
    return resolvedIdentities
}

func (ir *IdentityResolver) buildResolvedIdentity(cluster []*SocialEntity) ResolvedIdentity {
    identity := ResolvedIdentity{
        CanonicalID: ir.generateCanonicalID(cluster),
        Entities: cluster,
        Platforms: ir.extractPlatforms(cluster),
    }
    
    // Determine primary entity (most complete profile)
    identity.PrimaryEntity = ir.determinePrimaryEntity(cluster)
    
    // Calculate confidence
    identity.Confidence = ir.calculateIdentityConfidence(cluster)
    
    // Extract temporal information
    identity.FirstSeen, identity.LastSeen = ir.extractTemporalRange(cluster)
    
    // Gather evidence
    identity.Evidence = ir.gatherIdentityEvidence(cluster)
    
    return identity
}
