package repository_mock

import (
	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) CreateUserCSV(payload []model.User) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *UserRepositoryMock) Create(payload model.User) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *UserRepositoryMock) List(page int, size int) ([]model.User, sharedmodel.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]model.User), args.Get(1).(sharedmodel.Paging), args.Error(2)
}

func (m *UserRepositoryMock) GetUserByEmail(email string) (model.User, error) {
	args := m.Called(email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserRepositoryMock) GetUserByID(id string) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserRepositoryMock) GetUserByRole(role string, page int, size int) ([]model.User, sharedmodel.Paging, error) {
	args := m.Called(role, page, size)
	return args.Get(0).([]model.User), args.Get(1).(sharedmodel.Paging), args.Error(2)
}

func (m *UserRepositoryMock) Update(payload model.User) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *UserRepositoryMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
