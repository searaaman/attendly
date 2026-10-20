package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/searaaman/attendly/internal/config"
	"github.com/searaaman/attendly/internal/services"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		tokenString, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found || tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		claims, err := services.ParseToken(tokenString, config.Cfg.JWTSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		ctx.Set("user_id", claims["user_id"])
		ctx.Next()
	}
}
