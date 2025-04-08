package handlers

import (
	"akasha-api/pkg/config"
	"akasha-api/pkg/jwt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleLogin(c *gin.Context) {
	// TEMP LOG IN HANDLER MUST CHANGE AFTER GOOGLE OAUTH SETUP
	token, err := jwt.CreateToken(config.AkashaLearnUUID)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		c.Abort()
		return
	}

	if config.Dev {
		c.IndentedJSON(http.StatusOK, gin.H{"auth": token})
		return
	}
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("test", "test", 300, "/", "localhost", false, false)
	c.SetCookie("token", token, 86400, "/", "localhost", true, true)
	c.IndentedJSON(http.StatusOK, gin.H{"auth": "success"})
}
