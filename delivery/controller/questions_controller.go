package controller

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"enigmaCamp.com/instructor_led/delivery/middleware"
	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/shared/common"
	"enigmaCamp.com/instructor_led/usecase"
	"github.com/gin-gonic/gin"
)

type QuestionsController struct {
	questionsUC    usecase.QuestionsUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

// NewQuestionsController initializes a new QuestionsController.
func NewQuestionsController(questionsUC usecase.QuestionsUsecase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *QuestionsController {
	return &QuestionsController{
		questionsUC:    questionsUC,
		rg:             rg,
		authMiddleware: authMiddleware,
	}
}

func (q *QuestionsController) Route() {
	q.rg.GET("/questions", q.authMiddleware.RequireToken("ADMIN", "TRAINER"), q.GetQuestionsHandler)
	q.rg.GET("/questions/all", q.authMiddleware.RequireToken("ADMIN", "TRAINER"), q.ListQuestionsHandler)
	q.rg.POST("/questions", q.authMiddleware.RequireToken("ADMIN", "TRAINER", "PARTICIPANT"), q.CreateQuestionsHandler)
	q.rg.PATCH("/questions", q.authMiddleware.RequireToken("ADMIN", "TRAINER"), q.PatchQuestionsHandler)
	q.rg.DELETE("/questions/:id", q.authMiddleware.RequireToken("ADMIN"), q.DeleteQuestionsHandler)
}

func (q *QuestionsController) CreateQuestionsHandler(c *gin.Context) {
	var payload model.Question
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("QuestionsController.CreateQuestionsHandler:", err.Error())
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid json"+err.Error())
	}
	payloads, err := q.questionsUC.CreateQuestionsUC(payload)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, "failed to create questions"+err.Error())
	}
	common.SendSingleResponse(c, payloads, "questions created successfully")
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
		common.SendErrorResponse(c, http.StatusInternalServerError, "failed to list questions"+err.Error())
		return
	}
	common.SendSingleResponse(c, questions, "questions retrieved successfully")
}

func (q *QuestionsController) PatchQuestionsHandler(c *gin.Context) {
	var payload model.Question
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("QuestionsController.PatchQuestionsHandler:", err.Error())
		common.SendErrorResponse(c, http.StatusBadRequest, "invalid json"+err.Error())
		return
	}

	if payload.ID == "" || payload.Status == "" {
		common.SendErrorResponse(c, http.StatusBadRequest, "payload cannot be empty")
		return
	}

	err := q.questionsUC.UpdateQuestionStatus(payload)
	if err != nil {
		if err == sql.ErrNoRows {
			common.SendErrorResponse(c, http.StatusNotFound, "questions not found")
		} else {
			common.SendErrorResponse(c, http.StatusInternalServerError, "failed to update questions"+err.Error())
		}
		return
	}

	common.SendSingleResponse(c, "", "questions updated successfully")
}

func (q *QuestionsController) DeleteQuestionsHandler(c *gin.Context) {
	questionID := c.Param("id")
	if questionID == "" {
		common.SendErrorResponse(c, http.StatusBadRequest, "id cannot be empty")
		return
	}

	err := q.questionsUC.DeleteQuestion(questionID)
	if err != nil {
		if err == sql.ErrNoRows {
			common.SendErrorResponse(c, http.StatusNotFound, "questions not found")
		} else {
			common.SendErrorResponse(c, http.StatusInternalServerError, "failed to delete questions"+err.Error())
		}
		return
	}

	common.SendSingleResponse(c, "", "questions deleted successfully")
}
