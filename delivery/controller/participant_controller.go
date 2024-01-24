package controller

import (
	"net/http"
	"strconv"

	"enigmaCamp.com/instructor_led/shared/common"
	"enigmaCamp.com/instructor_led/usecase"
	"github.com/gin-gonic/gin"
)

type ParticipantController struct {
	participantUC usecase.ParticipantUseCase
	rg            *gin.RouterGroup
}

func NewParticipantController(participantUC usecase.ParticipantUseCase, rg *gin.RouterGroup) *ParticipantController {
	return &ParticipantController{
		participantUC: participantUC,
	}
}

func (p *ParticipantController) Route() {
	p.rg.GET("/schedules", p.FindAllScheduleHandler)
}

func (p *ParticipantController) FindAllScheduleHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	users, paging, err := p.participantUC.FindAllScheduleUC(page, size)
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	common.SendPagedResponse(c, users, paging, "success")
}
