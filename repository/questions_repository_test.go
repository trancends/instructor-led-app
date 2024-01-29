package repository

import (
	"database/sql"
	"testing"
	"time"

	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type QuestionsRepositoryTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    QuestionsRepository
}

var (
	currentTime      = time.Now()
	expectedQuestion = model.Question{
		ID:          "1",
		UserID:      "1",
		ScheduleID:  "1",
		Description: "test",
		Status:      "PROCESS",
		CreatedAt:   nil,
		UpdatedAt:   nil,
		DeletedAt:   nil,
	}
	expectedQuestionsError = model.Question{
		ID:          "1",
		ScheduleID:  "1",
		Description: "test",
		Status:      "PROCESS",
		CreatedAt:   nil,
		UpdatedAt:   nil,
		DeletedAt:   nil,
	}
	expectedQuestionsUpdate = model.Question{
		ID:          "1",
		UserID:      "1",
		ScheduleID:  "1",
		Description: "test",
		Status:      "FINISHED",
		CreatedAt:   &currentTime,
		UpdatedAt:   &currentTime,
		DeletedAt:   nil,
	}
	expectedQuestions = []model.Question{
		{
			ID:          "1",
			UserID:      "1",
			ScheduleID:  "1",
			Description: "test",
			Status:      "PROCESS",
			CreatedAt:   nil,
			UpdatedAt:   nil,
			DeletedAt:   nil,
		},
		{
			ID:          "2",
			UserID:      "2",
			ScheduleID:  "1",
			Description: "test",
			Status:      "FINISHED",
			CreatedAt:   nil,
			UpdatedAt:   nil,
			DeletedAt:   nil,
		},
	}

	expectedPaging = sharedmodel.Paging{
		Page:        1,
		RowsPerPage: 10,
		TotalRows:   len(expectedQuestions),
		TotalPages:  0,
	}
)

func (s *QuestionsRepositoryTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New()
	s.mockDB = db
	s.mockSql = mock
	s.repo = NewQuestionsRepository(s.mockDB)
}

func (s *QuestionsRepositoryTestSuite) TestGet() {

	s.mockSql.ExpectQuery("SELECT * FROM questions WHERE date = ?").
		WithArgs("2022-01-01").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "schedule_id", "description", "status", "created_at", "updated_at", "deleted_at"}).
				AddRow("1", "1", "1", "test", "PROCESS", nil, nil, nil).
				AddRow("2", "2", "1", "test", "FINISHED", nil, nil, nil),
		)
}
func (s *QuestionsRepositoryTestSuite) TestUpdate() {
	s.mockSql.ExpectQuery("UPDATE questions SET status = ? WHERE id = ?").
		WithArgs("FINISHED", "1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "schedule_id", "description", "status", "created_at", "updated_at", "deleted_at"}).AddRow("1", "1", "1", "test", "PROCESS", nil, nil, nil))
}
func (s *QuestionsRepositoryTestSuite) TestDelete() {
	s.mockSql.ExpectQuery("DELETE FROM questions WHERE id = ?").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "schedule_id", "description", "status", "created_at", "updated_at", "deleted_at"}).AddRow("1", "1", "1", "test", "PROCESS", nil, nil, nil))
}
func TestQuestionsRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(QuestionsRepositoryTestSuite))
}
