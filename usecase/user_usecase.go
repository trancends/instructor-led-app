package usecase

import (
	"log"

	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/repository"
)

type UserUsecase interface {
	CreateUser(payload model.User) error
}

type userUsecase struct {
	userRepository repository.UserRepository
}

func NewUserUsecase(userRepository repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}

func (u *userUsecase) CreateUser(payload model.User) error {
	log.Println("user usecase CreateUser")

	err := u.userRepository.Create(payload)
	if err != nil {
		return err
	}
	return nil
}
