package usecase_mock

import (
	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
	"github.com/stretchr/testify/mock"
)

type SchedulesUseCaseMock struct {
	mock.Mock
}

func (m *SchedulesUseCaseMock) FindScheduleByRole(page int, size int, role string) ([]model.Schedule, sharedmodel.Paging, error) {
	args := m.Called(page, size, role)
	return args.Get(0).([]model.Schedule), args.Get(1).(sharedmodel.Paging), args.Error(2)
}

func (m *SchedulesUseCaseMock) FindAllScheduleUC(page int, size int) ([]model.Schedule, sharedmodel.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]model.Schedule), args.Get(1).(sharedmodel.Paging), args.Error(2)
}

func (m *SchedulesUseCaseMock) CreateScheduledUC(payload model.Schedule) (model.Schedule, error) {
	args := m.Called(payload)
	return args.Get(0).(model.Schedule), args.Error(1)
}

func (m *SchedulesUseCaseMock) FindByIDUC(id string) (model.Schedule, error) {
	args := m.Called(id)
	return args.Get(0).(model.Schedule), args.Error(1)
}

func (m *SchedulesUseCaseMock) DeletedScheduleIDUC(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *SchedulesUseCaseMock) UpdateScheduleDocumentation(id string, pictureURL string) error {
	args := m.Called(id, pictureURL)
	return args.Error(0)
}

func NewSchedulesUseCaseMock() *SchedulesUseCaseMock {
	return &SchedulesUseCaseMock{}
}
