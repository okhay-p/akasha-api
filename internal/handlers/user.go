package handlers

import (
	"akasha-api/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
	}
	c.IndentedJSON(http.StatusOK, users)
}
