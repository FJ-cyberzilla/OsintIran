// tests/integration/platforms/platform-test-runner.js
const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

class PlatformTestRunner {
    constructor() {
        this.results = {
            passed: 0,
            failed: 0,
            skipped: 0,
            platforms: {}
        };
        this.startTime = null;
    }

    async runAllPlatformTests() {
        console.log('🚀 Starting Platform Integration Test Suite\n');
        this.startTime = Date.now();

        const platforms = [
            { name: 'Facebook', file: 'facebook.integration.test.js' },
            { name: 'Instagram', file: 'instagram.integration.test.js' },
            { name: 'Twitter', file: 'twitter.integration.test.js' },
            { name: 'LinkedIn', file: 'linkedin.integration.test.js' },
            { name: 'WhatsApp', file: 'whatsapp.integration.test.js' },
            { name: 'Telegram', file: 'telegram.integration.test.js' },
            { name: 'Rubika', file: 'iranian-platforms/rubika.integration.test.js' },
            { name: 'Eitaa', file: 'iranian-platforms/eitaa.integration.test.js' },
            { name: 'Soroush', file: 'iranian-platforms/soroush.integration.test.js' },
            { name: 'Bale', file: 'iranian-platforms/bale.integration.test.js' },
        ];

        for (const platform of platforms) {
            await this.runPlatformTest(platform);
        }

        this.generateReport();
    }

    async runPlatformTest(platform) {
        console.log(`\n🔍 Testing ${platform.name} Integration...`);
        
        const testFile = path.join(__dirname, platform.file);
        
        if (!fs.existsSync(testFile)) {
            console.log(`   ⚠️  Test file not found: ${testFile}`);
            this.results.skipped++;
            this.results.platforms[platform.name] = 'skipped';
            return;
        }

        return new Promise((resolve) => {
            const mochaProcess = spawn('npx', [
                'mocha',
                testFile,
                '--timeout', '120000',
                '--reporter', 'json'
            ], {
                stdio: ['pipe', 'pipe', 'pipe'],
                env: { ...process.env, NODE_ENV: 'test' }
            });

            let stdout = '';
            let stderr = '';

            mochaProcess.stdout.on('data', (data) => {
                stdout += data.toString();
            });

            mochaProcess.stderr.on('data', (data) => {
                stderr += data.toString();
            });

            mochaProcess.on('close', (code) => {
                try {
                    const result = JSON.parse(stdout);
                    
                    if (result.stats.failures === 0) {
                        console.log(`   ✅ ${platform.name}: PASSED (${result.stats.passes} tests)`);
                        this.results.passed++;
                        this.results.platforms[platform.name] = 'passed';
                    } else {
                        console.log(`   ❌ ${platform.name}: FAILED (${result.stats.failures} failures)`);
                        this.results.failed++;
                        this.results.platforms[platform.name] = 'failed';
                        
                        // Log failures
                        result.failures.forEach(failure => {
                            console.log(`      - ${failure.fullTitle}`);
                            console.log(`        ${failure.err.message}`);
                        });
                    }
                } catch (parseError) {
                    console.log(`   ❌ ${platform.name}: ERROR - ${parseError.message}`);
                    this.results.failed++;
                    this.results.platforms[platform.name] = 'error';
                }

                if (stderr) {
                    console.log(`   ⚠️  ${platform.name} stderr: ${stderr}`);
                }

                resolve();
            });

            // Timeout after 3 minutes
            setTimeout(() => {
                mochaProcess.kill();
                console.log(`   ⏰ ${platform.name}: TIMEOUT`);
                this.results.failed++;
                this.results.platforms[platform.name] = 'timeout';
                resolve();
            }, 180000);
        });
    }

    generateReport() {
        const duration = Date.now() - this.startTime;
        const minutes = Math.floor(duration / 60000);
        const seconds = ((duration % 60000) / 1000).toFixed(0);

        console.log('\n' + '='.repeat(60));
        console.log('📊 PLATFORM INTEGRATION TEST REPORT');
        console.log('='.repeat(60));
        
        console.log(`\n⏱️  Duration: ${minutes}m ${seconds}s`);
        console.log(`📈 Summary: ${this.results.passed} passed, ${this.results.failed} failed, ${this.results.skipped} skipped`);
        
        console.log('\n🔧 Platform Results:');
        Object.entries(this.results.platforms).forEach(([platform, status]) => {
            const icon = status === 'passed' ? '✅' : status === 'failed' ? '❌' : '⚠️';
            console.log(`   ${icon} ${platform}: ${status}`);
        });

        // Calculate success rate
        const total = this.results.passed + this.results.failed + this.results.skipped;
        const successRate = total > 0 ? (this.results.passed / total) * 100 : 0;

        console.log(`\n🎯 Success Rate: ${successRate.toFixed(1)}%`);

        if (this.results.failed === 0) {
            console.log('\n🎉 ALL PLATFORM INTEGRATION TESTS PASSED!');
        } else {
            console.log(`\n❌ ${this.results.failed} platform test(s) failed`);
            process.exit(1);
        }

        // Save report to file
        this.saveReportToFile(duration);
    }

    saveReportToFile(duration) {
        const report = {
            timestamp: new Date().toISOString(),
            duration: duration,
            summary: {
                passed: this.results.passed,
                failed: this.results.failed,
                skipped: this.results.skipped
            },
            platforms: this.results.platforms
        };

        const reportsDir = path.join(__dirname, '../../test-reports');
        if (!fs.existsSync(reportsDir)) {
            fs.mkdirSync(reportsDir, { recursive: true });
        }

        const reportFile = path.join(reportsDir, `platform-integration-${Date.now()}.json`);
        fs.writeFileSync(reportFile, JSON.stringify(report, null, 2));
        
        console.log(`\n📄 Detailed report saved to: ${reportFile}`);
    }
}

// Run if called directly
if (require.main === module) {
    const runner = new PlatformTestRunner();
    runner.runAllPlatformTests().catch(console.error);
}

module.exports = PlatformTestRunner;
