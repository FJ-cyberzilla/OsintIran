// admin-dashboard/src/components/Performance/OptimizationDashboard.tsx
import React, { useState } from 'react';
import { usePerformance } from '../../../hooks/usePerformance';

export const OptimizationDashboard: React.FC = () => {
  const { metrics, optimizations, applyOptimization } = usePerformance();
  const [activeTab, setActiveTab] = useState<'overview' | 'bottlenecks' | 'recommendations'>('overview');

  return (
    <div className="optimization-dashboard">
      <div className="dashboard-header">
        <h1>Performance Optimization</h1>
        <p>Monitor and optimize system performance in real-time</p>
      </div>

      {/* Performance Score */}
      <div className="performance-score">
        <div className="score-card">
          <div className="score-value">87%</div>
          <div className="score-label">Overall Performance</div>
          <div className="score-trend positive">+5% this week</div>
        </div>
        
        <div className="score-breakdown">
          <div className="breakdown-item">
            <span className="label">Proxy Health:</span>
            <span className="value positive">94%</span>
          </div>
          <div className="breakdown-item">
            <span className="label">Success Rate:</span>
            <span className="value positive">89%</span>
          </div>
          <div className="breakdown-item">
            <span className="label">Response Time:</span>
            <span className="value warning">1.2s</span>
          </div>
          <div className="breakdown-item">
            <span className="label">Cost Efficiency:</span>
            <span className="value positive">82%</span>
          </div>
        </div>
      </div>

      {/* Optimization Tabs */}
      <div className="optimization-tabs">
        <button 
          className={activeTab === 'overview' ? 'active' : ''}
          onClick={() => setActiveTab('overview')}
        >
          📊 Overview
        </button>
        <button 
          className={activeTab === 'bottlenecks' ? 'active' : ''}
          onClick={() => setActiveTab('bottlenecks')}
        >
          🚧 Bottlenecks
        </button>
        <button 
          className={activeTab === 'recommendations' ? 'active' : ''}
          onClick={() => setActiveTab('recommendations')}
        >
          💡 Recommendations
        </button>
      </div>

      {/* Overview Tab */}
      {activeTab === 'overview' && (
        <div className="overview-tab">
          <div className="metrics-grid">
            <div className="metric-widget">
              <h4>Request Throughput</h4>
              <div className="metric-chart">
                {/* Chart would go here */}
                <div className="chart-placeholder">Throughput Chart</div>
              </div>
              <div className="metric-stats">
                <span>Current: 1,250 req/min</span>
                <span className="positive">+15%</span>
              </div>
            </div>

            <div className="metric-widget">
              <h4>Response Times</h4>
              <div className="metric-chart">
                <div className="chart-placeholder">Response Time Chart</div>
              </div>
              <div className="metric-stats">
                <span>Average: 1.2s</span>
                <span className="warning">+0.2s</span>
              </div>
            </div>

            <div className="metric-widget">
              <h4>Error Rates</h4>
              <div className="metric-chart">
                <div className="chart-placeholder">Error Rate Chart</div>
              </div>
              <div className="metric-stats">
                <span>Current: 2.1%</span>
                <span className="positive">-0.5%</span>
              </div>
            </div>

            <div className="metric-widget">
              <h4>Cost per Request</h4>
              <div className="metric-chart">
                <div className="chart-placeholder">Cost Chart</div>
              </div>
              <div className="metric-stats">
                <span>Average: $0.0021</span>
                <span className="positive">-12%</span>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Bottlenecks Tab */}
      {activeTab === 'bottlenecks' && (
        <div className="bottlenecks-tab">
          <h3>System Bottlenecks</h3>
          
          <div className="bottleneck-list">
            <div className="bottleneck-item critical">
              <div className="bottleneck-header">
                <span className="severity critical">CRITICAL</span>
                <span className="title">Proxy Response Time Degradation</span>
              </div>
              <div className="bottleneck-details">
                <p>15% of proxies are responding slower than 3 seconds</p>
                <div className="bottleneck-metrics">
                  <span>Affected: 3,450 proxies</span>
                  <span>Impact: 12% success rate drop</span>
                </div>
              </div>
              <button className="btn-warning">Optimize Now</button>
            </div>

            <div className="bottleneck-item warning">
              <div className="bottleneck-header">
                <span className="severity warning">WARNING</span>
                <span className="title">Database Connection Pool Exhaustion</span>
              </div>
              <div className="bottleneck-details">
                <p>Database connections reaching 85% capacity during peak hours</p>
                <div className="bottleneck-metrics">
                  <span>Peak Usage: 185/200 connections</span>
                  <span>Impact: Potential timeouts</span>
                </div>
              </div>
              <button className="btn-warning">Increase Pool Size</button>
            </div>

            <div className="bottleneck-item info">
              <div className="bottleneck-header">
                <span className="severity info">INFO</span>
                <span className="title">Memory Usage Optimization</span>
              </div>
              <div className="bottleneck-details">
                <p>Browser workers using 2.1GB RAM on average</p>
                <div className="bottleneck-metrics">
                  <span>Current: 2.1GB/worker</span>
                  <span>Target: 1.5GB/worker</span>
                </div>
              </div>
              <button className="btn-info">Optimize Memory</button>
            </div>
          </div>
        </div>
      )}

      {/* Recommendations Tab */}
      {activeTab === 'recommendations' && (
        <div className="recommendations-tab">
          <h3>Optimization Recommendations</h3>
          
          <div className="recommendation-list">
            <div className="recommendation-item high-impact">
              <div className="recommendation-header">
                <span className="impact high">HIGH IMPACT</span>
                <span className="savings">Save: $1,200/month</span>
              </div>
              <h4>Implement Proxy Tiering</h4>
              <p>Use cheaper datacenter proxies for simple checks and reserve residential proxies for complex platforms</p>
              <div className="recommendation-actions">
                <button className="btn-success" onClick={() => applyOptimization('proxy_tiering')}>
                  Apply Optimization
                </button>
                <span className="effort">Effort: Medium</span>
              </div>
            </div>

            <div className="recommendation-item medium-impact">
              <div className="recommendation-header">
                <span className="impact medium">MEDIUM IMPACT</span>
                <span className="savings">Save: 350 req/min</span>
              </div>
              <h4>Optimize Request Batching</h4>
              <p>Batch similar requests to reduce overhead and improve throughput</p>
              <div className="recommendation-actions">
                <button className="btn-success" onClick={() => applyOptimization('request_batching')}>
                  Apply Optimization
                </button>
                <span className="effort">Effort: Low</span>
              </div>
            </div>

            <div className="recommendation-item low-impact">
              <div className="recommendation-header">
                <span className="impact low">LOW IMPACT</span>
                <span className="savings">Save: 0.1s avg response</span>
              </div>
              <h4>Enable Connection Keep-Alive</h4>
              <p>Reuse HTTP connections to reduce TLS handshake overhead</p>
              <div className="recommendation-actions">
                <button className="btn-success" onClick={() => applyOptimization('keep_alive')}>
                  Apply Optimization
                </button>
                <span className="effort">Effort: Low</span>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Real-time Optimization Log */}
      <div className="optimization-log">
        <h4>Optimization Activity Log</h4>
        <div className="log-entries">
          <div className="log-entry">
            <span className="timestamp">14:32:05</span>
            <span className="action">Applied: Proxy tiering optimization</span>
            <span className="result positive">Success rate improved by 8%</span>
          </div>
          <div className="log-entry">
            <span className="timestamp">14:15:22</span>
            <span className="action">Rotated: 1,200 slow proxies</span>
            <span className="result positive">Avg response time decreased by 0.3s</span>
          </div>
          <div className="log-entry">
            <span className="timestamp">13:58:41</span>
            <span className="action">Optimized: Database query indexes</span>
            <span className="result positive">Query time reduced by 45%</span>
          </div>
        </div>
      </div>
    </div>
  );
};
