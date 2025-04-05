package services

import (
	"akasha-api/internal/model"
	"akasha-api/pkg/db"
)

func GetAllUsers() ([]model.AlUser, error) {
	var users []model.AlUser

	if err := db.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
