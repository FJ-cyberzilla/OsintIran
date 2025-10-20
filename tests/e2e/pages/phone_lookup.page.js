// tests/e2e/pages/phone_lookup.page.js
class PhoneLookupPage {
  navigate() {
    cy.get('[data-testid="phone-lookup-tab"]').click();
    return this;
  }

  enterPhoneNumber(phoneNumber) {
    cy.get('[data-testid="phone-input"]').clear().type(phoneNumber);
    return this;
  }

  search() {
    cy.get('[data-testid="search-button"]').click();
    return this;
  }

  expectNormalization(expectedNormalized) {
    cy.get('[data-testid="normalized-phone"]')
      .should('contain', expectedNormalized);
    return this;
  }

  expectCarrier(expectedCarrier) {
    cy.get('[data-testid="carrier-info"]')
      .should('contain', expectedCarrier);
    return this;
  }

  discoverEmails() {
    cy.get('[data-testid="discover-emails-btn"]').click();
    return this;
  }

  waitForEmailDiscovery() {
    cy.get('[data-testid="email-discovery-progress"]', { timeout: 30000 })
      .should('not.exist');
    return this;
  }

  expectEmailResults(minCount = 1) {
    cy.get('[data-testid="email-results"]')
      .should('have.length.at.least', minCount);
    return this;
  }

  generateReport() {
    cy.get('[data-testid="generate-report-btn"]').click();
    return this;
  }

  waitForReport() {
    cy.get('[data-testid="report-status"]', { timeout: 45000 })
      .should('contain', 'Completed');
    return this;
  }
}

export default new PhoneLookupPage();
