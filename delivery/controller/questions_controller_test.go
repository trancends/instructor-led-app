package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"enigmaCamp.com/instructor_led/mock/middleware_mock"
	"enigmaCamp.com/instructor_led/mock/usecase_mock"
	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type QuestionsControllerSuite struct {
	suite.Suite
	rg                   *gin.RouterGroup
	QuestionsUsecaseMock *usecase_mock.QuestionUscaseMock // Fix the typo here
	AuthMiddlewareMock   *middleware_mock.AuthMiddlewareMock
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

func (s *QuestionsControllerSuite) SetupTest() {
	s.QuestionsUsecaseMock = new(usecase_mock.QuestionUscaseMock) // Fix the typo here
	s.AuthMiddlewareMock = new(middleware_mock.AuthMiddlewareMock)
	gin.SetMode(gin.TestMode)
	engine := gin.Default()
	rg := engine.Group("/api/v1")
	rg.Use(s.AuthMiddlewareMock.RequireToken("ADMIN"))
	s.rg = rg
}

func (s *QuestionsControllerSuite) TestGetQuestions() {
	s.QuestionsUsecaseMock.On("ListQuestions", expectedQuestions).Return(expectedQuestions, nil)
	questionController := NewQuestionsController(s.QuestionsUsecaseMock, s.rg, s.AuthMiddlewareMock)
	questionController.Route()
	req, err := http.NewRequest("GET", "/api/v1/questions", nil)
	s.NoError(err)
	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("expectedQuestion", expectedQuestions[0].ID)
	s.Equal(http.StatusOK, record.Code)
}

func (s *QuestionsControllerSuite) TestUpdateQuestion() {
	s.QuestionsUsecaseMock.On("UpdateQuestions", expectedQuestionsUpdate).Return(expectedQuestionsUpdate, nil)
	questionController := NewQuestionsController(s.QuestionsUsecaseMock, s.rg, s.AuthMiddlewareMock)
	questionController.Route()
	req, err := http.NewRequest("PUT", "/api/v1/questions", nil)
	s.NoError(err)
	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("expectedQuestion", expectedQuestionsUpdate.ID)
	s.Equal(http.StatusOK, record.Code)
}
func (s *QuestionsControllerSuite) TestDeleteQuestion() {
	s.QuestionsUsecaseMock.On("DeleteQuestions", expectedQuestionsUpdate).Return(expectedQuestionsUpdate, nil)
	questionController := NewQuestionsController(s.QuestionsUsecaseMock, s.rg, s.AuthMiddlewareMock)
	questionController.Route()
	req, err := http.NewRequest("DELETE", "/api/v1/questions", nil)
	s.NoError(err)
	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("expectedQuestion", expectedQuestionsUpdate.ID)
	s.Equal(http.StatusOK, record.Code)
}
func (s *QuestionsControllerSuite) TestCreateQuestion() {
	s.QuestionsUsecaseMock.On("CreateQuestions", expectedQuestion).Return(expectedQuestion, nil)
	questionController := NewQuestionsController(s.QuestionsUsecaseMock, s.rg, s.AuthMiddlewareMock)
	questionController.Route()
	req, err := http.NewRequest("POST", "/api/v1/questions", nil)
	s.NoError(err)
	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	s.Equal(http.StatusOK, record.Code)
}
func TestQustionsControllerSuite(t *testing.T) {
	suite.Run(t, new(QuestionsControllerSuite))
}
