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

func InsertNewUser(user *model.AlUser) (uuid.UUID, error) {
	if err := db.DB.Create(user).Error; err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

func GetUserUUIDByGoogleSub(sub string) (uuid.UUID, error) {
	var user model.AlUser

	if err := db.DB.First(&user, sub).Error; err != nil {
		return user.ID, err
	}
	return user.ID, nil
}
