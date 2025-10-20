// admin-dashboard/src/components/Performance/CostAnalytics.tsx
import React, { useState } from 'react';

export const CostAnalytics: React.FC = () => {
  const [timeRange, setTimeRange] = useState<'7d' | '30d' | '90d'>('30d');
  const [costBreakdown, setCostBreakdown] = useState<'category' | 'platform' | 'proxy'>('category');

  // Mock data - would come from API
  const costData = {
    totalCost: 12580.45,
    monthlyTrend: -12.5, // percentage
    byCategory: [
      { category: 'Proxy Services', cost: 8450.20, percentage: 67.2 },
      { category: 'Cloud Infrastructure', cost: 2850.75, percentage: 22.7 },
      { category: 'API Services', cost: 980.50, percentage: 7.8 },
      { category: 'Data Storage', cost: 299.00, percentage: 2.3 },
    ],
    byPlatform: [
      { platform: 'Facebook', cost: 3250.40, requests: 1250000, costPerRequest: 0.0026 },
      { platform: 'Instagram', cost: 2850.20, requests: 980000, costPerRequest: 0.0029 },
      { platform: 'LinkedIn', cost: 1980.75, requests: 450000, costPerRequest: 0.0044 },
      { platform: 'Twitter', cost: 1250.30, requests: 620000, costPerRequest: 0.0020 },
    ],
    dailyCosts: Array.from({ length: 30 }, (_, i) => ({
      date: new Date(Date.now() - (29 - i) * 24 * 60 * 60 * 1000),
      cost: 350 + Math.random() * 150,
      requests: 80000 + Math.random() * 40000
    }))
  };

  const calculateSavings = () => {
    const potentialSavings = {
      proxyOptimization: 2150,
      requestBatching: 890,
      cacheImplementation: 540,
      platformOptimization: 1230
    };
    
    return {
      totalPotential: Object.values(potentialSavings).reduce((a, b) => a + b, 0),
      breakdown: potentialSavings
    };
  };

  const savings = calculateSavings();

  return (
    <div className="cost-analytics">
      <div className="analytics-header">
        <h1>Cost Monitoring & Analytics</h1>
        <div className="header-controls">
          <select value={timeRange} onChange={(e) => setTimeRange(e.target.value as any)}>
            <option value="7d">Last 7 Days</option>
            <option value="30d">Last 30 Days</option>
            <option value="90d">Last 90 Days</option>
          </select>
          <select value={costBreakdown} onChange={(e) => setCostBreakdown(e.target.value as any)}>
            <option value="category">By Category</option>
            <option value="platform">By Platform</option>
            <option value="proxy">By Proxy Type</option>
          </select>
        </div>
      </div>

      {/* Cost Overview */}
      <div className="cost-overview">
        <div className="total-cost">
          <div className="amount">${costData.totalCost.toLocaleString()}</div>
          <div className="label">Total Cost ({timeRange})</div>
          <div className={`trend ${costData.monthlyTrend >= 0 ? 'positive' : 'negative'}`}>
            {costData.monthlyTrend >= 0 ? '↗' : '↘'} {Math.abs(costData.monthlyTrend)}%
          </div>
        </div>

        <div className="cost-metrics">
          <div className="metric">
            <span className="value">${(costData.totalCost / 30).toFixed(2)}</span>
            <span className="label">Daily Average</span>
          </div>
          <div className="metric">
            <span className="value">${(costData.totalCost / 1000000).toFixed(4)}</span>
            <span className="label">Cost per Request</span>
          </div>
          <div className="metric">
            <span className="value">{savings.totalPotential.toLocaleString()}</span>
            <span className="label">Potential Monthly Savings</span>
          </div>
        </div>
      </div>

      {/* Cost Breakdown */}
      <div className="cost-breakdown">
        <h3>Cost Breakdown</h3>
        
        {costBreakdown === 'category' && (
          <div className="breakdown-chart">
            {costData.byCategory.map(item => (
              <div key={item.category} className="breakdown-item">
                <div className="item-header">
                  <span className="category">{item.category}</span>
                  <span className="cost">${item.cost.toLocaleString()}</span>
                </div>
                <div className="progress-bar">
                  <div 
                    className="progress-fill" 
                    style={{ width: `${item.percentage}%` }}
                  />
                </div>
                <span className="percentage">{item.percentage}%</span>
              </div>
            ))}
          </div>
        )}

        {costBreakdown === 'platform' && (
          <div className="platform-costs">
            <table>
              <thead>
                <tr>
                  <th>Platform</th>
                  <th>Total Cost</th>
                  <th>Requests</th>
                  <th>Cost/Request</th>
                  <th>Efficiency</th>
                </tr>
              </thead>
              <tbody>
                {costData.byPlatform.map(platform => (
                  <tr key={platform.platform}>
                    <td>{platform.platform}</td>
                    <td>${platform.cost.toLocaleString()}</td>
                    <td>{platform.requests.toLocaleString()}</td>
                    <td>${platform.costPerRequest.toFixed(4)}</td>
                    <td>
                      <span className={`efficiency ${
                        platform.costPerRequest < 0.003 ? 'good' : 
                        platform.costPerRequest < 0.005 ? 'average' : 'poor'
                      }`}>
                        {platform.costPerRequest < 0.003 ? 'Good' : 
                         platform.costPerRequest < 0.005 ? 'Average' : 'Poor'}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Savings Opportunities */}
      <div className="savings-opportunities">
        <h3>Savings Opportunities</h3>
        <div className="opportunities-grid">
          <div className="opportunity-card">
            <h4>Proxy Optimization</h4>
            <div className="savings-amount">${savings.breakdown.proxyOptimization}/month</div>
            <p>Use cheaper proxies for non-critical platforms</p>
            <button className="btn-success">Implement</button>
          </div>

          <div className="opportunity-card">
            <h4>Request Batching</h4>
            <div className="savings-amount">${savings.breakdown.requestBatching}/month</div>
            <p>Batch similar requests to reduce overhead</p>
            <button className="btn-success">Implement</button>
          </div>

          <div className="opportunity-card">
            <h4>Cache Implementation</h4>
            <div className="savings-amount">${savings.breakdown.cacheImplementation}/month</div>
            <p>Cache frequent requests to avoid reprocessing</p>
            <button className="btn-success">Implement</button>
          </div>

          <div className="opportunity-card">
            <h4>Platform Optimization</h4>
            <div className="savings-amount">${savings.breakdown.platformOptimization}/month</div>
            <p>Optimize platform-specific request patterns</p>
            <button className="btn-success">Implement</button>
          </div>
        </div>
      </div>

      {/* Cost Projections */}
      <div className="cost-projections">
        <h3>Cost Projections</h3>
        <div className="projection-cards">
          <div className="projection-card">
            <h5>Current Monthly</h5>
            <div className="amount">${(costData.totalCost * 12 / 365 * 30).toFixed(2)}</div>
          </div>
          <div className="projection-card optimized">
            <h5>With Optimizations</h5>
            <div className="amount">${(costData.totalCost * 12 / 365 * 30 - savings.totalPotential).toFixed(2)}</div>
            <div className="savings">Save ${savings.totalPotential}</div>
          </div>
          <div className="projection-card">
            <h5>Projected Growth (6mo)</h5>
            <div className="amount">${(costData.totalCost * 12 / 365 * 30 * 1.5).toFixed(2)}</div>
          </div>
        </div>
      </div>
    </div>
  );
};
