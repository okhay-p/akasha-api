package middlewares

import (
	"akasha-api/pkg/config"
	"akasha-api/pkg/jwt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			token string
			err   error
		)

		if config.Dev {
			authHeader := c.GetHeader("Authorization")
			if len(authHeader) < 7 {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				c.Abort()
				return

			}
			token = authHeader[7:]
		} else {
			token, err = c.Cookie("token")
			if err != nil {
				log.Println("ERROR:", err)
				c.Abort()
				return
			}
		}
		claims, err := jwt.VerifyToken(token)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			log.Println(err)
			c.Abort()
			return

		}
		c.Set("UUID", claims.Subject)

		c.Next()
	}
}
