// controller/questions_controller.go
package controller

import (
	"log"
	"net/http"
	"time"

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
	q.rg.GET("/question/all", q.ListQuestionsHandler)
}

func (q *QuestionsController) GetQuestionsHandler(c *gin.Context) {
	date := c.Query("date")

	// Validasi format tanggal
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Printf("Invalid date format: %v\n", err)
		if date == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Date is required"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Expected format: 2006-01-02"})
		return
	}

	schedules, err := q.questionsUC.GetQuestion(date)
	log.Println(schedules)
	if err != nil {
		log.Printf("Error retrieving schedules for date %s: %v\n", date, err)
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
