package handler

import (
	"Finance/internal/response"
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
		response.Error(c, http.StatusBadRequest, "invalid registration request")
		return
	}
	user, err := h.userService.Register(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRegistration):
			response.Error(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrUsernameTaken):
			response.Error(c, http.StatusConflict, "username already exists")
		default:
			log.Printf("register user: %v", err)
			response.Error(c, http.StatusInternalServerError, "registration failed")
		}
		return
	}
	response.Success(c, http.StatusCreated, user)
}
