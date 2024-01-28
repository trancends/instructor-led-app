package usecase_mock

import (
	"enigmaCamp.com/instructor_led/model"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
	"github.com/stretchr/testify/mock"
)

type QuestionsUseCaseMock struct {
	mock.Mock
}

func (m *QuestionsUseCaseMock) CreateQuestionsUC(payload model.Question) (model.Question, error) {
	args := m.Called(payload)
	return args.Get(0).(model.Question), args.Error(1)
}

func (m *QuestionsUseCaseMock) GetQuestion(date string) ([]*model.Schedule, error) {
	args := m.Called(date)
	return args.Get(0).([]*model.Schedule), args.Error(1)
}

func (m *QuestionsUseCaseMock) FindAllQuestionsUC(page int, size int) ([]model.Question, sharedmodel.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]model.Question), args.Get(1).(sharedmodel.Paging), args.Error(2)
}

func (m *QuestionsUseCaseMock) ListQuestionsByScheduleID(scheduleID string) ([]model.Question, error) {
	args := m.Called(scheduleID)
	return args.Get(0).([]model.Question), args.Error(1)
}

func (m *QuestionsUseCaseMock) DeleteQuestion(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *QuestionsUseCaseMock) UpdateQuestionStatus(payload model.Question) error {
	args := m.Called(payload)
	return args.Error(0)
}
