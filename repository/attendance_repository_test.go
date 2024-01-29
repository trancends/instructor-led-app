package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"enigmaCamp.com/instructor_led/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type AttendanceRepositoryTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    AttendanceRepository
}

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

func (s *AttendanceRepositoryTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New()
	s.mockDB = db
	s.mockSql = mock
	s.repo = NewAttendanceRepository(s.mockDB)
}

func (s *AttendanceRepositoryTestSuite) TestCreate() {
	s.mockSql.
		ExpectQuery(regexp.QuoteMeta("INSERT INTO attendances (user_id, schedule_id) VALUES ($1, $2) RETURNING id")).
		WithArgs(expectedAttendance.UserID, expectedAttendance.ScheduleID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(expectedAttendance.ID))

	attendance, err := s.repo.Create(expectedAttendance.UserID, expectedAttendance.ScheduleID)

	s.Nil(err)
	s.NoError(err)
	s.Equal(expectedAttendance, attendance)
}

func (s *AttendanceRepositoryTestSuite) TestList() {
	s.mockSql.
		ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, schedule_id FROM attendances")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "schedule_id"}).
			AddRow(expectedAttendances[0].ID, expectedAttendances[0].UserID, expectedAttendances[0].ScheduleID).
			AddRow(expectedAttendances[1].ID, expectedAttendances[1].UserID, expectedAttendances[1].ScheduleID))

	attendances, err := s.repo.List()

	s.Nil(err)
	s.Equal(expectedAttendances, attendances)
}

func (s *AttendanceRepositoryTestSuite) TestGetByID() {
	s.mockSql.
		ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, schedule_id, created_at, updated_at, deleted_at FROM attendances WHERE user_id = $1 AND schedule_id = $2")).
		WithArgs(expectedAttendance.ID, expectedAttendance.ScheduleID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "schedule_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(expectedAttendanceByID.ID, expectedAttendanceByID.UserID, expectedAttendanceByID.ScheduleID, expectedAttendanceByID.CreatedAt, expectedAttendanceByID.UpdatedAt, expectedAttendanceByID.DeletedAt))

	attendance, err := s.repo.GetByID(expectedAttendance.UserID, expectedAttendance.ScheduleID)

	s.Nil(err)
	s.Equal(expectedAttendance, attendance)
}

func (s *AttendanceRepositoryTestSuite) TestGetAttendance() {
	s.mockSql.
		ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, schedule_id FROM attendances WHERE id = $1")).
		WithArgs(expectedAttendance.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "schedule_id"}).
			AddRow(expectedAttendanceByID.ID, expectedAttendanceByID.UserID, expectedAttendanceByID.ScheduleID))

	attendance, err := s.repo.GetAttendance(expectedAttendance.ID)

	s.Nil(err)
	s.Equal(expectedAttendance, attendance)
}

func (s *AttendanceRepositoryTestSuite) TestDeleteAttendance() {
	s.mockSql.
		ExpectExec(regexp.QuoteMeta("UPDATE attendances SET deleted_at = $1 WHERE id = $2")).
		WithArgs(customTimeMatcher(&currTime), expectedAttendance.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.DeleteAttendance(expectedAttendance.ID)

	s.Nil(err)
	s.NoError(err)
}

func TestAttendanceRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AttendanceRepositoryTestSuite))
}
