package controller

import (
	"log"
	"net/http"

	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/shared/common"
	"enigmaCamp.com/instructor_led/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUC usecase.UserUsecase
	rg     gin.RouterGroup
}

func NewUserController(rg gin.RouterGroup, userUC usecase.UserUsecase) *UserController {
	return &UserController{
		userUC: userUC,
		rg:     rg,
	}
}

func (u *UserController) Route() {
	u.rg.POST("/users", u.CreateUserHanlder)
}

func (u *UserController) CreateUserHanlder(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("invalid json at CreateUserHanlder")
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid json"+err.Error())
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" || user.Role == "" {
		common.SendErrorResponse(c, http.StatusBadRequest, "name, email, role, and password are required")
		return
	}

	err := u.userUC.CreateUser(user)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, "failed to create user"+err.Error())
		return
	}
	common.SendSingleResponse(c, nil, "user created successfully")
}
