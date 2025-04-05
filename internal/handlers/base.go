package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAPIStatus(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "API is up and running 🥳"})
}
