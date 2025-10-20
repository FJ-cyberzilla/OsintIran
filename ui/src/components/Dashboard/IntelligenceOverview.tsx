// ui/src/components/Dashboard/IntelligenceOverview.tsx
import React from 'react';
import { useIntelSearch } from '../../hooks/useIntelSearch';
import { useRealTimeData } from '../../hooks/useRealTimeData';

export const IntelligenceOverview: React.FC = () => {
  const { searchResults, isSearching } = useIntelSearch();
  const { systemHealth, agentStatus } = useRealTimeData();

  return (
    <div className="intelligence-dashboard">
      {/* System Status Header */}
      <div className="system-status">
        <div className="status-item">
          <span className="status-label">سیستم:</span>
          <span className={`status-value ${systemHealth.overall > 80 ? 'healthy' : 'warning'}`}>
            {systemHealth.overall}% سلامت
          </span>
        </div>
        <div className="status-item">
          <span className="status-label">عامل‌های فعال:</span>
          <span className="status-value">{agentStatus.active}/{agentStatus.total}</span>
        </div>
        <div className="status-item">
          <span className="status-label">پروکسی‌های سالم:</span>
          <span className="status-value">{systemHealth.healthyProxies}</span>
        </div>
      </div>

      {/* Search Interface */}
      <div className="search-interface">
        <h2>جستجوی هوشمند اطلاعات تلفن</h2>
        <PhoneSearchForm />
      </div>

      {/* Results Display */}
      {searchResults && (
        <div className="results-container">
          <IntelligenceTabs results={searchResults} />
          <PlatformResultsGrid platforms={searchResults.platforms} />
          <CrossPlatformAnalysis data={searchResults} />
          <RiskAssessmentPanel risks={searchResults.risks} />
        </div>
      )}

      {/* Real-time Agent Monitoring */}
      <div className="agent-monitor">
        <h3>مانیتورینگ عامل‌های هوشمند</h3>
        <AgentHierarchyView agents={agentStatus.details} />
        <PerformanceMetrics metrics={agentStatus.metrics} />
      </div>
    </div>
  );
};

// ui/src/components/AgentManager/MicroAgentGenerator.tsx
export const MicroAgentGenerator: React.FC = () => {
  const [agentType, setAgentType] = React.useState('');
  const [config, setConfig] = React.useState({});
  const { generateAgent, isGenerating } = useAgentManagement();

  const agentTemplates = [
    { value: 'social_scraper', label: 'عامل جمع‌آوری شبکه‌های اجتماعی' },
    { value: 'phone_analyzer', label: 'عامل تحلیل شماره تلفن' },
    { value: 'email_discoverer', label: 'عامل کشف ایمیل' },
    { value: 'behavior_analyzer', label: 'عامل تحلیل رفتار' },
    { value: 'risk_assessor', label: 'عامل ارزیابی ریسک' },
  ];

  const handleGenerate = async () => {
    const agent = await generateAgent(agentType, config);
    if (agent) {
      // Add to agent hierarchy
    }
  };

  return (
    <div className="agent-generator">
      <h3>تولید کننده عامل‌های میکرو</h3>
      
      <div className="generator-form">
        <select 
          value={agentType} 
          onChange={(e) => setAgentType(e.target.value)}
          className="agent-type-select"
        >
          <option value="">انتخاب نوع عامل</option>
          {agentTemplates.map(template => (
            <option key={template.value} value={template.value}>
              {template.label}
            </option>
          ))}
        </select>

        <div className="agent-config">
          <h4>پیکربندی عامل</h4>
          {/* Dynamic configuration based on agent type */}
          <AgentConfigurator 
            agentType={agentType} 
            config={config}
            onChange={setConfig}
          />
        </div>

        <button 
          onClick={handleGenerate}
          disabled={!agentType || isGenerating}
          className="generate-btn"
        >
          {isGenerating ? 'در حال تولید...' : 'تولید عامل جدید'}
        </button>
      </div>
    </div>
  );
};
