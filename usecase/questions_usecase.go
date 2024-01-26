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
	DeleteQuestion(id string) error
	UpdateQuestionStatus(payload model.Question) error
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

func (u *questionsUsecase) DeleteQuestion(id string) error {
	_, err := u.questionsRepository.GetByID(id)
	if err != nil {
		log.Println("QuestionsUsecase.DeleteQuestion:", err.Error())
		return err
	}
	err = u.questionsRepository.Delete(id)
	if err != nil {
		log.Println("QuestionsUsecase.DeleteQuestion:", err.Error())
		return err
	}
	return nil
}

func (u *questionsUsecase) UpdateQuestionStatus(payload model.Question) error {
	_, err := u.questionsRepository.GetByID(payload.ID)
	if err != nil {
		log.Println("QuestionsUsecase.DeleteQuestion:", err.Error())
		return err
	}
	err = u.questionsRepository.UpdateStatus(payload)
	if err != nil {
		log.Println("QuestionsUsecase.UpdateQuestionStatus:", err.Error())
		return err
	}
	return nil
}
