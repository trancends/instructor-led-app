package delivery

import (
	"database/sql"
	"fmt"
	"log"

	"enigmaCamp.com/instructor_led/config"
	"enigmaCamp.com/instructor_led/delivery/controller"
	"enigmaCamp.com/instructor_led/delivery/middleware"
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
	authMiddleware := middleware.NewAuthMiddleware(s.jwtService)
	controller.NewUserController(s.userUC, rg, authMiddleware).Route()
	controller.NewSchedulesController(s.scheduleUC, rg, authMiddleware).Route()
	controller.NewQuestionsController(s.questionsUC, rg, authMiddleware).Route()
	controller.NewAttandanceController(rg, s.attendanceUC, authMiddleware).Route()
	controller.NewAuthController(s.authUC, rg).Route()
}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.host); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v\n", err)
	}

	fmt.Println("Welcome to the Instructor Led App!")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)

	db, err := sql.Open(cfg.Driver, psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database connection: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v\n", err)
	}
	fmt.Println("Connected to Database")
	if err != nil {
		log.Fatalf("Config error: %v\n", err)
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
	authUC := usecase.NewAuthUseCase(userUseCase, jwtService)
	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	return &Server{
		scheduleUC:   schedulesUseCase,
		userUC:       userUseCase,
		questionsUC:  questionsUseCase,
		attendanceUC: attendanceUseCase,
		authUC:       authUC,
		jwtService:   jwtService,
		engine:       engine,
		host:         host,
	}
}
