package controller

import (
	"enigmaCamp.com/instructor_led/delivery/middleware"
	"enigmaCamp.com/instructor_led/mock/usecase_mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ScheduleControllerTestSuite struct {
	suite.Suite
	rg  *gin.RouterGroup
	sum *usecase_mock.SchedulesUseCaseMock
}

func (s *ScheduleControllerTestSuite) SetupTest() {
	s.sum = usecase_mock.NewSchedulesUseCaseMock()
	s.rg = gin.Default().Group("/schedules")
	controller.NewSchedulesController(s.sum, s.rg, middleware.AuthMiddleware{})
}
