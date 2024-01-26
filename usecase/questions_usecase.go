package usecase

import (
	"enigmaCamp.com/instructor_led/model"
	repository "enigmaCamp.com/instructor_led/repository"
)

type QuestionsUsecase interface {
	CreateQuestionsUC(payload model.Questions) (model.Questions, error)
}

type questionsUsecase struct {
	questionsRepository repository.QuestionsRepository
}

// CreateQuestionsUC implements QuestionsUsecase.
func (q *questionsUsecase) CreateQuestionsUC(payload model.Questions) (model.Questions, error) {
	questions, err := q.questionsRepository.CreateQuestions(payload)
	if err != nil {
		return questions, err
	}
	return questions, nil
}

func NewQuestionsUsecase(questionsRepository repository.QuestionsRepository) QuestionsUsecase {
	return &questionsUsecase{
		questionsRepository: questionsRepository,
	}

// questions/usecase/usecase.go
package usecase

import (
	"log"
	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/repository"
)

type QuestionsUsecase interface {
	GetQuestion(date string) ([]*model.Schedule, error)
	ListQuestions() ([]model.Question, error)
}

type questionsUsecase struct {
	questionsRepo repository.QuestionsRepository
}

// NewQuestionsUsecase initializes a new QuestionsUsecase.
func NewQuestionsUsecase(questionsRepo repository.QuestionsRepository) QuestionsUsecase {
	return &questionsUsecase{
		questionsRepo: questionsRepo,
	}
}

// GetQuestion returns a list of schedules based on the given date.
func (u *questionsUsecase) GetQuestion(date string) ([]*model.Schedule, error) {
	schedules, err := u.questionsRepo.Get(date)
	if err != nil {
		log.Println("QuestionsUsecase.GetQuestion:", err.Error())
		return nil, err
	}
	return schedules, nil
}

func (u *questionsUsecase) ListQuestions() ([]model.Question, error) {
	questions, err := u.questionsRepo.List()
	if err != nil {
		log.Println("QuestionsUsecase.ListQuestions:", err.Error())
		return nil, err
	}
	return questions, nil
}
