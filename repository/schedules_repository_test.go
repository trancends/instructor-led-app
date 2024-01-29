package repository

import (
	"database/sql"
	"errors"
	"math"
	"regexp"
	"testing"
	"time"

	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
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
	expectedPaging = sharedmodel.Paging{
		Page:        1,
		RowsPerPage: 10,
		TotalRows:   10,
		TotalPages:  0,
	}
	currentTime  = time.Now().Local()
	expectedUser = model.User{
		ID:       "1",
		Name:     "test",
		Email:    "test@example.com",
		Password: "test",
		Role:     "PARTICIPANT",
	}

	expectedUserError = model.User{
		ID:       "1",
		Name:     "test",
		Password: "test",
		Role:     "PARTICIPANT",
	}

	expectedUserUpdate = model.User{
		ID:        "1",
		Name:      "UpdatedName",
		Email:     "updatedemail@example.com",
		Password:  "updatedpassword",
		UpdatedAt: &currentTime,
	}

	expectedUsers = []model.User{
		{ID: "1", Name: "User1", Email: "user1@example.com", Role: "PARTICIPANT"},
		{ID: "2", Name: "User2", Email: "user2@example.com", Role: "PARTICIPANT"},
		// Add more expected users as needed
	}
)

func (s *ScheduleRepositoryTestSuite) SetupSuite() {
	db, mock, _ := sqlmock.New()
	s.mockDB = db
	s.mockSql = mock
	s.repo = NewSchedulesRepository(s.mockDB)
}

func (s *ScheduleRepositoryTestSuite) TestListScheduleByRole() {
	rows := sqlmock.NewRows([]string{"id", "user_id", "date", "start_time", "end_time", "documentation"})
	for _, schedules := range expectedSchedules {
		rows.AddRow(schedules.ID, schedules.UserID, schedules.Date, schedules.StartTime, schedules.EndTime, schedules.Documentation)
	}

	// Sesuaikan query dengan klausul JOIN dan parameter yang benar
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT s.id, s.user_id, s.date, s.start_time, s.end_time, s.documentation FROM schedules s JOIN users u ON s.user_id = u.id WHERE u.role = $3 AND s.deleted_at IS NULL ORDER BY s.created_at DESC LIMIT $1 OFFSET $2`)).
		WithArgs(expectedPaging.RowsPerPage, expectedPaging.RowsPerPage*(expectedPaging.Page-1), expectedUser.Role).
		WillReturnRows(rows)

	// Sesuaikan hasil paging dengan total baris yang sesuai
	expectedPaging.TotalRows = len(expectedSchedules)
	expectedPaging.TotalPages = int(math.Ceil(float64(expectedPaging.TotalRows) / float64(expectedPaging.RowsPerPage)))

	// Panggil ListScheduleByRole dengan parameter yang benar
	schedule, paging, err := s.repo.ListScheduleByRole(expectedPaging.RowsPerPage, expectedPaging.RowsPerPage*(expectedPaging.Page-1), expectedUser.Role)

	s.NoError(err)
	s.Equal(expectedSchedules, schedule)
	s.Equal(expectedPaging, paging)
}

func (s *ScheduleRepositoryTestSuite) TestListScheduled_succes() {}
func (s *ScheduleRepositoryTestSuite) TestCreateScheduled() {
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`INSERT INTO schedules (user_id, date, start_time, end_time, documentation) VALUES ($1, $2, $3, $4, $5) RETURNING id,user_id, date, start_time, end_time, documentation`)).WithArgs(expectedSchedule.UserID, expectedSchedule.Date, expectedSchedule.StartTime, expectedSchedule.EndTime).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "date", "start_time", "end_time", "documentation"}).AddRow(expectedSchedule.ID, expectedSchedule.UserID, expectedSchedule.Date, expectedSchedule.StartTime, expectedSchedule.EndTime, expectedSchedule.Documentation))

	_, err := s.repo.CreateScheduled(expectedSchedule)
	s.Error(err)
	s.NotNil(err)
}

func (s *ScheduleRepositoryTestSuite) TestListScheduled_failed() {
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`INSERT INTO schedules (user_id, date, start_time, end_time, documentation) VALUES ($1, $2, $3, $4, $5) RETURNING id,user_id, date, start_time, end_time, documentation`)).WillReturnError(errors.New("arguments do not match: expected 3, but got 4 arguments"))

	_, err := s.repo.CreateScheduled(expectedSchedule)
	s.Error(err)
	s.NotNil(err)
}
func (s *ScheduleRepositoryTestSuite) TestGetByID_succes() {
	rows := sqlmock.NewRows([]string{"id", "user_id", "date", "start_time", "end_time", "documentation"}).AddRow(expectedSchedule.ID, expectedSchedule.UserID, expectedSchedule.Date, expectedSchedule.StartTime, expectedSchedule.EndTime, expectedSchedule.Documentation)

	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, date, start_time, end_time, documentation FROM schedules WHERE id = $1`)).WithArgs(expectedSchedule.ID).WillReturnRows(rows)

	schedules, err := s.repo.GetByID(expectedSchedule.ID)
	s.Nil(err)
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

	// Ensure all expectations were met
	err = s.mockSql.ExpectationsWereMet()
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

	// Ensure all expectations were met
	err = s.mockSql.ExpectationsWereMet()
	s.NoError(err)
}

func TestScheduleRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ScheduleRepositoryTestSuite))
}
