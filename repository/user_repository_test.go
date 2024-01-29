package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    UserRepository
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

func (s *UserRepositoryTestSuite) TestGetUserByEmail_Success() {
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.Role)
	s.mockSql.
		ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, password, role FROM users WHERE email = $1 AND deleted_at IS NULL")).
		WithArgs(expectedUser.Email).
		WillReturnRows(rows)

	user, err := s.repo.GetUserByEmail(expectedUser.Email)
	s.Nil(err)
	s.Equal(expectedUser, user)
}

func (s *UserRepositoryTestSuite) TestGetUserByEmail_Failed() {
	s.mockSql.
		ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, password, role FROM users WHERE email = $1 AND deleted_at IS NULL")).
		WithArgs(expectedUser.Email).
		WillReturnError(errors.New("sql: no rows in result set"))
	_, err := s.repo.GetUserByEmail(expectedUser.Email)
	s.NotNil(err)
	s.Error(err)
}

func (s *UserRepositoryTestSuite) TestGetUserByID_Success() {
	rows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Role)
	s.mockSql.
		ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, role FROM users WHERE id = $1 AND deleted_at IS NULL")).
		WithArgs(expectedUser.ID).
		WillReturnRows(rows)

	user, err := s.repo.GetUserByID(expectedUser.ID)
	s.Nil(err)
	s.Equal(expectedUser.Email, user.Email)
}

func (s *UserRepositoryTestSuite) TestGetUserByID_Failed() {
	s.mockSql.
		ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, role FROM users WHERE id = $1 AND deleted_at IS NULL")).
		WithArgs(expectedUser.ID).
		WillReturnError(errors.New("sql: no rows in result set"))
	_, err := s.repo.GetUserByID(expectedUser.ID)
	s.NotNil(err)
	s.Error(err)
}

func (s *UserRepositoryTestSuite) TestListSuccess() {
	rows := sqlmock.NewRows([]string{"id", "name", "email", "role"})
	for _, user := range expectedUsers {
		rows.AddRow(user.ID, user.Name, user.Email, user.Role)
	}

	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, role FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`)).
		WithArgs(expectedPaging.RowsPerPage, (expectedPaging.Page-1)*expectedPaging.RowsPerPage).
		WillReturnRows(rows)

	s.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM users")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(expectedUsers)))

	users, paging, err := s.repo.List(expectedPaging.Page, expectedPaging.RowsPerPage)

	s.NoError(err)
	s.Equal(expectedUsers, users)
	s.Equal(expectedPaging, paging)
}

func (s *UserRepositoryTestSuite) TestGetUserByRole() {
	rows := sqlmock.NewRows([]string{"id", "name", "email", "role"})
	for _, user := range expectedUsers {
		rows.AddRow(user.ID, user.Name, user.Email, user.Role)
	}

	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, role FROM users WHERE role = $3 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`)).
		WithArgs("PARTICIPANT", expectedPaging.RowsPerPage, (expectedPaging.Page-1)*expectedPaging.RowsPerPage).
		WillReturnRows(rows)

	s.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM users WHERE role = $1")).
		WithArgs("PARTICIPANT").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(expectedUsers)))

	users, paging, err := s.repo.GetUserByRole(expectedUser.Role, expectedPaging.Page, expectedPaging.RowsPerPage)

	s.NoError(err)
	s.Equal(expectedUsers, users)
	s.Equal(expectedPaging, paging)
}

func (s *UserRepositoryTestSuite) TestGetUserByRoleNoRows() {
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, role FROM users WHERE role = $3 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`)).
		WithArgs(expectedUser.Role, expectedPaging.RowsPerPage, (expectedPaging.Page-1)*expectedPaging.RowsPerPage).
		WillReturnError(errors.New("sql: no rows in result set"))

	s.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM users WHERE role = $1")).
		WithArgs(expectedUser.Role).
		WillReturnError(errors.New("sql: no rows in result set"))

	_, _, err := s.repo.GetUserByRole(expectedUser.Role, expectedPaging.Page, expectedPaging.RowsPerPage)
	s.Error(err)
	s.NotNil(err)
}

func (s *UserRepositoryTestSuite) TestListFailed() {
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, role FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`)).
		WithArgs(expectedPaging.RowsPerPage, (expectedPaging.Page-1)*expectedPaging.RowsPerPage).
		WillReturnError(errors.New("sql: no rows in result set"))

	s.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM users")).
		WillReturnError(errors.New("sql: no rows in result set"))

	_, _, err := s.repo.List(expectedPaging.Page, expectedPaging.RowsPerPage)
	s.Error(err)
	s.NotNil(err)
}

func (s *UserRepositoryTestSuite) TestUpdate() {
	// Mock the database query and expected result
	s.mockSql.ExpectExec(regexp.QuoteMeta("UPDATE users SET name = $1, email = $2, password = $3, updated_at = $4 WHERE id = $5")).
		WithArgs(expectedUserUpdate.Name, expectedUserUpdate.Email, expectedUserUpdate.Password, customTimeMatcher(expectedUserUpdate.UpdatedAt), expectedUserUpdate.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the Update method
	err := s.repo.Update(expectedUserUpdate)

	// Verify the result
	s.NoError(err)
	s.Nil(err)
}

func (s *UserRepositoryTestSuite) TestUpdateFailed() {
	// Mock the database query and expected result
	s.mockSql.ExpectExec(regexp.QuoteMeta("UPDATE users SET name = $1, email = $2, password = $3, updated_at = $4 WHERE id = $5")).
		WithArgs(expectedUserUpdate.Name, expectedUserUpdate.Email, expectedUserUpdate.Password, customTimeMatcher(expectedUserUpdate.UpdatedAt), expectedUserUpdate.ID).
		WillReturnError(errors.New("sql: no rows in result set"))

	// Call the Update method
	err := s.repo.Update(expectedUserUpdate)
	s.Error(err)
}

func (s *UserRepositoryTestSuite) TestDelete() {
	// Mock the database query and expected result
	s.mockSql.ExpectExec(regexp.QuoteMeta("UPDATE users SET deleted_at = $1 WHERE id = $2")).
		WithArgs(customTimeMatcher(expectedUser.DeletedAt), expectedUser.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the Delete method
	err := s.repo.Delete(expectedUser.ID)
	s.NoError(err)
	s.Nil(err)
}

func customTimeMatcher(expected *time.Time) interface{} {
	return sqlmock.AnyArg()
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
