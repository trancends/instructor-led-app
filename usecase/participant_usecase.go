package usecase

import (
	"fmt"

	"enigmaCamp.com/instructor_led/model"
	repository "enigmaCamp.com/instructor_led/repository/participant"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
)

type ParticipantUseCase interface {
	FindAllScheduleUC(page int, size int) ([]model.User, sharedmodel.Paging, error)
}

type participantUseCase struct {
	participantRepository repository.ParticipantRepository
}

// FindAllSchedule implements ParticipantUseCase.
func (p *participantUseCase) FindAllScheduleUC(page int, size int) ([]model.User, sharedmodel.Paging, error) {
	var user model.User

	if user.Role != "participant" {
		return nil, sharedmodel.Paging{}, fmt.Errorf("unauthorized")
	}
	users, paging, err := p.participantRepository.ListScheduled(page, size)
	if err != nil {
		return nil, sharedmodel.Paging{}, err
	}
	return users, paging, nil
}

func NewParticipantUseCase(participantRepository repository.ParticipantRepository) ParticipantUseCase {
	return &participantUseCase{
		participantRepository: participantRepository,
	}
}
