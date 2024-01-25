package controller

import (
	"log"
	"net/http"
	"net/mail"
	"strconv"

	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/shared/common"
	"enigmaCamp.com/instructor_led/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUC usecase.UserUsecase
	rg     *gin.RouterGroup
}

func NewUserController(userUC usecase.UserUsecase, rg *gin.RouterGroup) *UserController {
	return &UserController{
		userUC: userUC,
		rg:     rg,
	}
}

func (u *UserController) Route() {
	u.rg.POST("/users", u.CreateUserHanlder)
}

func (u *UserController) CreateUserHanlder(c *gin.Context) {
	csv := c.Query("csv")
	if csv != "" {
		return
	}
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

	log.Println("calling user usecase CreateUser")
	err := u.userUC.CreateUser(user)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, "failed to create user"+err.Error())
		return
	}
	common.SendSingleResponse(c, nil, "user created successfully")
}

func (u *UserController) GetAllUserHandler(c *gin.Context) {
	userRole := c.Query("role")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid page"+err.Error())
		return
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid size"+err.Error())
		return
	}
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}

	if userRole != "" {
		if userRole == "TRAINER" || userRole == "PARTICIPANT" {
			log.Println("calling user usecase GetUserByRole")
			users, paging, err := u.userUC.GetUserByRole(userRole, page, size)
			if err != nil {
				common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
				return
			}
			common.SendPagedResponse(c, users, paging, "success")
			return
		}
	}

	log.Println("calling user usecase ListAllUsers")
	users, paging, err := u.userUC.ListAllUsers(page, size)
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	common.SendPagedResponse(c, users, paging, "success")
}

func (u *UserController) GetUserByEmailHandler(c *gin.Context) {
	userEmail := c.Param("email")
	_, err := mail.ParseAddress(userEmail)
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid email"+err.Error())
		return
	}
	log.Println("calling user usecase GetUserByEmail")
	user, err := u.userUC.GetUserByEmail(userEmail)
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "user not found"+err.Error())
		return
	}

	common.SendSingleResponse(c, user, "success")
}
