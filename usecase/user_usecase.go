package usecase

import (
	"log"

	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/repository"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
)

type UserUsecase interface {
	CreateUser(payload model.User) error
	ListAllUsers(page int, size int) ([]model.User, sharedmodel.Paging, error)
	GetUserByEmail(email string) (model.User, error)
	GetUserByRole(role string, page int, size int) ([]model.User, sharedmodel.Paging, error)
	UpdateUser(payload model.User) error
	DeleteUser(id string) error
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
	log.Println("calling user repo Create")
	err := u.userRepository.Create(payload)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) ListAllUsers(page int, size int) ([]model.User, sharedmodel.Paging, error) {
	log.Println("calling user repo List")
	users, paging, err := u.userRepository.List(page, size)
	if err != nil {
		return nil, sharedmodel.Paging{}, err
	}
	return users, paging, nil
}

func (u *userUsecase) GetUserByEmail(email string) (model.User, error) {
	log.Println("calling user repo Get")
	user, err := u.userRepository.GetUserByEmail(email)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *userUsecase) GetUserByRole(role string, page int, size int) ([]model.User, sharedmodel.Paging, error) {
	log.Println("calling user repo GetUserByRole")
	users, paging, err := u.userRepository.GetUserByRole(role, page, size)
	if err != nil {
		return nil, sharedmodel.Paging{}, err
	}
	return users, paging, nil
}

func (u *userUsecase) UpdateUser(payload model.User) error {
	log.Println("calling user repo Update")
	err := u.userRepository.Update(payload)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) DeleteUser(id string) error {
	log.Println("calling user repo Delete")
	err := u.userRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
