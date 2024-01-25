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

type UserRepository interface {
	Create(payload model.User) error
	List(page int, size int) ([]model.User, sharedmodel.Paging, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Create(payload model.User) error {
	var err error
	user := payload
	currTime := time.Now().Local()
	user.CreatedAt = currTime
	user.UpdatedAt = currTime

	err = u.db.QueryRow(config.InsertUser, user.Name, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		log.Println("err at inserting user", err)
		return err
	}

	return nil
}

func (u *userRepository) List(page int, size int) ([]model.User, sharedmodel.Paging, error) {
	var users []model.User
	offset := (page - 1) * size
	query := config.SelectUserPagination

	rows, err := u.db.Query(query, size, offset)
	if err != nil {
		log.Println("userRepository.Query:", err.Error())
		return nil, sharedmodel.Paging{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role)
		if err != nil {
			log.Println("userRepository.Scan:", err.Error())
			return nil, sharedmodel.Paging{}, err
		}
		users = append(users, user)
	}

	totalRows := 0
	if err := u.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalRows); err != nil {
		log.Println("userRepository select count:", err.Error())
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
