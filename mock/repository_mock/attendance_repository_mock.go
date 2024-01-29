package repository_mock

import (
	"enigmaCamp.com/instructor_led/model"
	"github.com/stretchr/testify/mock"
)

type AttendanceRepositoryMock struct {
	mock.Mock
}

func (m *AttendanceRepositoryMock) GetAttendance(id string) (model.Attendance, error) {
	args := m.Called(id)
	return args.Get(0).(model.Attendance), args.Error(1)
}

func (m *AttendanceRepositoryMock) List() ([]model.Attendance, error) {
	args := m.Called()
	return args.Get(0).([]model.Attendance), args.Error(1)
}

func (m *AttendanceRepositoryMock) Create(user_id string, schedule_id string) (model.Attendance, error) {
	args := m.Called(user_id, schedule_id)
	return args.Get(0).(model.Attendance), args.Error(1)
}

func (m *AttendanceRepositoryMock) GetByID(user_id string, schedule_id string) (model.Attendance, error) {
	args := m.Called(user_id, schedule_id)
	return args.Get(0).(model.Attendance), args.Error(1)
}

func (m *AttendanceRepositoryMock) DeleteAttendance(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
