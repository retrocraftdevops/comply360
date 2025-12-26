package services

import (
	"testing"

	"github.com/comply360/auth-service/internal/repository"
	"github.com/comply360/shared/models"
	testhelpers "github.com/comply360/shared/testing"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	// Setup
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)
	tdb.CreateTestTables(t)

	tredis := testhelpers.SetupTestRedis(t)
	defer tredis.Cleanup(t)

	userRepo := repository.NewUserRepository(tdb.DB)
	authService := NewAuthService(userRepo, tredis.Client, "test_jwt_secret")

	// Test: Successful registration
	req := &models.RegisterRequest{
		Email:     "newuser@example.com",
		Password:  "SecurePassword123!",
		FirstName: "John",
		LastName:  "Doe",
	}

	user, err := authService.Register(tdb.TenantID, req)
	testhelpers.AssertNoError(t, err, "Failed to register user")
	testhelpers.AssertNotNil(t, user, "User should not be nil")
	testhelpers.AssertEqual(t, req.Email, user.Email, "Email mismatch")
	testhelpers.AssertFalse(t, user.EmailVerified, "Email should not be verified initially")

	// Verify password was hashed
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	testhelpers.AssertNoError(t, err, "Password hash verification failed")

	// Test: Duplicate email
	_, err = authService.Register(tdb.TenantID, req)
	testhelpers.AssertError(t, err, "Should fail to register duplicate email")
}

func TestAuthService_Login(t *testing.T) {
	// Setup
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)
	tdb.CreateTestTables(t)

	tredis := testhelpers.SetupTestRedis(t)
	defer tredis.Cleanup(t)

	userRepo := repository.NewUserRepository(tdb.DB)
	authService := NewAuthService(userRepo, tredis.Client, "test_jwt_secret")

	// Create a user first
	password := "SecurePassword123!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	firstName := "John"
	lastName := "Doe"

	user := &models.User{
		TenantID:     tdb.TenantID,
		Email:        "login@example.com",
		PasswordHash: string(hashedPassword),
		FirstName:    &firstName,
		LastName:     &lastName,
		Status:       models.UserStatusActive,
	}
	err := userRepo.Create(user)
	testhelpers.AssertNoError(t, err)

	// Assign role
	err = userRepo.AssignRole(user.ID, models.RoleClient, nil)
	testhelpers.AssertNoError(t, err)

	// Test: Successful login
	loginReq := &models.LoginRequest{
		Email:    "login@example.com",
		Password: password,
	}

	authResp, err := authService.Login(tdb.TenantID, loginReq)
	testhelpers.AssertNoError(t, err, "Login failed")
	testhelpers.AssertNotNil(t, authResp, "Auth response should not be nil")
	testhelpers.AssertNotEqual(t, "", authResp.AccessToken, "Access token should be set")
	testhelpers.AssertNotEqual(t, "", authResp.RefreshToken, "Refresh token should be set")

	// Test: Invalid password
	loginReq.Password = "WrongPassword"
	_, err = authService.Login(tdb.TenantID, loginReq)
	testhelpers.AssertError(t, err, "Should fail with wrong password")

	// Test: Non-existent user
	loginReq.Email = "nonexistent@example.com"
	loginReq.Password = password
	_, err = authService.Login(tdb.TenantID, loginReq)
	testhelpers.AssertError(t, err, "Should fail with non-existent user")
}

func TestAuthService_LoginAccountLocking(t *testing.T) {
	// Setup
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)
	tdb.CreateTestTables(t)

	tredis := testhelpers.SetupTestRedis(t)
	defer tredis.Cleanup(t)

	userRepo := repository.NewUserRepository(tdb.DB)
	authService := NewAuthService(userRepo, tredis.Client, "test_jwt_secret")

	// Create a user
	password := "SecurePassword123!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	firstName := "John"

	user := &models.User{
		TenantID:     tdb.TenantID,
		Email:        "locktest@example.com",
		PasswordHash: string(hashedPassword),
		FirstName:    &firstName,
		Status:       models.UserStatusActive,
	}
	err := userRepo.Create(user)
	testhelpers.AssertNoError(t, err)

	err = userRepo.AssignRole(user.ID, models.RoleClient, nil)
	testhelpers.AssertNoError(t, err)

	// Test: Attempt login with wrong password multiple times
	loginReq := &models.LoginRequest{
		Email:    "locktest@example.com",
		Password: "WrongPassword",
	}

	// Attempt 5 failed logins (maxFailedAttempts)
	for i := 0; i < 5; i++ {
		_, err = authService.Login(tdb.TenantID, loginReq)
		testhelpers.AssertError(t, err, "Should fail with wrong password")
	}

	// Verify account is locked
	locked, err := userRepo.GetByEmail(tdb.TenantID, "locktest@example.com")
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertNotNil(t, locked.LockedUntil, "Account should be locked")

	// Test: Even correct password should fail when account is locked
	loginReq.Password = password
	_, err = authService.Login(tdb.TenantID, loginReq)
	testhelpers.AssertError(t, err, "Should fail when account is locked")
}

func TestAuthService_ValidateToken(t *testing.T) {
	// Setup
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)
	tdb.CreateTestTables(t)

	tredis := testhelpers.SetupTestRedis(t)
	defer tredis.Cleanup(t)

	userRepo := repository.NewUserRepository(tdb.DB)
	authService := NewAuthService(userRepo, tredis.Client, "test_jwt_secret")

	// Create a user and login
	password := "SecurePassword123!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	firstName := "Jane"

	user := &models.User{
		TenantID:     tdb.TenantID,
		Email:        "tokentest@example.com",
		PasswordHash: string(hashedPassword),
		FirstName:    &firstName,
		Status:       models.UserStatusActive,
	}
	err := userRepo.Create(user)
	testhelpers.AssertNoError(t, err)

	err = userRepo.AssignRole(user.ID, models.RoleTenantAdmin, nil)
	testhelpers.AssertNoError(t, err)

	loginReq := &models.LoginRequest{
		Email:    "tokentest@example.com",
		Password: password,
	}

	authResp, err := authService.Login(tdb.TenantID, loginReq)
	testhelpers.AssertNoError(t, err)

	// Test: Validate access token
	claims, err := authService.ValidateToken(authResp.AccessToken)
	testhelpers.AssertNoError(t, err, "Failed to validate token")
	testhelpers.AssertNotNil(t, claims, "Claims should not be nil")
	testhelpers.AssertEqual(t, user.Email, claims.Email, "Email mismatch in claims")
	testhelpers.AssertEqual(t, user.ID, claims.UserID, "UserID mismatch in claims")

	// Test: Invalid token
	_, err = authService.ValidateToken("invalid.token.here")
	testhelpers.AssertError(t, err, "Should fail to validate invalid token")

	// Test: Token with wrong secret
	wrongSecretService := NewAuthService(userRepo, tredis.Client, "wrong_secret")
	_, err = wrongSecretService.ValidateToken(authResp.AccessToken)
	testhelpers.AssertError(t, err, "Should fail to validate token with wrong secret")
}
