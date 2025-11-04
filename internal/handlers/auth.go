package handlers

import (
	"errors"
	"net/http"

	"github.com/Stenoliv/didlydoodash_api/internal/config"
	"github.com/Stenoliv/didlydoodash_api/internal/dto"
	"github.com/Stenoliv/didlydoodash_api/internal/services"
	"github.com/Stenoliv/didlydoodash_api/pkg/logging"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
	cfg     *config.EnvConfig
}

func NewAuthHandler(service *services.AuthService, cfg *config.EnvConfig) *AuthHandler {
	return &AuthHandler{
		service: service,
		cfg:     cfg,
	}
}

func (h *AuthHandler) Routes(router *gin.RouterGroup) {
	auth := router.Group("/auth")

	auth.POST("/signin", h.SignIn)
	auth.POST("/signup", h.SignUp)
	auth.POST("/refresh", h.Refresh)
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var body dto.SignInRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(utils.NewError(http.StatusBadRequest, "invalid input", err))
		return
	}

	// SignIn in service layer
	user, tokens, err := h.service.SignIn(c.Request.Context(), body)
	if err != nil {
		c.Error(err)
		return
	}

	// Response
	c.JSON(http.StatusOK, dto.AuthResponse{
		User: dto.UserResponse{
			ID:       user.ID,
			Username: user.Email,
			Email:    user.Email,
		},
		Tokens: *tokens,
	})
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var body dto.SignUpRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(utils.NewError(http.StatusBadRequest, "invalid input", err))
		return
	}

	// SignUp in service layer
	user, tokens, err := h.service.SignUp(c.Request.Context(), body)
	if err != nil {
		c.Error(err)
		return
	}

	// Response
	c.JSON(http.StatusCreated, dto.AuthResponse{
		User: dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
		Tokens: *tokens,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logging.WithLayer(ctx, "handler")

	logger.Info("extracting token from context")
	token := utils.ExtractToken(c)
	if token == "" {
		logger.Warn("refresh token not found in request")
		c.Error(utils.NewError(http.StatusUnauthorized, "no token provided", errors.New("no token provided in request")))
		return
	}

	logger.Debug("sending token to service layer to validate and generate new tokens")
	tokens, err := h.service.Refresh(ctx, dto.RefreshRequest{
		Token: token,
	})
	if err != nil {
		c.Error(err)
		return
	}

	// Respond with new tokens
	logger.WithField("user_id", tokens.UserID).Info("successfully refreshed user's tokens")
	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
