package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/user/config"
	handlers "github.com/tanush-128/openzo_backend/user/internal/api"
	"github.com/tanush-128/openzo_backend/user/internal/middlewares"
	userpb "github.com/tanush-128/openzo_backend/user/internal/pb"
	"github.com/tanush-128/openzo_backend/user/internal/repository"
	"github.com/tanush-128/openzo_backend/user/internal/service"
)

type Server struct {
	userpb.UserServiceServer
	userRepository repository.UserRepository
	userService    service.UserService
}

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
	}

	db, err := connectToDB(cfg) // Implement database connection logic
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to database: %w", err))
	}

	userRepository := repository.NewUserRepository(db)
	userDataRepository := repository.NewUserDataRepository(db)

	userService := service.NewUserService(userRepository)
	userDataService := service.NewUserDataService(userDataRepository)

	go service.GrpcServer(cfg, &service.Server{UserRepository: userRepository, UserService: userService})
	// Initialize HTTP server with Gin
	router := gin.Default()
	handler := handlers.NewHandler(&userService, &userDataService)

	router.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong v3 trying to deploy",
		})
	})
	router.POST("/", handler.CreateUser)
	router.GET("/:id", handler.GetUserByID)
	router.GET("/email/:email", handler.GetUserByEmail)
	router.PUT("/", handler.UpdateUser)
	router.POST("/signin", handler.UserSignIn)
	router.POST("/userdata", handler.CreateUserData)
	router.GET("/userdata/:id", handler.GetUserDataByID)
	router.PUT("/userdata", handler.UpdateUserData)
	router.DELETE("/userdata/:id", handler.DeleteUserData)
	router.Use(middlewares.JwtMiddleware)
	router.GET("/jwt", handler.GetUserWithJWT)

	router.Run(fmt.Sprintf(":%s", cfg.HTTPPort))

}
