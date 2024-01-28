package controller

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"strconv"

	"enigmaCamp.com/instructor_led/delivery/middleware"
	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/shared/common"
	"enigmaCamp.com/instructor_led/shared/utils"
	"enigmaCamp.com/instructor_led/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUC         usecase.UserUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func NewUserController(userUC usecase.UserUsecase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *UserController {
	return &UserController{
		userUC:         userUC,
		rg:             rg,
		authMiddleware: authMiddleware,
	}
}

func (u *UserController) Route() {
	u.rg.POST("/users/csv", u.CreateUserCSVHandler)
	u.rg.POST("/users", u.CreateUserHanlder)
	u.rg.GET("/users", u.GetAllUserHandler)
	u.rg.GET("/users/:email", u.GetUserByEmailHandler)
	u.rg.PUT("/users/:id", u.UpdateUserHandler)
	u.rg.DELETE("/users/:id", u.authMiddleware.RequireToken("ADMIN"), u.DeleteUserHandler)
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
		log.Println("error at calling user usecase CreateUser", err)
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
				log.Println("error at calling userUC getUserByRole", err)
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
		log.Println("error at calling userUC ListAllUsers", err)
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
		log.Println("error at user usecase GetUserByEmail", err)
		common.SendErrorResponse(c, http.StatusBadRequest, "user not found"+err.Error())
		return
	}

	common.SendSingleResponse(c, user, "success")
}

func (u *UserController) UpdateUserHandler(c *gin.Context) {
	userId := c.Param("id")
	var user model.User
	_, err := u.userUC.GetUserByID(userId)
	if err == sql.ErrNoRows {
		log.Println("error at calling user usecase GetUserByID", err)
		common.SendErrorResponse(c, http.StatusBadRequest, "user not found"+err.Error())
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("invalid json at UpdateUserHandler")
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid json"+err.Error())
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		common.SendErrorResponse(c, http.StatusBadRequest, "name, email, and password are required")
		return
	}

	user.Password, _ = utils.GetHashPassword(user.Password)
	user.ID = userId
	err = u.userUC.UpdateUser(user)
	if err != nil {
		log.Println("error at calling user usecase UpdateUser", err)
		common.SendErrorResponse(c, http.StatusInternalServerError, "failed to update user"+err.Error())
		return
	}
	common.SendSingleResponse(c, user, "user updated successfully")
}

func (u *UserController) DeleteUserHandler(c *gin.Context) {
	userId := c.Param("id")
	_, err := u.userUC.GetUserByID(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			common.SendErrorResponse(c, http.StatusBadRequest, "user not found"+err.Error())
			return
		}
		log.Println("error at calling user usecase GetUserByID", err)
		common.SendErrorResponse(c, http.StatusInternalServerError, "failed to get user"+err.Error())
		return
	}

	err = u.userUC.DeleteUser(userId)
	if err != nil {
		fmt.Println("error at calling user usecase DeleteUser", err)
		common.SendErrorResponse(c, http.StatusInternalServerError, "failed to delete user"+err.Error())
		return
	}

	common.SendSingleResponse(c, userId, "user deleted successfully")
}

func (u *UserController) CreateUserCSVHandler(c *gin.Context) {
	// parse the form
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid form"+err.Error())
		return
	}
	// get the file from the form
	file, _, err := c.Request.FormFile("csv")
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid file"+err.Error())
		return
	}
	defer file.Close()

	// parse the CSV file
	reader := csv.NewReader(file)
	var users []model.User
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			common.SendErrorResponse(c, http.StatusInternalServerError, "failed to read csv"+err.Error())
			return
		}
		var user model.User
		user.Name = record[0]
		user.Email = record[1]
		user.Password, err = utils.GetHashPassword(record[2])
		if err != nil {
			common.SendErrorResponse(c, http.StatusInternalServerError, "failed to hash password"+err.Error())
			return
		}
		user.Role = record[3]

		users = append(users, user)
	}

	err = u.userUC.CreateUserCSV(users)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, "failed to create user"+err.Error())
		return
	}
	common.SendSingleResponse(c, "", "users created successfully")
}
