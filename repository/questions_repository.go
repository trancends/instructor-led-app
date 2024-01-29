package repository

import (
	"database/sql"
	"log"
	"time"

	"enigmaCamp.com/instructor_led/config"
	"enigmaCamp.com/instructor_led/model"
)

type QuestionsRepository interface {
	CreateQuestions(payload model.Question) (model.Question, error)
	Get(date string) ([]*model.Schedule, error)
	GetByID(id string) (model.Question, error)
	GetByScheduleID(scheduleID string) ([]model.Question, error)
	List() ([]model.Question, error)
	Delete(id string) error
	UpdateStatus(payload model.Question) error
}

type questionsRepository struct {
	db *sql.DB
}

// CreateQuestions implements QuestionsRepository.
func (q *questionsRepository) CreateQuestions(payload model.Question) (model.Question, error) {
	questions := payload
	rows, err := q.db.Query(config.InsertQuestions, payload.UserID, payload.ScheduleID, payload.Description)
	if err != nil {
		log.Println("questionsRepository.Query:", err.Error())
		return questions, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&questions.ID)
		if err != nil {
			log.Println("questionsRepository.Scan:", err.Error())
		}
	}
	return questions, nil
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

// GetByID implements QuestionsRepository.
func (q *questionsRepository) GetByID(id string) (model.Question, error) {
	var question model.Question
	err := q.db.QueryRow(config.SelectQuestionsByID, id).Scan(&question.ID, &question.ScheduleID, &question.Description, &question.Status)
	if err != nil {
		log.Println("questionsRepository.GetByID:", err.Error())
		return question, err
	}
	return question, nil
}

// GetByScheduleID
func (q *questionsRepository) GetByScheduleID(scheduleID string) ([]model.Question, error) {
	var questions []model.Question
	rows, err := q.db.Query(config.SelectQuestionsByScheduleID, scheduleID)
	if err != nil {
		log.Println("questionsRepository.GetByScheduleID:", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var question model.Question
		err := rows.Scan(&question.ID, &question.ScheduleID, &question.Description, &question.Status)
		if err != nil {
			log.Println("questionsRepository.GetByScheduleID:", err.Error())
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

// Delete implements QuestionsRepository.
// Soft Delete
func (q *questionsRepository) Delete(id string) error {
	deletedAt := time.Now().Local()
	_, err := q.db.Exec(config.DeleteQuestions, deletedAt, id)
	if err != nil {
		log.Println("questionsRepository.Delete:", err.Error())
		return err
	}
	return nil
}

// Update implements QuestionsRepository.
func (q *questionsRepository) UpdateStatus(payload model.Question) error {
	updated_at := time.Now().Local()
	_, err := q.db.Exec(config.UpdateQuestions, payload.Status, updated_at, payload.ID)
	if err != nil {
		log.Println("questionsRepository.Update:", err.Error())
		return err
	}
	return nil
}

func NewQuestionsRepository(db *sql.DB) QuestionsRepository {
	return &questionsRepository{
		db: db,
	}
}
