package repository

import (
	"database/sql"
	"log"

	"enigmaCamp.com/instructor_led/config"
	"enigmaCamp.com/instructor_led/model"
)

type QuestionsRepository interface {
	CreateQuestions(payload model.Questions) (model.Questions, error)
}

type questionsRepository struct {
	db *sql.DB
}

// CreateQuestions implements QuestionsRepository.
func (q *questionsRepository) CreateQuestions(payload model.Questions) (model.Questions, error) {
	var questions model.Questions
	rows, err := q.db.Query(config.InsertQuestions)
	if err != nil {
		log.Println("questionsRepository.Query:", err.Error())
		return questions, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&questions.ID, &questions.ScheduleID, &questions.Description, &questions.CreatedAt, &questions.UpdatedAt)
		if err != nil {
			log.Println("questionsRepository.Scan:", err.Error())
		}
	}
	return questions, nil

}

func NewQuestionsRepository(db *sql.DB) QuestionsRepository {
	return &questionsRepository{
		db: db,
	}
}
