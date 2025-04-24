package handlers

import (
	"akasha-api/internal/model"
	"akasha-api/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FeedbackValue struct {
	Feedback string `json:"feedback"`
}

func CreateFeedback(c *gin.Context) {

	var fb model.AlFeedback

	var feedbackText FeedbackValue
	if err := c.BindJSON(&feedbackText); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		c.Abort()
		return
	}

	fb.Text = feedbackText.Feedback
	err := services.InsertFeedback(&fb)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		log.Println(err)
		return
	}
	c.Status(http.StatusCreated)
}
