package middleware

import (
	"strings"

	"github.com/Alfian57/ruang-tenang-api/pkg/response"
	"github.com/Alfian57/ruang-tenang-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.AbortUnauthorized(c, "Header otorisasi diperlukan")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.AbortUnauthorized(c, "Format header otorisasi tidak valid")
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			response.AbortUnauthorized(c, "Token tidak valid atau telah kedaluwarsa")
			return
		}
		if !claims.PhoneVerified {
			response.AbortUnauthorized(c, "Verifikasi nomor WhatsApp diperlukan")
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		// Enforce ban/block/suspension on every request, not just at login,
		// so moderation actions take effect immediately for existing tokens.
		if !enforceAccountStatus(c) {
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware checks for token but doesn't require it
// Useful for endpoints that can be accessed by both guests and logged-in users
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.Next()
			return
		}
		if !claims.PhoneVerified {
			c.Next()
			return
		}

		// Set user info in context if token is valid
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		if !enforceAccountStatus(c) {
			return
		}

		c.Next()
	}
}

// AdminMiddleware checks if user has admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != "admin" {
			response.AbortForbidden(c, "Akses admin diperlukan")
			return
		}
		c.Next()
	}
}

// UserMiddleware checks if user has user role.
func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != "user" {
			response.AbortForbidden(c, "Akses pengguna diperlukan")
			return
		}
		c.Next()
	}
}

// MemberMiddleware remains for backward compatibility and maps to user role.
func MemberMiddleware() gin.HandlerFunc {
	return UserMiddleware()
}

// MitraMiddleware checks if user has mitra role.
func MitraMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != "mitra" {
			response.AbortForbidden(c, "Akses mitra diperlukan")
			return
		}
		c.Next()
	}
}

// GetUserID helper to get user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetUserRole helper to get user role from context
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	roleStr, ok := role.(string)
	if !ok {
		return "", false
	}
	return roleStr, true
}
