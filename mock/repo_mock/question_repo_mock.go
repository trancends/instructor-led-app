package repo_mock

import (
	"enigmaCamp.com/instructor_led/model"
	"github.com/stretchr/testify/mock"
)

type QuestionRepositoryMock struct {
	mock.Mock
}

// Get implements repository.QuestionsRepository.
func (q *QuestionRepositoryMock) Get(date string) ([]*model.Schedule, error) {
	args := q.Called(date)
	return args.Get(0).([]*model.Schedule), args.Error(1)
}

// UpdateStatus implements repository.QuestionsRepository.
func (q *QuestionRepositoryMock) UpdateStatus(payload model.Question) error {
	args := q.Called(payload)
	return args.Error(0)
}

// CreateQuestions implements repository.QuestionsRepository.
func (q *QuestionRepositoryMock) CreateQuestions(payload model.Question) (model.Question, error) {
	args := q.Called(payload)
	return args.Get(0).(model.Question), args.Error(1)
}

// Delete implements repository.QuestionsRepository.
func (q *QuestionRepositoryMock) Delete(id string) error {
	args := q.Called(id)
	return args.Error(0)
}

// Get implements repository.QuestionsRepository
// GetByID implements repository.QuestionsRepository.
func (q *QuestionRepositoryMock) GetByID(id string) (model.Question, error) {
	args := q.Called(id)
	return args.Get(0).(model.Question), args.Error(1)
}

// List implements repository.QuestionsRepository.
func (q *QuestionRepositoryMock) List() ([]model.Question, error) {
	args := q.Called()
	return args.Get(0).([]model.Question), args.Error(1)
}
