package middleware

import (
	"api-gateway/internal/core/ports"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(validator ports.TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		claims, err := validator.Validate(token)
		if err != nil {
			slog.Error("Error to validate token", err, token)
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
