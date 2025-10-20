// admin-dashboard/src/components/Jobs/BulkJobManager.tsx
import React, { useState, useRef } from 'react';
import { useBulkOperations } from '../../../hooks/useBulkOperations';

export const BulkJobManager: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'upload' | 'template' | 'history'>('upload');
  const [selectedPlatforms, setSelectedPlatforms] = useState<string[]>([]);
  const [jobPriority, setJobPriority] = useState<'low' | 'normal' | 'high'>('normal');
  const fileInputRef = useRef<HTMLInputElement>(null);

  const {
    uploadFile,
    createBulkJob,
    isUploading,
    uploadProgress,
    jobStatus
  } = useBulkOperations();

  const platforms = [
    { id: 'facebook', name: 'Facebook', category: 'social' },
    { id: 'instagram', name: 'Instagram', category: 'social' },
    { id: 'telegram', name: 'Telegram', category: 'messaging' },
    { id: 'whatsapp', name: 'WhatsApp', category: 'messaging' },
    { id: 'linkedin', name: 'LinkedIn', category: 'professional' },
    { id: 'twitter', name: 'Twitter', category: 'social' },
    // ... 100+ platforms
  ];

  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    try {
      const result = await uploadFile(file);
      if (result.success) {
        await createBulkJob({
          phoneNumbers: result.data.phoneNumbers,
          platforms: selectedPlatforms,
          priority: jobPriority,
          name: `Bulk Job - ${new Date().toLocaleString('fa-IR')}`
        });
      }
    } catch (error) {
      console.error('Bulk job creation failed:', error);
    }
  };

  const downloadTemplate = () => {
    const template = "phone_number,notes\n+989123456789,Sample number\n+989123456790,Another number";
    const blob = new Blob([template], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'phone_intelligence_template.csv';
    a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="bulk-job-manager">
      <div className="manager-header">
        <h1>Bulk Job Management</h1>
        <p>Process thousands of phone numbers efficiently</p>
      </div>

      {/* Tab Navigation */}
      <div className="tab-navigation">
        <button 
          className={activeTab === 'upload' ? 'active' : ''}
          onClick={() => setActiveTab('upload')}
        >
          📤 Upload File
        </button>
        <button 
          className={activeTab === 'template' ? 'active' : ''}
          onClick={() => setActiveTab('template')}
        >
          📋 Download Template
        </button>
        <button 
          className={activeTab === 'history' ? 'active' : ''}
          onClick={() => setActiveTab('history')}
        >
          📊 Job History
        </button>
      </div>

      {/* Upload Section */}
      {activeTab === 'upload' && (
        <div className="upload-section">
          <div className="upload-area" onClick={() => fileInputRef.current?.click()}>
            <input
              type="file"
              ref={fileInputRef}
              onChange={handleFileUpload}
              accept=".csv,.xlsx,.txt"
              style={{ display: 'none' }}
            />
            <div className="upload-content">
              <div className="upload-icon">📁</div>
              <h3>Upload Phone Numbers File</h3>
              <p>Supported formats: CSV, Excel, Text (one number per line)</p>
              <p className="file-size">Max file size: 100MB</p>
            </div>
          </div>

          {isUploading && (
            <div className="upload-progress">
              <div className="progress-bar">
                <div 
                  className="progress-fill" 
                  style={{ width: `${uploadProgress}%` }}
                />
              </div>
              <span>{uploadProgress}%</span>
            </div>
          )}

          {/* Platform Selection */}
          <div className="platform-selection">
            <h4>Select Target Platforms</h4>
            <div className="platform-categories">
              {['social', 'messaging', 'professional'].map(category => (
                <div key={category} className="platform-category">
                  <h5>{category.charAt(0).toUpperCase() + category.slice(1)}</h5>
                  <div className="platform-checkboxes">
                    {platforms
                      .filter(p => p.category === category)
                      .map(platform => (
                        <label key={platform.id} className="checkbox-label">
                          <input
                            type="checkbox"
                            checked={selectedPlatforms.includes(platform.id)}
                            onChange={(e) => {
                              if (e.target.checked) {
                                setSelectedPlatforms([...selectedPlatforms, platform.id]);
                              } else {
                                setSelectedPlatforms(selectedPlatforms.filter(p => p !== platform.id));
                              }
                            }}
                          />
                          {platform.name}
                        </label>
                      ))
                    }
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Job Configuration */}
          <div className="job-configuration">
            <h4>Job Configuration</h4>
            <div className="config-options">
              <div className="config-group">
                <label>Priority:</label>
                <select 
                  value={jobPriority} 
                  onChange={(e) => setJobPriority(e.target.value as any)}
                >
                  <option value="low">Low</option>
                  <option value="normal">Normal</option>
                  <option value="high">High</option>
                </select>
              </div>
              
              <div className="config-group">
                <label>Concurrent Workers:</label>
                <input type="number" min="1" max="100" defaultValue="10" />
              </div>

              <div className="config-group">
                <label>Rate Limit (req/min):</label>
                <input type="number" min="1" max="1000" defaultValue="100" />
              </div>
            </div>
          </div>

          <button 
            className="btn-primary large"
            disabled={selectedPlatforms.length === 0 || isUploading}
            onClick={() => fileInputRef.current?.click()}
          >
            {isUploading ? 'Processing...' : 'Start Bulk Job'}
          </button>
        </div>
      )}

      {/* Template Section */}
      {activeTab === 'template' && (
        <div className="template-section">
          <div className="template-info">
            <h3>Download Template File</h3>
            <p>Use our template to ensure proper formatting for bulk uploads.</p>
            
            <div className="template-format">
              <h4>Expected Format:</h4>
              <pre>
                {`phone_number,notes,country_code
+989123456789,Sample Iranian number,IR
+14155552671,Sample US number,US
447123456789,Sample UK number,GB`}
              </pre>
            </div>

            <button className="btn-success" onClick={downloadTemplate}>
              📥 Download CSV Template
            </button>
          </div>
        </div>
      )}

      {/* History Section */}
      {activeTab === 'history' && (
        <div className="history-section">
          <BulkJobHistory />
        </div>
      )}
    </div>
  );
};

// Bulk Job History Component
const BulkJobHistory: React.FC = () => {
  const { bulkJobs, loading, error } = useBulkOperations();

  if (loading) return <div className="loading">Loading job history...</div>;
  if (error) return <div className="error">Error loading history: {error}</div>;

  return (
    <div className="bulk-job-history">
      <div className="history-header">
        <h3>Bulk Job History</h3>
        <div className="history-filters">
          <select>
            <option value="all">All Status</option>
            <option value="completed">Completed</option>
            <option value="processing">Processing</option>
            <option value="failed">Failed</option>
          </select>
          <input type="date" placeholder="From date" />
          <input type="date" placeholder="To date" />
        </div>
      </div>

      <div className="jobs-table">
        <table>
          <thead>
            <tr>
              <th>Job ID</th>
              <th>Name</th>
              <th>Numbers</th>
              <th>Platforms</th>
              <th>Status</th>
              <th>Created</th>
              <th>Progress</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {bulkJobs.map(job => (
              <tr key={job.id}>
                <td>{job.id.slice(0, 8)}...</td>
                <td>{job.name}</td>
                <td>{job.totalNumbers.toLocaleString()}</td>
                <td>{job.platforms.slice(0, 3).join(', ')} {job.platforms.length > 3 && `+${job.platforms.length - 3} more`}</td>
                <td>
                  <span className={`status-badge status-${job.status}`}>
                    {job.status}
                  </span>
                </td>
                <td>{new Date(job.createdAt).toLocaleDateString('fa-IR')}</td>
                <td>
                  <div className="progress-container">
                    <div 
                      className="progress-bar" 
                      style={{ width: `${job.progress}%` }}
                    />
                    <span>{job.progress}%</span>
                  </div>
                </td>
                <td>
                  <button className="btn-sm">View</button>
                  <button className="btn-sm">Export</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
