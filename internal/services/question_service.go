package services

import (
	"akasha-api/internal/model"
	"akasha-api/pkg/db"
	"log"

	"github.com/google/uuid"
)

func GetQuestionByUUID(qid uuid.UUID) (model.AlQuestion, error) {
	var ques model.AlQuestion

	if err := db.DB.First(&ques, qid).Error; err != nil {
		log.Println(err)
		return ques, err
	}

	return ques, nil

}
