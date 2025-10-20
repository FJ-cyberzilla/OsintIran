import React from 'react';
import { ResolvedIdentity, CorrelationResult } from '../../../types/analytics';
import { PlatformIcon } from '../../common/PlatformIcon';

interface CrossPlatformViewProps {
    identities: ResolvedIdentity[];
    correlations: CorrelationResult[];
}

export const CrossPlatformView: React.FC<CrossPlatformViewProps> = ({
    identities,
    correlations
}) => {
    return (
        <div className="cross-platform-view">
            <div className="identities-section">
                <h3>Resolved Identities ({identities.length})</h3>
                <div className="identity-cards">
                    {identities.map(identity => (
                        <IdentityCard key={identity.canonicalID} identity={identity} />
                    ))}
                </div>
            </div>
            
            <div className="correlations-section">
                <h3>Entity Correlations ({correlations.length})</h3>
                <div className="correlation-matrix">
                    {correlations.map(correlation => (
                        <CorrelationRow key={`${correlation.entity1.id}-${correlation.entity2.id}`} 
                                      correlation={correlation} />
                    ))}
                </div>
            </div>
        </div>
    );
};

const IdentityCard: React.FC<{ identity: ResolvedIdentity }> = ({ identity }) => (
    <div className="identity-card">
        <div className="identity-header">
            <h4>{identity.primaryEntity?.displayName || 'Unknown Identity'}</h4>
            <span className="confidence-badge">
                {Math.round(identity.confidence * 100)}%
            </span>
        </div>
        
        <div className="platforms">
            {identity.platforms.map(platform => (
                <PlatformIcon key={platform} platform={platform} size="small" />
            ))}
        </div>
        
        <div className="entity-count">
            {identity.entities.length} linked entities
        </div>
        
        <div className="temporal-info">
            First seen: {identity.firstSeen.toLocaleDateString()}
        </div>
    </div>
);
