package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/comply360/shared/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	// Context keys
	UserIDKey    = "user_id"
	UserEmailKey = "user_email"
	UserRolesKey = "user_roles"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, errors.Unauthorized("Missing authorization header"))
			c.Abort()
			return
		}

		// Check for Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, errors.Unauthorized("Invalid authorization format. Expected 'Bearer {token}'"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, errors.NewAPIError(errors.ErrInvalidToken, "Invalid or expired token"))
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, errors.NewAPIError(errors.ErrInvalidToken, "Token is not valid"))
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, errors.Unauthorized("Invalid token claims"))
			c.Abort()
			return
		}

		// Extract user ID
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, errors.Unauthorized("Invalid user ID in token"))
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errors.Unauthorized("Invalid user ID format in token"))
			c.Abort()
			return
		}

		// Extract tenant ID (if present)
		if tenantIDStr, ok := claims["tenant_id"].(string); ok {
			tenantID, err := uuid.Parse(tenantIDStr)
			if err == nil {
				c.Set(TenantIDKey, tenantID)
			}
		}

		// Extract email
		if email, ok := claims["email"].(string); ok {
			c.Set(UserEmailKey, email)
		}

		// Extract roles
		if rolesInterface, ok := claims["roles"]; ok {
			// Roles can be either []interface{} or []string
			var roles []string
			switch v := rolesInterface.(type) {
			case []interface{}:
				for _, role := range v {
					if roleStr, ok := role.(string); ok {
						roles = append(roles, roleStr)
					}
				}
			case []string:
				roles = v
			}
			c.Set(UserRolesKey, roles)
		}

		// Set user context
		c.Set(UserIDKey, userID)

		c.Next()
	}
}

// RequireRole middleware ensures user has one of the required roles
func RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user roles from context
		rolesInterface, exists := c.Get(UserRolesKey)
		if !exists {
			c.JSON(http.StatusForbidden, errors.Forbidden("No roles found for user"))
			c.Abort()
			return
		}

		userRoles, ok := rolesInterface.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, errors.Forbidden("Invalid roles format"))
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRequiredRole := false
		for _, requiredRole := range requiredRoles {
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					hasRequiredRole = true
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}

		if !hasRequiredRole {
			c.JSON(http.StatusForbidden, errors.NewAPIError(
				errors.ErrInsufficientPermissions,
				fmt.Sprintf("This operation requires one of the following roles: %v", requiredRoles),
			))
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID retrieves user ID from Gin context
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user ID type in context")
	}

	return id, nil
}

// GetUserRoles retrieves user roles from Gin context
func GetUserRoles(c *gin.Context) ([]string, error) {
	roles, exists := c.Get(UserRolesKey)
	if !exists {
		return nil, fmt.Errorf("user roles not found in context")
	}

	r, ok := roles.([]string)
	if !ok {
		return nil, fmt.Errorf("invalid roles type in context")
	}

	return r, nil
}
