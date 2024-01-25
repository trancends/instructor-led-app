// questions/repository/repository.go
package repository

import (
	"database/sql"
	"log"

	"enigmaCamp.com/instructor_led/model"
)

type QuestionsRepository interface {
	Get(id string) (model.Question, error)
	List() ([]model.Question, error)
}

type questionsRepository struct {
	db *sql.DB
}

// Get implements QuestionsRepository.
func (q *questionsRepository) Get(id string) (model.Question, error) {
	var question model.Question

	row := q.db.QueryRow(`
		SELECT
			schedules.id AS schedule_id,
			schedules.date,
			schedules.start_time,
			schedules.end_time,
			questions.id AS question_id,
			questions.description,
			questions.status
		FROM
			schedules
		JOIN
			questions ON schedules.id = questions.schedule_id
		WHERE
			questions.id = $1;
	`, id)

	err := row.Scan(
		&question.ScheduleID,
		&question.Date,
		&question.ID,
		&question.Description,
		&question.Status,
	)

	if err != nil {
		log.Println("questionsRepository.Get:", err.Error())
		return model.Question{}, err
	}

	return question, nil
}

// List implements QuestionsRepository.
func (q *questionsRepository) List() ([]model.Question, error) {
	var questions []model.Question

	rows, err := q.db.Query(`
		SELECT
			schedules.id AS schedule_id,
			schedules.date,
			schedules.start_time,
			schedules.end_time,
			questions.id AS question_id,
			questions.description,
			questions.status
		FROM
			schedules
		JOIN
			questions ON schedules.id = questions.schedule_id;
	`)
	if err != nil {
		log.Println("questionsRepository.List:", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var question model.Question

		err := rows.Scan(
			&question.ScheduleID,
			&question.Date,
			&question.ID,
			&question.Description,
			&question.Status,
			&question,
		)
		if err != nil {
			log.Println("questionsRepository.List:", err.Error())
			return nil, err
		}

		questions = append(questions, question)
	}

	if err := rows.Err(); err != nil {
		log.Println("questionsRepository.List:", err.Error())
		return nil, err
	}

	return questions, nil
}

func NewQuestionsRepository(db *sql.DB) QuestionsRepository {
	return &questionsRepository{
		db: db,
	}
}
