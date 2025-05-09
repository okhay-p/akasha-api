// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const TableNameAlQuestion = "al_question"

// AlQuestion mapped from table <al_question>
type AlQuestion struct {
	ID            uuid.UUID `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	LessonID      uuid.UUID `gorm:"column:lesson_id" json:"lesson_id"`
	QuestionText  string `gorm:"column:question_text" json:"question_text"`
	OrderNumber   int32  `gorm:"column:order_number;not null" json:"order_number"`
	Options       pq.StringArray `gorm:"column:options;not null;type:text[]" json:"options"`
	CorrectAnswer int32  `gorm:"column:correct_answer;not null" json:"correct_answer"`
}

// TableName AlQuestion's table name
func (*AlQuestion) TableName() string {
	return TableNameAlQuestion
}
