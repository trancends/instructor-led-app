// questions/usecase/usecase.go
package usecase

import (
	"log"

	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/repository"
)

type QuestionsUsecase interface {
	GetQuestion(id string) (model.Question, error)
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

// GetQuestion returns a single question based on the given ID.
func (u *questionsUsecase) GetQuestion(id string) (model.Question, error) {
	question, err := u.questionsRepo.Get(id)
	if err != nil {
		log.Println("QuestionsUsecase.GetQuestion:", err.Error())
		return model.Question{}, err
	}
	return question, nil
}

// ListQuestions returns a list of questions.
func (u *questionsUsecase) ListQuestions() ([]model.Question, error) {
	questions, err := u.questionsRepo.List()
	if err != nil {
		log.Println("QuestionsUsecase.ListQuestions:", err.Error())
		return nil, err
	}
	return questions, nil
}
