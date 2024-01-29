package usecase

import (
	"errors"
	"testing"
	"time"

	"enigmaCamp.com/instructor_led/mock/repository_mock"
	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
	"github.com/stretchr/testify/suite"
)

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

type UserUseCaseTestSuite struct {
	suite.Suite
	userRepoMock *repository_mock.UserRepositoryMock
	userUseCase  UserUsecase
}

func (s *UserUseCaseTestSuite) TestCreateUser() {
	s.userRepoMock.On("Create", expectedUser).Return(nil)

	err := s.userUseCase.CreateUser(expectedUser)

	s.NoError(err)
	s.Nil(err)
}

func (s *UserUseCaseTestSuite) TestCreateUser_Error() {
	s.userRepoMock.On("Create", expectedUserError).Return(errors.New("error creating user"))

	err := s.userUseCase.CreateUser(expectedUserError)

	s.Error(err)
	s.NotNil(err)
}

func (s *UserUseCaseTestSuite) TestListAllUsers() {
	s.userRepoMock.On("List", 1, 10).Return(expectedUsers, expectedPaging, nil)

	users, paging, err := s.userUseCase.ListAllUsers(1, 10)

	s.NoError(err)
	s.Equal(expectedUsers, users)
	s.Equal(expectedPaging, paging)
}

func (s *UserUseCaseTestSuite) TestListAllUsers_Error() {
	s.userRepoMock.On("List", 1, 10).Return([]model.User{}, expectedPaging, errors.New("error listing users"))

	users, paging, err := s.userUseCase.ListAllUsers(1, 10)

	s.Error(err)
	s.NotNil(err)
	s.Nil(users)
	s.Zero(paging)
}

func (s *UserUseCaseTestSuite) TestGetUserByEmail() {
	s.userRepoMock.On("GetUserByEmail", "test@example.com").Return(expectedUser, nil)

	user, err := s.userUseCase.GetUserByEmail("test@example.com")

	s.NoError(err)
	s.Equal(expectedUser, user)
}

func (s *UserUseCaseTestSuite) TestGetUserByEmail_Error() {
	s.userRepoMock.On("GetUserByEmail", "test@example.com").Return(model.User{}, errors.New("error getting user"))

	user, err := s.userUseCase.GetUserByEmail("test@example.com")

	s.Error(err)
	s.Equal(model.User{}, user)
}

func (s *UserUseCaseTestSuite) TestGetUserByID() {
	s.userRepoMock.On("GetUserByID", "1").Return(expectedUser, nil)

	user, err := s.userUseCase.GetUserByID("1")

	s.NoError(err)
	s.Equal(expectedUser, user)
}

func (s *UserUseCaseTestSuite) TestGetUserByID_Error() {
	s.userRepoMock.On("GetUserByID", "1").Return(model.User{}, errors.New("error getting user"))

	user, err := s.userUseCase.GetUserByID("1")

	s.Error(err)
	s.Equal(model.User{}, user)
}

func (s *UserUseCaseTestSuite) TestGetUserByRole() {
	s.userRepoMock.On("GetUserByRole", "PARTICIPANT", 1, 10).Return(expectedUsers, expectedPaging, nil)

	users, paging, err := s.userUseCase.GetUserByRole("PARTICIPANT", 1, 10)

	s.NoError(err)
	s.Equal(expectedUsers, users)
	s.Equal(expectedPaging, paging)
}

func (s *UserUseCaseTestSuite) TestGetUserByRole_Error() {
	s.userRepoMock.On("GetUserByRole", "PARTICIPANT", 1, 10).Return([]model.User{}, expectedPaging, errors.New("error getting users"))

	users, paging, err := s.userUseCase.GetUserByRole("PARTICIPANT", 1, 10)

	s.Error(err)
	s.Nil(users)
	s.Zero(paging)
}

func (s *UserUseCaseTestSuite) TestUpdateUser() {
	s.userRepoMock.On("Update", expectedUser).Return(nil)

	err := s.userUseCase.UpdateUser(expectedUser)

	s.NoError(err)
}

func (s *UserUseCaseTestSuite) TestUpdateUser_Error() {
	s.userRepoMock.On("Update", expectedUserError).Return(errors.New("error updating user"))

	err := s.userUseCase.UpdateUser(expectedUserError)

	s.Error(err)
}

func (s *UserUseCaseTestSuite) TestDeleteUser() {
	s.userRepoMock.On("Delete", "1").Return(nil)

	err := s.userUseCase.DeleteUser("1")

	s.NoError(err)
}

func (s *UserUseCaseTestSuite) TestDeleteUser_Error() {
	s.userRepoMock.On("Delete", "1").Return(errors.New("error deleting user"))

	err := s.userUseCase.DeleteUser("1")

	s.Error(err)
}

func (s *UserUseCaseTestSuite) SetupTest() {
	s.userRepoMock = new(repository_mock.UserRepositoryMock)
	s.userUseCase = NewUserUsecase(s.userRepoMock)
}

func TestUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseTestSuite))
}
