package handlers

import (
	"akasha-api/internal/model"
	"akasha-api/internal/req_structs"
	"akasha-api/internal/services"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateTopic(c *gin.Context) {
	user_id, ok := c.Get("UUID")

	if !ok {
		log.Println("user_id not found")
		c.Abort()
		return
	}

	uuid, err := uuid.Parse(user_id.(string))

	// TODO: Might need validation
	var newInput req_structs.GenerateTopicReqBody

	content := c.PostForm("content")
	isPublicStr := c.PostForm("is_public")       // Gin returns strings, so handle bool carefully
	onlyContentStr := c.PostForm("only_content") // Same as above
	numOfLessonsStr := c.PostForm("num_of_lessons")

	// Handle boolean values
	isPublic := false
	if isPublicStr == "true" {
		isPublic = true
	}

	onlyContent := false
	if onlyContentStr == "true" {
		onlyContent = true
	}

	//Handle integer value
	numOfLessons := 0
	_, err = fmt.Sscan(numOfLessonsStr, &numOfLessons)
	if err != nil {
		log.Printf("Error converting num_of_lessons to integer: %s\n", err)
		numOfLessons = 3 // Provide a default value or handle the error as needed
	}
	newInput.Content = content
	newInput.IsPublic = isPublic
	newInput.OnlyContent = onlyContent
	newInput.NumOfLessons = uint8(numOfLessons)

	file, err := c.FormFile("attachment")
	if err != nil {
		log.Println(err)
	}
	if file != nil && file.Size > 15*1024*1024 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "File size too large"})
		c.Abort()
		return
	}

	// Get lesson plan from AI
	lessonPlan, err := services.GenerateLessonPlan(newInput, file)
	if err != nil {

		if len(err.Error()) > 11 && err.Error()[:9] == "req_error" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()[11:]})
			c.Abort()
			return
		}
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		c.Abort()
		return
	}

	// Get username from uuid
	user, err := services.GetUserByUUID(uuid)

	// Store the topic in DB
	var topic model.AlTopic

	topic.Title = lessonPlan.MainTitle
	topic.Emoji = lessonPlan.Emoji
	topic.CreatedBy = user.Username
	topic.IsPublic = newInput.IsPublic
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

	var prg model.AlUserTopicProgress

	prg.TopicID = topicId
	prg.UserID = uuid
	prg.CurrentLesson = 0

	_, err = services.GetOrInsertTopicProgress(&prg)

	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": topicId})

}

func FirstOrCreateTopicProgress(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println(c.Param("id"))
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	user_id, ok := c.Get("UUID")

	if !ok {
		log.Println("user_id not found")
		c.Abort()
		return
	}

	uId, err := uuid.Parse(user_id.(string))

	var prg model.AlUserTopicProgress

	prg.TopicID = id
	prg.UserID = uId
	prg.CurrentLesson = 0

	res, err := services.GetOrInsertTopicProgress(&prg)
	c.IndentedJSON(http.StatusOK, res)
}

func UpdateTopicProgressCurrentLesson(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))
	order, err := strconv.Atoi(c.Param("order"))
	if err != nil {
		log.Println(c.Param("id"))
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	user_id, ok := c.Get("UUID")

	if !ok {
		log.Println("user_id not found")
		c.Abort()
		return
	}

	uId, err := uuid.Parse(user_id.(string))
	prg, err := services.GetTopicProgress(uId, id)

	if err == gorm.ErrRecordNotFound {
		c.Status(http.StatusNotFound)
		c.Abort()
		return
	}

	err = services.UpdateTopicProgress(&prg, int32(order)+1)
	if err != nil {

		log.Println(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)

}

func DeleteProgress(c *gin.Context) {
	user_id, ok := c.Get("UUID")
	if !ok {
		log.Println("user_id not found")
		c.Abort()
		return
	}

	uId, err := uuid.Parse(user_id.(string))
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		log.Println("Missing IDs for deleting progress")
		return
	}

	prg, err := services.GetTopicProgress(uId, id)
	err = services.DeleteTopicProgress(&prg)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		log.Println(err)
		return
	}

	c.Status(http.StatusNoContent)

}

func GetTopicsL1(c *gin.Context) {

	user_id, ok := c.Get("UUID")

	if !ok {
		log.Println("user_id not found")
		c.Abort()
		return
	}

	uId, err := uuid.Parse(user_id.(string))
	user, err := services.GetUserByUUID(uId)

	topics, err := services.GetAllPublicTopics()
	userTopics, err := services.GetAllUserGeneratedTopics(user.Username)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"public_topics": topics, "user_topics": userTopics})

}

func GetTopicByUUID(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println(c.Param("id"))
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	topic, err := services.GetTopicByUUID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Status(http.StatusNotFound)
			c.Abort()
			return

		}
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, topic)
}

func GetFullTopicDetails(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println(c.Param("id"))
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	topic, err := services.GetTopicFullDetailsByUUID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Status(http.StatusNotFound)
			c.Abort()
			return

		}
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, topic)
}

func GetTopicsRelatedToUser(c *gin.Context) {
	user_id, ok := c.Get("UUID")

	if !ok {
		log.Println("user_id not found")
		c.Abort()
		return
	}

	uId, _ := uuid.Parse(user_id.(string))
	topics, err := services.GetTopicsRelatedToUser(uId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, topics)
}
