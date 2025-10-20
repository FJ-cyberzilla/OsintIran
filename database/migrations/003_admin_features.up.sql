-- database/migrations/003_admin_features.up.sql

-- Bulk Jobs Table
CREATE TABLE bulk_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    total_numbers INTEGER NOT NULL,
    processed_numbers INTEGER DEFAULT 0,
    platforms JSONB NOT NULL,
    priority VARCHAR(50) DEFAULT 'normal',
    status VARCHAR(50) DEFAULT 'pending',
    file_path VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Cost Analytics Table
CREATE TABLE cost_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date DATE NOT NULL,
    category VARCHAR(100) NOT NULL,
    platform VARCHAR(100),
    proxy_type VARCHAR(50),
    cost DECIMAL(10, 4) NOT NULL,
    requests INTEGER NOT NULL,
    successful_requests INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Performance Metrics Table
CREATE TABLE performance_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMP NOT NULL,
    metric_type VARCHAR(100) NOT NULL,
    value DECIMAL(10, 4) NOT NULL,
    tags JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Proxy Acquisition History
CREATE TABLE proxy_acquisitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider VARCHAR(100) NOT NULL,
    proxy_type VARCHAR(50) NOT NULL,
    count INTEGER NOT NULL,
    cost DECIMAL(10, 2) NOT NULL,
    countries JSONB,
    acquired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active'
);

-- Create indexes for better performance
CREATE INDEX idx_bulk_jobs_status ON bulk_jobs(status);
CREATE INDEX idx_bulk_jobs_user_id ON bulk_jobs(user_id);
CREATE INDEX idx_cost_analytics_date ON cost_analytics(date);
CREATE INDEX idx_performance_metrics_timestamp ON performance_metrics(timestamp);
CREATE INDEX idx_proxy_acquisitions_provider ON proxy_acquisitions(provider);
