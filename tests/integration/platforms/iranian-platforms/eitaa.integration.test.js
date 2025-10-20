// tests/integration/platforms/iranian-platforms/eitaa.integration.test.js
const { describe, it, before, after } = require('mocha');
const { expect } = require('chai');
const EitaaAgent = require('../../../../browser-workers/src/platforms/iranian/eitaa-agent');

describe('Eitaa Platform Integration Tests', function() {
    this.timeout(90000);
    
    let eitaaAgent;

    before(async () => {
        console.log('🚀 Initializing Eitaa integration tests...');
        
        const ProxyManager = require('../../../../browser-workers/src/proxy/proxy-manager');
        const proxyManager = new ProxyManager();
        const iranianProxy = await proxyManager.getResidentialProxy('IR');
        
        eitaaAgent = new EitaaAgent();
        await eitaaAgent.initialize({
            proxy: iranianProxy,
            userAgent: 'Mozilla/5.0 (Linux; Android 11; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.115 Mobile Safari/537.36',
            language: 'fa-IR'
        });
    });

    after(async () => {
        if (eitaaAgent) {
            await eitaaAgent.cleanup();
        }
    });

    describe('Eitaa Phone Search', () => {
        it('should search for Iranian phone numbers on Eitaa', async () => {
            const testNumbers = [
                '+989123456789',
                '+989351234567'
            ];

            for (const phoneNumber of testNumbers) {
                const result = await eitaaAgent.searchByPhone(phoneNumber);
                
                expect(result).to.have.property('success').that.is.a('boolean');
                expect(result).to.have.property('platform', 'eitaa');
                expect(result).to.have.property('profiles').that.is.an('array');
                
                if (result.profiles.length > 0) {
                    const profile = result.profiles[0];
                    expect(profile).to.have.property('username').that.is.a('string');
                    expect(profile).to.have.property('isChannel').that.is.a('boolean');
                    expect(profile).to.have.property('membersCount').that.is.a('number');
                }
                
                await new Promise(resolve => setTimeout(resolve, 2000));
            }
        });
    });
});
