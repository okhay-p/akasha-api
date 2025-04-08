package services

import (
	"akasha-api/internal/model"
	"akasha-api/pkg/db"
	"encoding/json"
	"log"
	"time"

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

func GetOrInsertTopicProgress(prg *model.AlUserTopicProgress) (model.AlUserTopicProgress, error) {
	var res model.AlUserTopicProgress

	if err := db.DB.FirstOrCreate(&res, model.AlUserTopicProgress{TopicID: prg.TopicID, UserID: prg.UserID}).Error; err != nil {
		return res, err
	}

	// db.DB.Create(prg)

	return res, nil
}

func GetTopicProgress(uid uuid.UUID, tid uuid.UUID) (model.AlUserTopicProgress, error) {
	var progress model.AlUserTopicProgress

	err := db.DB.Where("user_id = ? AND topic_id = ?", uid, tid).First(&progress).Error

	return progress, err
}

func UpdateTopicProgress(prg *model.AlUserTopicProgress, newLesson int32) error {
	return db.DB.Model(prg).Update("current_lesson", newLesson).Error
}

func GetAllTopics() ([]model.AlTopic, error) {
	var topics []model.AlTopic

	if err := db.DB.Find(&topics).Error; err != nil {
		return topics, err
	}

	return topics, nil
}

func GetTopicByUUID(uuid uuid.UUID) (model.AlTopic, error) {
	var topic model.AlTopic

	if err := db.DB.First(&topic, uuid).Error; err != nil {
		log.Println(err)
		return topic, err
	}
	return topic, nil
}

// Question struct to represent a single question
type TopicQuestion struct {
	QuestionID    string   `json:"question_id"`
	QuestionText  string   `json:"question_text"`
	QuestionOrder int      `json:"question_order"`
	Options       []string `json:"options"`
}

// Lesson struct to represent a single lesson
type TopicLesson struct {
	LessonID    string          `json:"lesson_id"`
	LessonTitle string          `json:"lesson_title"`
	Objectives  []string        `json:"objectives"`
	Content     []string        `json:"content"`
	LessonOrder int             `json:"lesson_order"`
	Questions   []TopicQuestion `json:"questions"`
}

// Topic struct to represent the entire topic
type TopicFullDetails struct {
	TopicID    string        `json:"topic_id"`
	TopicTitle string        `json:"topic_title"`
	IsPublic   bool          `json:"is_public"`
	Emoji      string        `json:"emoji"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	CreatedBy  string        `json:"created_by"`
	Lessons    []TopicLesson `json:"lessons"`
}

type TopicFullDetailsResult struct {
	TopicData []uint8 `json:"topic_data"`
}

func GetTopicFullDetailsByUUID(uuid uuid.UUID) (TopicFullDetails, error) {

	stmt := `
	WITH LessonQuestions AS (
    SELECT
        l.topic_id,
        l.id AS lesson_id,
        l.title AS lesson_title,
        l.objectives,
        l.content,
        l.order_number AS lesson_order,
        COALESCE(json_agg(
            json_build_object(
                'question_id', q.id,
                'question_text', q.question_text,
                'question_order', q.order_number,
                'options', q.options
            )
        ), '[]'::json) AS questions
    FROM
        al_lesson AS l
    LEFT JOIN
        al_question AS q ON l.id = q.lesson_id
    GROUP BY
        l.id, l.title, l.objectives, l.content, l.order_number
),
TopicLessons AS (
    SELECT
        t.id AS topic_id,
        t.title AS topic_title,
        t.is_public,
        t.emoji,
        t.created_at,
        t.updated_at,
        t.created_by,
        COALESCE(json_agg(
            json_build_object(
                'lesson_id', lq.lesson_id,
                'lesson_title', lq.lesson_title,
                'objectives', lq.objectives,
                'content', lq.content,
                'lesson_order', lq.lesson_order,
                'questions', lq.questions
            )
        ), '[]'::json) AS lessons
    FROM
        al_topic AS t
    LEFT JOIN
        LessonQuestions AS lq ON t.id = lq.topic_id
    WHERE
        t.id = ?
    GROUP BY
        t.id, t.title, t.is_public, t.emoji, t.created_at, t.updated_at, t.created_by
)
SELECT
    json_build_object(
        'topic_id', tl.topic_id,
        'topic_title', tl.topic_title,
        'is_public', tl.is_public,
        'emoji', tl.emoji,
        'created_at', tl.created_at,
        'updated_at', tl.updated_at,
        'created_by', tl.created_by,
        'lessons', tl.lessons
    ) AS topic_data
FROM
    TopicLessons AS tl;`

	var result TopicFullDetailsResult
	var topic TopicFullDetails
	if err := db.DB.Raw(stmt, uuid).Scan(&result).Error; err != nil {
		log.Println(err)
		return topic, err
	}

	err := json.Unmarshal(result.TopicData, &topic)
	if err != nil {
		log.Println(err)
		return topic, err
	}

	return topic, nil
}
