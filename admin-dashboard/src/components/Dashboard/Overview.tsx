// admin-dashboard/src/components/Dashboard/Overview.tsx
import React from 'react';
import { useRealTimeData } from '../../../hooks/useRealTimeData';
import { useWebSocket } from '../../../hooks/useWebSocket';

export const DashboardOverview: React.FC = () => {
  const { systemHealth, proxyStats, jobMetrics, performance } = useRealTimeData();
  const { connected } = useWebSocket();

  return (
    <div className="dashboard-overview">
      {/* Status Header */}
      <div className="status-header">
        <div className="system-status">
          <span className={`status-indicator ${connected ? 'connected' : 'disconnected'}`}>
            {connected ? '🟢' : '🔴'} {connected ? 'Connected' : 'Disconnected'}
          </span>
          <span className="last-update">
            Last update: {new Date().toLocaleTimeString('fa-IR')}
          </span>
        </div>
      </div>

      {/* Key Metrics Grid */}
      <div className="metrics-grid">
        <MetricCard
          title="Active Jobs"
          value={jobMetrics.active}
          change={jobMetrics.change}
          icon="📊"
          color="blue"
        />
        <MetricCard
          title="Healthy Proxies"
          value={proxyStats.healthy}
          total={proxyStats.total}
          icon="🌐"
          color="green"
        />
        <MetricCard
          title="Success Rate"
          value={`${performance.successRate}%`}
          change={performance.successRateChange}
          icon="✅"
          color="purple"
        />
        <MetricCard
          title="Cost Today"
          value={`$${performance.costToday}`}
          change={performance.costChange}
          icon="💰"
          color="orange"
        />
      </div>

      {/* Quick Actions */}
      <div className="quick-actions">
        <h3>Quick Actions</h3>
        <div className="action-buttons">
          <button className="btn-primary" onClick={() => window.open('/admin/jobs/bulk', '_self')}>
            📦 Create Bulk Job
          </button>
          <button className="btn-secondary" onClick={() => window.open('/admin/proxies', '_self')}>
            🔧 Manage Proxies
          </button>
          <button className="btn-success" onClick={() => window.open('/admin/analytics', '_self')}>
            📈 View Analytics
          </button>
          <button className="btn-warning" onClick={() => window.open('/admin/performance', '_self')}>
            ⚡ Performance
          </button>
        </div>
      </div>

      {/* Real-time Activity Feed */}
      <div className="activity-feed">
        <h3>Real-time Activity</h3>
        <ActivityFeed />
      </div>
    </div>
  );
};

const MetricCard: React.FC<{
  title: string;
  value: string | number;
  total?: number;
  change?: number;
  icon: string;
  color: string;
}> = ({ title, value, total, change, icon, color }) => (
  <div className={`metric-card metric-${color}`}>
    <div className="metric-header">
      <span className="metric-icon">{icon}</span>
      <span className="metric-title">{title}</span>
    </div>
    <div className="metric-value">{value}</div>
    {total && (
      <div className="metric-total">of {total} total</div>
    )}
    {change !== undefined && (
      <div className={`metric-change ${change >= 0 ? 'positive' : 'negative'}`}>
        {change >= 0 ? '↗' : '↘'} {Math.abs(change)}%
      </div>
    )}
  </div>
);
