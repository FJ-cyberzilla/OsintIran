// tests/integration/platforms/facebook.test.js
const FacebookAgent = require('../../../browser-workers/src/platforms/facebook-agent');
const ProxyManager = require('../../../browser-workers/src/proxy/proxy-manager');
const { describe, it, before, after } = require('mocha');
const { expect } = require('chai');
const sinon = require('sinon');
const { URL } = require('url');
describe('Facebook Platform Integration', function() {
  this.timeout(60000); // 60 second timeout for browser operations
  
  let facebookAgent;
  let proxyManager;
  let browser;

  before(async () => {
    proxyManager = new ProxyManager();
    facebookAgent = new FacebookAgent();
    
    // Initialize with test proxy
    await facebookAgent.initialize({
      proxy: await proxyManager.getTestProxy(),
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
    });
  });

  after(async () => {
    await facebookAgent.cleanup();
  });

  describe('Phone Number Search', () => {
    it('should search for phone number and return profiles', async () => {
      const testPhone = '+989123456789';
      
      const result = await facebookAgent.searchByPhone(testPhone);
      
      expect(result).to.have.property('success', true);
      expect(result).to.have.property('platform', 'facebook');
      expect(result).to.have.property('profiles').that.is.an('array');
      
      if (result.profiles.length > 0) {
        const profile = result.profiles[0];
        expect(profile).to.have.property('username').that.is.a('string');
        expect(profile).to.have.property('profileUrl').that.is.a('string');
        const allowedFacebookHosts = ['facebook.com', 'www.facebook.com'];
        const urlHost = (new URL(profile.profileUrl)).hostname;
        expect(allowedFacebookHosts).to.include(urlHost);
        expect(profile).to.have.property('verified').that.is.a('boolean');
      }
    });

    it('should handle non-existent phone numbers gracefully', async () => {
      const nonExistentPhone = '+989999999999';
      
      const result = await facebookAgent.searchByPhone(nonExistentPhone);
      
      expect(result).to.have.property('success', true);
      expect(result.profiles).to.be.an('array').that.is.empty;
    });

    it('should respect rate limiting', async () => {
      const phones = Array.from({ length: 5 }, (_, i) => `+9891234567${i}`);
      
      const results = await Promise.all(
        phones.map(phone => facebookAgent.searchByPhone(phone))
      );
      
      // Check if any requests were rate limited
      const rateLimited = results.filter(r => r.rateLimited);
      expect(rateLimited.length).to.be.lessThan(3); // Should not hit hard limits
    });
  });

  describe('Profile Analysis', () => {
    it('should extract profile information correctly', async () => {
      const testUsername = 'testuser'; // Mock username for testing
      
      const profile = await facebookAgent.analyzeProfile(testUsername);
      
      expect(profile).to.have.property('username', testUsername);
      expect(profile).to.have.property('fullName').that.is.a('string');
      expect(profile).to.have.property('followerCount').that.is.a('number');
      expect(profile).to.have.property('friendsCount').that.is.a('number');
      expect(profile).to.have.property('isVerified').that.is.a('boolean');
      expect(profile).to.have.property('lastActive').that.is.a('string');
    });

    it('should detect profile visibility settings', async () => {
      const testUsername = 'publicuser';
      
      const profile = await facebookAgent.analyzeProfile(testUsername);
      
      expect(profile).to.have.property('isPublic').that.is.a('boolean');
      expect(profile).to.have.property('canMessage').that.is.a('boolean');
      expect(profile).to.have.property('canFriendRequest').that.is.a('boolean');
    });
  });

  describe('CAPTCHA Handling', () => {
    it('should detect CAPTCHA challenges', async () => {
      // Force CAPTCHA by making rapid requests
      const rapidRequests = Array.from({ length: 10 }, (_, i) => 
        facebookAgent.searchByPhone(`+9891234567${i}`)
      );
      
      const results = await Promise.all(rapidRequests);
      const captchaResults = results.filter(r => r.captchaRequired);
      
      // Should detect CAPTCHA at some point
      expect(captchaResults.length).to.be.greaterThan(0);
    });

    it('should solve CAPTCHA challenges when detected', async () => {
      const captchaResult = {
        captchaRequired: true,
        captchaImage: 'base64-encoded-image',
        captchaId: 'test-captcha-id'
      };
      
      const solved = await facebookAgent.solveCaptcha(captchaResult);
      
      expect(solved).to.have.property('success').that.is.a('boolean');
      if (solved.success) {
        expect(solved).to.have.property('solution').that.is.a('string');
      }
    });
  });
});
