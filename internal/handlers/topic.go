package handlers

import (
	"akasha-api/internal/model"
	"akasha-api/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type inputContent struct {
	Content string `json:"content"`
}

func CreateTopic(c *gin.Context) {
	user_id, ok := c.Get("UUID")

	if !ok {
		log.Println("user_id not found")
	}

	// TODO: Might need validation
	var newInput inputContent
	if err := c.BindJSON(&newInput); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		c.Abort()
		return
	}

	// Get lesson plan from AI
	lessonPlan, err := services.GenerateLessonPlan(newInput.Content)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		c.Abort()
		return
	}

	// Store the topic in DB
	var topic model.AlTopic

	topic.Title = lessonPlan.MainTitle
	topic.Emoji = lessonPlan.Emoji
	topic.CreatedBy = user_id.(string)
	topic.IsPublic = true
	topic.StatusID = 1

	topicId, err := services.InsertTopic(&topic)

	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		c.Abort()
		return
	}

	// Iterate over lessons and store each lesson
	for i, l := range lessonPlan.Lessons {

		var lesson model.AlLesson
		lesson.TopicID = topicId
		lesson.Title = l.Title
		lesson.Objectives = l.Objectives
		lesson.Content = l.Content
		lesson.OrderNumber = int32(i)

		lessonId, err := services.InsertLesson(&lesson)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			c.Abort()
			return
		}

		// Iterate over questions of each lesson
		log.Println(lessonId)
		for ii, q := range l.Questions {

			var question model.AlQuestion

			question.LessonID = lessonId
			question.QuestionText = q.QuestionText
			question.Options = q.Options
			question.CorrectAnswer = int32(q.CorrectAnswer)
			question.OrderNumber = int32(ii)

			_, err = services.InsertQuestion(&question)

			if err != nil {
				log.Println(err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
				c.Abort()
				return
			}

		}

	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": topicId})

}
