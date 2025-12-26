package repository

import (
	"database/sql"
	"fmt"

	"github.com/comply360/shared/models"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (
			tenant_id, email, password_hash, first_name, last_name, status
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		user.TenantID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Status,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(tenantID, userID uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, status,
			email_verified, email_verified_at, mfa_enabled, mfa_method, mfa_secret,
			failed_login_attempts, locked_until, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`

	user := &models.User{}
	var mfaSecret sql.NullString

	err := r.db.QueryRow(query, userID, tenantID).Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Status,
		&user.EmailVerified,
		&user.EmailVerifiedAt,
		&user.MFAEnabled,
		&user.MFAMethod,
		&mfaSecret,
		&user.FailedLoginAttempts,
		&user.LockedUntil,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Set MFASecret if valid
	if mfaSecret.Valid {
		user.MFASecret = mfaSecret.String
	}

	// Load roles from user_roles table
	user.Roles, _ = r.GetUserRoles(user.ID)

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(tenantID uuid.UUID, email string) (*models.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, status,
			email_verified, email_verified_at, mfa_enabled, mfa_method, mfa_secret,
			failed_login_attempts, locked_until, last_login_at, created_at, updated_at
		FROM users
		WHERE email = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`

	user := &models.User{}
	var mfaSecret sql.NullString

	err := r.db.QueryRow(query, email, tenantID).Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Status,
		&user.EmailVerified,
		&user.EmailVerifiedAt,
		&user.MFAEnabled,
		&user.MFAMethod,
		&mfaSecret,
		&user.FailedLoginAttempts,
		&user.LockedUntil,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Set MFASecret if valid
	if mfaSecret.Valid {
		user.MFASecret = mfaSecret.String
	}

	// Load roles from user_roles table
	user.Roles, _ = r.GetUserRoles(user.ID)

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users SET
			email = $1,
			first_name = $2,
			last_name = $3,
			status = $4,
			email_verified = $5,
			email_verified_at = $6,
			mfa_enabled = $7,
			mfa_method = $8,
			mfa_secret = $9,
			failed_login_attempts = $10,
			locked_until = $11,
			last_login_at = $12,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $13 AND tenant_id = $14
	`

	result, err := r.db.Exec(
		query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Status,
		user.EmailVerified,
		user.EmailVerifiedAt,
		user.MFAEnabled,
		user.MFAMethod,
		user.MFASecret,
		user.FailedLoginAttempts,
		user.LockedUntil,
		user.LastLoginAt,
		user.ID,
		user.TenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(tenantID, userID uuid.UUID, passwordHash string) error {
	query := `
		UPDATE users SET
			password_hash = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3
	`

	result, err := r.db.Exec(query, passwordHash, userID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// IncrementFailedLoginAttempts increments the failed login attempts counter
func (r *UserRepository) IncrementFailedLoginAttempts(tenantID uuid.UUID, email string) error {
	query := `
		UPDATE users SET
			failed_login_attempts = failed_login_attempts + 1,
			updated_at = CURRENT_TIMESTAMP
		WHERE email = $1 AND tenant_id = $2
	`

	_, err := r.db.Exec(query, email, tenantID)
	return err
}

// ResetFailedLoginAttempts resets the failed login attempts counter
func (r *UserRepository) ResetFailedLoginAttempts(tenantID, userID uuid.UUID) error {
	query := `
		UPDATE users SET
			failed_login_attempts = 0,
			locked_until = NULL,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND tenant_id = $2
	`

	_, err := r.db.Exec(query, userID, tenantID)
	return err
}

// LockAccount locks a user account
func (r *UserRepository) LockAccount(tenantID uuid.UUID, email string, lockedUntil sql.NullTime) error {
	query := `
		UPDATE users SET
			locked_until = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE email = $2 AND tenant_id = $3
	`

	_, err := r.db.Exec(query, lockedUntil, email, tenantID)
	return err
}

// VerifyEmail marks a user's email as verified
func (r *UserRepository) VerifyEmail(tenantID, userID uuid.UUID) error {
	query := `
		UPDATE users SET
			email_verified = true,
			email_verified_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND tenant_id = $2
	`

	result, err := r.db.Exec(query, userID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// GetUserRoles retrieves all roles for a user
func (r *UserRepository) GetUserRoles(userID uuid.UUID) ([]string, error) {
	query := `
		SELECT role FROM user_roles
		WHERE user_id = $1 AND (expires_at IS NULL OR expires_at > NOW())
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// AssignRole assigns a role to a user
func (r *UserRepository) AssignRole(userID uuid.UUID, role string, grantedBy *uuid.UUID) error {
	query := `
		INSERT INTO user_roles (user_id, role, granted_by)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, role) DO NOTHING
	`

	_, err := r.db.Exec(query, userID, role, grantedBy)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

// RemoveRole removes a role from a user
func (r *UserRepository) RemoveRole(userID uuid.UUID, role string) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role = $2`

	result, err := r.db.Exec(query, userID, role)
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found for user")
	}

	return nil
}

// Delete soft deletes a user by setting deleted_at timestamp
func (r *UserRepository) Delete(tenantID, userID uuid.UUID) error {
	query := `
		UPDATE users SET
			deleted_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, userID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}
