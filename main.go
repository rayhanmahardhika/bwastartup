package main

import (
	"bwastartup/auth"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	// running router
	router.Run()
}

// fungsi middleware dibungkus agar bisa passing data
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		// get nilai header Authorization Bearer TOken
		authHeader := c.GetHeader("Authorization")
		// check bentuk hedaer auth yang dikirim
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse(
				"Unauthorized",
				http.StatusUnauthorized,
				"error",
				nil,
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response) // agar tidak dilanjutkan request nya
			return
		}

		// menghilangkan bearer
		tokenJWT := ""
		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) == 2 {
			tokenJWT = tokenString[1]
		}

		// validate token menggunakan service
		token, err := authService.ValidateToken(tokenJWT)
		if err != nil {
			response := helper.APIResponse(
				"Unauthorized",
				http.StatusUnauthorized,
				"error",
				nil,
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response) // agar tidak dilanjutkan request nya
			return
		}

		// check token claim
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse(
				"Unauthorized",
				http.StatusUnauthorized,
				"error",
				nil,
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response) // agar tidak dilanjutkan request nya
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse(
				"Unauthorized",
				http.StatusUnauthorized,
				"error",
				nil,
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response) // agar tidak dilanjutkan request nya
			return
		}

		// jika semua aman, akan disimpan dalam context
		c.Set("currentUser", user)
	}
}
