package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"enigmaCamp.com/instructor_led/mock/usecase_mock"
	"enigmaCamp.com/instructor_led/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// var mockSchedules = []model.Schedule{
// 	{
// 		ID:            "1",
// 		UserID:        "test-user-id",
// 		Date:          "2022-01-01",
// 		StartTime:     "08:00",
// 		EndTime:       "09:00",
// 		Documentation: "Documentation",
// 		CreatedAt:     time.Now(),
// 		UpdatedAt:     &time.Time{},
// 		DeletedAt:     &time.Time{},
// 		Questions:     []model.Question{},
// 	},
// }

type ScheduleControllerTestSuite struct {
	suite.Suite
	rg  *gin.RouterGroup
	sum *usecase_mock.SchedulesUseCaseMock
}

func (s *ScheduleControllerTestSuite) SetupTest() {
	s.sum = new(usecase_mock.SchedulesUseCaseMock)
	r := gin.Default()
	gin.SetMode(gin.TestMode)
	rg := r.Group("/api/v1")
	s.rg = rg
}

func (s *ScheduleControllerTestSuite) TestCreateScheduleHandler_success() {
	payloadSuccess := model.Schedule{
		UserID:        "test-user-id",
		Date:          "2022-01-01",
		StartTime:     "08:00",
		EndTime:       "09:00",
		Documentation: "Documentation",
	}
	s.sum.On("CreateScheduledUC", payloadSuccess).Return(model.Schedule{ID: "1", Date: "2022-01-01", StartTime: "08:00", EndTime: "09:00", Documentation: "Documentation"}, nil)

	ScheduleController := NewSchedulesController(s.sum, s.rg, nil)
	ScheduleController.Route()

	// Convert payload to JSON and set it as the request body
	requestBody, err := json.Marshal(payloadSuccess)
	s.NoError(err)

	req, err := http.NewRequest("POST", "/api/v1/schedules", bytes.NewBuffer(requestBody))
	s.NoError(err)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	ScheduleController.CreateScheduleHandler(ctx)
	s.Equal(http.StatusOK, record.Code)
}

func (s *ScheduleControllerTestSuite) TestCreateScheduleHandler_fail() {
	// Create an expectation with a valid model.Schedule object
	s.sum.On("CreateScheduledUC", mock.AnythingOfType("model.Schedule")).Return(model.Schedule{}, errors.New("some error"))

	ScheduleController := NewSchedulesController(s.sum, s.rg, nil)
	ScheduleController.Route()

	req, err := http.NewRequest("POST", "/api/v1/schedules", nil)
	s.NoError(err)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req

	ScheduleController.CreateScheduleHandler(ctx)
	s.Equal(http.StatusBadRequest, record.Code)
}

func TestScheduleControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ScheduleControllerTestSuite))
}
