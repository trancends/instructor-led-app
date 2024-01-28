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

type UserRepositoryTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    UserRepository
}

var currentTime = time.Now().Local()

var expectedUser = model.User{
	ID:       "1",
	Name:     "test",
	Email:    "test@example.com",
	Password: "test",
	Role:     "PARTICIPANT",
}

var expectedUserError = model.User{
	ID:       "1",
	Name:     "test",
	Password: "test",
	Role:     "PARTICIPANT",
}

func (s *UserRepositoryTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New()
	s.mockDB = db
	s.mockSql = mock
	s.repo = NewUserRepository(s.mockDB)
}

func (s *UserRepositoryTestSuite) TestCreate_Success() {
	s.mockSql.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (name,email,password,role) VALUES ($1,$2,$3,$4) RETURNING id`)).
		WithArgs(expectedUser.Name, expectedUser.Email,
			expectedUser.Password, expectedUser.Role).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(expectedUser.ID))

	err := s.repo.Create(expectedUser)
	s.Nil(err)
	s.NoError(err)
}

func (s *UserRepositoryTestSuite) TestCreate_Failed() {
	s.mockSql.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (name,email,password,role) VALUES ($1,$2,$3,$4) RETURNING id`)).
		WithArgs(expectedUser.Name, expectedUser.Email,
			expectedUser.Password).
		WillReturnError(errors.New("arguments do not match: expected 3, but got 4 arguments"))
	err := s.repo.Create(expectedUserError)
	s.NotNil(err)
	s.Error(err)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
