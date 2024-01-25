package controller

import (
	"net/http"

	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/shared/common"
	"enigmaCamp.com/instructor_led/usecase"
	"github.com/gin-gonic/gin"
)

type QuestionsController struct {
	questionsUC usecase.QuestionsUsecase
	rg          *gin.RouterGroup
}

func NewQuestionsController(questionsUC usecase.QuestionsUsecase, rg *gin.RouterGroup) *QuestionsController {
	return &QuestionsController{
		questionsUC: questionsUC,
		rg:          rg,
	}
}

func (q *QuestionsController) Route() {
	q.rg.POST("/questions", q.CreateQuestionsHandler)
}

func (q *QuestionsController) CreateQuestionsHandler(c *gin.Context) {
	var payload model.Questions
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid json"+err.Error())
	}
	payloads, err := q.questionsUC.CreateQuestionsUC(payload)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, "failed to create questions"+err.Error())
	}
	common.SendSingleResponse(c, payloads, "questions created successfully")
}
