package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"enigmaCamp.com/instructor_led/mock/middleware_mock"
	"enigmaCamp.com/instructor_led/mock/usecase_mock"
	"enigmaCamp.com/instructor_led/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AttendanceControllerTestSuite struct {
	suite.Suite
	rg                    *gin.RouterGroup
	AttendanceUsecaseMock *usecase_mock.AttendanceUsecaseMock
	AuthMiddlewareMock    *middleware_mock.AuthMiddlewareMock
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
	newUUID                  = uuid.New()
	expectedAttendanceDelete = model.Attendance{
		ID:         newUUID.String(),
		UserID:     "1",
		ScheduleID: "1",
		DeletedAt:  &currTime,
	}
)

func (s *AttendanceControllerTestSuite) SetupTest() {
	s.AttendanceUsecaseMock = new(usecase_mock.AttendanceUsecaseMock)
	s.AuthMiddlewareMock = new(middleware_mock.AuthMiddlewareMock)
	r := gin.Default()
	gin.SetMode(gin.TestMode)
	rg := r.Group("/api/v1")
	rg.Use(s.AuthMiddlewareMock.RequireToken("ADMIN", "TRAINER", "PARTICIPANT"))
	s.rg = rg
}

func (s *AttendanceControllerTestSuite) TestGetAllAttendanceHandler() {
	s.AttendanceUsecaseMock.On("ListAttendance").Return(expectedAttendances, nil)

	attendanceController := NewAttandanceController(s.rg, s.AttendanceUsecaseMock, s.AuthMiddlewareMock)
	attendanceController.Route()

	req, err := http.NewRequest("GET", "/api/v1/attendances", nil)
	s.NoError(err)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("attendances", expectedAttendances)

	s.Equal(http.StatusOK, record.Code)
}

func (s *AttendanceControllerTestSuite) TestGetAttendanceByIDHandler() {
	s.AttendanceUsecaseMock.On("GetAttendanceByID", expectedAttendance.ID).Return(expectedAttendanceByID, nil)

	attendanceController := NewAttandanceController(s.rg, s.AttendanceUsecaseMock, s.AuthMiddlewareMock)
	attendanceController.Route()

	req, err := http.NewRequest("GET", "/api/v1/attendances/1", nil)
	s.NoError(err)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("attendance", expectedAttendanceByID)

	s.Equal(http.StatusOK, record.Code)
}

func (s *AttendanceControllerTestSuite) TestAddAttendanceHandler() {
	s.AttendanceUsecaseMock.On("AddAttendance", expectedAttendance.ID, expectedAttendance.ScheduleID).Return(expectedAttendance, nil)

	attendanceController := NewAttandanceController(s.rg, s.AttendanceUsecaseMock, s.AuthMiddlewareMock)
	attendanceController.Route()

	req, err := http.NewRequest("POST", "/api/v1/attendances", nil)
	s.NoError(err)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("attendance", expectedAttendance)

	s.Equal(http.StatusOK, record.Code)
}

func (s *AttendanceControllerTestSuite) TestDeleteAttendanceHandler() {
	s.AttendanceUsecaseMock.On("DeleteAttendance", expectedAttendanceDelete.ID).Return(expectedAttendanceDelete, nil)

	attendanceController := NewAttandanceController(s.rg, s.AttendanceUsecaseMock, s.AuthMiddlewareMock)
	attendanceController.Route()

	req, err := http.NewRequest("DELETE", "/api/v1/attendances/1", nil)
	s.NoError(err)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("attendance", expectedAttendanceDelete)

	s.Equal(http.StatusOK, record.Code)
}

func TestAttendanceControllerSuite(t *testing.T) {
	suite.Run(t, new(AttendanceControllerTestSuite))
}
