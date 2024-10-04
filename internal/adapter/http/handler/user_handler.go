package handler

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/service"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetAll(c *gin.Context)
	Create(c *gin.Context)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *userHandler) Create(c *gin.Context) {
	var input dto.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]string, len(ve))
			for i, fe := range ve {
				out[i] = fmt.Sprintf("Field '%s' failed on the '%s' tag", fe.Field(), fe.Tag())
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": out})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map DTO to domain entity
	user := entity.User{
		Name:  input.Name,
		Email: input.Email,
		Age:   input.Age,
	}

	createdUser, err := h.userService.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}
