package repo_mock

import (
	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
	"github.com/stretchr/testify/mock"
)

type ScheduleRepositoryMock struct {
	mock.Mock
}

func (m *ScheduleRepositoryMock) ListScheduleByRole(page int, size int, role string) ([]model.Schedule, sharedmodel.Paging, error) {
	args := m.Called(page, size, role)
	return args.Get(0).([]model.Schedule), args.Get(1).(sharedmodel.Paging), args.Error(2)
}
func (m *ScheduleRepositoryMock) CreateScheduled(payload model.Schedule) (model.Schedule, error) {
	args := m.Called(payload)
	return args.Get(0).(model.Schedule), args.Error(1)
}

func (m *ScheduleRepositoryMock) ListScheduled(page int, size int) ([]model.Schedule, sharedmodel.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]model.Schedule), args.Get(1).(sharedmodel.Paging), args.Error(2)
}

func (m *ScheduleRepositoryMock) GetByID(id string) (model.Schedule, error) {
	args := m.Called(id)
	return args.Get(0).(model.Schedule), args.Error(1)
}

func (m *ScheduleRepositoryMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *ScheduleRepositoryMock) UpdateDocumentation(id string, pictureURL string) error {
	args := m.Called(id, pictureURL)
	return args.Error(0)
}
