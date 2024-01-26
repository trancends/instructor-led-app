package controller

import (
	"net/http"

	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/usecase"
	"github.com/gin-gonic/gin"
)

type AttandanceController struct {
	attendanceUC usecase.AttendanceUsecase
	rg           *gin.RouterGroup
}

func NewAttandanceController(rg *gin.RouterGroup, attendanceUC usecase.AttendanceUsecase) *AttandanceController {
	return &AttandanceController{
		attendanceUC: attendanceUC,
		rg:           rg,
	}
}

func (a *AttandanceController) Route() {
	a.rg.GET("/attendance", a.GetAllAttendanceHandler)
	a.rg.GET("/attendance/:id", a.GetAttendanceByID)
	a.rg.POST("/attendance", a.AddAttendanceHandler)
}

func (a *AttandanceController) GetAllAttendanceHandler(c *gin.Context) {
	// Retrieve data from usecase
	attendanceList, err := a.attendanceUC.ListAttendances()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendance"})
		return
	}

	// Respond with the retrieved data
	c.JSON(http.StatusOK, gin.H{"attendance": attendanceList})
}

func (a *AttandanceController) GetAttendanceByID(c *gin.Context) {
	id := c.Param("id")
	attendance, err := a.attendanceUC.GetAttendance(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attendance": attendance})
}

func (a *AttandanceController) AddAttendanceHandler(c *gin.Context) {
	var attendance model.Attendance
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdAttendance, err := a.attendanceUC.AddAttendance(attendance.UserID, attendance.ScheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attendance": createdAttendance})
}
