package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"enigmaCamp.com/instructor_led/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type ScheduleRepositoryTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    ScheduleRepository
}

var (
	expectedSchedules = []model.Schedule{
		{
			ID:            "1",
			UserID:        "test-user-id",
			Date:          "2022-01-01",
			StartTime:     "08:00",
			EndTime:       "09:00",
			Documentation: "Documentation",
			CreatedAt:     time.Now(),
			UpdatedAt:     &time.Time{},
			DeletedAt:     &time.Time{},
			Questions: []model.Question{
				{
					ID:          "1",
					UserID:      "test-user-id",
					ScheduleID:  "1",
					Description: "Test",
					Status:      "PROCESS",
					CreatedAt:   &time.Time{},
					UpdatedAt:   &time.Time{},
					DeletedAt:   &time.Time{},
				},
			},
		},
	}
	expectedSchedule = model.Schedule{
		ID:            "1",
		UserID:        "test-user-id",
		Date:          "2022-01-01",
		StartTime:     "08:00",
		EndTime:       "09:00",
		Documentation: "Documentation",
		CreatedAt:     time.Now(),
		UpdatedAt:     &time.Time{},
		DeletedAt:     &time.Time{},
		Questions: []model.Question{
			{
				ID:          "1",
				UserID:      "test-user-id",
				ScheduleID:  "1",
				Description: "Test",
				Status:      "PROCESS",
				CreatedAt:   &time.Time{},
				UpdatedAt:   &time.Time{},
				DeletedAt:   &time.Time{},
			},
		},
	}
)

func (s *ScheduleRepositoryTestSuite) SetupSuite() {
	db, mock, _ := sqlmock.New()
	s.mockDB = db
	s.mockSql = mock
	s.repo = NewSchedulesRepository(s.mockDB)
}

// func (s *ScheduleRepositoryTestSuite) TestCreateScheduled() {
// 	s.mockSql.ExpectQuery(regexp.QuoteMeta(`INSERT INTO schedules (user_id, date, start_time, end_time, documentation) VALUES ($1, $2, $3, $4, $5) RETURNING id,user_id, date, start_time, end_time, documentation`)).WithArgs(expectedSchedule.UserID, expectedSchedule.Date, expectedSchedule.StartTime, expectedSchedule.EndTime).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "date", "start_time", "end_time", "documentation"}).AddRow(expectedSchedule.ID, expectedSchedule.UserID, expectedSchedule.Date, expectedSchedule.StartTime, expectedSchedule.EndTime, expectedSchedule.Documentation))

// 	_, err := s.repo.CreateScheduled(expectedSchedule)
// 	s.Error(err)
// 	s.NotNil(err)
// }

func (s *ScheduleRepositoryTestSuite) TestListScheduled_failed() {
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`INSERT INTO schedules (user_id, date, start_time, end_time, documentation) VALUES ($1, $2, $3, $4, $5) RETURNING id,user_id, date, start_time, end_time, documentation`)).WillReturnError(errors.New("arguments do not match: expected 3, but got 4 arguments"))

	_, err := s.repo.CreateScheduled(expectedSchedule)
	s.Error(err)
	s.NotNil(err)
}

func (s *ScheduleRepositoryTestSuite) TestGetByID_succes() {
	rows := sqlmock.NewRows([]string{"id", "user_id", "date", "start_time", "end_time", "documentation"}).AddRow(expectedSchedule.ID, expectedSchedule.UserID, expectedSchedule.Date, expectedSchedule.StartTime, expectedSchedule.EndTime, expectedSchedule.Documentation)

	s.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, date, start_time, end_time, documentation FROM schedules WHERE id = $1")).WithArgs(expectedSchedule.ID).WillReturnRows(rows)

	schedules, _ := s.repo.GetByID(expectedSchedule.ID)
	// s.Nil(err)
	s.Equal(expectedSchedule.ID, schedules.ID)
}

func (s *ScheduleRepositoryTestSuite) TestGetByID_failed() {
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, date, start_time, end_time, documentation FROM schedules WHERE id = $1`)).WithArgs(expectedSchedule.ID).WillReturnError(errors.New("sql: no rows in result set"))

	_, err := s.repo.GetByID(expectedSchedule.ID)
	s.Error(err)
	s.NotNil(err)
}

func (s *ScheduleRepositoryTestSuite) TestDelete_success() {
	// Prepare the expected SQL query
	query := "UPDATE schedules SET deleted_at = $1 WHERE id = $2"

	// Expectation: mock query execution
	s.mockSql.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), expectedSchedule.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the actual method
	err := s.repo.Delete(expectedSchedule.ID)

	// Assertions
	s.NoError(err)
}

func (s *ScheduleRepositoryTestSuite) TestUpdateDocumentation() {
	// Prepare the expected SQL query
	query := "UPDATE schedules SET documentation = $1, updated_at = $2 WHERE id = $3"

	// Expectation: mock query execution
	s.mockSql.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), expectedSchedule.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the actual method
	err := s.repo.UpdateDocumentation(expectedSchedule.ID, "http://example.com/picture.jpg")

	// Assertions
	s.NoError(err)
}

func TestScheduleRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ScheduleRepositoryTestSuite))
}
