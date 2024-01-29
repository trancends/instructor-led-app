package usecase_mock

import (
	"enigmaCamp.com/instructor_led/model"
	"github.com/stretchr/testify/mock"
)

type QuestionUscaseMock struct {
	mock.Mock
}

// CreateQuestionsUC implements usecase.QuestionsUsecase.
func (*QuestionUscaseMock) CreateQuestionsUC(payload model.Question) (model.Question, error) {
	panic("unimplemented")
}

// ListQuestionsByScheduleID implements usecase.QuestionsUsecase.
func (*QuestionUscaseMock) ListQuestionsByScheduleID(scheduleID string) ([]model.Question, error) {
	panic("unimplemented")
}

// CreateQuestionsUC implements usecase.QuestionsUsecase.

func (q *QuestionUscaseMock) CreateQuestion(payload model.Question) (model.Question, error) {
	args := q.Called(payload)
	return args.Get(0).(model.Question), args.Error(1)
}
func (q *QuestionUscaseMock) ListQuestions() ([]model.Question, error) {
	args := q.Called()
	return args.Get(0).([]model.Question), args.Error(1)
}
func (q *QuestionUscaseMock) GetQuestion(date string) ([]*model.Schedule, error) {
	args := q.Called(date)
	return args.Get(0).([]*model.Schedule), args.Error(1)
}
func (q *QuestionUscaseMock) DeleteQuestion(id string) error {
	args := q.Called(id)
	return args.Error(0)
}
func (q *QuestionUscaseMock) UpdateQuestionStatus(payload model.Question) error {
	args := q.Called(payload)
	return args.Error(0)
}
