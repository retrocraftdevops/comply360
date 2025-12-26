# Authentication and Authorization System - Implementation Tasks

**Spec:** Authentication and Authorization System
**Created:** December 27, 2025
**Total Estimated Time:** 2 weeks (L)
**Phase:** 1 - Foundation

---

## Priority 1: Database Schema and Models (2 days)

### 1.1 Create Database Schema (1 day)
- [ ] Create `users` table with all required fields
- [ ] Create `user_roles` table for role assignments
- [ ] Create `oauth_accounts` table for OAuth integrations
- [ ] Create `password_reset_tokens` table
- [ ] Create `email_verification_tokens` table
- [ ] Add all necessary indexes (email, status, tenant_id)
- [ ] Create Prisma schema definitions
- [ ] Create database migrations
- [ ] Test schema in local PostgreSQL

### 1.2 Create Domain Models (1 day)
- [ ] Define User model (Go struct)
- [ ] Define UserRole model
- [ ] Define OAuthAccount model
- [ ] Define Session model
- [ ] Define JWT claims structs
- [ ] Create model validation functions
- [ ] Write unit tests for models
- [ ] Document model relationships

---

## Priority 2: User Registration and Email Verification (3 days)

### 2.1 Registration API (2 days)
- [ ] Create `POST /api/auth/register` endpoint
- [ ] Implement email/password validation
- [ ] Check for duplicate email addresses
- [ ] Hash passwords using bcrypt (cost factor 12)
- [ ] Create user record in database
- [ ] Assign default role (agent or client)
- [ ] Generate email verification token
- [ ] Send verification email (SendGrid)
- [ ] Write API tests
- [ ] Add rate limiting (5 registrations/hour per IP)

### 2.2 Email Verification (1 day)
- [ ] Create `POST /api/auth/verify-email` endpoint
- [ ] Validate verification token
- [ ] Check token expiration (24 hours)
- [ ] Mark email as verified
- [ ] Update user record
- [ ] Send welcome email
- [ ] Create `POST /api/auth/resend-verification` endpoint
- [ ] Write integration tests

---

## Priority 3: Login and JWT Token Management (3 days)

### 3.1 Login API (2 days)
- [ ] Create `POST /api/auth/login` endpoint
- [ ] Validate email and password
- [ ] Fetch user by email and tenant
- [ ] Verify password hash using bcrypt
- [ ] Check account status (active/suspended/locked)
- [ ] Check failed login attempts
- [ ] Implement account lockout (5 attempts, 30 min)
- [ ] Generate JWT access token (15 min expiry)
- [ ] Generate JWT refresh token (7 day expiry)
- [ ] Create session in Redis
- [ ] Return tokens and user info
- [ ] Write comprehensive tests

### 3.2 JWT Token Generation (1 day)
- [ ] Generate RSA key pair for signing
- [ ] Implement access token generation function
- [ ] Implement refresh token generation function
- [ ] Add claims: user ID, tenant ID, roles, permissions
- [ ] Sign tokens using RS256 algorithm
- [ ] Store refresh token in session
- [ ] Write token generation tests
- [ ] Test token expiration handling

---

## Priority 4: Token Validation and Middleware (2 days)

### 4.1 JWT Validation Middleware (1 day)
- [ ] Create JWT validation middleware
- [ ] Extract token from Authorization header
- [ ] Validate token signature
- [ ] Check token expiration
- [ ] Extract claims (user ID, tenant ID, roles)
- [ ] Inject user context into request
- [ ] Handle invalid/expired tokens
- [ ] Write middleware tests

### 4.2 Token Refresh (1 day)
- [ ] Create `POST /api/auth/refresh` endpoint
- [ ] Validate refresh token
- [ ] Check session validity
- [ ] Generate new access token
- [ ] Return new token
- [ ] Update session last activity
- [ ] Write tests for token refresh flow

---

## Priority 5: Role-Based Access Control (RBAC) (3 days)

### 5.1 Casbin Setup (1 day)
- [ ] Install Casbin library
- [ ] Create RBAC model configuration
- [ ] Create policy rules CSV file
- [ ] Initialize Casbin enforcer
- [ ] Load policies from database
- [ ] Test policy enforcement
- [ ] Document RBAC model

### 5.2 Permission System (1 day)
- [ ] Define permission constants
- [ ] Create permission-to-role mappings
- [ ] Implement permission checking function
- [ ] Store policies in database
- [ ] Create policy management functions
- [ ] Write permission tests
- [ ] Document permission model

### 5.3 Authorization Middleware (1 day)
- [ ] Create authorization middleware
- [ ] Integrate with Casbin enforcer
- [ ] Check user permissions for resources
- [ ] Handle forbidden access (403)
- [ ] Log authorization failures
- [ ] Write middleware tests
- [ ] Add to protected routes

---

## Priority 6: Multi-Factor Authentication (MFA) (3 days)

### 6.1 TOTP Implementation (2 days)
- [ ] Install TOTP library (google-authenticator)
- [ ] Create `POST /api/auth/mfa/enable` endpoint
- [ ] Generate TOTP secret
- [ ] Generate QR code for Google Authenticator
- [ ] Store secret in user record
- [ ] Create `POST /api/auth/mfa/verify` endpoint
- [ ] Validate TOTP codes (30-second window)
- [ ] Generate backup codes
- [ ] Create `POST /api/auth/mfa/disable` endpoint
- [ ] Write MFA tests

### 6.2 SMS/Email MFA (1 day)
- [ ] Implement SMS code generation
- [ ] Send SMS codes via Twilio
- [ ] Implement email code generation
- [ ] Send email codes via SendGrid
- [ ] Create code verification endpoints
- [ ] Set code expiration (10 minutes)
- [ ] Add rate limiting (3 attempts)
- [ ] Write tests

---

## Priority 7: OAuth Integration (2 days)

### 7.1 OAuth Setup (1 day)
- [ ] Register OAuth apps (Google, Microsoft, GitHub)
- [ ] Configure OAuth redirect URIs
- [ ] Create `GET /api/auth/oauth/:provider` endpoint
- [ ] Generate OAuth state parameter
- [ ] Redirect to OAuth provider
- [ ] Create `GET /api/auth/oauth/:provider/callback` endpoint
- [ ] Exchange code for access token
- [ ] Fetch user info from provider

### 7.2 OAuth Account Linking (1 day)
- [ ] Find or create user from OAuth data
- [ ] Link OAuth account to user record
- [ ] Store OAuth tokens in database
- [ ] Generate JWT tokens
- [ ] Handle OAuth errors
- [ ] Write OAuth integration tests
- [ ] Document OAuth flow

---

## Priority 8: Password Management (2 days)

### 8.1 Password Reset (1 day)
- [ ] Create `POST /api/auth/password/forgot` endpoint
- [ ] Generate password reset token
- [ ] Store token in database (1 hour expiry)
- [ ] Send password reset email
- [ ] Create `POST /api/auth/password/reset` endpoint
- [ ] Validate reset token
- [ ] Update password
- [ ] Invalidate token
- [ ] Write tests

### 8.2 Password Change (1 day)
- [ ] Create `PUT /api/auth/password/change` endpoint
- [ ] Require authentication
- [ ] Validate current password
- [ ] Validate new password
- [ ] Update password hash
- [ ] Revoke all sessions
- [ ] Send notification email
- [ ] Write tests

---

## Priority 9: Session Management (1 day)

### 9.1 Session Operations (1 day)
- [ ] Implement session creation in Redis
- [ ] Store session metadata (IP, user agent)
- [ ] Create `GET /api/auth/sessions` endpoint
- [ ] List all user sessions
- [ ] Create `DELETE /api/auth/sessions/:id` endpoint
- [ ] Revoke specific session
- [ ] Create `DELETE /api/auth/sessions` endpoint
- [ ] Revoke all sessions
- [ ] Track session activity
- [ ] Write session tests

---

## Priority 10: User Management (2 days)

### 10.1 User CRUD Operations (1 day)
- [ ] Create `GET /api/users` endpoint (list users)
- [ ] Add pagination and filtering
- [ ] Create `GET /api/users/:id` endpoint
- [ ] Create `PUT /api/users/:id` endpoint
- [ ] Create `DELETE /api/users/:id` endpoint (soft delete)
- [ ] Add authorization checks (tenant_admin only)
- [ ] Write API tests

### 10.2 User Status Management (1 day)
- [ ] Create `POST /api/users/:id/suspend` endpoint
- [ ] Create `POST /api/users/:id/activate` endpoint
- [ ] Send notification emails
- [ ] Revoke sessions on suspension
- [ ] Create `GET /api/users/me` endpoint
- [ ] Create `PUT /api/users/me` endpoint
- [ ] Write tests

---

## Priority 11: Role Management (1 day)

### 11.1 Role Assignment (1 day)
- [ ] Create `GET /api/users/:id/roles` endpoint
- [ ] Create `POST /api/users/:id/roles` endpoint
- [ ] Validate role exists
- [ ] Assign role to user
- [ ] Create `DELETE /api/users/:id/roles/:role` endpoint
- [ ] Remove role from user
- [ ] Update Casbin policies
- [ ] Write tests

---

## Priority 12: Testing and Security (2 days)

### 12.1 Security Testing (1 day)
- [ ] Test brute force protection
- [ ] Test SQL injection prevention
- [ ] Test XSS prevention
- [ ] Test CSRF protection
- [ ] Test token manipulation
- [ ] Test session hijacking prevention
- [ ] Document security measures

### 12.2 Integration Testing (1 day)
- [ ] Test complete registration flow
- [ ] Test login flow with MFA
- [ ] Test OAuth flow
- [ ] Test password reset flow
- [ ] Test role assignment flow
- [ ] Test session management
- [ ] Achieve 80%+ code coverage

---

## Priority 13: Documentation (1 day)

### 13.1 API Documentation (0.5 day)
- [ ] Generate OpenAPI/Swagger documentation
- [ ] Document all endpoints
- [ ] Add request/response examples
- [ ] Document error codes
- [ ] Publish API docs

### 13.2 Developer Documentation (0.5 day)
- [ ] Document authentication flow
- [ ] Document authorization model
- [ ] Create integration guide
- [ ] Document environment variables
- [ ] Create troubleshooting guide

---

## Definition of Done

- [ ] All code passes AI validation (80%+ score)
- [ ] Unit tests written and passing (80%+ coverage)
- [ ] Integration tests written and passing
- [ ] Security tests passed
- [ ] Code reviewed and approved
- [ ] API documentation published
- [ ] Developer documentation complete
- [ ] Deployed to staging environment
- [ ] Performance benchmarks met
- [ ] Product owner acceptance

---

## Technical Notes

### Password Hashing
- Use bcrypt with cost factor 12
- Salt automatically handled by bcrypt
- Never log or store plain-text passwords

### JWT Signing
- Use RS256 algorithm (asymmetric)
- Store private key securely (environment variable or secrets manager)
- Rotate keys every 90 days
- Public key for validation can be cached

### Session Storage
- Store in Redis for fast access
- Set TTL to match refresh token expiry
- Use session ID as Redis key
- Invalidate on logout

### Account Lockout
- Lock after 5 failed attempts
- Lock duration: 30 minutes
- Reset counter on successful login
- Send email notification on lockout

### Rate Limiting
- Registration: 5/hour per IP
- Login: 10/minute per IP
- MFA verification: 3/minute per user
- Password reset: 3/hour per email

### Error Handling
- Never reveal user existence in login errors
- Use generic "Invalid credentials" message
- Log security events for monitoring
- Return 401 for authentication failures
- Return 403 for authorization failures

### Performance Targets
- Login: < 500ms
- Token validation: < 50ms
- Authorization check: < 10ms
- Session lookup: < 20ms

---

**Next Steps:** Begin with Priority 1 (Database Schema and Models)
