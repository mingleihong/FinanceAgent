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
	Password string `json:"password" binding:"required"`
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid login request")
		return
	}
	result, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			response.Error(c, http.StatusUnauthorized, "invalid username or password")
		default:
			log.Printf("login user: %v", err)
			response.Error(c, http.StatusInternalServerError, "login failed")
		}
		return
	}
	c.Header("Cache-Control", "no-store")
	response.Success(c, http.StatusOK, result)
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
