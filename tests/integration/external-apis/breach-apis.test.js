// tests/integration/external-apis/breach-apis.test.js
const { describe, it, before } = require('mocha');
const { expect } = require('chai');
const axios = require('axios');
const BreachChecker = require('../../../intelligence-engine/internal/email_discovery/breach-checker');

describe('External Breach API Integration Tests', function() {
    this.timeout(30000);
    
    let breachChecker;

    before(() => {
        breachChecker = new BreachChecker();
    });

    describe('Have I Been Pwned API', () => {
        it('should integrate with HIBP API successfully', async () => {
            const testEmails = [
                'test@example.com',
                'user@gmail.com'
            ];

            for (const email of testEmails) {
                const result = await breachChecker.checkHIBP(email);
                
                expect(result).to.have.property('email', email);
                expect(result).to.have.property('found').that.is.a('boolean');
                expect(result).to.have.property('breaches').that.is.an('array');
                
                if (result.found) {
                    expect(result.breaches.length).to.be.greaterThan(0);
                    expect(result).to.have.property('firstBreached').that.is.a('string');
                }
                
                // Respect API rate limits
                await new Promise(resolve => setTimeout(resolve, 1600)); // HIBP: 1 request/1.6s
            }
        });

        it('should handle HIBP API errors gracefully', async () => {
            // Test with invalid API key to trigger error
            const originalApiKey = process.env.HIBP_API_KEY;
            process.env.HIBP_API_KEY = 'invalid-key';
            
            try {
                const result = await breachChecker.checkHIBP('test@example.com');
                expect(result.found).to.be.false;
            } catch (error) {
                expect(error).to.have.property('message').that.includes('API');
            } finally {
                process.env.HIBP_API_KEY = originalApiKey;
            }
        });
    });

    describe('Dehashed API Integration', () => {
        it('should query Dehashed API for breach data', async () => {
            const testQuery = 'domain:example.com';
            
            const result = await breachChecker.checkDehashed(testQuery);
            
            expect(result).to.have.property('success').that.is.a('boolean');
            expect(result).to.have.property('results').that.is.an('array');
            expect(result).to.have.property('total').that.is.a('number');
            
            if (result.results.length > 0) {
                const breach = result.results[0];
                expect(breach).to.have.property('email').that.is.a('string');
                expect(breach).to.have.property('password').that.is.a('string');
            }
        });
    });

    describe('Breach Data Correlation', () => {
        it('should correlate data across multiple breach APIs', async () => {
            const testPhone = '+989123456789';
            const generatedEmails = [
                '989123456789@gmail.com',
                '09123456789@yahoo.com'
            ];

            const allResults = [];
            
            for (const email of generatedEmails) {
                const hibpResult = await breachChecker.checkHIBP(email);
                const dehashedResult = await breachChecker.checkDehashed(`email:${email}`);
                
                allResults.push({
                    email,
                    hibp: hibpResult,
                    dehashed: dehashedResult
                });
                
                await new Promise(resolve => setTimeout(resolve, 2000));
            }

            // Analyze correlation results
            const foundInAny = allResults.some(r => r.hibp.found || r.dehashed.success);
            console.log(`📊 Found in breaches: ${foundInAny}`);
            
            expect(allResults).to.have.lengthOf(generatedEmails.length);
        });
    });
});
