package usecase

import (
	"fmt"
	"log"

	"enigmaCamp.com/instructor_led/model/dto"
	"enigmaCamp.com/instructor_led/shared/service"
	"enigmaCamp.com/instructor_led/shared/utils"
	"golang.org/x/crypto/bcrypt"
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
	user, err := a.userUC.GetUserByEmail(payload.Email)
	if err != nil {
		return dto.AuthResponseDTO{}, err
	}
	log.Println("authUC.FindAuthorByEmail", user)

	test := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	log.Println("bcrypt returned: ", test)
	if !utils.CheckPassword(user.Password, payload.Password) {
		return dto.AuthResponseDTO{}, fmt.Errorf("wrong password")
	}
	log.Println("password correct")

	// TODO generate jwt
	tokenDto, err := a.jwtService.GenerateToken(user)
	log.Println("TokenDto", tokenDto)
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
