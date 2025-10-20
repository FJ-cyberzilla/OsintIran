# tests/load/locustfile.py
from locust import HttpUser, task, between, TaskSet
import random
import json

class PhoneIntelligenceTasks(TaskSet):
    
    def on_start(self):
        """Initialize user with test data"""
        self.phone_numbers = [
            "+989123456789", "+989123456790", "+989123456791",
            "+989123456792", "+989123456793", "+989123456794"
        ]
        self.platforms = ["facebook", "instagram", "twitter", "linkedin"]
        self.tenant_id = f"tenant-{random.randint(1, 100)}"
        
    @task(3)
    def phone_lookup(self):
        """Test phone lookup endpoint"""
        phone = random.choice(self.phone_numbers)
        payload = {
            "phone_numbers": [phone],
            "platforms": random.sample(self.platforms, 2)
        }
        
        headers = {
            "Authorization": f"Bearer {self.get_api_key()}",
            "Content-Type": "application/json",
            "X-Tenant-ID": self.tenant_id
        }
        
        with self.client.post(
            "/api/v1/intelligence/phone-lookup",
            json=payload,
            headers=headers,
            catch_response=True,
            name="Phone Lookup"
        ) as response:
            if response.status_code == 200:
                response.success()
                self.job_id = response.json().get("job_id")
            elif response.status_code == 429:
                response.failure("Rate limit exceeded")
            else:
                response.failure(f"Failed with status {response.status_code}")
    
    @task(2)
    def email_discovery(self):
        """Test email discovery endpoint"""
        if hasattr(self, 'job_id'):
            payload = {"job_id": self.job_id}
            
            headers = {
                "Authorization": f"Bearer {self.get_api_key()}",
                "Content-Type": "application/json"
            }
            
            with self.client.post(
                "/api/v1/intelligence/email-discovery",
                json=payload,
                headers=headers,
                catch_response=True,
                name="Email Discovery"
            ) as response:
                if response.status_code == 200:
                    response.success()
                else:
                    response.failure(f"Failed with status {response.status_code}")
    
    @task(1)
    def bulk_operations(self):
        """Test bulk operations endpoint"""
        payload = {
            "phone_numbers": random.sample(self.phone_numbers, 3),
            "platforms": self.platforms,
            "priority": "normal"
        }
        
        headers = {
            "Authorization": f"Bearer {self.get_api_key()}",
            "Content-Type": "application/json"
        }
        
        with self.client.post(
            "/api/v1/intelligence/bulk-operations",
            json=payload,
            headers=headers,
            name="Bulk Operations"
        ) as response:
            if response.status_code == 202:
                response.success()
            else:
                response.failure(f"Failed with status {response.status_code}")
    
    @task(5)
    def health_check(self):
        """Test health endpoint"""
        self.client.get("/api/v1/health", name="Health Check")
    
    def get_api_key(self):
        """Get API key for tenant"""
        # In real implementation, this would rotate keys
        return "test-api-key"

class WebsiteUser(HttpUser):
    tasks = [PhoneIntelligenceTasks]
    wait_time = between(1, 5)  # Wait 1-5 seconds between tasks
    
    def on_start(self):
        self.client.headers = {
            "User-Agent": "Locust Load Test",
            "Accept": "application/json"
        }
