// tests/load/k6/smoke_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// Custom metrics
const failureRate = new Rate('failed_requests');
const requestDuration = new Trend('request_duration');
const totalRequests = new Counter('total_requests');

export const options = {
  stages: [
    { duration: '2m', target: 10 },  // Ramp up to 10 users
    { duration: '5m', target: 10 },  // Stay at 10 users
    { duration: '2m', target: 0 },   // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% of requests under 2s
    http_req_failed: ['rate<0.01'],    // Less than 1% failures
    failed_requests: ['rate<0.05'],    // Custom failure rate
  },
};

const BASE_URL = __ENV.TEST_URL || 'http://localhost:8080';
const API_KEY = __ENV.API_KEY;

export default function () {
  const params = {
    headers: {
      'Authorization': `Bearer ${API_KEY}`,
      'Content-Type': 'application/json',
      'X-Tenant-ID': 'test-tenant',
    },
    tags: { name: 'phone_lookup' },
  };

  // Test phone lookup endpoint
  const phoneLookupPayload = JSON.stringify({
    phone_numbers: ['+989123456789'],
    platforms: ['facebook', 'instagram'],
  });

  const responses = http.batch([
    ['POST', `${BASE_URL}/api/v1/intelligence/phone-lookup`, phoneLookupPayload, params],
    ['GET', `${BASE_URL}/api/v1/health`, null, params],
    ['GET', `${BASE_URL}/api/v1/admin/usage`, null, params],
  ]);

  // Check phone lookup response
  const phoneLookupResp = responses[0];
  totalRequests.add(1);
  
  const phoneLookupSuccess = check(phoneLookupResp, {
    'phone lookup status is 200': (r) => r.status === 200,
    'phone lookup has job_id': (r) => r.json('job_id') !== undefined,
    'phone lookup response time acceptable': (r) => r.timings.duration < 5000,
  });

  failureRate.add(!phoneLookupSuccess);
  requestDuration.add(phoneLookupResp.timings.duration);

  // Check health endpoint
  check(responses[1], {
    'health check status is 200': (r) => r.status === 200,
  });

  sleep(1);
}
