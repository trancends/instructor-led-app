package usecase

import (
	"errors"
	"testing"
	"time"

	"enigmaCamp.com/instructor_led/mock/repository_mock"
	"enigmaCamp.com/instructor_led/model"
	"github.com/stretchr/testify/suite"
)

var (
	expectedAttendances = []model.Attendance{
		{
			ID:         "1",
			UserID:     "1",
			ScheduleID: "1",
		},
		{
			ID:         "2",
			UserID:     "2",
			ScheduleID: "2",
		},
	}
	expectedAttendance = model.Attendance{
		ID:         "1",
		UserID:     "1",
		ScheduleID: "1",
	}
	currTime               = time.Now().Local()
	expectedAttendanceByID = model.Attendance{
		ID:         "1",
		UserID:     "1",
		ScheduleID: "1",
		CreatedAt:  nil,
		UpdatedAt:  nil,
		DeletedAt:  nil,
	}
)

type AttendanceUsecaseTestSuite struct {
	suite.Suite
	attendanceRepoMock *repository_mock.AttendanceRepositoryMock
	attendanceUsecase  AttendanceUsecase
}

func (s *AttendanceUsecaseTestSuite) SetupTest() {
	s.attendanceRepoMock = new(repository_mock.AttendanceRepositoryMock)
	s.attendanceUsecase = NewAttendanceUsecase(s.attendanceRepoMock)
}

func (s *AttendanceUsecaseTestSuite) TestAddAttendance() {
	s.attendanceRepoMock.On("GetByID", expectedAttendance.UserID, expectedAttendance.ScheduleID).Return(model.Attendance{}, nil)
	s.attendanceRepoMock.On("Create", expectedAttendance.UserID, expectedAttendance.ScheduleID).Return(expectedAttendance, nil)

	attendance, err := s.attendanceUsecase.AddAttendance(expectedAttendance.UserID, expectedAttendance.ScheduleID)

	s.Nil(err)
	s.Equal(expectedAttendance, attendance)
}

func (s *AttendanceUsecaseTestSuite) TestGetAddAttendance_Error() {
	s.attendanceRepoMock.On("GetByID", expectedAttendance.UserID, expectedAttendance.ScheduleID).Return(expectedAttendance, errors.New("err"))
	s.attendanceRepoMock.On("Create", expectedAttendance.UserID, expectedAttendance.ScheduleID).Return(expectedAttendance, nil)

	_, err := s.attendanceUsecase.AddAttendance(expectedAttendance.UserID, expectedAttendance.ScheduleID)

	s.Error(err)
}

func (s *AttendanceUsecaseTestSuite) TestGetAttendance() {
	s.attendanceRepoMock.On("GetAttendance", expectedAttendance.ID).Return(expectedAttendance, nil)

	attendance, err := s.attendanceUsecase.GetAttendance(expectedAttendance.ID)

	s.Nil(err)
	s.Equal(expectedAttendance, attendance)
}

func (s *AttendanceUsecaseTestSuite) TestGetAttendance_Error() {
	s.attendanceRepoMock.On("GetAttendance", expectedAttendance.ID).Return(model.Attendance{}, errors.New("err"))

	_, err := s.attendanceUsecase.GetAttendance(expectedAttendance.ID)

	s.Error(err)
}

func (s *AttendanceUsecaseTestSuite) TestListAttendances() {
	s.attendanceRepoMock.On("List").Return(expectedAttendances, nil)

	attendances, err := s.attendanceUsecase.ListAttendances()

	s.Nil(err)
	s.Equal(expectedAttendances, attendances)
}

func (s *AttendanceUsecaseTestSuite) TestListAttendances_Error() {
	s.attendanceRepoMock.On("List").Return([]model.Attendance{}, errors.New("err"))

	_, err := s.attendanceUsecase.ListAttendances()

	s.Error(err)
}

func TestAttendanceUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(AttendanceUsecaseTestSuite))
}
