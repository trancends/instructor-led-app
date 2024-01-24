package repository

import (
	"database/sql"
	"log"
	"math"

	"enigmaCamp.com/instructor_led/config"
	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
)

type ParticipantRepository interface {
	ListScheduled(page int, size int) ([]model.Schedule, sharedmodel.Paging, error)
}

type participantRepository struct {
	db *sql.DB
}

// List implements ParticipantRepository.
func (p *participantRepository) ListScheduled(page int, size int) ([]model.Schedule, sharedmodel.Paging, error) {
	var schedules []model.Schedule
	offset := (page - 1) * size
	query := config.SelectSchedulePagination

	rows, err := p.db.Query(query, size, offset)
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
	if err := p.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalRows); err != nil {
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

func NewSchedulesRepository(db *sql.DB) ParticipantRepository {
	return &participantRepository{
		db: db,
	}
}
