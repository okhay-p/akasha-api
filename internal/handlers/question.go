package handlers

import (
	"akasha-api/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CheckAnswer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	option, err := strconv.Atoi(c.Param("option"))
	ques, err := services.GetQuestionByUUID(id)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"answer": ques.CorrectAnswer == int32(option)})

}
