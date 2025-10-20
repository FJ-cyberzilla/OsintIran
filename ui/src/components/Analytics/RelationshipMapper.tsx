// ui/src/components/Analytics/RelationshipMapper.tsx
import React, { useRef, useEffect } from 'react';
import * as d3 from 'd3';
import { SocialGraph, SocialEntity, SocialConnection } from '../../../types/analytics';

interface SocialNetworkGraphProps {
    graph: SocialGraph;
    width?: number;
    height?: number;
    onNodeClick?: (entity: SocialEntity) => void;
    onConnectionClick?: (connection: SocialConnection) => void;
}

export const SocialNetworkGraph: React.FC<SocialNetworkGraphProps> = ({
    graph,
    width = 800,
    height = 600,
    onNodeClick,
    onConnectionClick
}) => {
    const svgRef = useRef<SVGSVGElement>(null);
    
    useEffect(() => {
        if (!svgRef.current || !graph) return;
        
        const svg = d3.select(svgRef.current);
        svg.selectAll("*").remove(); // Clear previous
        
        // Create force-directed graph
        const simulation = d3.forceSimulation(graph.nodes)
            .force("link", d3.forceLink(graph.connections).id((d: any) => d.id))
            .force("charge", d3.forceManyBody().strength(-1000))
            .force("center", d3.forceCenter(width / 2, height / 2));
            
        // Draw connections
        const link = svg.append("g")
            .selectAll("line")
            .data(graph.connections)
            .enter().append("line")
            .attr("stroke-width", (d: any) => Math.sqrt(d.strength) * 3)
            .attr("stroke", "#999")
            .on("click", (event, d) => onConnectionClick?.(d));
            
        // Draw nodes
        const node = svg.append("g")
            .selectAll("circle")
            .data(graph.nodes)
            .enter().append("circle")
            .attr("r", (d: any) => Math.sqrt(d.influence) * 5 + 5)
            .attr("fill", (d: any) => getNodeColor(d.type))
            .on("click", (event, d) => onNodeClick?.(d));
            
        // Add labels
        const label = svg.append("g")
            .selectAll("text")
            .data(graph.nodes)
            .enter().append("text")
            .text((d: any) => d.username)
            .attr("font-size", "10px")
            .attr("dx", 12)
            .attr("dy", 4);
            
        // Update positions
        simulation.on("tick", () => {
            link
                .attr("x1", (d: any) => d.source.x)
                .attr("y1", (d: any) => d.source.y)
                .attr("x2", (d: any) => d.target.x)
                .attr("y2", (d: any) => d.target.y);
                
            node
                .attr("cx", (d: any) => d.x)
                .attr("cy", (d: any) => d.y);
                
            label
                .attr("x", (d: any) => d.x)
                .attr("y", (d: any) => d.y);
        });
        
    }, [graph, width, height]);
    
    return (
        <div className="social-network-graph">
            <svg ref={svgRef} width={width} height={height} />
        </div>
    );
};

const getNodeColor = (type: string): string => {
    const colors = {
        person: '#4CAF50',
        organization: '#2196F3',
        phone: '#FF9800',
        email: '#9C27B0'
    };
    return colors[type] || '#607D8B';
};
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
