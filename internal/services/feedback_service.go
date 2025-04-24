package services

import (
	"akasha-api/internal/model"
	"akasha-api/pkg/db"
)

func InsertFeedback(fb *model.AlFeedback) error {
	return db.DB.Create(fb).Error
}
