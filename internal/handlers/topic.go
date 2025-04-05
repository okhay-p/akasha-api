package handlers

import (
	"akasha-api/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type inputContent struct {
	Content string `json:"content"`
}

func CreateTopic(c *gin.Context) {
	var newInput inputContent
	if err := c.BindJSON(&newInput); err != nil {
		log.Fatal(err)
		return
	}

	lessonPlan, err := services.GenerateLessonPlan(newInput.Content)
	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
	}

	log.Println(lessonPlan.MainTitle)

	c.IndentedJSON(http.StatusOK, gin.H{"message": newInput.Content})

}
