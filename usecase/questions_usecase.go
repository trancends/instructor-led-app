package usecase

import (
	"log"

	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/repository"
)

type QuestionsUsecase interface {
	CreateQuestionsUC(payload model.Question) (model.Question, error)
	GetQuestion(date string) ([]*model.Schedule, error)
	ListQuestions() ([]model.Question, error)
}

type questionsUsecase struct {
	questionsRepository repository.QuestionsRepository
}

func NewQuestionsUsecase(questionsRepository repository.QuestionsRepository) QuestionsUsecase {
	return &questionsUsecase{
		questionsRepository: questionsRepository,
	}
}

// CreateQuestionsUC implements QuestionsUsecase.
func (q *questionsUsecase) CreateQuestionsUC(payload model.Question) (model.Question, error) {
	questions, err := q.questionsRepository.CreateQuestions(payload)
	if err != nil {
		return questions, err
	}
	return questions, nil
}

// GetQuestion returns a list of schedules based on the given date.
func (u *questionsUsecase) GetQuestion(date string) ([]*model.Schedule, error) {
	schedules, err := u.questionsRepository.Get(date)
	if err != nil {
		log.Println("QuestionsUsecase.GetQuestion:", err.Error())
		return nil, err
	}
	return schedules, nil
}

func (u *questionsUsecase) ListQuestions() ([]model.Question, error) {
	questions, err := u.questionsRepository.List()
	if err != nil {
		log.Println("QuestionsUsecase.ListQuestions:", err.Error())
		return nil, err
	}
	return questions, nil
}
