// controller/questions_controller.go
package controller

import (
	"net/http"

	"enigmaCamp.com/instructor_led/usecase"
	"github.com/gin-gonic/gin"
)

type QuestionsController struct {
	questionsUC usecase.QuestionsUsecase
	rg          *gin.RouterGroup
}

// NewQuestionsController initializes a new QuestionsController.
func NewQuestionsController(questionsUC usecase.QuestionsUsecase, rg *gin.RouterGroup) *QuestionsController {
	return &QuestionsController{
		questionsUC: questionsUC,
		rg:          rg,
	}
}

func (q *QuestionsController) Route() {
	q.rg.GET("/question", q.GetQuestionsHandler)
}

func (q *QuestionsController) GetQuestionsHandler(c *gin.Context) {
	// Call the use case to retrieve a list of questions
	questions, err := q.questionsUC.ListQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve questions"})
		return
	}

	// Return the list of questions as JSON
	c.JSON(http.StatusOK, gin.H{"questions": questions})
}
