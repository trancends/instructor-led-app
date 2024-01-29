package repository

import (
	"database/sql"
	"log"
	"math"
	"time"

	"enigmaCamp.com/instructor_led/config"
	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
)

type ScheduleRepository interface {
	ListScheduleByRole(page int, size int, role string) ([]model.Schedule, sharedmodel.Paging, error)
	ListScheduled(page int, size int) ([]model.Schedule, sharedmodel.Paging, error)
	CreateScheduled(payload model.Schedule) (model.Schedule, error)
	GetByID(id string) (model.Schedule, error)
	Delete(id string) error
	UpdateDocumentation(id string, pictureURL string) error
}

type scheduleRepository struct {
	db *sql.DB
}

// CreateScheduled implements ScheduleRepository.
func (s *scheduleRepository) CreateScheduled(payload model.Schedule) (model.Schedule, error) {
	var schedule model.Schedule

	// Assuming s.db is a *sql.DB connection
	err := s.db.QueryRow(
		config.InsertSchedule,
		payload.UserID,
		payload.Date,
		payload.StartTime,
		payload.EndTime,
		payload.Documentation,
	).Scan(
		&schedule.ID,
		&schedule.UserID,
		&schedule.Date,
		&schedule.StartTime,
		&schedule.EndTime,
		&schedule.Documentation,
	)
	if err != nil {
		log.Println("scheduleRepository.CreateScheduled:", err.Error())
		return schedule, err
	}

	return schedule, nil
}

// ListScheduleByRole implements ScheduleRepository.
func (s *scheduleRepository) ListScheduleByRole(page int, size int, role string) ([]model.Schedule, sharedmodel.Paging, error) {
	var schedules []model.Schedule
	offset := (page - 1) * size
	query := config.SelectScheduleByRole
	rows, err := s.db.Query(query, size, offset, role)
	if err != nil {
		log.Println("scheduleRepository.Query:", err.Error())
		return nil, sharedmodel.Paging{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var schedule model.Schedule
		err := rows.Scan(&schedule.ID, &schedule.UserID, &schedule.Date, &schedule.StartTime, &schedule.EndTime, &schedule.Documentation)
		if err != nil {
			log.Println("scheduleRepository.Scan:", err.Error())
			return nil, sharedmodel.Paging{}, err
		}
		schedules = append(schedules, schedule)
	}

	totalRows := 0
	totalQuery := `SELECT COUNT(schedule_id) FROM schedules s JOIN users u ON s.user_id = u.id WHERE u.role = $1 AND s.deleted_at IS NULL`
	if err := s.db.QueryRow(totalQuery, role).Scan(&totalRows); err != nil {
		log.Println("scheduleRepository.GetScheduleByRole select count:", err.Error())
		return nil, sharedmodel.Paging{}, err
	}

	paging := sharedmodel.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return schedules, paging, nil
}

// List implements ParticipantRepository.
func (s *scheduleRepository) ListScheduled(page int, size int) ([]model.Schedule, sharedmodel.Paging, error) {
	var schedules []model.Schedule
	offset := (page - 1) * size
	query := config.SelectSchedulePagination

	rows, err := s.db.Query(query, size, offset)
	if err != nil {
		log.Println("scheduleRepository.Query:", err.Error())
		return nil, sharedmodel.Paging{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var schedule model.Schedule
		err := rows.Scan(&schedule.ID, &schedule.UserID, &schedule.Date, &schedule.StartTime, &schedule.EndTime, &schedule.Documentation)
		if err != nil {
			log.Println("scheduleRepository.Scan:", err.Error())
			return nil, sharedmodel.Paging{}, err
		}
		schedules = append(schedules, schedule)
	}

	totalRows := 0
	if err := s.db.QueryRow("SELECT COUNT(*) FROM schedules").Scan(&totalRows); err != nil {
		return nil, sharedmodel.Paging{}, err
	}

	paging := sharedmodel.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}
	return schedules, paging, nil
}

// GetByID implements ScheduleRepository.
func (s *scheduleRepository) GetByID(id string) (model.Schedule, error) {
	var schedule model.Schedule
	schedule.ID = id
	rows, err := s.db.Query(config.SelectScheduleByID, id)
	if err != nil {
		log.Println("scheduleRepository.Query:", err.Error())
		return schedule, err
	}
	defer rows.Close()

	if !rows.Next() {
		return schedule, sql.ErrNoRows
	}

	err = rows.Scan(&schedule.ID, &schedule.UserID, &schedule.Date, &schedule.StartTime, &schedule.EndTime, &schedule.Documentation)
	if err != nil {
		log.Println("scheduleRepository.Scan:", err.Error())
		return schedule, err
	}
	return schedule, nil
}

func (s *scheduleRepository) Delete(id string) error {
	deletedAt := time.Now().Local()
	result, err := s.db.Exec(config.DeleteSchedule, deletedAt, id)
	if err != nil {
		log.Println("scheduleRepository.Delete:", err.Error())
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("scheduleRepository.Delete:", err.Error())
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *scheduleRepository) UpdateDocumentation(id string, pictureURL string) error {
	currentTime := time.Now().Local()
	_, err := s.db.Exec(config.UpdateScheduleDoc, pictureURL, currentTime, id)
	if err != nil {
		log.Println("scheduleRepository.UpdateDocumentation:", err.Error())
		return err
	}

	return nil
}

func NewSchedulesRepository(db *sql.DB) ScheduleRepository {
	return &scheduleRepository{
		db: db,
	}
}
