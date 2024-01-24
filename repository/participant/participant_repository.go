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
	ListScheduled(page int, size int) ([]model.User, sharedmodel.Paging, error)
}

type participantRepository struct {
	db *sql.DB
}

// List implements ParticipantRepository.
func (p *participantRepository) ListScheduled(page int, size int) ([]model.User, sharedmodel.Paging, error) {
	var users []model.User
	offset := (page - 1) * size
	query := config.SelectTaskPagination

	rows, err := p.db.Query(query, size, offset)
	if err != nil {
		log.Println("taskRepository.Query:", err.Error())
		return nil, sharedmodel.Paging{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			log.Println("taskRepository.Scan:", err.Error())
			return nil, sharedmodel.Paging{}, err
		}
		users = append(users, user)
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
	return users, paging, nil
}

func NewParticipantRepository(db *sql.DB) ParticipantRepository {
	return &participantRepository{
		db: db,
	}
}
