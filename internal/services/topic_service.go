package services

import (
	"akasha-api/internal/model"
	"akasha-api/pkg/db"

	"github.com/google/uuid"
)

func InsertTopic(topic *model.AlTopic) (uuid.UUID, error) {
	if err := db.DB.Create(topic).Error; err != nil {
		return uuid.Nil, err
	}
	return topic.ID, nil
}

func InsertLesson(lesson *model.AlLesson) (uuid.UUID, error) {
	if err := db.DB.Create(lesson).Error; err != nil {
		return uuid.Nil, err
	}
	return lesson.ID, nil
}

func InsertQuestion(question *model.AlQuestion) (uuid.UUID, error) {
	if err := db.DB.Create(question).Error; err != nil {
		return uuid.Nil, err
	}
	return question.ID, nil
}

func GetTopicByUUID(uuid uuid.UUID) (model.AlTopic, error) {
	var topic model.AlTopic

	if err := db.DB.Find(&topic, uuid).Error; err != nil {
		return topic, err
	}
	return topic, nil
}
