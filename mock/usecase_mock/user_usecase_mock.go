package usecase_mock

import (
	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
	"github.com/stretchr/testify/mock"
)

type UserUsecaseMock struct {
	mock.Mock
}

func (m *UserUsecaseMock) CreateUserCSV(payload []model.User) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *UserUsecaseMock) CreateUser(payload model.User) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *UserUsecaseMock) ListAllUsers(page int, size int) ([]model.User, sharedmodel.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]model.User), args.Get(1).(sharedmodel.Paging), args.Error(2)
}

func (m *UserUsecaseMock) GetUserByEmail(email string) (model.User, error) {
	args := m.Called(email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserUsecaseMock) GetUserByID(id string) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserUsecaseMock) GetUserByRole(role string, page int, size int) ([]model.User, sharedmodel.Paging, error) {
	args := m.Called(role, page, size)
	return args.Get(0).([]model.User), args.Get(1).(sharedmodel.Paging), args.Error(2)
}

func (m *UserUsecaseMock) UpdateUser(payload model.User) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *UserUsecaseMock) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
