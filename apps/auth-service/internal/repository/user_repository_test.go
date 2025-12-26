package repository

import (
	"testing"

	"github.com/comply360/shared/models"
	testhelpers "github.com/comply360/shared/testing"
	"github.com/google/uuid"
)

func TestUserRepository_Create(t *testing.T) {
	// Setup test database
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)
	tdb.CreateTestTables(t)

	repo := NewUserRepository(tdb.DB)

	// Test data
	user := &models.User{
		TenantID:     tdb.TenantID,
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Status:       models.UserStatusActive,
	}

	// Test: Create user
	err := repo.Create(user)
	testhelpers.AssertNoError(t, err, "Failed to create user")
	testhelpers.AssertNotEqual(t, uuid.Nil, user.ID, "User ID should be set after creation")

	// Test: Verify user was created
	retrieved, err := repo.GetByID(tdb.TenantID, user.ID)
	testhelpers.AssertNoError(t, err, "Failed to retrieve user")
	testhelpers.AssertEqual(t, user.Email, retrieved.Email, "Email mismatch")
	testhelpers.AssertEqual(t, user.Status, retrieved.Status, "Status mismatch")
}

func TestUserRepository_GetByEmail(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)
	tdb.CreateTestTables(t)

	repo := NewUserRepository(tdb.DB)

	// Create test user
	user := &models.User{
		TenantID:     tdb.TenantID,
		Email:        "getbyemail@example.com",
		PasswordHash: "hashed_password",
		Status:       models.UserStatusActive,
	}
	err := repo.Create(user)
	testhelpers.AssertNoError(t, err)

	// Test: Get by email
	retrieved, err := repo.GetByEmail(tdb.TenantID, "getbyemail@example.com")
	testhelpers.AssertNoError(t, err, "Failed to get user by email")
	testhelpers.AssertNotNil(t, retrieved, "User should be found")
	testhelpers.AssertEqual(t, user.ID, retrieved.ID, "User ID mismatch")
	testhelpers.AssertEqual(t, user.Email, retrieved.Email, "Email mismatch")

	// Test: Non-existent email
	notFound, err := repo.GetByEmail(tdb.TenantID, "nonexistent@example.com")
	testhelpers.AssertError(t, err, "Should return error for non-existent email")
	testhelpers.AssertNil(t, notFound, "Should return nil for non-existent user")
}

func TestUserRepository_Update(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)
	tdb.CreateTestTables(t)

	repo := NewUserRepository(tdb.DB)

	// Create test user
	user := &models.User{
		TenantID:     tdb.TenantID,
		Email:        "update@example.com",
		PasswordHash: "hashed_password",
		Status:       models.UserStatusActive,
	}
	err := repo.Create(user)
	testhelpers.AssertNoError(t, err)

	// Test: Update user
	firstName := "John"
	lastName := "Doe"
	user.FirstName = &firstName
	user.LastName = &lastName
	user.EmailVerified = true

	err = repo.Update(user)
	testhelpers.AssertNoError(t, err, "Failed to update user")

	// Verify update
	retrieved, err := repo.GetByID(tdb.TenantID, user.ID)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, firstName, *retrieved.FirstName, "First name not updated")
	testhelpers.AssertEqual(t, lastName, *retrieved.LastName, "Last name not updated")
	testhelpers.AssertTrue(t, retrieved.EmailVerified, "EmailVerified not updated")
}

func TestUserRepository_AssignRole(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)
	tdb.CreateTestTables(t)

	repo := NewUserRepository(tdb.DB)

	// Create test user
	user := &models.User{
		TenantID:     tdb.TenantID,
		Email:        "role@example.com",
		PasswordHash: "hashed_password",
		Status:       models.UserStatusActive,
	}
	err := repo.Create(user)
	testhelpers.AssertNoError(t, err)

	// Test: Assign role
	err = repo.AssignRole(user.ID, models.RoleTenantAdmin, nil)
	testhelpers.AssertNoError(t, err, "Failed to assign role")

	// Test: Assign duplicate role (should be idempotent - no error)
	err = repo.AssignRole(user.ID, models.RoleTenantAdmin, nil)
	testhelpers.AssertNoError(t, err, "Duplicate role assignment should be idempotent")

	// Verify only one role exists
	roles, err := repo.GetUserRoles(user.ID)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, 1, len(roles), "Should have exactly 1 role after duplicate assignment")
}

func TestUserRepository_IncrementFailedLoginAttempts(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)
	tdb.CreateTestTables(t)

	repo := NewUserRepository(tdb.DB)

	// Create test user
	user := &models.User{
		TenantID:     tdb.TenantID,
		Email:        "failedlogin@example.com",
		PasswordHash: "hashed_password",
		Status:       models.UserStatusActive,
	}
	err := repo.Create(user)
	testhelpers.AssertNoError(t, err)

	// Test: Increment failed login attempts
	err = repo.IncrementFailedLoginAttempts(tdb.TenantID, user.Email)
	testhelpers.AssertNoError(t, err, "Failed to increment failed login attempts")

	// Verify increment
	retrieved, err := repo.GetByEmail(tdb.TenantID, user.Email)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, 1, retrieved.FailedLoginAttempts, "Failed login attempts not incremented")

	// Increment again
	err = repo.IncrementFailedLoginAttempts(tdb.TenantID, user.Email)
	testhelpers.AssertNoError(t, err)

	retrieved, err = repo.GetByEmail(tdb.TenantID, user.Email)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, 2, retrieved.FailedLoginAttempts, "Failed login attempts should be 2")
}

func TestUserRepository_Delete(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)
	tdb.CreateTestTables(t)

	repo := NewUserRepository(tdb.DB)

	// Create test user
	user := &models.User{
		TenantID:     tdb.TenantID,
		Email:        "delete@example.com",
		PasswordHash: "hashed_password",
		Status:       models.UserStatusActive,
	}
	err := repo.Create(user)
	testhelpers.AssertNoError(t, err)

	// Test: Delete user
	err = repo.Delete(tdb.TenantID, user.ID)
	testhelpers.AssertNoError(t, err, "Failed to delete user")

	// Verify deletion
	_, err = repo.GetByID(tdb.TenantID, user.ID)
	testhelpers.AssertError(t, err, "Should return error for deleted user")
}
