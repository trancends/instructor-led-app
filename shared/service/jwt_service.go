package service

import (
	"fmt"
	"log"
	"time"

	"enigmaCamp.com/instructor_led/config"
	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/model/dto"
	sharedmodel "enigmaCamp.com/instructor_led/shared/shared_model"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	GenerateToken(user model.User) (dto.AuthResponseDTO, error)
	ParseToken(tokenHeadr string) (jwt.MapClaims, error)
	GetKey() []byte
}

type jwtService struct {
	cfg config.TokenConfig
}

func (j *jwtService) ParseToken(tokenHeader string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return j.cfg.JwtSignatureKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("oops, failed to verify token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("oops, failed to claim token")
	}
	return claims, nil
}

func (j *jwtService) GetKey() []byte {
	return j.cfg.JwtSignatureKey
}

func (j *jwtService) GenerateToken(user model.User) (dto.AuthResponseDTO, error) {
	claims := sharedmodel.CustomClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.IssuerName,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.JwtExpiresTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.cfg.JwtSignatureKey)
	log.Println("tokenString:")
	log.Println(tokenString)
	if err != nil {
		return dto.AuthResponseDTO{}, err
	}
	return dto.AuthResponseDTO{
		Token: tokenString,
	}, nil
}

func NewJwtService(cfg config.TokenConfig) JwtService {
	return &jwtService{
		cfg: cfg,
	}
}
