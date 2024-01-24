package controller

import (
	"net/http"
	"strconv"

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
}

func (s *SchedulesController) FindAllScheduleHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	users, paging, err := s.schedulesUC.FindAllScheduleUC(page, size)
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	common.SendPagedResponse(c, users, paging, "success")
}
