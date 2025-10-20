// tests/load/k6/load_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Trend, Rate } from 'k6/metrics';

export const options = {
  stages: [
    { duration: '5m', target: 50 },   // Ramp up to 50 users
    { duration: '10m', target: 50 },  // Stay at 50 users
    { duration: '3m', target: 100 },  // Spike to 100 users
    { duration: '10m', target: 100 }, // Stay at 100 users  
    { duration: '5m', target: 0 },    // Ramp down to 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<3000', 'p(99)<5000'],
    http_req_failed: ['rate<0.02'],
    'http_reqs{endpoint:phone-lookup}': ['count>1000'],
  },
  ext: {
    loadimpact: {
      projectID: 123456,
      name: 'Phone Intelligence Load Test'
    }
  }
};

const BASE_URL = __ENV.TEST_URL || 'http://localhost:8080';
const API_KEYS = JSON.parse(__ENV.API_KEYS || '["key1", "key2"]');

// Test data
const phoneNumbers = [
  '+989123456789', '+989123456790', '+989123456791',
  '+989123456792', '+989123456793', '+989123456794',
];

export default function () {
  const randomApiKey = API_KEYS[Math.floor(Math.random() * API_KEYS.length)];
  const randomPhone = phoneNumbers[Math.floor(Math.random() * phoneNumbers.length)];
  
  const params = {
    headers: {
      'Authorization': `Bearer ${randomApiKey}`,
      'Content-Type': 'application/json',
      'X-Tenant-ID': `tenant-${Math.floor(Math.random() * 10)}`,
    },
    tags: { endpoint: 'phone-lookup' },
  };

  // Test 1: Phone lookup
  const lookupResp = http.post(
    `${BASE_URL}/api/v1/intelligence/phone-lookup`,
    JSON.stringify({
      phone_numbers: [randomPhone],
      platforms: ['facebook', 'instagram', 'twitter'],
    }),
    params
  );

  check(lookupResp, {
    'phone lookup successful': (r) => r.status === 200,
    'response has job id': (r) => r.json('job_id') !== undefined,
  });

  // Test 2: Email discovery (if lookup was successful)
  if (lookupResp.status === 200) {
    const jobId = lookupResp.json('job_id');
    
    sleep(2); // Wait for processing
    
    const emailResp = http.post(
      `${BASE_URL}/api/v1/intelligence/email-discovery`,
      JSON.stringify({ job_id: jobId }),
      params
    );

    check(emailResp, {
      'email discovery successful': (r) => r.status === 200,
    });
  }

  // Test 3: Bulk operations (every 10th request)
  if (__ITER % 10 === 0) {
    const bulkPayload = JSON.stringify({
      phone_numbers: phoneNumbers.slice(0, 3),
      platforms: ['facebook', 'linkedin'],
      priority: 'high'
    });

    const bulkResp = http.post(
      `${BASE_URL}/api/v1/intelligence/bulk-operations`,
      bulkPayload,
      params
    );

    check(bulkResp, {
      'bulk operation successful': (r) => r.status === 202,
    });
  }

  sleep(Math.random() * 2 + 1); // Random sleep between 1-3 seconds
}
