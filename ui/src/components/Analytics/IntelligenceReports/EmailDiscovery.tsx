// ui/src/components/Analytics/IntelligenceReports/EmailDiscovery.tsx
import React, { useState } from 'react';
import { useIntelligenceAPI } from '../../../../hooks/useIntelligenceAPI';

export const EmailDiscovery: React.FC = () => {
  const [phoneNumber, setPhoneNumber] = useState('');
  const [discoveryResult, setDiscoveryResult] = useState<EmailDiscoveryResult | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  
  const { discoverEmails } = useIntelligenceAPI();

  const handleDiscoverEmails = async () => {
    if (!phoneNumber) return;
    
    setIsLoading(true);
    try {
      const result = await discoverEmails(phoneNumber);
      setDiscoveryResult(result);
    } catch (error) {
      console.error('Email discovery failed:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="email-discovery">
      <div className="discovery-header">
        <h2>📧 Email Discovery Engine</h2>
        <p>Find email addresses associated with phone numbers using advanced intelligence</p>
      </div>

      <div className="discovery-input">
        <input
          type="tel"
          value={phoneNumber}
          onChange={(e) => setPhoneNumber(e.target.value)}
          placeholder="Enter phone number (e.g., +989123456789)"
          className="phone-input"
        />
        <button 
          onClick={handleDiscoverEmails}
          disabled={!phoneNumber || isLoading}
          className="btn-primary"
        >
          {isLoading ? 'Discovering...' : 'Discover Emails'}
        </button>
      </div>

      {discoveryResult && (
        <div className="discovery-results">
          <div className="results-summary">
            <h3>Discovery Results</h3>
            <div className="summary-stats">
              <span className="stat">
                <strong>{discoveryResult.totalFound}</strong> emails found
              </span>
              <span className="stat">
                <strong>{discoveryResult.confidence}%</strong> confidence
              </span>
              <span className="stat">
                <strong>{discoveryResult.breachMatches}</strong> in breaches
              </span>
            </div>
          </div>

          <div className="emails-list">
            {discoveryResult.emails.map((email, index) => (
              <EmailCard key={index} email={email} />
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

const EmailCard: React.FC<{ email: DiscoveredEmail }> = ({ email }) => (
  <div className={`email-card confidence-${Math.floor(email.confidence * 10)}`}>
    <div className="email-header">
      <span className="email-address">{email.email}</span>
      <span className={`confidence-badge ${email.confidence > 0.8 ? 'high' : email.confidence > 0.5 ? 'medium' : 'low'}`}>
        {Math.round(email.confidence * 100)}%
      </span>
    </div>
    
    <div className="email-details">
      <div className="detail">
        <span className="label">Source:</span>
        <span className="value">{email.source}</span>
      </div>
      
      <div className="detail">
        <span className="label">Pattern:</span>
        <span className="value">{email.patternType}</span>
      </div>
      
      {email.foundInBreach && (
        <div className="detail warning">
          <span className="label">⚠️ Breach:</span>
          <span className="value">Found in data breaches</span>
        </div>
      )}
      
      {email.socialProfiles.length > 0 && (
        <div className="detail">
          <span className="label">Profiles:</span>
          <div className="social-profiles">
            {email.socialProfiles.map((profile, idx) => (
              <span key={idx} className="social-profile">
                {profile.platform}: {profile.username}
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  </div>
);
