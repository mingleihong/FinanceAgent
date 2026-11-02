package handler

import (
	"Finance/internal/service"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": 1, "message": "invalid registration request"})
		return
	}
	user, err := h.userService.Register(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRegistration):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		case errors.Is(err, service.ErrUsernameTaken):
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"code": 1, "message": "username already exists"})
		default:
			log.Printf("register user: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "registration failed"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": user})
}
