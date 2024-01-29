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

var expectedSchedules = []model.Schedule{
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
var expectedPaging = sharedmodel.Paging{
	Page:        1,
	RowsPerPage: 10,
	TotalRows:   len(expectedSchedules),
	TotalPages:  0,
}

var (
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

type ScheduleControllerTestSuite struct {
	suite.Suite
	rg                 *gin.RouterGroup
	sum                *usecase_mock.SchedulesUseCaseMock
	AuthMiddlewareMock *middleware_mock.AuthMiddlewareMock
}

func (s *ScheduleControllerTestSuite) SetupTest() {
	s.sum = new(usecase_mock.SchedulesUseCaseMock)
	s.AuthMiddlewareMock = new(middleware_mock.AuthMiddlewareMock)
	r := gin.Default()
	gin.SetMode(gin.TestMode)
	rg := r.Group("/api/v1")
	rg.Use(s.AuthMiddlewareMock.RequireToken("ADMIN", "TRAINER", "PARTICIPANT"))
	s.rg = rg
}

func (s *ScheduleControllerTestSuite) TestCreateScheduleHandler_success() {
	s.sum.Mock.On("CreateScheduledUC", expectedSchedules).Return(expectedSchedules, nil)
	schedulesController := NewSchedulesController(s.sum, s.rg, s.AuthMiddlewareMock)
	schedulesController.Route()
	req, err := http.NewRequest("POST", "/api/v1/schedules", nil)
	s.NoError(err)
	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("user", expectedSchedules[0].UserID)
	s.Equal(http.StatusOK, record.Code)
}

// func (s *ScheduleControllerTestSuite) TestCreateScheduleHandler_fail() {
// 	s.sum.Mock.On("CreateScheduledUC", expectedSchedules).Return(nil, errors.New("error"))
// 	schedulesController := NewSchedulesController(s.sum, s.rg, s.AuthMiddlewareMock)
// 	schedulesController.Route()
// 	req, err := http.NewRequest("POST", "/api/v1/schedules", nil)
// 	s.NoError(err)
// 	record := httptest.NewRecorder()
// 	ctx, _ := gin.CreateTestContext(record)
// 	ctx.Request = req
// 	ctx.Set("user", expectedSchedules[0].UserID)
// 	s.Equal(http.StatusBadRequest, record.Code)
// }

func (s *ScheduleControllerTestSuite) TestFindAllScheduleHandler_success() {
	s.sum.Mock.On("FindAllScheduleUC", expectedPaging.Page, expectedPaging.RowsPerPage).Return(expectedSchedules, expectedPaging, nil)

	schedulesController := NewSchedulesController(s.sum, s.rg, s.AuthMiddlewareMock)
	schedulesController.Route()

	req, err := http.NewRequest("GET", "/api/v1/schedules?page=1&size=10", nil)
	s.NoError(err)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	schedulesController.FindAllScheduleHandler(ctx)

	s.Equal(http.StatusOK, record.Code)
}

// func (s *ScheduleControllerTestSuite) TestFindAllScheduleHandler_fail() {
// 	s.sum.Mock.On("FindAllScheduleUC", mockSchedules).Return(nil, errors.New("error"))
// 	schedulesController := NewSchedulesController(s.sum, s.rg, s.AuthMiddlewareMock)
// 	schedulesController.Route()
// 	req, err := http.NewRequest("GET", "/api/v1/schedules?page=1&size=10", nil)
// 	s.NoError(err)
// 	record := httptest.NewRecorder()
// 	ctx, _ := gin.CreateTestContext(record)
// 	ctx.Request = req
// 	schedulesController.FindAllScheduleHandler(ctx)
// 	s.Equal(http.StatusBadRequest, record.Code)
// }

func (s *ScheduleControllerTestSuite) TestFindSchedulesByRoleHandler_success() {
	s.sum.Mock.On("FindScheduleByRole", expectedPaging.Page, expectedPaging.RowsPerPage, expectedUsers[0].Role).Return(expectedSchedules, expectedPaging, nil)
	schedulesController := NewSchedulesController(s.sum, s.rg, s.AuthMiddlewareMock)
	schedulesController.Route()
	req, err := http.NewRequest("GET", "/api/v1/schedules/role?role=PARTICIPANT&size=10&page=1", nil) // Provide valid role, size, and page parameters
	s.NoError(err)
	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	schedulesController.FindSchedulesByRoleHandler(ctx)
	s.Equal(http.StatusOK, record.Code) // Change to the expected success status code
}

//	func (s *ScheduleControllerTestSuite) TestFindSchedulesByRoleHandler_fail() {
//		s.sum.Mock.On("FindSchedulesByRoleUC", mockSchedules).Return(nil, errors.New("error"))
//		SchedulesController := NewSchedulesController(s.sum, s.rg, s.AuthMiddlewareMock)
//		SchedulesController.Route()
//		req, err := http.NewRequest("GET", "/api/v1/schedules/role", nil)
//		s.NoError(err)
//		record := httptest.NewRecorder()
//		ctx, _ := gin.CreateTestContext(record)
//		ctx.Request = req
//		SchedulesController.FindSchedulesByRoleHandler(ctx)
//		s.Equal(http.StatusBadRequest, record.Code)
//	}
func (s *ScheduleControllerTestSuite) TestFindByIDScheduleHandler_success() {
	s.sum.Mock.On("FindByIDUC", "1").Return(expectedSchedules, nil)
	SchedulesController := NewSchedulesController(s.sum, s.rg, s.AuthMiddlewareMock)
	SchedulesController.Route()
	req, err := http.NewRequest("GET", "/api/v1/schedules/1", nil)
	s.NoError(err)
	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	SchedulesController.FindByIDScheduleHandler(ctx)
	// Add assertions here
}

// func (s *ScheduleControllerTestSuite) TestFindByIDScheduleHandler_fail() {
// 	s.sum.Mock.On("FindByIDSchedulesUC", mockSchedules).Return(nil, errors.New("error"))
// 	SchedulesController := NewSchedulesController(s.sum, s.rg, s.AuthMiddlewareMock)
// 	SchedulesController.Route()
// 	req, err := http.NewRequest("GET", "/api/v1/schedules/1", nil)
// 	s.NoError(err)
// 	record := httptest.NewRecorder()
// 	ctx, _ := gin.CreateTestContext(record)
// 	ctx.Request = req
// 	SchedulesController.FindByIDScheduleHandler(ctx)
// 	s.Equal(http.StatusBadRequest, record.Code)
// }

func (s *ScheduleControllerTestSuite) TestDeleteScheduleHandler_success() {
	s.sum.Mock.On("FindByIDUC", "1").Return(expectedSchedules, nil)
	s.sum.Mock.On("DeletedScheduleIDUC", "1").Return(nil)
	SchedulesController := NewSchedulesController(s.sum, s.rg, s.AuthMiddlewareMock)
	SchedulesController.Route()
	req, err := http.NewRequest("DELETE", "/api/v1/schedules/1", nil)
	s.NoError(err)
	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	SchedulesController.DeleteScheduleHandler(ctx)
	s.Equal(http.StatusOK, record.Code)
}

// func (s *ScheduleControllerTestSuite) TestDeleteScheduleHandler_fail() {
// 	s.sum.Mock.On("DeleteSchedulesUC", mockSchedules).Return(nil, errors.New("error"))
// 	SchedulesController := NewSchedulesController(s.sum, s.rg, s.AuthMiddlewareMock)
// 	SchedulesController.Route()
// 	req, err := http.NewRequest("DELETE", "/api/v1/schedules/1", nil)
// 	s.NoError(err)
// 	record := httptest.NewRecorder()
// 	ctx, _ := gin.CreateTestContext(record)
// 	ctx.Request = req
// 	SchedulesController.DeleteScheduleHandler(ctx)
// 	s.Equal(http.StatusBadRequest, record.Code)
// }

func (s *ScheduleControllerTestSuite) TestUpdateScheduleHandler_success() {
	s.sum.Mock.On("FindByIDUC", "1").Return(expectedSchedules, nil)
	s.sum.Mock.On("UpdateScheduleDocumentation", "1", "example.jpg").Return(nil)
	SchedulesController := NewSchedulesController(s.sum, s.rg, s.AuthMiddlewareMock)
	SchedulesController.Route()
	req, err := http.NewRequest("PATCH", "/api/v1/schedules/1", nil)
	s.NoError(err)
	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("picture", "example.jpg")
	SchedulesController.UploadDocumentationHandler(ctx)
	s.Equal(http.StatusOK, record.Code)
}

// func (s *ScheduleControllerTestSuite) TestUpdateScheduleHandler_fail() {
// 	s.sum.Mock.On("UpdateSchedulesUC", mockSchedules).Return(nil, errors.New("error"))
// 	SchedulesController := NewSchedulesController(s.sum, s.rg, s.AuthMiddlewareMock)
// 	SchedulesController.Route()
// 	req, err := http.NewRequest("PATCH", "/api/v1/schedules/1", nil)
// 	s.NoError(err)
// 	record := httptest.NewRecorder()
// 	ctx, _ := gin.CreateTestContext(record)
// 	ctx.Request = req
// 	SchedulesController.UploadDocumentationHandler(ctx)
// 	s.Equal(http.StatusBadRequest, record.Code)
// }

func TestScheduleControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ScheduleControllerTestSuite))
}
