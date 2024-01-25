package repository

import (
	"database/sql"
	"log"
	"math"

	"enigmaCamp.com/instructor_led/config"
	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
)

type ScheduleRepository interface {
	ListScheduled(page int, size int) ([]model.Schedule, sharedmodel.Paging, error)
	CreateScheduled(payload model.Schedule) (model.Schedule, error)
	GetByID(id string) (model.Schedule, error)
}

type scheduleRepository struct {
	db *sql.DB
}

// CreateScheduled implements ScheduleRepository.
func (s *scheduleRepository) CreateScheduled(payload model.Schedule) (model.Schedule, error) {
	var schedule model.Schedule
	rows, err := s.db.Query(config.InsertSchedule)
	if err != nil {
		log.Println("scheduleRepository.Query:", err.Error())
		return schedule, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&schedule.ID, &schedule.UserID, &schedule.Date, &schedule.StartTime, &schedule.EndTime, &schedule.Documentation)
		if err != nil {
			log.Println("scheduleRepository.Scan:", err.Error())
		}
	}
	return schedule, nil
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
		err := rows.Scan(&schedule.ID, &schedule.UserID, &schedule.Date, &schedule.StartTime, &schedule.EndTime, &schedule.Documentation, &schedule.CreatedAt, &schedule.UpdatedAt)
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
	rows, err := s.db.Query(config.SelectScheduleByID, id)
	if err != nil {
		log.Println("scheduleRepository.Query:", err.Error())
		return schedule, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&schedule.ID, &schedule.UserID, &schedule.Date, &schedule.StartTime, &schedule.EndTime, &schedule.Documentation)
		if err != nil {
			log.Println("scheduleRepository.Scan:", err.Error())
			return schedule, err
		}
	}
	return schedule, nil
}

func NewSchedulesRepository(db *sql.DB) ScheduleRepository {
	return &scheduleRepository{
		db: db,
	}
}
