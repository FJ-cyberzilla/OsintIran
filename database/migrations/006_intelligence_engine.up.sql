-- database/migrations/006_intelligence_engine.up.sql

-- Email Discovery Results
CREATE TABLE email_discovery_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    confidence DECIMAL(3,2) NOT NULL,
    source VARCHAR(100) NOT NULL,
    pattern_type VARCHAR(100),
    found_in_breach BOOLEAN DEFAULT FALSE,
    social_profiles JSONB,
    verified BOOLEAN DEFAULT FALSE,
    discovered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(phone_number, email)
);

-- Social Graph Data
CREATE TABLE social_graphs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(50) NOT NULL,
    graph_data JSONB NOT NULL,
    node_count INTEGER NOT NULL,
    edge_count INTEGER NOT NULL,
    density DECIMAL(4,3),
    central_node VARCHAR(255),
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Behavioral Patterns
CREATE TABLE behavioral_patterns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_identifier VARCHAR(255) NOT NULL,
    platform VARCHAR(100) NOT NULL,
    activity_type VARCHAR(100) NOT NULL,
    pattern_data JSONB NOT NULL,
    confidence DECIMAL(3,2) NOT NULL,
    anomaly_score DECIMAL(3,2),
    detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Risk Assessments
CREATE TABLE risk_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(50) NOT NULL,
    overall_score DECIMAL(3,2) NOT NULL,
    risk_level VARCHAR(20) NOT NULL,
    factors JSONB NOT NULL,
    confidence DECIMAL(3,2) NOT NULL,
    assessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Intelligence Reports
CREATE TABLE intelligence_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id VARCHAR(100) UNIQUE NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    report_type VARCHAR(50) NOT NULL,
    report_data JSONB NOT NULL,
    executive_summary JSONB,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_email_discovery_phone ON email_discovery_results(phone_number);
CREATE INDEX idx_email_discovery_email ON email_discovery_results(email);
CREATE INDEX idx_social_graphs_phone ON social_graphs(phone_number);
CREATE INDEX idx_behavioral_patterns_user ON behavioral_patterns(user_identifier);
CREATE INDEX idx_risk_assessments_phone ON risk_assessments(phone_number);
CREATE INDEX idx_intelligence_reports_phone ON intelligence_reports(phone_number);
CREATE INDEX idx_intelligence_reports_created ON intelligence_reports(generated_at);
