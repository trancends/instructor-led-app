package controller

import (
	"database/sql"
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
	if err != nil {
		log.Println("error at loginHandler - Unable to bind JSON payload", err)
		common.SendErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// send payload to usecase
	response, err := a.authUC.Login(payload)
	if err != nil {
		log.Println("error at loginHandler - AuthUseCase login failed", err)

		if err == sql.ErrNoRows {
			common.SendErrorResponse(c, http.StatusBadRequest, "User not found")
			return
		}

		common.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Login failed: %v", err))
		return
	}

	common.SendSingleResponse(c, response, "Login success")
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
