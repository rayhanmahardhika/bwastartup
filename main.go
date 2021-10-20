package main

import (
	"bwastartup/auth"
	"bwastartup/handler"
	"bwastartup/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// connect DB
	dsn := "root:@tcp(127.0.0.1:3306)/bwa_startup_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	// inisiasi User Repo, Service, dan Handler
	// REPOSITORIES
	userRepository := user.NewRepository(db)

	// SERVICES
	userService := user.NewService(userRepository)
	authService := auth.NewService()

	// HANDLERS
	userHandler := handler.NewUserHandler(userService, authService)

	// inisiasi router
	router := gin.Default()

	// Routing
	// api versioning
	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email-checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)

	// running router
	router.Run()
}
