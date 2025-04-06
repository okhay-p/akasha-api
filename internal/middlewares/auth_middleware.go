package middlewares

import (
	"akasha-api/pkg/config"
	"akasha-api/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if config.Dev {
			authHeader := c.GetHeader("Authorization")
			if len(authHeader) < 7 {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				c.Abort()
				return
			}
			token := authHeader[7:]
			claims, err := jwt.VerifyToken(token)
			if err != nil {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				c.Abort()
				return
			}
			c.Set("UUID", claims.Subject)
		}

		c.Next()
	}
}
