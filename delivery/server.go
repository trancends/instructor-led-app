package delivery

import (
	"database/sql"
	"fmt"
	"log"

	"enigmaCamp.com/instructor_led/config"
	"enigmaCamp.com/instructor_led/delivery/controller"
	repository "enigmaCamp.com/instructor_led/repository"
	"enigmaCamp.com/instructor_led/shared/service"
	"enigmaCamp.com/instructor_led/usecase"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type Server struct {
	scheduleUC   usecase.ShecdulesUseCase
	userUC       usecase.UserUsecase
	questionsUC  usecase.QuestionsUsecase
	attendanceUC usecase.AttendanceUsecase
	authUC       usecase.AuthUseCase
	jwtService   service.JwtService
	engine       *gin.Engine
	host         string
}

func (s *Server) initRoute() {
	log.Println("init route")
	route := s.engine
	route.Static("/scheduleImages", "./scheduleImages")
	rg := route.Group("/api/v1")
	controller.NewUserController(s.userUC, rg).Route()
	controller.NewSchedulesController(s.scheduleUC, rg).Route()
	controller.NewQuestionsController(s.questionsUC, rg).Route()
	controller.NewAttandanceController(rg, s.attendanceUC).Route()
	controller.NewAuthController(s.authUC, rg).Route()
}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("failed to start server %v", err))
	}
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Welcome to the Instructor Led App!")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)

	db, err := sql.Open(cfg.Driver, psqlInfo)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Database")
	if err != nil {
		log.Fatal(fmt.Errorf("config error: %v", err))
	}

	userRepository := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUsecase(userRepository)
	schedulesRepository := repository.NewSchedulesRepository(db)
	schedulesUseCase := usecase.NewSchedulesUseCase(schedulesRepository)
	questionRepository := repository.NewQuestionsRepository(db)
	questionsUseCase := usecase.NewQuestionsUsecase(questionRepository)
	attendanceRepository := repository.NewAttendanceRepository(db)
	attendanceUseCase := usecase.NewAttendanceUsecase(attendanceRepository)
	jwtService := service.NewJwtService(cfg.TokenConfig)
	authUseCase := usecase.NewAuthUseCase(userUseCase, jwtService)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	return &Server{
		scheduleUC:   schedulesUseCase,
		userUC:       userUseCase,
		questionsUC:  questionsUseCase,
		attendanceUC: attendanceUseCase,
		authUC:       authUseCase,
		engine:       engine,
		host:         host,
	}
}
