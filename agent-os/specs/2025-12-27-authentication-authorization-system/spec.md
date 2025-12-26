# Authentication and Authorization System - Specification

**Version:** 1.0.0
**Date:** December 27, 2025
**Author:** Comply360 Development Team
**Status:** Planning
**Phase:** 1 - Foundation

---

## Executive Summary

This specification defines the authentication and authorization system for Comply360, providing secure multi-tenant user management, role-based access control (RBAC), session management, and OAuth integration. This system is foundational for all platform security.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                  AUTHENTICATION FLOW                         │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  1. User Login                                               │
│     │                                                         │
│     ▼                                                         │
│  2. Credential Validation                                    │
│     │                                                         │
│     ▼                                                         │
│  3. Multi-Factor Authentication (if enabled)                 │
│     │                                                         │
│     ▼                                                         │
│  4. JWT Token Generation                                     │
│     │                                                         │
│     ▼                                                         │
│  5. Session Creation                                         │
│     │                                                         │
│     ▼                                                         │
│  6. Return Tokens (Access + Refresh)                         │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                  AUTHORIZATION FLOW                          │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  1. Request with JWT Token                                   │
│     │                                                         │
│     ▼                                                         │
│  2. Token Validation                                         │
│     │                                                         │
│     ▼                                                         │
│  3. Extract User & Tenant Context                            │
│     │                                                         │
│     ▼                                                         │
│  4. Check Permissions (RBAC via Casbin)                      │
│     │                                                         │
│     ▼                                                         │
│  5. Grant/Deny Access                                        │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

---

## Core Components

### 1. User Authentication

#### 1.1 Email/Password Authentication

**Flow:**
1. User submits email and password
2. Hash password using bcrypt
3. Compare with stored hash
4. Validate password strength requirements
5. Check account status (active/suspended/locked)
6. Generate JWT tokens on success

**Password Requirements:**
- Minimum 8 characters
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 number
- At least 1 special character
- Not in common password list

**Implementation (Go):**
```go
type AuthService struct {
    db    *sql.DB
    redis *redis.Client
}

func (s *AuthService) Login(email, password, tenantID string) (*AuthResponse, error) {
    // 1. Fetch user by email and tenant
    user, err := s.db.GetUserByEmail(tenantID, email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    // 2. Check account status
    if user.Status != "active" {
        return nil, ErrAccountSuspended
    }

    // 3. Verify password
    if !bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) {
        // Increment failed login attempts
        s.incrementFailedLogins(user.ID)
        return nil, ErrInvalidCredentials
    }

    // 4. Check if MFA is required
    if user.MFAEnabled {
        // Return challenge requiring MFA verification
        return s.initiateMFAChallenge(user.ID)
    }

    // 5. Generate JWT tokens
    accessToken, refreshToken, err := s.generateTokens(user)
    if err != nil {
        return nil, err
    }

    // 6. Create session
    session := &Session{
        UserID:       user.ID,
        TenantID:     user.TenantID,
        RefreshToken: refreshToken,
        ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
    }
    s.createSession(session)

    // 7. Reset failed login attempts
    s.resetFailedLogins(user.ID)

    return &AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        User:         user,
    }, nil
}
```

#### 1.2 Multi-Factor Authentication (MFA)

**Supported Methods:**
- TOTP (Time-Based One-Time Password) via Google Authenticator, Authy
- SMS verification via Twilio
- Email verification codes

**MFA Flow:**
```go
func (s *AuthService) VerifyMFA(userID, code string) (*AuthResponse, error) {
    user, err := s.db.GetUserByID(userID)
    if err != nil {
        return nil, err
    }

    valid := false
    switch user.MFAMethod {
    case "totp":
        valid = totp.Validate(code, user.MFASecret)
    case "sms":
        valid = s.verifySMSCode(userID, code)
    case "email":
        valid = s.verifyEmailCode(userID, code)
    }

    if !valid {
        return nil, ErrInvalidMFACode
    }

    // Generate tokens after successful MFA
    return s.generateTokens(user)
}
```

#### 1.3 OAuth Integration

**Supported Providers:**
- Google OAuth
- Microsoft Azure AD
- GitHub (for developers)

**OAuth Flow:**
```go
func (s *AuthService) HandleOAuthCallback(provider, code, state string) (*AuthResponse, error) {
    // 1. Validate state parameter
    if !s.validateState(state) {
        return nil, ErrInvalidState
    }

    // 2. Exchange code for token
    oauthToken, err := s.exchangeCodeForToken(provider, code)
    if err != nil {
        return nil, err
    }

    // 3. Fetch user info from provider
    userInfo, err := s.fetchOAuthUserInfo(provider, oauthToken)
    if err != nil {
        return nil, err
    }

    // 4. Find or create user
    user, created := s.findOrCreateOAuthUser(userInfo)

    // 5. Link OAuth account
    s.linkOAuthAccount(user.ID, provider, userInfo.ProviderID)

    // 6. Generate JWT tokens
    return s.generateTokens(user)
}
```

### 2. JWT Token Management

#### 2.1 Token Structure

**Access Token (Short-lived: 15 minutes):**
```json
{
  "sub": "user_uuid",
  "tenant_id": "tenant_uuid",
  "email": "user@example.com",
  "roles": ["agent", "admin"],
  "permissions": ["registrations:read", "registrations:write"],
  "exp": 1735300800,
  "iat": 1735300000,
  "iss": "comply360.com",
  "aud": "comply360-api"
}
```

**Refresh Token (Long-lived: 7 days):**
```json
{
  "sub": "user_uuid",
  "tenant_id": "tenant_uuid",
  "session_id": "session_uuid",
  "exp": 1735905600,
  "iat": 1735300800,
  "iss": "comply360.com",
  "type": "refresh"
}
```

#### 2.2 Token Generation

```go
func (s *AuthService) generateTokens(user *User) (string, string, error) {
    // Access Token
    accessClaims := jwt.MapClaims{
        "sub":        user.ID,
        "tenant_id":  user.TenantID,
        "email":      user.Email,
        "roles":      user.Roles,
        "permissions": s.getUserPermissions(user),
        "exp":        time.Now().Add(15 * time.Minute).Unix(),
        "iat":        time.Now().Unix(),
        "iss":        "comply360.com",
        "aud":        "comply360-api",
    }

    accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
    accessTokenString, err := accessToken.SignedString(s.privateKey)
    if err != nil {
        return "", "", err
    }

    // Refresh Token
    sessionID := uuid.New().String()
    refreshClaims := jwt.MapClaims{
        "sub":        user.ID,
        "tenant_id":  user.TenantID,
        "session_id": sessionID,
        "exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
        "iat":        time.Now().Unix(),
        "iss":        "comply360.com",
        "type":       "refresh",
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString(s.privateKey)
    if err != nil {
        return "", "", err
    }

    return accessTokenString, refreshTokenString, nil
}
```

#### 2.3 Token Refresh

```go
func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
    // 1. Validate refresh token
    claims, err := s.validateJWT(refreshToken)
    if err != nil {
        return nil, ErrInvalidToken
    }

    // 2. Check if token type is refresh
    if claims["type"] != "refresh" {
        return nil, ErrInvalidTokenType
    }

    // 3. Verify session exists and is valid
    sessionID := claims["session_id"].(string)
    session, err := s.getSession(sessionID)
    if err != nil || session.Revoked {
        return nil, ErrSessionExpired
    }

    // 4. Fetch user
    userID := claims["sub"].(string)
    user, err := s.db.GetUserByID(userID)
    if err != nil {
        return nil, err
    }

    // 5. Generate new tokens
    return s.generateTokens(user)
}
```

### 3. Role-Based Access Control (RBAC)

#### 3.1 Role Hierarchy

**Global Roles:**
- `super_admin` - Platform administrator (access to all tenants)
- `global_support` - Support staff (read-only access to all tenants)

**Tenant Roles:**
- `tenant_admin` - Tenant administrator (full access within tenant)
- `tenant_manager` - Tenant manager (limited admin access)
- `agent` - Corporate service agent (core features)
- `agent_assistant` - Agent assistant (limited features)
- `client` - End client (self-service portal)

#### 3.2 Permission Model

**Resource-Action Permissions:**
```
registrations:read
registrations:write
registrations:delete
clients:read
clients:write
documents:read
documents:write
commissions:read
settings:read
settings:write
users:read
users:write
reports:read
```

#### 3.3 Casbin Implementation

**Model Configuration (RBAC with domains):**
```conf
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
```

**Policy Rules:**
```csv
p, agent, tenant_abc, registrations, read
p, agent, tenant_abc, registrations, write
p, agent, tenant_abc, clients, read
p, agent, tenant_abc, clients, write
p, tenant_admin, tenant_abc, *, *
g, user_123, agent, tenant_abc
```

**Authorization Middleware:**
```go
func (s *AuthService) Authorize(resource, action string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        tenantID := c.GetString("tenant_id")

        // Check permission using Casbin
        allowed, err := s.enforcer.Enforce(userID, tenantID, resource, action)
        if err != nil || !allowed {
            c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 4. Session Management

#### 4.1 Session Storage

**Redis-based session storage:**
```go
type Session struct {
    ID           string    `json:"id"`
    UserID       string    `json:"user_id"`
    TenantID     string    `json:"tenant_id"`
    RefreshToken string    `json:"refresh_token"`
    CreatedAt    time.Time `json:"created_at"`
    ExpiresAt    time.Time `json:"expires_at"`
    LastActivity time.Time `json:"last_activity"`
    IPAddress    string    `json:"ip_address"`
    UserAgent    string    `json:"user_agent"`
    Revoked      bool      `json:"revoked"`
}

func (s *AuthService) createSession(session *Session) error {
    key := fmt.Sprintf("session:%s", session.ID)
    data, _ := json.Marshal(session)
    return s.redis.Set(ctx, key, data, 7*24*time.Hour).Err()
}
```

#### 4.2 Session Revocation

```go
func (s *AuthService) RevokeSession(sessionID string) error {
    key := fmt.Sprintf("session:%s", sessionID)
    session, err := s.getSession(sessionID)
    if err != nil {
        return err
    }

    session.Revoked = true
    data, _ := json.Marshal(session)
    return s.redis.Set(ctx, key, data, 7*24*time.Hour).Err()
}

func (s *AuthService) RevokeAllUserSessions(userID string) error {
    // Find all sessions for user
    sessions := s.getUserSessions(userID)
    for _, session := range sessions {
        s.RevokeSession(session.ID)
    }
    return nil
}
```

---

## Database Schema

### Users Table (Per Tenant Schema)

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- active, suspended, locked
    email_verified BOOLEAN NOT NULL DEFAULT false,
    email_verified_at TIMESTAMP,
    mfa_enabled BOOLEAN NOT NULL DEFAULT false,
    mfa_method VARCHAR(20), -- totp, sms, email
    mfa_secret VARCHAR(255),
    failed_login_attempts INT NOT NULL DEFAULT 0,
    locked_until TIMESTAMP,
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, email)
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
```

### User Roles Table (Per Tenant Schema)

```sql
CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, role)
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
```

### OAuth Accounts Table (Per Tenant Schema)

```sql
CREATE TABLE oauth_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- google, microsoft, github
    provider_user_id VARCHAR(255) NOT NULL,
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_oauth_accounts_user_id ON oauth_accounts(user_id);
```

### Password Reset Tokens Table (Per Tenant Schema)

```sql
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_password_reset_tokens_token ON password_reset_tokens(token);
CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
```

---

## API Endpoints

### Authentication Endpoints

```
POST   /api/auth/register               # Register new user
POST   /api/auth/login                  # Login with email/password
POST   /api/auth/logout                 # Logout (revoke session)
POST   /api/auth/refresh                # Refresh access token
POST   /api/auth/verify-email           # Verify email address
POST   /api/auth/resend-verification    # Resend verification email

POST   /api/auth/password/forgot        # Request password reset
POST   /api/auth/password/reset         # Reset password with token
PUT    /api/auth/password/change        # Change password (authenticated)

POST   /api/auth/mfa/enable             # Enable MFA
POST   /api/auth/mfa/verify             # Verify MFA code
POST   /api/auth/mfa/disable            # Disable MFA

GET    /api/auth/oauth/:provider        # Initiate OAuth flow
GET    /api/auth/oauth/:provider/callback # OAuth callback
```

### User Management Endpoints

```
GET    /api/users                       # List users (tenant admin)
GET    /api/users/:id                   # Get user details
PUT    /api/users/:id                   # Update user
DELETE /api/users/:id                   # Delete user (soft delete)
POST   /api/users/:id/suspend           # Suspend user
POST   /api/users/:id/activate          # Activate user
GET    /api/users/me                    # Get current user
PUT    /api/users/me                    # Update current user
```

### Role Management Endpoints

```
GET    /api/users/:id/roles             # Get user roles
POST   /api/users/:id/roles             # Assign role to user
DELETE /api/users/:id/roles/:role       # Remove role from user
GET    /api/roles                       # List available roles
GET    /api/roles/:role/permissions     # Get role permissions
```

---

## Security Requirements

### 1. Password Security
- Bcrypt hashing with cost factor 12
- Salted hashes
- Enforce password complexity requirements
- Prevent common passwords

### 2. Account Lockout
- Lock account after 5 failed login attempts
- Lock duration: 30 minutes
- Send notification on account lockout

### 3. Token Security
- RS256 algorithm for JWT signing
- Rotate signing keys every 90 days
- Short-lived access tokens (15 minutes)
- Secure refresh token storage
- Token revocation support

### 4. Session Security
- Session invalidation on password change
- Concurrent session limits (optional)
- Session activity tracking
- IP address validation (optional)

### 5. MFA Security
- TOTP with 30-second time window
- Backup codes for account recovery
- Rate limiting on MFA verification attempts

---

## Performance Requirements

- **Login:** < 500ms response time
- **Token Validation:** < 50ms
- **Authorization Check:** < 10ms
- **Session Lookup:** < 20ms (cached)

---

## Testing Requirements

### Unit Tests
- Password hashing and validation
- JWT token generation and validation
- MFA code generation and verification
- Role and permission checking

### Integration Tests
- Complete login flow
- OAuth integration flow
- Password reset flow
- Session management

### Security Tests
- Brute force protection
- Token manipulation prevention
- SQL injection prevention
- XSS prevention

---

## Success Criteria

1. ✅ Users can register and login securely
2. ✅ JWT tokens generated and validated correctly
3. ✅ RBAC enforces permissions across all resources
4. ✅ MFA works for all supported methods
5. ✅ OAuth integration functional for all providers
6. ✅ Account lockout prevents brute force attacks
7. ✅ Sessions tracked and revocable
8. ✅ Performance targets met

---

**Next Steps:** See `tasks.md` for implementation task breakdown.
