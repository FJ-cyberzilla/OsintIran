// ui/src/components/Analytics/RelationshipMapper.tsx
import React from 'react';
import { ForceGraph2D } from 'react-force-graph';

export const RelationshipMapper: React.FC<{ socialGraph: SocialGraph }> = ({ socialGraph }) => {
  const graphData = convertToGraphData(socialGraph);

  return (
    <div className="relationship-mapper">
      <div className="mapper-header">
        <h2>🕸️ Relationship Mapping</h2>
        <p>Visualize social connections and relationship networks</p>
      </div>

      <div className="graph-container">
        <ForceGraph2D
          graphData={graphData}
          nodeLabel="name"
          nodeColor={node => getNodeColor(node)}
          linkColor={link => getLinkColor(link)}
          nodeRelSize={6}
          linkWidth={link => link.strength * 3}
          onNodeClick={node => handleNodeClick(node)}
        />
      </div>

      <div className="graph-legend">
        <div className="legend-item">
          <div className="node-sample phone"></div>
          <span>Phone Number</span>
        </div>
        <div className="legend-item">
          <div className="node-sample email"></div>
          <span>Email Address</span>
        </div>
        <div className="legend-item">
          <div className="node-sample social"></div>
          <span>Social Profile</span>
        </div>
        <div className="legend-item">
          <div className="link-sample strong"></div>
          <span>Strong Connection</span>
        </div>
        <div className="legend-item">
          <div className="link-sample weak"></div>
          <span>Weak Connection</span>
        </div>
      </div>

      <div className="network-metrics">
        <h4>Network Analysis</h4>
        <div className="metrics-grid">
          <div className="metric">
            <span className="value">{socialGraph.nodes.length}</span>
            <span className="label">Total Nodes</span>
          </div>
          <div className="metric">
            <span className="value">{socialGraph.edges.length}</span>
            <span className="label">Connections</span>
          </div>
          <div className="metric">
            <span className="value">{socialGraph.density.toFixed(3)}</span>
            <span className="label">Network Density</span>
          </div>
          <div className="metric">
            <span className="value">{socialGraph.centralNode?.value}</span>
            <span className="label">Most Central</span>
          </div>
        </div>
      </div>
    </div>
  );
};
