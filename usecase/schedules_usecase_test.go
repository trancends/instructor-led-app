package usecase

import (
	"errors"
	"testing"
	"time"

	repo_mock "enigmaCamp.com/instructor_led/mock/repository_mock"
	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
	"github.com/stretchr/testify/suite"
)

var expectedSchedule = model.Schedule{
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
			ScheduleID:  "1",
			Description: "Test",
			Status:      "PROCESS",
			CreatedAt:   &time.Time{},
			UpdatedAt:   &time.Time{},
			DeletedAt:   &time.Time{},
		},
	},
}

var mockPaging = sharedmodel.Paging{
	Page:        1,
	RowsPerPage: 1,
	TotalRows:   5,
	TotalPages:  1,
}

var mockData = []model.Schedule{expectedSchedule}

type ScheduleUseCaseTestSuite struct {
	suite.Suite
	srm *repo_mock.ScheduleRepositoryMock
	suc ShecdulesUseCase
}

func (s *ScheduleUseCaseTestSuite) SetupTest() {
	s.srm = new(repo_mock.ScheduleRepositoryMock)
	s.suc = NewSchedulesUseCase(s.srm)
}

func (s *ScheduleUseCaseTestSuite) TestCreateScheduledUC_Success() {
	s.srm.On("CreateScheduled", expectedSchedule).Return(expectedSchedule, nil)
	actual, err := s.suc.CreateScheduledUC(expectedSchedule)
	s.NoError(err)
	s.Nil(err)
	s.Equal(expectedSchedule, actual)
}

func (s *ScheduleUseCaseTestSuite) TestCreateScheduledUC_Fail() {
	s.srm.On("CreateScheduled", expectedSchedule).Return(model.Schedule{}, errors.New("some error"))
	_, err := s.suc.CreateScheduledUC(expectedSchedule)
	s.Error(err)
	s.NotNil(err)
}

func (s *ScheduleUseCaseTestSuite) TestFindScheduleByRole_Success() {
	s.srm.On("ListScheduleByRole", mockPaging.Page, mockPaging.RowsPerPage, "admin").Return(mockData, mockPaging, nil)
	actual, paging, err := s.suc.FindScheduleByRole(mockPaging.Page, mockPaging.RowsPerPage, "admin")
	s.NoError(err)
	s.Nil(err)
	s.Equal(mockData, actual)
	s.Equal(mockPaging, paging)
}

func (s *ScheduleUseCaseTestSuite) TestFindScheduleByRole_Fail() {
	s.srm.On("ListScheduleByRole", mockPaging.Page, mockPaging.RowsPerPage, "admin").Return([]model.Schedule{}, sharedmodel.Paging{}, errors.New("some error"))
	_, _, err := s.suc.FindScheduleByRole(mockPaging.Page, mockPaging.RowsPerPage, "admin")
	s.Error(err)
	s.NotNil(err)
	s.Equal(errors.New("some error"), err)
}

func (s *ScheduleUseCaseTestSuite) TestFindAllScheduleUC_Success() {
	s.srm.On("ListScheduled", mockPaging.Page, mockPaging.RowsPerPage).Return(mockData, mockPaging, nil)
	actual, paging, err := s.suc.FindAllScheduleUC(mockPaging.Page, mockPaging.RowsPerPage)
	s.NoError(err)
	s.Nil(err)
	s.Equal(mockData, actual)
	s.Equal(mockPaging, paging)
}

func (s *ScheduleUseCaseTestSuite) TestFindAllScheduleUC_Fail() {
	s.srm.On("ListScheduled", mockPaging.Page, mockPaging.RowsPerPage).Return([]model.Schedule{}, sharedmodel.Paging{}, errors.New("some error"))
	_, _, err := s.suc.FindAllScheduleUC(mockPaging.Page, mockPaging.RowsPerPage)
	s.Error(err)
	s.NotNil(err)
	s.Equal(errors.New("some error"), err)
}

func (s *ScheduleUseCaseTestSuite) TestFindByIDUC_Success() {
	s.srm.On("GetByID", expectedSchedule.ID).Return(expectedSchedule, nil)
	actual, err := s.suc.FindByIDUC(expectedSchedule.ID)
	s.NoError(err)
	s.Nil(err)
	s.Equal(expectedSchedule, actual)
}
func (s *ScheduleUseCaseTestSuite) TestFindByIDUC_Fail() {
	s.srm.On("GetByID", expectedSchedule.ID).Return(model.Schedule{}, errors.New("some error"))
	_, err := s.suc.FindByIDUC(expectedSchedule.ID)
	s.Error(err)
	s.NotNil(err)
	s.Equal(errors.New("some error"), err)
}

func (s *ScheduleUseCaseTestSuite) TestDeletedScheduleIDUC_Success() {
	s.srm.On("Delete", expectedSchedule.ID).Return(nil)
	err := s.suc.DeletedScheduleIDUC(expectedSchedule.ID)
	s.NoError(err)
	s.Nil(err)
	s.Equal(nil, err)
}

func (s *ScheduleUseCaseTestSuite) TestDeletedScheduleIDUC_Fail() {
	s.srm.On("Delete", expectedSchedule.ID).Return(errors.New("some error"))
	err := s.suc.DeletedScheduleIDUC(expectedSchedule.ID)
	s.Error(err)
	s.NotNil(err)
	s.Equal(errors.New("some error"), err)
}

func (s *ScheduleUseCaseTestSuite) TestUpdateScheduleDocumentation_Success() {
	s.srm.On("UpdateDocumentation", expectedSchedule.ID, expectedSchedule.Documentation).Return(nil)
	err := s.suc.UpdateScheduleDocumentation(expectedSchedule.ID, expectedSchedule.Documentation)
	s.NoError(err)
	s.Nil(err)
	s.Equal(nil, err)
}

func TestScheduleUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ScheduleUseCaseTestSuite))
}
