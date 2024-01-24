package repository

import (
	"database/sql"
	"log"
	"time"

	"enigmaCamp.com/instructor_led/config"
	"enigmaCamp.com/instructor_led/model"
)

type UserRepository interface {
	Create(payload model.User) error
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
	payload.CreatedAt = currTime
	payload.UpdatedAt = currTime

	err = u.db.QueryRow(config.InsertUser, user.Name, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		log.Println("err at inserting user", err)
		return err
	}

	return nil
}
