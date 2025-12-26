package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user within a tenant
type User struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	TenantID             uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	Email                string     `json:"email" db:"email"`
	PasswordHash         string     `json:"-" db:"password_hash"` // Never expose password hash
	FirstName            *string    `json:"first_name,omitempty" db:"first_name"`
	LastName             *string    `json:"last_name,omitempty" db:"last_name"`
	Phone                *string    `json:"phone,omitempty" db:"phone"`
	Mobile               *string    `json:"mobile,omitempty" db:"mobile"`
	Status               string     `json:"status" db:"status"`
	EmailVerified        bool       `json:"email_verified" db:"email_verified"`
	EmailVerifiedAt      *time.Time `json:"email_verified_at,omitempty" db:"email_verified_at"`
	MFAEnabled           bool       `json:"mfa_enabled" db:"mfa_enabled"`
	MFAMethod            *string    `json:"mfa_method,omitempty" db:"mfa_method"`
	MFASecret            string     `json:"-" db:"mfa_secret"` // Never expose MFA secret
	FailedLoginAttempts  int        `json:"failed_login_attempts" db:"failed_login_attempts"`
	LockedUntil          *time.Time `json:"locked_until,omitempty" db:"locked_until"`
	LastLoginAt          *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	LastLoginIP          *string    `json:"last_login_ip,omitempty" db:"last_login_ip"`
	PasswordChangedAt    *time.Time `json:"password_changed_at,omitempty" db:"password_changed_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	Roles                []string   `json:"roles,omitempty" db:"-"` // Populated from user_roles table
}

// UserStatus constants
const (
	UserStatusActive    = "active"
	UserStatusSuspended = "suspended"
	UserStatusLocked    = "locked"
	UserStatusDeleted   = "deleted"
)

// UserRole constants
const (
	RoleTenantAdmin     = "tenant_admin"
	RoleTenantManager   = "tenant_manager"
	RoleAgent           = "agent"
	RoleAgentAssistant  = "agent_assistant"
	RoleClient          = "client"
)

// MFAMethod constants
const (
	MFAMethodTOTP  = "totp"
	MFAMethodSMS   = "sms"
	MFAMethodEmail = "email"
)

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive && u.DeletedAt == nil
}

// IsLocked checks if user account is locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// FullName returns the user's full name
func (u *User) FullName() string {
	if u.FirstName != nil && u.LastName != nil {
		return *u.FirstName + " " + *u.LastName
	}
	if u.FirstName != nil {
		return *u.FirstName
	}
	if u.LastName != nil {
		return *u.LastName
	}
	return u.Email
}

// HasRole checks if user has a specific role
func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // seconds
	User         *User  `json:"user"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
