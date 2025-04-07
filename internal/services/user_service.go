package services

import (
	"akasha-api/internal/model"
	"akasha-api/pkg/db"

	"github.com/google/uuid"
)

func GetAllUsers() ([]model.AlUser, error) {
	var users []model.AlUser

	if err := db.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByUUID(uuid uuid.UUID) (model.AlUser, error) {
	var user model.AlUser

	if err := db.DB.First(&user, uuid).Error; err != nil {
		return user, err
	}
	return user, nil
}
