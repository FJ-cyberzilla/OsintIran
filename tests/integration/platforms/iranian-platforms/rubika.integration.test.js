// tests/integration/platforms/iranian-platforms/rubika.integration.test.js
const { describe, it, before, after } = require('mocha');
const { expect } = require('chai');
const RubikaAgent = require('../../../../browser-workers/src/platforms/iranian/rubika-agent');

describe('Rubika Platform Integration Tests', function() {
    this.timeout(90000); // 90 seconds for Iranian platforms
    
    let rubikaAgent;
    let iranianProxy;

    before(async () => {
        console.log('🚀 Initializing Rubika integration tests...');
        
        // Get Iranian residential proxy
        const ProxyManager = require('../../../../browser-workers/src/proxy/proxy-manager');
        const proxyManager = new ProxyManager();
        iranianProxy = await proxyManager.getResidentialProxy('IR');
        
        // Initialize Rubika agent with Iranian configuration
        rubikaAgent = new RubikaAgent();
        await rubikaAgent.initialize({
            proxy: iranianProxy,
            userAgent: 'Mozilla/5.0 (Linux; Android 10; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36',
            language: 'fa-IR',
            timezone: 'Asia/Tehran'
        });
    });

    after(async () => {
        if (rubikaAgent) {
            await rubikaAgent.cleanup();
        }
    });

    describe('Rubika Phone Search', () => {
        it('should search for Iranian phone numbers on Rubika', async () => {
            const iranianNumbers = [
                '+989123456789',
                '+989351234567',
                '+989191234567'
            ];

            for (const phoneNumber of iranianNumbers) {
                console.log(`🔍 Searching Rubika for: ${phoneNumber}`);
                
                const result = await rubikaAgent.searchByPhone(phoneNumber);
                
                // Validate Rubika-specific response
                expect(result).to.have.property('success').that.is.a('boolean');
                expect(result).to.have.property('platform', 'rubika');
                expect(result).to.have.property('profiles').that.is.an('array');
                
                if (result.profiles.length > 0) {
                    const profile = result.profiles[0];
                    expect(profile).to.have.property('username').that.is.a('string');
                    expect(profile).to.have.property('profileUrl').that.includes('rubika');
                    expect(profile).to.have.property('isVerified').that.is.a('boolean');
                    
                    // Rubika-specific fields
                    expect(profile).to.have.property('lastSeen').that.is.a('string');
                    expect(profile).to.have.property('isContact').that.is.a('boolean');
                }
                
                await new Promise(resolve => setTimeout(resolve, 3000));
            }
        });

        it('should handle Rubika-specific authentication', async () => {
            const authResult = await rubikaAgent.authenticate();
            
            expect(authResult).to.have.property('authenticated').that.is.a('boolean');
            expect(authResult).to.have.property('requires2FA').that.is.a('boolean');
            
            if (authResult.requires2FA) {
                expect(authResult).to.have.property('authMethod').that.is.a('string');
            }
        });
    });
});
