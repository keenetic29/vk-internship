package handlers

import (
	"VK/internal/domain"
	"VK/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
    Register(username, password string) (*domain.User, error)
    Login(username, password string) (string, error)
    ValidateToken(token string) (uint, error)
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	logger.Log.Info("Register request received",
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
		"client_ip", c.ClientIP(),
	)

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid registration request",
			"error", err.Error(),
			"username", req.Username,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Debug("Attempting to register user",
		"username", req.Username,
	)

	user, err := h.authService.Register(req.Username, req.Password)
	if err != nil {
		logger.Log.Error("Registration failed",
			"error", err.Error(),
			"username", req.Username,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("User registered successfully",
		"user_id", user.ID,
		"username", user.Username,
	)

	c.JSON(http.StatusCreated, user)
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	logger.Log.Info("Login request received",
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
		"client_ip", c.ClientIP(),
	)

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid login request",
			"error", err.Error(),
			"username", req.Username,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Debug("Attempting to authenticate user",
		"username", req.Username,
	)

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		logger.Log.Warn("Login failed",
			"error", err.Error(),
			"username", req.Username,
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("User logged in successfully",
		"username", req.Username,
		"token_prefix", token[:10]+"...", // логирую только префикс токена
	)

	c.Header("Authorization", token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}