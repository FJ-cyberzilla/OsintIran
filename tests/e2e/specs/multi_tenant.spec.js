// tests/e2e/specs/multi_tenant.spec.js
describe('Multi-tenant Administration', () => {
  let testTenant;
  let testUser;

  before(() => {
    // Create test tenant and user
    cy.createTestTenant().then(tenant => {
      testTenant = tenant;
      return cy.createTestUser(tenant.id);
    }).then(user => {
      testUser = user;
    });
  });

  after(() => {
    // Cleanup test data
    cy.cleanupTestData();
  });

  it('should isolate tenant data properly', () => {
    // Login as tenant admin
    cy.loginAsTenantAdmin(testTenant.slug, testUser.email);
    
    // Create some test data
    cy.createTestPhoneLookup('+989111111111');
    cy.createTestPhoneLookup('+989222222222');
    
    // Verify data is isolated to tenant
    cy.get('[data-testid="lookup-history"]')
      .should('have.length', 2)
      .and('contain', '+989111111111')
      .and('contain', '+989222222222');

    // Login as different tenant admin
    cy.loginAsDifferentTenant();
    
    // Verify no data leakage
    cy.get('[data-testid="lookup-history"]')
      .should('not.contain', '+989111111111')
      .and('not.contain', '+989222222222');
  });

  it('should enforce tenant quotas', () => {
    cy.loginAsTenantAdmin(testTenant.slug, testUser.email);
    
    // Exhaust quota with multiple requests
    const requests = Array.from({ length: 105 }, (_, i) => 
      `+9891234567${i.toString().padStart(2, '0')}`
    );

    requests.forEach(phone => {
      cy.get('[data-testid="phone-input"]').clear().type(phone);
      cy.get('[data-testid="search-button"]').click();
      
      // Check for quota exceeded message after 100 requests
      if (requests.indexOf(phone) >= 100) {
        cy.get('[data-testid="quota-exceeded"]', { timeout: 5000 })
          .should('be.visible');
        return false; // Break loop
      }
      
      // Small delay between requests
      cy.wait(100);
    });
  });

  it('should manage users within tenant', () => {
    cy.loginAsTenantAdmin(testTenant.slug, testUser.email);
    cy.visit('/admin/users');
    
    // Add new user
    cy.get('[data-testid="add-user-btn"]').click();
    cy.get('[data-testid="user-email"]').type('newuser@test.com');
    cy.get('[data-testid="user-role"]').select('viewer');
    cy.get('[data-testid="save-user-btn"]').click();
    
    // Verify user was added
    cy.get('[data-testid="users-table"]')
      .should('contain', 'newuser@test.com')
      .and('contain', 'viewer');
    
    // Test user permissions
    cy.loginAsUser('newuser@test.com', 'temp-password');
    
    // Viewer should not see admin sections
    cy.get('[data-testid="admin-tab"]').should('not.exist');
    
    // Viewer should be able to view reports
    cy.visit('/reports');
    cy.get('[data-testid="reports-list"]').should('be.visible');
  });
});
