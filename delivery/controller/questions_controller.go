// controller/questions_controller.go
package controller

import (
	"log"
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
	q.rg.GET("/questions", q.ListQuestionsHandler)
}

func (q *QuestionsController) GetQuestionsHandler(c *gin.Context) {

	date := c.Query("date")

	// Call the use case to retrieve a list of schedules based on the given date
	schedules, err := q.questionsUC.GetQuestion(date)
	log.Println(schedules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve questions"})
		return
	}

	// Return the list of schedules as JSON
	c.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

func (q *QuestionsController) ListQuestionsHandler(c *gin.Context) {
	questions, err := q.questionsUC.ListQuestions()
	log.Println(questions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve questions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"questions": questions})
}
