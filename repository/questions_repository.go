// questions/repository/repository.go
package repository

import (
	"database/sql"
	"log"

	"enigmaCamp.com/instructor_led/model"
)

type QuestionsRepository interface {
	Get(date string) ([]*model.Schedule, error)
	List() ([]model.Question, error)
}

type questionsRepository struct {
	db *sql.DB
}

// List implements QuestionsRepository.
func (q *questionsRepository) List() ([]model.Question, error) {
	rows, err := q.db.Query("SELECT id, description, status FROM questions")
	if err != nil {
		log.Println("Error retrieving questions:", err)
		return nil, err
	}
	defer rows.Close()

	var questions []model.Question

	for rows.Next() {
		var question model.Question
		err := rows.Scan(&question.ID,
			&question.Description,
			&question.Status)
		if err != nil {
			log.Println("Error scanning question row:", err)
			return nil, err
		}

		questions = append(questions, question)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over question rows:", err)
		return nil, err
	}

	return questions, nil
}

// Get implements QuestionsRepository.
func (q *questionsRepository) Get(date string) ([]*model.Schedule, error) {
	rows, err := q.db.Query(`
		SELECT
			schedules.id AS schedule_id,
			schedules.user_id,
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
			schedules.date = $1 ;
	`, date)
	if err != nil {
		log.Println("questionsRepository.Get:", err.Error())
		return nil, err
	}
	defer rows.Close()

	ScheduleSlice := make([]*model.Schedule, 0)

	for rows.Next() {
		var question model.Question
		var schedule model.Schedule

		err := rows.Scan(
			&schedule.ID,
			&schedule.UserID,
			&schedule.Date,
			&schedule.StartTime,
			&schedule.EndTime,
			&question.ID,
			&question.Description,
			&question.Status,
		)

		if err != nil {
			log.Println("questionsRepository.Get:", err.Error())
			return nil, err
		}

		var found bool
		for _, existingSchedule := range ScheduleSlice {
			if existingSchedule.ID == schedule.ID {
				found = true
				existingSchedule.Questions = append(existingSchedule.Questions, question)
				break
			}
		}

		// If schedule doesn't exist, add it to ScheduleSlice
		if !found {
			schedule.Questions = append(schedule.Questions, question)
			ScheduleSlice = append(ScheduleSlice, &schedule)
		}
	}

	return ScheduleSlice, nil
}

func NewQuestionsRepository(db *sql.DB) QuestionsRepository {
	return &questionsRepository{
		db: db,
	}
}
