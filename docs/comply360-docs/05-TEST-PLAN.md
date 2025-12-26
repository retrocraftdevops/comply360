# Comply360: Test Plan

**Document Version:** 1.0  
**Date:** December 26, 2025  
**QA Lead:** TBD  
**Status:** Ready for Implementation

---

## 1. Testing Strategy

### 1.1 Test Pyramid

```
           ┌─────────┐
          /  E2E (10%) \
         /───────────────\
        /  Integration (30%)\
       /─────────────────────\
      /    Unit Tests (60%)    \
     /───────────────────────────\
```

**Rationale**: Heavy emphasis on unit tests for speed and reliability, fewer E2E tests which are slower and more brittle.

### 1.2 Test Types

| Test Type | Coverage | Tools | Frequency |
|-----------|----------|-------|-----------|
| Unit Tests | 80%+ | Vitest, Go testing | Every commit |
| Integration Tests | Key flows | Testcontainers, Supertest | Every PR |
| E2E Tests | Critical paths | Playwright | Before deployment |
| Performance Tests | Load scenarios | k6, Artillery | Weekly |
| Security Tests | Vulnerabilities | OWASP ZAP, Snyk | Daily scan |
| UAT | User acceptance | Manual + Beta users | Pre-release |

---

## 2. Unit Testing

### 2.1 Frontend Unit Tests (Vitest + React Testing Library)

**Coverage Requirements**:
- **Minimum**: 80% overall coverage
- **Critical components**: 90%+ coverage (forms, auth, payments)

**Test Structure**:
```typescript
// Example: NameReservationForm.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import NameReservationForm from './NameReservationForm';

describe('NameReservationForm', () => {
  it('should search for name availability', async () => {
    const mockSearch = vi.fn().mockResolvedValue({
      available: true,
      alternatives: []
    });
    
    render(<NameReservationForm onSearch={mockSearch} />);
    
    const input = screen.getByLabelText('Company Name');
    fireEvent.change(input, { target: { value: 'Test Company' } });
    
    await waitFor(() => {
      expect(mockSearch).toHaveBeenCalledWith('Test Company');
    });
  });

  it('should show validation error for invalid name', async () => {
    render(<NameReservationForm />);
    
    const input = screen.getByLabelText('Company Name');
    fireEvent.change(input, { target: { value: 'Bank' } }); // Prohibited word
    
    await waitFor(() => {
      expect(screen.getByText(/prohibited word/i)).toBeInTheDocument();
    });
  });
  
  it('should disable submit button while loading', async () => {
    render(<NameReservationForm isSubmitting={true} />);
    
    const button = screen.getByRole('button', { name: /reserve/i });
    expect(button).toBeDisabled();
  });
});
```

**What to Test**:
- Component rendering
- User interactions (click, type, submit)
- Form validation
- Conditional rendering
- Error states
- Loading states
- Accessibility (ARIA labels, keyboard navigation)

**What NOT to Test**:
- Implementation details (internal state, methods)
- Third-party libraries (React, Zod)
- Trivial code (getters, simple renders)

---

### 2.2 Backend Unit Tests (Go testing)

**Coverage Requirements**:
- **Minimum**: 85% overall coverage
- **Business logic**: 95%+ coverage (commission calculations, validations)

**Test Structure**:
```go
// Example: commission_test.go
package service

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalculateCommission(t *testing.T) {
	tests := []struct {
		name           string
		registrationFee float64
		serviceFee     float64
		rate           float64
		expectedAmount float64
	}{
		{
			name:           "Standard commission calculation",
			registrationFee: 1000.00,
			serviceFee:     500.00,
			rate:           0.15,
			expectedAmount: 225.00, // (1000 + 500) * 0.15
		},
		{
			name:           "Zero service fee",
			registrationFee: 1000.00,
			serviceFee:     0.00,
			rate:           0.15,
			expectedAmount: 150.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateCommission(tt.registrationFee, tt.serviceFee, tt.rate)
			assert.Equal(t, tt.expectedAmount, result)
		})
	}
}

func TestCreateRegistration_ValidData(t *testing.T) {
	// Arrange
	mockDB := setupMockDB(t)
	service := NewRegistrationService(mockDB)
	
	input := &CreateRegistrationInput{
		Type:       "PTY_LTD",
		TenantID:   "tenant_123",
		ClientName: "John Doe",
		// ... other fields
	}

	// Act
	registration, err := service.CreateRegistration(context.Background(), input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, registration)
	assert.Equal(t, "DRAFT", registration.Status)
	assert.NotEmpty(t, registration.ID)
}

func TestCreateRegistration_InvalidTenant(t *testing.T) {
	mockDB := setupMockDB(t)
	service := NewRegistrationService(mockDB)
	
	input := &CreateRegistrationInput{
		TenantID: "invalid_tenant",
		// ...
	}

	registration, err := service.CreateRegistration(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, registration)
	assert.Contains(t, err.Error(), "tenant not found")
}
```

**What to Test**:
- Business logic (calculations, transformations)
- Validation rules
- Error handling
- Database operations (using mocks)
- API endpoints (request/response)
- Authorization checks

---

## 3. Integration Testing

### 3.1 API Integration Tests

**Tools**: Supertest (Node.js), httptest (Go)

**Test Structure**:
```typescript
// Example: registrations.integration.test.ts
import request from 'supertest';
import { app } from '../app';
import { setupTestDB, teardownTestDB } from './helpers';

describe('Registration API Integration', () => {
  beforeAll(async () => {
    await setupTestDB();
  });

  afterAll(async () => {
    await teardownTestDB();
  });

  describe('POST /api/v1/registrations', () => {
    it('should create a new registration', async () => {
      const response = await request(app)
        .post('/api/v1/registrations')
        .set('Authorization', `Bearer ${validToken}`)
        .send({
          type: 'PTY_LTD',
          clientName: 'Test Client',
          clientEmail: 'test@example.com',
          jurisdiction: 'ZA'
        });

      expect(response.status).toBe(201);
      expect(response.body.success).toBe(true);
      expect(response.body.data).toHaveProperty('id');
      expect(response.body.data.status).toBe('DRAFT');
    });

    it('should reject request without authentication', async () => {
      const response = await request(app)
        .post('/api/v1/registrations')
        .send({ type: 'PTY_LTD' });

      expect(response.status).toBe(401);
    });

    it('should reject request with invalid tenant', async () => {
      const response = await request(app)
        .post('/api/v1/registrations')
        .set('Authorization', `Bearer ${invalidTenantToken}`)
        .send({ type: 'PTY_LTD' });

      expect(response.status).toBe(403);
    });
  });
});
```

**Test Scenarios**:
- Full API request/response cycle
- Database persistence
- Multi-tenant isolation (tenant A cannot access tenant B data)
- Authentication and authorization
- Validation errors
- Rate limiting
- External API mocking (CIPC, payment gateways)

---

### 3.2 Database Integration Tests

**Tools**: Testcontainers (Docker-based PostgreSQL)

```go
func TestRegistration_DatabaseIntegration(t *testing.T) {
	// Start PostgreSQL container
	postgres, err := testcontainers.NewContainer(/* ... */)
	require.NoError(t, err)
	defer postgres.Terminate(context.Background())

	// Run migrations
	db := setupTestDatabase(postgres)

	// Test
	registration := &Registration{
		TenantID:   "tenant_123",
		Type:       "PTY_LTD",
		Status:     "DRAFT",
	}

	err = db.Create(registration).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, registration.ID)

	// Verify RLS works (tenant isolation)
	db = SetTenantContext(db, "different_tenant")
	var found Registration
	err = db.First(&found, registration.ID).Error
	assert.Error(t, err) // Should not find - different tenant
}
```

---

## 4. End-to-End Testing

### 4.1 E2E Test Framework (Playwright)

**Critical User Journeys to Test**:

1. **Complete Registration Flow**
   - Agent logs in
   - Switches to South Africa jurisdiction
   - Searches for company name
   - Reserves name
   - Completes Pty Ltd registration wizard (all 7 steps)
   - Uploads documents
   - Processes payment
   - Receives confirmation
   - Downloads certificate (when approved)

2. **Multi-Tenant Isolation**
   - Create two tenant accounts
   - Verify Tenant A cannot see Tenant B data
   - Verify subdomains work correctly

3. **Commission Tracking**
   - Complete registration
   - Verify commission calculated correctly
   - Check commission dashboard updates

**Example E2E Test**:
```typescript
// tests/e2e/registration-flow.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Complete Registration Flow', () => {
  test('should allow agent to complete Pty Ltd registration', async ({ page }) => {
    // Login
    await page.goto('https://demo.comply360.com');
    await page.fill('[name=email]', 'agent@demo.com');
    await page.fill('[name=password]', 'password123');
    await page.click('button:text("Sign In")');
    await expect(page).toHaveURL(/\/dashboard/);

    // Start registration
    await page.click('button:text("New Registration")');
    await page.click('text=Private Company (Pty Ltd)');

    // Step 1: Name Search
    await page.fill('[name=companyName]', 'Test Trading (Pty) Ltd');
    await page.click('button:text("Check Availability")');
    await expect(page.locator('.availability-status')).toContainText('Available');
    await page.click('button:text("Continue")');

    // Step 2: Company Details
    await page.fill('[name=businessAddress]', '123 Main Road, Johannesburg');
    await page.selectOption('[name=businessActivity]', 'Retail Trade');
    await page.click('button:text("Continue")');

    // Step 3: Directors
    await page.click('button:text("Add Director")');
    await page.fill('[name="directors[0].fullName"]', 'John Doe');
    await page.fill('[name="directors[0].idNumber"]', '8001015009087');
    await page.fill('[name="directors[0].email"]', 'john@example.com');
    await page.click('button:text("Continue")');

    // ... (continue for all steps)

    // Final step: Review and Submit
    await expect(page.locator('.registration-summary')).toBeVisible();
    await page.click('button:text("Submit Registration")');

    // Payment
    await page.fill('[name=cardNumber]', '4242424242424242'); // Stripe test card
    await page.fill('[name=expiry]', '12/25');
    await page.fill('[name=cvc]', '123');
    await page.click('button:text("Pay Now")');

    // Confirmation
    await expect(page.locator('.success-message')).toContainText('Registration Submitted');
    await expect(page.locator('.registration-number')).toBeVisible();
  });
});
```

**E2E Test Configuration**:
```typescript
// playwright.config.ts
export default defineConfig({
  testDir: './tests/e2e',
  timeout: 60000, // 60 seconds per test
  retries: 2, // Retry failed tests
  workers: 4, // Parallel execution
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'mobile',
      use: { ...devices['iPhone 13'] },
    },
  ],
});
```

---

## 5. Performance Testing

### 5.1 Load Testing (k6)

**Test Scenarios**:

1. **Normal Load**: 100 concurrent users, 10-minute duration
2. **Peak Load**: 500 concurrent users, 5-minute duration
3. **Spike Test**: 0 → 1000 users in 1 minute, then back to 0
4. **Endurance Test**: 200 concurrent users, 2-hour duration

**Example k6 Script**:
```javascript
// load-test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 100 }, // Ramp up
    { duration: '5m', target: 100 }, // Steady state
    { duration: '2m', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% of requests < 2s
    http_req_failed: ['rate<0.01'],    // <1% failure rate
  },
};

export default function() {
  // Login
  let loginRes = http.post('https://api.comply360.com/v1/auth/login', JSON.stringify({
    email: 'loadtest@example.com',
    password: 'password123',
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  check(loginRes, {
    'login successful': (r) => r.status === 200,
  });

  let token = loginRes.json('accessToken');

  // Get registrations
  let regsRes = http.get('https://api.comply360.com/v1/registrations', {
    headers: { Authorization: `Bearer ${token}` },
  });

  check(regsRes, {
    'registrations fetched': (r) => r.status === 200,
    'response time OK': (r) => r.timings.duration < 2000,
  });

  sleep(1); // Think time
}
```

**Performance Targets**:
- **Response Time**: p95 < 2 seconds, p99 < 5 seconds
- **Throughput**: 100 requests/second sustained
- **Error Rate**: < 1%
- **Concurrent Users**: 10,000+

---

## 6. Security Testing

### 6.1 Automated Security Scans

**Tools**:
- **OWASP ZAP**: Dynamic application security testing (DAST)
- **Snyk**: Dependency vulnerability scanning
- **Trivy**: Container image scanning
- **SonarQube**: Static code analysis (SAST)

**Scan Frequency**:
- **Daily**: Dependency scans (Snyk)
- **Weekly**: Full DAST scan (OWASP ZAP)
- **On PR**: Static analysis (SonarQube)
- **On Deploy**: Container scan (Trivy)

**Critical Vulnerabilities**:
- SQL Injection
- XSS (Cross-Site Scripting)
- CSRF (Cross-Site Request Forgery)
- Authentication bypass
- Authorization flaws
- Sensitive data exposure
- Insufficient logging

---

### 6.2 Manual Penetration Testing

**Frequency**: Quarterly

**Test Scenarios**:
- Multi-tenant data leakage attempts
- Privilege escalation
- API fuzzing
- Session hijacking
- File upload vulnerabilities
- Business logic flaws

---

## 7. User Acceptance Testing (UAT)

### 7.1 UAT Process

**Participants**:
- 10 beta agent partners
- 5 internal stakeholders
- Product owner

**Duration**: 2 weeks before production launch

**Test Scenarios**:
1. Complete onboarding process
2. Register 5 different company types
3. Process payments end-to-end
4. Generate and download reports
5. Test all notification channels
6. Mobile responsiveness testing
7. Usability feedback

**Exit Criteria**:
- NPS > 40
- CSAT > 4.0/5.0
- <5 critical bugs
- 90%+ task completion rate
- All P0 feedback addressed

---

## 8. Test Data Management

### 8.1 Test Environments

| Environment | Purpose | Data | Refresh |
|-------------|---------|------|---------|
| Dev | Active development | Synthetic | On demand |
| Staging | Pre-production | Anonymized prod copy | Weekly |
| UAT | User acceptance | Synthetic + beta | Manual |
| Production | Live | Real | N/A |

### 8.2 Test Data Strategy

**Synthetic Data Generation**:
- Use Faker.js for frontend tests
- Use factory patterns for backend tests
- Maintain test fixtures in Git

**Data Anonymization** (for staging):
- Hash emails: `user@example.com` → `u9f2k@example.com`
- Mask ID numbers: `8001015009087` → `8001****9087`
- Remove sensitive documents
- Retain data relationships

---

## 9. CI/CD Test Integration

### 9.1 GitHub Actions Workflow

```yaml
name: Test Pipeline

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Frontend Unit Tests
        run: |
          npm install
          npm run test:unit
      - name: Run Backend Unit Tests
        run: |
          cd services/api
          go test ./... -cover -coverprofile=coverage.out

  integration-tests:
    runs-on: ubuntu-latest
    needs: unit-tests
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
      redis:
        image: redis:7
    steps:
      - uses: actions/checkout@v3
      - name: Run Integration Tests
        run: npm run test:integration

  e2e-tests:
    runs-on: ubuntu-latest
    needs: integration-tests
    steps:
      - uses: actions/checkout@v3
      - name: Run E2E Tests
        run: npx playwright test

  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Snyk Scan
        run: npx snyk test --severity-threshold=high
```

---

## 10. Test Metrics & Reporting

### 10.1 KPIs

| Metric | Target | Current |
|--------|--------|---------|
| Unit Test Coverage | 80%+ | TBD |
| Integration Test Pass Rate | 100% | TBD |
| E2E Test Pass Rate | 95%+ | TBD |
| Build Success Rate | 90%+ | TBD |
| Mean Time to Detect (MTTD) | <10 min | TBD |
| Mean Time to Resolve (MTTR) | <4 hours | TBD |

### 10.2 Test Reporting

**Tools**: Jest HTML Reporter, Allure, Playwright HTML Report

**Distribution**:
- Test results published to GitHub PR comments
- Coverage reports uploaded to CodeCov
- Daily test summary emailed to team
- Dashboard (Grafana) with test metrics

---

**Test Plan Approval:**
- QA Lead: __________
- Development Lead: __________
- Product Owner: __________
- Date: __________

