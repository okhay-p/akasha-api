package services

import (
	"akasha-api/internal/model"
	"akasha-api/pkg/db"
)

func InsertTopic(topic *model.AlTopic) (string, error) {
	if err := db.DB.Create(topic).Error; err != nil {
		return "", err
	}
	return topic.ID, nil
}

func InsertLesson(lesson *model.AlLesson) (string, error) {
	if err := db.DB.Create(lesson).Error; err != nil {
		return "", err
	}
	return lesson.ID, nil
}

func InsertQuestion(question *model.AlQuestion) (string, error) {
	if err := db.DB.Create(question).Error; err != nil {
		return "", err
	}
	return question.ID, nil
}
