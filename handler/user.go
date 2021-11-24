package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

// menangkap API Register
func (h *userHandler) RegisterUser(c *gin.Context) {
	// tangkap input dari user
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValdiationError(err)
		errorMessage := gin.H{"errors": errors} // gin.H adalah map

		response := helper.APIResponse(
			"Account registere failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
		)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse(
			"Account register failed",
			http.StatusBadRequest,
			"error",
			nil,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse(
			"Account register failed",
			http.StatusBadRequest,
			"error",
			nil,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)

	response := helper.APIResponse(
		"Account successfuly registered",
		http.StatusOK,
		"success",
		formatter,
	)

	c.JSON(http.StatusOK, response)
}

// menangani API Login
func (h *userHandler) Login(c *gin.Context) {
	// input login
	//menangkap input ke handler
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValdiationError(err)
		errorMessage := gin.H{"errors": errors} // gin.H adalah map

		response := helper.APIResponse(
			"Login failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
		)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()} // gin.H adalah map
		response := helper.APIResponse(
			"Login failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
		)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedUser.ID)
	if err != nil {
		response := helper.APIResponse(
			"Login failed",
			http.StatusBadRequest,
			"error",
			nil,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedUser, token)
	response := helper.APIResponse(
		"Account successfuly logged",
		http.StatusOK,
		"success",
		formatter,
	)

	c.JSON(http.StatusOK, response)
}

// menangani email availability
func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValdiationError(err)
		errorMessage := gin.H{"errors": errors} // gin.H adalah map

		response := helper.APIResponse(
			"Email already registered.",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
		)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server Error"}
		response := helper.APIResponse(
			"Email already registered.",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
		)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isAvailable,
	}

	metaMessage := "Email already registered"
	if isAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(
		metaMessage,
		http.StatusOK,
		"success",
		data,
	)
	c.JSON(http.StatusOK, response)

}

// menangani API upload avatar + middleware
func (h *userHandler) UploadAvatar(c *gin.Context) {
	// get file dari form data
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse(
			"Failed to upload avatar Image",
			http.StatusBadRequest,
			"error",
			data,
		)

		c.JSON(http.StatusBadRequest, response)
		return
	}
	// buat file path + nama lalu save ke local server
	// get contenxt + di casting ke type User
	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID
	path := fmt.Sprintf("images/%d-%s", userId, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse(
			"Failed to upload avatar Image",
			http.StatusBadRequest,
			"error",
			data,
		)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userId, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse(
			"Failed to upload avatar Image",
			http.StatusBadRequest,
			"error",
			data,
		)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse(
		"Avatar successfuly uploaded",
		http.StatusOK,
		"success",
		data,
	)

	c.JSON(http.StatusOK, response)
}
