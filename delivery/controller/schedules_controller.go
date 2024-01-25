package controller

import (
	"log"
	"net/http"
	"strconv"

	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/shared/common"
	"enigmaCamp.com/instructor_led/usecase"
	"github.com/gin-gonic/gin"
)

type SchedulesController struct {
	schedulesUC usecase.ShecdulesUseCase
	rg          *gin.RouterGroup
}

func NewSchedulesController(schedulesUC usecase.ShecdulesUseCase, rg *gin.RouterGroup) *SchedulesController {
	return &SchedulesController{
		schedulesUC: schedulesUC,
		rg:          rg,
	}
}

func (s *SchedulesController) Route() {
	s.rg.GET("/schedules", s.FindAllScheduleHandler)
	s.rg.POST("/schedules", s.CreateScheduleHandler)
	s.rg.GET("/schedules/:id", s.FindByIDScheduleHandler)
}

func (s *SchedulesController) CreateScheduleHandler(c *gin.Context) {
	var payload model.Schedule
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid json"+err.Error())
	}
	payloads, err := s.schedulesUC.CreateScheduledUC(payload)
	if err != nil {
		log.Println("SchedulesController.CreateScheduleHandler:", err.Error())
		common.SendErrorResponse(c, http.StatusInternalServerError, "failed to create schedule"+err.Error())
	}
	common.SendSingleResponse(c, payloads, "schedule created successfully")
}

func (s *SchedulesController) FindAllScheduleHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	users, paging, err := s.schedulesUC.FindAllScheduleUC(page, size)
	if err != nil {
		log.Println("SchedulesController.FindAllScheduleHandler:", err.Error())
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	common.SendPagedResponse(c, users, paging, "success")
}

func (s *SchedulesController) FindByIDScheduleHandler(c *gin.Context) {
	var payload model.Schedule
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid json"+err.Error())
	}
	payloads, err := s.schedulesUC.FindByIDUC(payload.ID)
	if err != nil {
		log.Println("SchedulesController.FindByIDScheduleHandler:", err.Error())
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	common.SendSingleResponse(c, payloads, "success")
}

func (s *SchedulesController) DeleteScheduleHandler(c *gin.Context) {
	id := c.Param("id")
	err := s.schedulesUC.DeletedScheduleIDUC(id)
	if err != nil {
		log.Println("SchedulesController.DeleteScheduleHandler:", err.Error())
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	common.SendSingleResponse(c, nil, "success")
}
