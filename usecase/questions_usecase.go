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
}
