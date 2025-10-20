// tests/e2e/pages/login.page.js
class LoginPage {
  navigate() {
    cy.visit('/login');
    return this;
  }

  fillCredentials(email, password) {
    cy.get('[data-testid="email-input"]').clear().type(email);
    cy.get('[data-testid="password-input"]').clear().type(password);
    return this;
  }

  submit() {
    cy.get('[data-testid="login-btn"]').click();
    return this;
  }

  loginAsAdmin() {
    return this.navigate()
      .fillCredentials(Cypress.env('admin_email'), Cypress.env('admin_password'))
      .submit();
  }

  loginAsTenant(tenantSlug, email, password) {
    return this.navigate()
      .fillCredentials(`${email}@${tenantSlug}`, password)
      .submit();
  }

  expectSuccess() {
    cy.url().should('include', '/dashboard');
    cy.get('[data-testid="welcome-message"]').should('be.visible');
    return this;
  }

  expectError() {
    cy.get('[data-testid="error-message"]').should('be.visible');
    return this;
  }
}

export default new LoginPage();
