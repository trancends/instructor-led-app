package usecase

import (
	"enigmaCamp.com/instructor_led/model"
	repository "enigmaCamp.com/instructor_led/repository"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
)

type ShecdulesUseCase interface {
	FindAllScheduleUC(page int, size int) ([]model.Schedule, sharedmodel.Paging, error)
}

type schedulesUseCase struct {
	participantRepository repository.ParticipantRepository
}

// FindAllScheduleUC implements ShecdulesUseCase.
func (s *schedulesUseCase) FindAllScheduleUC(page int, size int) ([]model.Schedule, sharedmodel.Paging, error) {
	users, paging, err := s.participantRepository.ListScheduled(page, size)
	if err != nil {
		return nil, sharedmodel.Paging{}, err
	}
	return users, paging, nil
}

// FindAllSchedule implements ParticipantUseCase.

func NewSchedulesUseCase(participantRepository repository.ParticipantRepository) ShecdulesUseCase {
	return &schedulesUseCase{participantRepository}
}
