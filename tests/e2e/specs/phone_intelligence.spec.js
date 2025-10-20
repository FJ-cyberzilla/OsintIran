// tests/e2e/specs/phone_intelligence.spec.js
describe('Phone Intelligence Workflow', () => {
  beforeEach(() => {
    cy.loginAsAdmin();
    cy.visit('/dashboard');
  });

  it('should complete full phone intelligence workflow', () => {
    // 1. Navigate to phone lookup
    cy.get('[data-testid="phone-lookup-tab"]').click();
    
    // 2. Enter phone number and search
    const testPhone = '+989123456789';
    cy.get('[data-testid="phone-input"]').type(testPhone);
    cy.get('[data-testid="search-button"]').click();

    // 3. Verify normalization
    cy.get('[data-testid="normalized-phone"]')
      .should('contain', '+989123456789')
      .and('contain', 'Iran');

    // 4. Verify carrier detection
    cy.get('[data-testid="carrier-info"]')
      .should('contain', 'MCI');

    // 5. Start email discovery
    cy.get('[data-testid="discover-emails-btn"]').click();
    
    // 6. Wait for email discovery to complete
    cy.get('[data-testid="email-discovery-progress"]', { timeout: 30000 })
      .should('not.exist');
    
    // 7. Verify email results
    cy.get('[data-testid="email-results"]')
      .should('have.length.at.least', 1)
      .first()
      .within(() => {
        cy.get('[data-testid="email-address"]').should('contain', '@');
        cy.get('[data-testid="confidence-score"]').should('contain', '%');
      });

    // 8. Generate intelligence report
    cy.get('[data-testid="generate-report-btn"]').click();
    
    // 9. Verify report generation
    cy.get('[data-testid="report-status"]', { timeout: 45000 })
      .should('contain', 'Completed');
    
    // 10. Download report
    cy.get('[data-testid="download-report-btn"]').click();
    
    // 11. Verify file download
    cy.verifyDownload('intelligence-report.pdf', { timeout: 10000 });

    // 12. Check report content in preview
    cy.get('[data-testid="report-preview"]')
      .should('contain', 'Executive Summary')
      .and('contain', 'Risk Assessment')
      .and('contain', 'Social Graph');
  });

  it('should handle invalid phone numbers gracefully', () => {
    cy.get('[data-testid="phone-lookup-tab"]').click();
    
    // Test invalid number
    cy.get('[data-testid="phone-input"]').type('invalid-phone');
    cy.get('[data-testid="search-button"]').click();
    
    cy.get('[data-testid="error-message"]')
      .should('be.visible')
      .and('contain', 'valid phone number');
  });

  it('should respect rate limiting', () => {
    cy.get('[data-testid="phone-lookup-tab"]').click();
    
    // Make rapid consecutive requests
    for (let i = 0; i < 15; i++) {
      cy.get('[data-testid="phone-input"]').clear().type(`+9891234567${i}`);
      cy.get('[data-testid="search-button"]').click();
      
      if (i >= 10) {
        // Should hit rate limit around 10th request
        cy.get('[data-testid="rate-limit-message"]').should('be.visible');
        break;
      }
    }
  });
});
