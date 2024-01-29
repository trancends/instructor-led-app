package controller

import (
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
	rg                   *gin.Engine
	QuestionsUsecaseMock *usecase_mock.QuestionUscaseMock
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
	s.QuestionsUsecaseMock = new(usecase_mock.QuestionUscaseMock)
	s.AuthMiddlewareMock = new(middleware_mock.AuthMiddlewareMock)
	s.rg = gin.Default()
	gin.SetMode(gin.TestMode)
	rg := s.rg.Group("/api/v1") // Fix the variable name
	rg.Use(s.AuthMiddlewareMock.RequireToken("ADMIN"))
	s.rg = rg
}
