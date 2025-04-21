package handlers

import (
	"akasha-api/internal/model"
	"akasha-api/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
	}
	c.IndentedJSON(http.StatusOK, users)
}

func GetUserByUUID(c *gin.Context) {
	user_id, ok := c.Get("UUID")

	if !ok {
		log.Println("user_id not found")
		c.Abort()
		return
	}

	uuid, _ := uuid.Parse(user_id.(string))

	log.Println("GetUserByUUID:", uuid)
	var user model.AlUser
	var err error
	user, err = services.GetUserByUUID(uuid)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error getting user profile"})
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}
