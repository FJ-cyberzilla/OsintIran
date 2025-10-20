import React from 'react';
import { SocialGraph, InfluenceMetrics } from '../../../types/analytics';
import { SocialNetworkGraph } from './SocialNetworkGraph';
import { MetricsPanel } from '../common/MetricsPanel';

interface InfluenceNetworkProps {
    graph: SocialGraph;
    influenceMetrics: Map<string, InfluenceMetrics>;
}

export const InfluenceNetwork: React.FC<InfluenceNetworkProps> = ({
    graph,
    influenceMetrics
}) => {
    const [selectedNode, setSelectedNode] = React.useState<string | null>(null);
    
    const handleNodeClick = (entity: any) => {
        setSelectedNode(entity.id);
    };
    
    const selectedMetrics = selectedNode ? influenceMetrics.get(selectedNode) : null;
    
    return (
        <div className="influence-network">
            <div className="network-container">
                <SocialNetworkGraph
                    graph={graph}
                    onNodeClick={handleNodeClick}
                    width={1000}
                    height={700}
                />
            </div>
            
            {selectedMetrics && (
                <div className="metrics-sidebar">
                    <MetricsPanel metrics={selectedMetrics} />
                </div>
            )}
        </div>
    );
};
