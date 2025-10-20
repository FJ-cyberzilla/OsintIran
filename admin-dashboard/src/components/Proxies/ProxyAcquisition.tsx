// admin-dashboard/src/components/Proxies/ProxyAcquisition.tsx
import React, { useState } from 'react';

export const ProxyAcquisition: React.FC = () => {
  const [acquisitionMethod, setAcquisitionMethod] = useState<'manual' | 'api' | 'file'>('manual');
  const [selectedProvider, setSelectedProvider] = useState<string>('');
  const [proxyCount, setProxyCount] = useState<number>(1000);

  const providers = [
    { id: 'luminati', name: 'Luminati', type: 'residential', costPerGB: 15 },
    { id: 'oxylabs', name: 'Oxylabs', type: 'residential', costPerGB: 12 },
    { id: 'smartproxy', name: 'Smartproxy', type: 'residential', costPerGB: 10 },
    { id: 'geosurf', name: 'GeoSurf', type: 'mobile', costPerGB: 25 },
    { id: 'stormproxies', name: 'Storm Proxies', type: 'datacenter', costPerGB: 2 },
  ];

  const calculateCost = (providerId: string, count: number) => {
    const provider = providers.find(p => p.id === providerId);
    if (!provider) return 0;
    
    // Estimated cost calculation (simplified)
    const estimatedGB = count * 0.01; // 10MB per proxy
    return estimatedGB * provider.costPerGB;
  };

  const handleAcquireProxies = async () => {
    // Implementation for proxy acquisition
    console.log('Acquiring proxies:', {
      method: acquisitionMethod,
      provider: selectedProvider,
      count: proxyCount,
      cost: calculateCost(selectedProvider, proxyCount)
    });
  };

  return (
    <div className="proxy-acquisition">
      <div className="acquisition-header">
        <h1>Proxy Acquisition</h1>
        <p>Purchase and integrate new proxies into your rotation pool</p>
      </div>

      {/* Acquisition Methods */}
      <div className="acquisition-methods">
        <div className="method-tabs">
          <button 
            className={acquisitionMethod === 'manual' ? 'active' : ''}
            onClick={() => setAcquisitionMethod('manual')}
          >
            🔧 Manual Entry
          </button>
          <button 
            className={acquisitionMethod === 'api' ? 'active' : ''}
            onClick={() => setAcquisitionMethod('api')}
          >
            🔌 API Integration
          </button>
          <button 
            className={acquisitionMethod === 'file' ? 'active' : ''}
            onClick={() => setAcquisitionMethod('file')}
          >
            📁 File Upload
          </button>
        </div>

        {/* Manual Entry */}
        {acquisitionMethod === 'manual' && (
          <div className="manual-entry">
            <h3>Manual Proxy Entry</h3>
            <div className="entry-form">
              <textarea 
                placeholder="Enter proxies in format: ip:port:username:password&#10;192.168.1.1:8080:user:pass&#10;192.168.1.2:8080:user:pass"
                rows={10}
              />
              <button className="btn-primary">Add to Pool</button>
            </div>
          </div>
        )}

        {/* API Integration */}
        {acquisitionMethod === 'api' && (
          <div className="api-integration">
            <h3>Provider API Integration</h3>
            
            <div className="provider-selection">
              <label>Select Provider:</label>
              <select 
                value={selectedProvider} 
                onChange={(e) => setSelectedProvider(e.target.value)}
              >
                <option value="">Choose a provider...</option>
                {providers.map(provider => (
                  <option key={provider.id} value={provider.id}>
                    {provider.name} ({provider.type}) - ${provider.costPerGB}/GB
                  </option>
                ))}
              </select>
            </div>

            {selectedProvider && (
              <div className="provider-config">
                <div className="config-group">
                  <label>Number of Proxies:</label>
                  <input 
                    type="number" 
                    value={proxyCount}
                    onChange={(e) => setProxyCount(parseInt(e.target.value))}
                    min="1" 
                    max="10000" 
                  />
                </div>

                <div className="config-group">
                  <label>Target Countries:</label>
                  <select multiple>
                    <option value="IR">Iran</option>
                    <option value="US">United States</option>
                    <option value="GB">United Kingdom</option>
                    <option value="DE">Germany</option>
                    <option value="FR">France</option>
                  </select>
                </div>

                <div className="cost-calculation">
                  <h4>Cost Estimation</h4>
                  <div className="cost-breakdown">
                    <div>Proxies: {proxyCount.toLocaleString()}</div>
                    <div>Estimated Usage: {(proxyCount * 0.01).toFixed(2)} GB</div>
                    <div className="total-cost">
                      Total: ${calculateCost(selectedProvider, proxyCount).toFixed(2)}
                    </div>
                  </div>
                </div>

                <button className="btn-success" onClick={handleAcquireProxies}>
                  🛒 Purchase Proxies
                </button>
              </div>
            )}
          </div>
        )}

        {/* File Upload */}
        {acquisitionMethod === 'file' && (
          <div className="file-upload">
            <h3>Upload Proxy List</h3>
            <div className="upload-area">
              <input type="file" accept=".csv,.txt,.json" />
              <p>Upload a file containing proxy list in CSV, TXT, or JSON format</p>
            </div>
            
            <div className="format-examples">
              <h4>Supported Formats:</h4>
              <pre>
                {`CSV: ip,port,username,password,type,country
192.168.1.1,8080,user,pass,http,US

TXT: ip:port:username:password
192.168.1.1:8080:user:pass

JSON: [{"ip": "192.168.1.1", "port": 8080, ...}]`}
              </pre>
            </div>
          </div>
        )}
      </div>

      {/* Recent Acquisitions */}
      <div className="recent-acquisitions">
        <h3>Recent Proxy Acquisitions</h3>
        <AcquisitionHistory />
      </div>
    </div>
  );
};

const AcquisitionHistory: React.FC = () => {
  const acquisitions = [
    { id: 1, provider: 'Luminati', count: 5000, cost: 750, date: '2024-01-15', status: 'active' },
    { id: 2, provider: 'GeoSurf', count: 2000, cost: 500, date: '2024-01-10', status: 'active' },
    { id: 3, provider: 'Manual', count: 100, cost: 0, date: '2024-01-08', status: 'inactive' },
  ];

  return (
    <div className="acquisition-history">
      <table>
        <thead>
          <tr>
            <th>Date</th>
            <th>Provider</th>
            <th>Proxies</th>
            <th>Cost</th>
            <th>Status</th>
            <th>Health</th>
          </tr>
        </thead>
        <tbody>
          {acquisitions.map(acq => (
            <tr key={acq.id}>
              <td>{acq.date}</td>
              <td>{acq.provider}</td>
              <td>{acq.count.toLocaleString()}</td>
              <td>${acq.cost}</td>
              <td>
                <span className={`status-badge status-${acq.status}`}>
                  {acq.status}
                </span>
              </td>
              <td>
                <div className="health-indicator">
                  <div className="health-bar">
                    <div className="health-fill" style={{ width: '85%' }} />
                  </div>
                  <span>85%</span>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
