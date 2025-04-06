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

	uuid, ok := c.Get("UUID")
	if ok {
		log.Println("User: ", uuid)
	}

	// TODO: Might need validation
	var newInput inputContent
	if err := c.BindJSON(&newInput); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	// Get lesson plan from AI
	lessonPlan, err := services.GenerateLessonPlan(newInput.Content)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	log.Println(lessonPlan.MainTitle)

	c.IndentedJSON(http.StatusOK, gin.H{"message": newInput.Content})

}
