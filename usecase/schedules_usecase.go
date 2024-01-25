package usecase

import (
	"log"

	"enigmaCamp.com/instructor_led/model"
	repository "enigmaCamp.com/instructor_led/repository"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
)

type ShecdulesUseCase interface {
	FindAllScheduleUC(page int, size int) ([]model.Schedule, sharedmodel.Paging, error)
	CreateScheduledUC(payload model.Schedule) (model.Schedule, error)
	FindByIDUC(id string) (model.Schedule, error)
}

type schedulesUseCase struct {
	scheduleRepository repository.ScheduleRepository
}

// CreateScheduledUC implements ShecdulesUseCase.
func (s *schedulesUseCase) CreateScheduledUC(payload model.Schedule) (model.Schedule, error) {
	schedule, err := s.scheduleRepository.CreateScheduled(payload)
	if err != nil {
		log.Println("schedulesUseCase.CreateScheduledUC:", err.Error())
		return schedule, err
	}
	return schedule, nil
}

// FindAllScheduleUC implements ShecdulesUseCase.
func (s *schedulesUseCase) FindAllScheduleUC(page int, size int) ([]model.Schedule, sharedmodel.Paging, error) {
	users, paging, err := s.scheduleRepository.ListScheduled(page, size)
	if err != nil {
		log.Println("schedulesUseCase.FindAllScheduleUC:", err.Error())
		return nil, sharedmodel.Paging{}, err
	}
	return users, paging, nil
}

// FindByIDUC implements ShecdulesUseCase.
func (s *schedulesUseCase) FindByIDUC(id string) (model.Schedule, error) {
	schdule, err := s.scheduleRepository.GetByID(id)
	if err != nil {
		log.Println("schedulesUseCase.FindByIDUC:", err.Error())
		return schdule, err
	}
	return schdule, nil
}

// FindAllSchedule implements ParticipantUseCase.

func NewSchedulesUseCase(scheduleRepository repository.ScheduleRepository) ShecdulesUseCase {
	return &schedulesUseCase{scheduleRepository}
}
