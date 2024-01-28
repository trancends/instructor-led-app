package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"enigmaCamp.com/instructor_led/delivery/middleware"
	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/shared/common"
	"enigmaCamp.com/instructor_led/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SchedulesController struct {
	schedulesUC    usecase.ShecdulesUseCase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func NewSchedulesController(schedulesUC usecase.ShecdulesUseCase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *SchedulesController {
	return &SchedulesController{
		schedulesUC:    schedulesUC,
		rg:             rg,
		authMiddleware: authMiddleware,
	}
}

func (s *SchedulesController) Route() {
	s.rg.GET("/schedules", s.authMiddleware.RequireToken("ADMIN", "TRAINER"), s.FindAllScheduleHandler)
	s.rg.POST("/schedules", s.authMiddleware.RequireToken("ADMIN", "PARTICIPANT"), s.CreateScheduleHandler)
	s.rg.GET("/schedules/:id", s.authMiddleware.RequireToken("ADMIN", "TRAINER", "PARTICIPANT"), s.FindByIDScheduleHandler)
	s.rg.DELETE("/schedules/:id", s.authMiddleware.RequireToken("ADMIN"), s.DeleteScheduleHandler)
	s.rg.PATCH("/schedules/upload/:id", s.UploadDocumentationHandler)
}

func (s *SchedulesController) CreateScheduleHandler(c *gin.Context) {
	var payload model.Schedule
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("SchedulesController.CreateScheduleHandler:", err.Error())
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid json"+err.Error())
		return
	}

	if payload.UserID == "" || payload.Date == "" || payload.StartTime == "" || payload.EndTime == "" {
		common.SendErrorResponse(c, http.StatusBadRequest, "payload cannot be empty")
		return
	}

	schedule, err := s.schedulesUC.CreateScheduledUC(payload)
	if err != nil {
		log.Println("SchedulesController.CreateScheduleHandler:", err.Error())
		common.SendErrorResponse(c, http.StatusInternalServerError, "failed to create schedule"+err.Error())
		return
	}

	common.SendSingleResponse(c, schedule, "schedule created successfully")
}

func (s *SchedulesController) FindAllScheduleHandler(c *gin.Context) {
	pageQuery := c.Query("page")
	sizeQuery := c.Query("size")

	if pageQuery == "" || sizeQuery == "" {
		common.SendErrorResponse(c, http.StatusBadRequest, "pageParam or sizeParam cant be empty")
		return
	}

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid page param")
		return
	}
	size, err := strconv.Atoi(sizeQuery)
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid psize param")
		return
	}

	log.Println("calling user usecase FindAllScheduleUC")
	users, paging, err := s.schedulesUC.FindAllScheduleUC(page, size)
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	common.SendPagedResponse(c, users, paging, "success")
}

func (s *SchedulesController) FindByIDScheduleHandler(c *gin.Context) {
	id := c.Param("id")

	log.Println("calling user usecase FindByIDUC")
	schedule, err := s.schedulesUC.FindByIDUC(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			common.SendErrorResponse(c, http.StatusBadRequest, "schedule not found")
			return
		}
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	common.SendSingleResponse(c, schedule, "success")
}

func (s *SchedulesController) DeleteScheduleHandler(c *gin.Context) {
	id := c.Param("id")

	_, err := s.schedulesUC.FindByIDUC(id)
	if err != nil {
		common.SendErrorResponse(c, http.StatusNotFound, "schedule not found")
		return
	}

	err = s.schedulesUC.DeletedScheduleIDUC(id)
	if err != nil {
		log.Println("SchedulesController.DeleteScheduleHandler:", err.Error())
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	common.SendSingleResponse(c, nil, "success")
}

func (s *SchedulesController) UploadDocumentationHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.SendErrorResponse(c, http.StatusBadRequest, "id cannot be empty")
		return
	}
	log.Println(id)

	_, err := s.schedulesUC.FindByIDUC(id)
	if err != nil {
		log.Println("SchedulesController.UploadDocumentationHandler err: ", err)
		common.SendErrorResponse(c, http.StatusNotFound, "schedule not found")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid file"+err.Error())
		return
	}
	// max image size is 3 MB
	maxSize := 3 * 1024 * 1024 // 3 MB in bytes
	if file.Size > int64(maxSize) {
		common.SendErrorResponse(c, http.StatusBadRequest, "image size exceeds the maximum limit of 3 MB")
		return
	}

	// Define the file path to save the upload image.
	uuid := uuid.New()
	stringUUID := uuid.String()
	path := fmt.Sprintf("./scheduleImages/%s", stringUUID+"-"+file.Filename)

	// Save the uploaded image to the defined path
	err = saveUploadFile(file, path)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, "failed to save image"+err.Error())
		return
	}

	// Construct the URL for the saved picture.
	baseURL := "http://localhost:8080"
	pictureURL := fmt.Sprintf("%s/scheduleImages/%s", baseURL, stringUUID+"-"+file.Filename)

	// Update the schedule's documentation picture URL.
	err = s.schedulesUC.UpdateScheduleDocumentation(id, pictureURL)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, "failed to update schedule"+err.Error())
		return
	}
	common.SendSingleResponse(c, nil, "success")
}

func saveUploadFile(file *multipart.FileHeader, path string) error {
	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create a destination file for the uploaded content.
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
