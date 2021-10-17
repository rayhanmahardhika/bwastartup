package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

// menangkap API
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
		errors := helper.FormatValdiationError(err)
		errorMessage := gin.H{"errors": errors} // gin.H adalah map

		response := helper.APIResponse(
			"Account registere failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, "tokentoken")

	response := helper.APIResponse(
		"Account successfuly registered",
		http.StatusOK,
		"success",
		formatter,
	)

	c.JSON(http.StatusOK, response)
}
