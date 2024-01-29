package usecase

import (
	"testing"
	"time"

	"enigmaCamp.com/instructor_led/mock/repository_mock"
	"enigmaCamp.com/instructor_led/model"

	"github.com/stretchr/testify/suite"
)

var expectedQuestions = model.Question{
	ID:          "1",
	UserID:      "1",
	ScheduleID:  "1",
	Description: "test",
	Status:      "PROCESS",
	CreatedAt:   nil,
	UpdatedAt:   nil,
	DeletedAt:   nil,
}

var expectedSchedulepointer = []*model.Schedule{
	{
		ID:            "1",
		UserID:        "1",
		Date:          "2022-01-01",
		StartTime:     "11:00",
		EndTime:       "12:00",
		Documentation: "test",
		CreatedAt:     time.Now(),
		UpdatedAt:     nil,
		DeletedAt:     nil,
		Questions:     []model.Question{},
	},
}

type QuestionUsecaseSuite struct {
	suite.Suite
	qrm *repository_mock.QuestionRepositoryMock
	quc QuestionsUsecase
}

func (s *QuestionUsecaseSuite) SetupTest() {
	s.qrm = new(repository_mock.QuestionRepositoryMock)
	s.quc = NewQuestionsUsecase(s.qrm)
}

func (s *QuestionUsecaseSuite) TestCreateQuestion() {
	s.qrm.On("CreateQuestions", expectedQuestions).Return(expectedQuestions, nil)
	question, err := s.quc.CreateQuestionsUC(expectedQuestions)
	s.Equal(expectedQuestions, question)
	s.Nil(err)
	s.Equal(expectedQuestions.ID, question.ID)
}

func (s *QuestionUsecaseSuite) TestListQuestions() {
	s.qrm.On("List").Return([]model.Question{expectedQuestions}, nil)
	questions, err := s.quc.ListQuestions()
	s.Equal(expectedQuestions, questions[0])
	s.Nil(err)
}

func (s *QuestionUsecaseSuite) TestDeleteQuestion() {
	s.qrm.On("Delete", "1").Return(nil)
	s.qrm.On("GetByID", "1").Return(expectedQuestions, nil)
	err := s.quc.DeleteQuestion("1")
	s.Nil(err)
}

func (s *QuestionUsecaseSuite) TestUpdateQuestionStatus() {
	s.qrm.On("UpdateStatus", expectedQuestions).Return(nil)
	s.qrm.On("GetByID", "1").Return(expectedQuestions, nil)
	err := s.quc.UpdateQuestionStatus(expectedQuestions)
	s.Nil(err)
}

func (s *QuestionUsecaseSuite) TestGetByDate() {
	s.qrm.On("Get", "2022-01-01").Return(expectedSchedulepointer, nil)
	schedule, err := s.quc.GetQuestion("2022-01-01")
	s.Equal(expectedSchedulepointer, schedule)
	s.Nil(err)
}

func TestQuestionUsecaseSuite(t *testing.T) {
	suite.Run(t, new(QuestionUsecaseSuite))
}
