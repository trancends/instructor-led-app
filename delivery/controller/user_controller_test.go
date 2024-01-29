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

type UserControllerTestSuite struct {
	suite.Suite
	rg                 *gin.RouterGroup
	UserUsecaseMock    *usecase_mock.UserUsecaseMock
	AuthMiddlewareMock *middleware_mock.AuthMiddlewareMock
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

	expectedPaging = sharedmodel.Paging{
		Page:        1,
		RowsPerPage: 10,
		TotalRows:   len(expectedUsers),
		TotalPages:  1,
	}
)

func (s *UserControllerTestSuite) SetupTest() {
	s.UserUsecaseMock = new(usecase_mock.UserUsecaseMock)
	s.AuthMiddlewareMock = new(middleware_mock.AuthMiddlewareMock)
	r := gin.Default()
	gin.SetMode(gin.TestMode)
	rg := r.Group("/api/v1")
	rg.Use(s.AuthMiddlewareMock.RequireToken("ADMIN"))
	s.rg = rg
}

func (s *UserControllerTestSuite) TestCreateUserHandler() {
	s.UserUsecaseMock.On("CreateUser", expectedUser).Return(expectedUser, nil)

	userController := NewUserController(s.UserUsecaseMock, s.rg, s.AuthMiddlewareMock)
	userController.Route()

	req, err := http.NewRequest("POST", "/api/v1/users", nil)
	s.NoError(err)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("user", expectedUser.ID)
	s.Equal(http.StatusOK, record.Code)
}

func (s *UserControllerTestSuite) TestGetAllUserHandler() {
	s.UserUsecaseMock.On("GetAllUser", expectedPaging).Return(expectedUsers, nil)

	userController := NewUserController(s.UserUsecaseMock, s.rg, s.AuthMiddlewareMock)
	userController.Route()

	req, err := http.NewRequest("GET", "/api/v1/users", nil)
	s.NoError(err)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("user", expectedUsers)

	s.Equal(http.StatusOK, record.Code)
}

func (s *UserControllerTestSuite) TestGetUserByEmailHandler() {
	s.UserUsecaseMock.On("GetUserByEmail", expectedUser.Email).Return(expectedUser, nil)

	userController := NewUserController(s.UserUsecaseMock, s.rg, s.AuthMiddlewareMock)
	userController.Route()

	req, err := http.NewRequest("GET", "/api/v1/users/test@example.com", nil)
	s.NoError(err)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("user", expectedUser)

	s.Equal(http.StatusOK, record.Code)
}

func (s *UserControllerTestSuite) TestUpdateUserHandler() {
	s.UserUsecaseMock.On("UpdateUser", expectedUserUpdate).Return(expectedUserUpdate, nil)

	userController := NewUserController(s.UserUsecaseMock, s.rg, s.AuthMiddlewareMock)
	userController.Route()

	req, err := http.NewRequest("PUT", "/api/v1/users/1", nil)
	s.NoError(err)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("user", expectedUserUpdate)

	s.Equal(http.StatusOK, record.Code)
}

func (s *UserControllerTestSuite) TestDeleteUserHandler() {
	s.UserUsecaseMock.On("DeleteUser", expectedUser.ID).Return(nil)

	userController := NewUserController(s.UserUsecaseMock, s.rg, s.AuthMiddlewareMock)
	userController.Route()

	req, err := http.NewRequest("DELETE", "/api/v1/users/1", nil)
	s.NoError(err)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("user", expectedUser)

	s.Equal(http.StatusOK, record.Code)
}

func (s *UserControllerTestSuite) TestCreateUserCSVHandler() {
	s.UserUsecaseMock.On("CreateUserCSV", expectedUsers).Return(nil)

	userController := NewUserController(s.UserUsecaseMock, s.rg, s.AuthMiddlewareMock)
	userController.Route()

	req, err := http.NewRequest("POST", "/api/v1/users/csv", nil)
	s.NoError(err)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = req
	ctx.Set("users", expectedUsers)

	s.Equal(http.StatusOK, record.Code)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
