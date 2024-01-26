package usecase

import (
	"log"

	"enigmaCamp.com/instructor_led/model/dto"
	"enigmaCamp.com/instructor_led/shared/service"
)

type AuthUseCase interface {
	Login(payload dto.AuthRequestDTO) (dto.AuthResponseDTO, error)
	// GetKey() []byte
}

type authUseCase struct {
	userUC     UserUsecase
	jwtService service.JwtService
}

// func (a *authUseCase) GetKey() []byte {
// 	return a.jwtService.GetKey()
// }

func (a *authUseCase) Login(payload dto.AuthRequestDTO) (dto.AuthResponseDTO, error) {
	author, err := a.userUC.GetUserByEmail(payload.Email)
	if err != nil {
		return dto.AuthResponseDTO{}, err
	}
	log.Println("authUC.FindAuthorByEmail", author)
	log.Println(payload.Password == author.Password)

	if payload.Password != author.Password {
		return dto.AuthResponseDTO{}, err
	}
	log.Println("password correct")

	// TODO generate jwt
	tokenDto, err := a.jwtService.GenerateToken(author)
	log.Println("TokenDto")
	log.Println(tokenDto)
	if err != nil {
		return dto.AuthResponseDTO{}, err
	}
	return tokenDto, nil
}

func NewAuthUseCase(userUC UserUsecase, jwtService service.JwtService) AuthUseCase {
	return &authUseCase{
		userUC:     userUC,
		jwtService: jwtService,
	}
}
