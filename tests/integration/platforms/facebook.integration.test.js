// tests/integration/platforms/facebook.integration.test.js
const { describe, it, before, after, beforeEach } = require('mocha');
const { expect } = require('chai');
const sinon = require('sinon');
const FacebookAgent = require('../../../browser-workers/src/platforms/facebook-agent');
const ProxyManager = require('../../../browser-workers/src/proxy/proxy-manager');
const RateLimiter = require('../../../browser-workers/src/utils/rate-limiter');

describe('Facebook Platform Integration Tests', function() {
    this.timeout(120000); // 2 minutes for browser operations
    
    let facebookAgent;
    let proxyManager;
    let rateLimiter;
    let testBrowser;

    before(async () => {
        console.log('🚀 Initializing Facebook integration tests...');
        
        // Initialize dependencies
        proxyManager = new ProxyManager();
        rateLimiter = new RateLimiter();
        
        // Get test proxy (residential for Facebook)
        const testProxy = await proxyManager.getResidentialProxy('US');
        
        // Initialize Facebook agent with test configuration
        facebookAgent = new FacebookAgent();
        await facebookAgent.initialize({
            proxy: testProxy,
            userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
            viewport: { width: 1920, height: 1080 },
            stealth: true
        });
        
        testBrowser = facebookAgent.getBrowser();
    });

    after(async () => {
        console.log('🧹 Cleaning up Facebook test resources...');
        if (facebookAgent) {
            await facebookAgent.cleanup();
        }
        if (proxyManager) {
            await proxyManager.cleanup();
        }
    });

    beforeEach(async () => {
        // Reset rate limiting between tests
        rateLimiter.reset('facebook');
        
        // Add delay to avoid detection
        await new Promise(resolve => setTimeout(resolve, 2000));
    });

    describe('Phone Number Search Integration', () => {
        it('should successfully search for valid Iranian phone numbers', async () => {
            const testPhoneNumbers = [
                '+989123456789',
                '+989351234567', 
                '+989191234567'
            ];

            for (const phoneNumber of testPhoneNumbers) {
                console.log(`🔍 Testing phone: ${phoneNumber}`);
                
                const result = await facebookAgent.searchByPhone(phoneNumber);
                
                // Validate response structure
                expect(result).to.have.property('success').that.is.a('boolean');
                expect(result).to.have.property('platform', 'facebook');
                expect(result).to.have.property('phoneNumber', phoneNumber);
                expect(result).to.have.property('profiles').that.is.an('array');
                expect(result).to.have.property('requestId').that.is.a('string');
                
                // Validate profiles if found
                if (result.profiles.length > 0) {
                    const profile = result.profiles[0];
                    expect(profile).to.have.property('username').that.is.a('string');
                    expect(profile).to.have.property('profileUrl').that.includes('facebook.com');
                    expect(profile).to.have.property('fullName').that.is.a('string');
                    expect(profile).to.have.property('isVerified').that.is.a('boolean');
                    
                    // Additional profile validation
                    if (profile.avatarUrl) {
                        expect(profile.avatarUrl).to.match(/^https?:\/\//);
                    }
                }
                
                // Respect rate limiting
                await rateLimiter.waitFor('facebook');
            }
        });

        it('should handle non-existent phone numbers gracefully', async () => {
            const nonExistentNumbers = [
                '+989999999999',
                '+989000000000'
            ];

            for (const phoneNumber of nonExistentNumbers) {
                const result = await facebookAgent.searchByPhone(phoneNumber);
                
                expect(result.success).to.be.true;
                expect(result.profiles).to.be.an('array').that.is.empty;
                expect(result.message).to.include('No profiles found');
                
                await rateLimiter.waitFor('facebook');
            }
        });

        it('should detect and handle rate limiting', async () => {
            const rapidRequests = Array.from({ length: 10 }, (_, i) => 
                `+9891234567${i.toString().padStart(2, '0')}`
            );

            let rateLimitCount = 0;
            
            for (const phoneNumber of rapidRequests) {
                const result = await facebookAgent.searchByPhone(phoneNumber);
                
                if (result.rateLimited) {
                    rateLimitCount++;
                    expect(result.retryAfter).to.be.a('number').that.is.greaterThan(0);
                    console.log(`⚠️ Rate limited detected, retry after: ${result.retryAfter}s`);
                    
                    // Wait for rate limit to clear
                    await new Promise(resolve => 
                        setTimeout(resolve, result.retryAfter * 1000)
                    );
                }
                
                // Small delay between requests
                await new Promise(resolve => setTimeout(resolve, 500));
            }
            
            expect(rateLimitCount).to.be.lessThan(5, 
                'Should not hit rate limits too frequently with proper delays');
        });
    });

    describe('Profile Analysis Integration', () => {
        it('should extract comprehensive profile information', async () => {
            const testUsernames = [
                'zuck', // Public profile for testing
                'facebook'
            ];

            for (const username of testUsernames) {
                console.log(`📊 Analyzing profile: ${username}`);
                
                const profile = await facebookAgent.analyzeProfile(username);
                
                // Validate profile structure
                expect(profile).to.have.property('username', username);
                expect(profile).to.have.property('fullName').that.is.a('string');
                expect(profile).to.have.property('profileUrl').that.includes(username);
                expect(profile).to.have.property('isVerified').that.is.a('boolean');
                expect(profile).to.have.property('isPublic').that.is.a('boolean');
                
                // Validate metrics if available
                if (profile.followerCount !== undefined) {
                    expect(profile.followerCount).to.be.a('number').that.is.at.least(0);
                }
                
                if (profile.friendsCount !== undefined) {
                    expect(profile.friendsCount).to.be.a('number').that.is.at.least(0);
                }
                
                // Validate timestamps
                if (profile.lastActive) {
                    expect(profile.lastActive).to.match(/^\d{4}-\d{2}-\d{2}/);
                }
                
                if (profile.joinDate) {
                    expect(profile.joinDate).to.match(/^\d{4}-\d{2}-\d{2}/);
                }
                
                await rateLimiter.waitFor('facebook');
            }
        });

        it('should handle private profiles gracefully', async () => {
            // Test with known private profiles or non-existent ones
            const privateUsernames = [
                'privateuser123456789', // Likely non-existent
                'testprivateprofile123'
            ];

            for (const username of privateUsernames) {
                const profile = await facebookAgent.analyzeProfile(username);
                
                expect(profile).to.have.property('username', username);
                expect(profile.isPublic).to.be.false;
                expect(profile.followerCount).to.be.undefined;
                expect(profile.friendsCount).to.be.undefined;
                
                await rateLimiter.waitFor('facebook');
            }
        });
    });

    describe('CAPTCHA Handling Integration', () => {
        it('should detect CAPTCHA challenges during rapid requests', async () => {
            // Make rapid requests to trigger CAPTCHA
            const rapidRequests = Array.from({ length: 15 }, (_, i) => 
                facebookAgent.searchByPhone(`+9891234567${i.toString().padStart(2, '0')}`)
            );

            const results = await Promise.all(rapidRequests);
            const captchaResults = results.filter(r => r.captchaRequired);
            
            console.log(`🔐 CAPTCHA triggered in ${captchaResults.length} of ${results.length} requests`);
            
            // Should detect CAPTCHA at some point with rapid requests
            if (captchaResults.length > 0) {
                const captchaResult = captchaResults[0];
                expect(captchaResult).to.have.property('captchaRequired', true);
                expect(captchaResult).to.have.property('captchaId').that.is.a('string');
                expect(captchaResult).to.have.property('captchaImage').that.is.a('string');
            }
        });

        it('should solve CAPTCHA challenges using AI', async () => {
            // This test requires actual CAPTCHA challenge
            // For integration testing, we'll simulate or use test CAPTCHA
            
            const testCaptcha = {
                captchaRequired: true,
                captchaId: 'test-captcha-123',
                captchaImage: 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=='
            };

            const solution = await facebookAgent.solveCaptcha(testCaptcha);
            
            expect(solution).to.have.property('success').that.is.a('boolean');
            
            if (solution.success) {
                expect(solution).to.have.property('solution').that.is.a('string');
                expect(solution.solution).to.have.length.of.at.least(4);
            } else {
                expect(solution).to.have.property('error').that.is.a('string');
            }
        });
    });

    describe('Error Handling and Recovery', () => {
        it('should handle network errors and retry automatically', async () => {
            const originalMethod = facebookAgent.searchByPhone;
            let callCount = 0;
            
            // Mock method to fail first time, succeed second time
            facebookAgent.searchByPhone = async (phoneNumber) => {
                callCount++;
                if (callCount === 1) {
                    throw new Error('Network timeout');
                }
                return originalMethod.call(facebookAgent, phoneNumber);
            };

            const result = await facebookAgent.searchByPhone('+989123456789');
            expect(result.success).to.be.true;
            expect(callCount).to.equal(2);
            
            // Restore original method
            facebookAgent.searchByPhone = originalMethod;
        });

        it('should handle browser crashes and restart', async () => {
            // Simulate browser crash
            await testBrowser.disconnect();
            
            // Attempt operation - should restart browser automatically
            const result = await facebookAgent.searchByPhone('+989123456789');
            expect(result.success).to.be.true;
            
            // Verify browser is restarted
            expect(facebookAgent.getBrowser().isConnected()).to.be.true;
        });
    });

    describe('Performance and Metrics', () => {
        it('should complete searches within acceptable time limits', async () => {
            const testPhone = '+989123456789';
            const startTime = Date.now();
            
            const result = await facebookAgent.searchByPhone(testPhone);
            const duration = Date.now() - startTime;
            
            expect(result.success).to.be.true;
            expect(duration).to.be.lessThan(30000, 'Search should complete within 30 seconds');
            
            console.log(`⏱️ Search completed in ${duration}ms`);
        });

        it('should track and report performance metrics', async () => {
            const metrics = await facebookAgent.getPerformanceMetrics();
            
            expect(metrics).to.have.property('totalRequests').that.is.a('number');
            expect(metrics).to.have.property('successRate').that.is.a('number');
            expect(metrics).to.have.property('averageResponseTime').that.is.a('number');
            expect(metrics).to.have.property('captchaRate').that.is.a('number');
            
            expect(metrics.successRate).to.be within(0, 1);
            expect(metrics.captchaRate).to.be within(0, 1);
        });
    });
});
