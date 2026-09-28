package middleware

import (
	"be-golang-poliklinik/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(authService auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak, token tidak ditemukan"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := authService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("clean_token", tokenString)
		c.Next()
	}
}

func RoleBlockMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetString("clean_token")
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("rahasia_super_poliklinik"), nil
		})

		claims, _ := token.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)

		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Anda tidak memiliki izin."})
			c.Abort()
			return
		}
		c.Next()
	}
}
