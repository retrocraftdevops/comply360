package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/comply360/auth-service/internal/repository"
	"github.com/comply360/shared/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

const (
	maxFailedAttempts = 5
	accountLockDuration = 30 * time.Minute
	accessTokenDuration = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

type AuthService struct {
	userRepo   *repository.UserRepository
	redis      *redis.Client
	jwtSecret  string
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Email    string    `json:"email"`
	Roles    []string  `json:"roles"`
}

func NewAuthService(userRepo *repository.UserRepository, redis *redis.Client, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		redis:     redis,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account
func (s *AuthService) Register(tenantID uuid.UUID, req *models.RegisterRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(tenantID, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	firstName := req.FirstName
	lastName := req.LastName
	user := &models.User{
		TenantID:      tenantID,
		Email:         req.Email,
		PasswordHash:  string(hashedPassword),
		FirstName:     &firstName,
		LastName:      &lastName,
		Status:        models.UserStatusActive,
		EmailVerified: false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign default role (client)
	if err := s.userRepo.AssignRole(user.ID, models.RoleClient, nil); err != nil {
		return nil, fmt.Errorf("failed to assign default role: %w", err)
	}

	// TODO: Send verification email

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(tenantID uuid.UUID, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(tenantID, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if account is locked
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return nil, fmt.Errorf("account is locked until %v", user.LockedUntil)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// Increment failed login attempts
		s.userRepo.IncrementFailedLoginAttempts(tenantID, req.Email)

		// Lock account if max attempts reached
		if user.FailedLoginAttempts+1 >= maxFailedAttempts {
			lockedUntil := sql.NullTime{Time: time.Now().Add(accountLockDuration), Valid: true}
			s.userRepo.LockAccount(tenantID, req.Email, lockedUntil)
			return nil, fmt.Errorf("account locked due to too many failed login attempts")
		}

		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if MFA is enabled
	if user.MFAEnabled {
		// TODO: Implement full MFA flow
		return nil, fmt.Errorf("MFA authentication not yet implemented")
	}

	// Reset failed login attempts
	s.userRepo.ResetFailedLoginAttempts(tenantID, user.ID)

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in Redis
	ctx := context.Background()
	s.redis.Set(ctx, fmt.Sprintf("refresh_token:%s", refreshToken), user.ID.String(), refreshTokenDuration)

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(accessTokenDuration.Seconds()),
		User:         user,
	}, nil
}

// RefreshToken generates a new access token using a refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	ctx := context.Background()

	// Check if refresh token exists in Redis
	userIDStr, err := s.redis.Get(ctx, fmt.Sprintf("refresh_token:%s", refreshToken)).Result()
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// TODO: Get tenant ID from context
	// For now, we'll need to parse it from the refresh token
	claims := &jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	tenantIDStr, ok := (*claims)["tenant_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID")
	}

	// Get user
	user, err := s.userRepo.GetByID(tenantID, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &models.AuthResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(accessTokenDuration.Seconds()),
		User:        user,
	}, nil
}

// VerifyEmail verifies a user's email using a verification token
func (s *AuthService) VerifyEmail(tenantID uuid.UUID, token string) error {
	ctx := context.Background()

	// Get user ID from Redis
	userIDStr, err := s.redis.Get(ctx, fmt.Sprintf("email_verification:%s", token)).Result()
	if err != nil {
		return fmt.Errorf("invalid or expired verification token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	// Verify email
	if err := s.userRepo.VerifyEmail(tenantID, userID); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	// Delete verification token
	s.redis.Del(ctx, fmt.Sprintf("email_verification:%s", token))

	return nil
}

// SetupMFA generates MFA secret and QR code for user
func (s *AuthService) SetupMFA(tenantID, userID uuid.UUID, method string) (string, error) {
	// Get user
	user, err := s.userRepo.GetByID(tenantID, userID)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	// Generate TOTP secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Comply360",
		AccountName: user.Email,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	// Update user with MFA secret
	user.MFASecret = key.Secret()
	user.MFAMethod = &method

	if err := s.userRepo.Update(user); err != nil {
		return "", fmt.Errorf("failed to update user: %w", err)
	}

	// Return QR code URL
	return key.URL(), nil
}

// VerifyMFA verifies MFA code and enables MFA for user
func (s *AuthService) VerifyMFA(tenantID, userID uuid.UUID, code string) error {
	// Get user
	user, err := s.userRepo.GetByID(tenantID, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.MFASecret == "" {
		return fmt.Errorf("MFA not set up")
	}

	// Verify TOTP code
	valid := totp.Validate(code, user.MFASecret)
	if !valid {
		return fmt.Errorf("invalid MFA code")
	}

	// Enable MFA
	user.MFAEnabled = true
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to enable MFA: %w", err)
	}

	return nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Parse claims into TokenClaims struct
	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	tenantIDStr, ok := claims["tenant_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid tenant ID in token")
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID format: %w", err)
	}

	email, _ := claims["email"].(string)

	// Parse roles
	var roles []string
	if rolesInterface, ok := claims["roles"].([]interface{}); ok {
		for _, role := range rolesInterface {
			if roleStr, ok := role.(string); ok {
				roles = append(roles, roleStr)
			}
		}
	}

	return &TokenClaims{
		UserID:   userID,
		TenantID: tenantID,
		Email:    email,
		Roles:    roles,
	}, nil
}

// generateAccessToken generates a JWT access token
func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":       user.ID.String(),
		"tenant_id": user.TenantID.String(),
		"email":     user.Email,
		"roles":     user.Roles,
		"exp":       time.Now().Add(accessTokenDuration).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// generateRefreshToken generates a JWT refresh token
func (s *AuthService) generateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":       user.ID.String(),
		"tenant_id": user.TenantID.String(),
		"exp":       time.Now().Add(refreshTokenDuration).Unix(),
		"iat":       time.Now().Unix(),
		"type":      "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// generateMFAToken generates a temporary token for MFA verification
func (s *AuthService) generateMFAToken(userID, tenantID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub":       userID.String(),
		"tenant_id": tenantID.String(),
		"exp":       time.Now().Add(5 * time.Minute).Unix(),
		"iat":       time.Now().Unix(),
		"type":      "mfa",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
