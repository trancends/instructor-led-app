package controller

import (
	"fmt"
	"log"
	"net/http"

	"enigmaCamp.com/instructor_led/model/dto"
	"enigmaCamp.com/instructor_led/shared/common"
	"enigmaCamp.com/instructor_led/usecase"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUC usecase.AuthUseCase
	rg     *gin.RouterGroup
}

func (a *AuthController) loginHandler(c *gin.Context) {
	log.Println("hit loginHandler")
	var payload dto.AuthRequestDTO
	err := c.ShouldBindJSON(&payload)
	fmt.Println(payload)
	if err != nil {
		log.Println("error at loginHandler", err)
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// send payload to usecase
	response, err := a.authUC.Login(payload)
	fmt.Println(response)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.SendSingleResponse(c, response, "login success")
}

func (a *AuthController) Route() {
	a.rg.POST("/login", a.loginHandler)
}

func NewAuthController(authUC usecase.AuthUseCase, rg *gin.RouterGroup) *AuthController {
	return &AuthController{
		authUC: authUC,
		rg:     rg,
	}
}
