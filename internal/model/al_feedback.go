package model

import "time"

const TableNameAlFeedback = "al_feedback"

// AlLesson mapped from table <al_feedback>
type AlFeedback struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	Text      string    `gorm:"column:text;not null" json:"text"`
	CreatedAt time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
}

func (*AlFeedback) TableName() string {
	return TableNameAlFeedback
}
