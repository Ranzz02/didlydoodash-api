package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/Stenoliv/didlydoodash_api/internal/config"
	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/Stenoliv/didlydoodash_api/internal/dto"
	"github.com/Stenoliv/didlydoodash_api/internal/repositories"
	"github.com/Stenoliv/didlydoodash_api/pkg/logging"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   *repositories.UserRepository
	tx     *repositories.TxManager
	cfg    *config.EnvConfig
	logger *logrus.Logger
}

func NewAuthService(repo *repositories.UserRepository, tx *repositories.TxManager, cfg *config.EnvConfig, logger *logrus.Logger) *AuthService {
	return &AuthService{
		repo:   repo,
		tx:     tx,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *AuthService) generateTokens(ctx context.Context, user repository.User, remember bool) (*dto.Tokens, error) {
	params := utils.TokenParams{
		UserID:     user.ID,
		RememberMe: remember,
	}

	// Generate access token
	access, err := utils.GenerateAccessToken(s.cfg, params)
	if err != nil {
		return nil, utils.NewError(http.StatusInternalServerError, "failed to generate access token", err)
	}

	// Generate refresh token
	refresh, err := utils.GenerateRefreshToken(s.cfg, params)
	if err != nil {
		return nil, utils.NewError(http.StatusInternalServerError, "failed to generate refresh token", err)
	}

	tokens := dto.Tokens{
		Access:  access,
		Refresh: refresh,
	}

	return &tokens, nil
}

func (s *AuthService) SignIn(ctx context.Context, params dto.SignInRequest) (*repository.User, *dto.Tokens, error) {
	logger := logging.WithLayer(ctx, "service").WithField("user_email", params.Email)
	logger.Info("attempting sign-in")

	user, err := s.repo.GetByEmail(ctx, params.Email)
	if err != nil {
		logger.WithError(err).Error("failed to retrieve user by email")
		return nil, nil, utils.NewError(http.StatusForbidden, "invalid email or password", err)
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		logger.WithError(err).Error("password does not match")
		return nil, nil, utils.NewError(http.StatusForbidden, "invalid email or password", err)
	}

	// Generate tokens
	tokens, err := s.generateTokens(ctx, user, params.Remember)
	if err != nil {
		logger.WithError(err).Error("token generation failed")
		return nil, nil, err
	}

	logger.WithField("user_id", user.ID).Info("sign-in successful")
	return &user, tokens, nil
}

func (s *AuthService) SignUp(ctx context.Context, params dto.SignUpRequest) (*repository.User, *dto.Tokens, error) {
	logger := logging.WithLayer(ctx, "service").WithField("new_user", params.Username)
	logger.Info("attempting sign-up")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.WithError(err).Error("failed to hash password")
		return nil, nil, utils.NewError(http.StatusInternalServerError, "failed to hash password", err)
	}

	args := repository.CreateUserParams{
		ID:       gonanoid.Must(),
		Email:    params.Email,
		Password: string(hashedPassword),
		Username: params.Username,
	}

	// Save user to database
	user, err := s.repo.CreateUser(ctx, args)
	if err != nil {
		logger.WithError(err).Error("failed to create user in database")
		return nil, nil, utils.NewError(http.StatusInternalServerError, "failed to save user to database", err)
	}

	// Generate tokens
	tokens, err := s.generateTokens(ctx, user, params.Remember)
	if err != nil {
		logger.WithError(err).Error("token generation failed")
		return nil, nil, err
	}

	logger.WithField("user_id", user.ID).Info("sign-up successful")
	return &user, tokens, nil
}

func (s *AuthService) Refresh(ctx context.Context, params dto.RefreshRequest) (*dto.Tokens, error) {
	logger := logging.WithLayer(ctx, "service")

	logger.Debug("validating refresh token")
	claims, err := utils.ValidateToken(s.cfg, params.Token, utils.RefreshToken)
	if err != nil {
		logger.WithError(err).Warn("invalid or expired refresh token")
		return nil, utils.NewError(http.StatusUnauthorized, "invalid or expired refresh token", err)
	}

	// Extract user ID (sub claim)
	sub, ok := (*claims)["sub"].(string)
	if !ok || sub == "" {
		logger.Warn("missing or invalid 'sub' claim in refresh token")
		return nil, utils.NewError(http.StatusUnauthorized, "invalid refresh token", errors.New("missing subject in token"))
	}
	logger = logger.WithField("user_id", sub)

	// Determine if the token had "rememberMe" flag
	rememberMe, _ := (*claims)["rememberMe"].(bool)

	logger.Debug("generating new tokens for user")
	tokens, err := s.generateTokens(ctx, repository.User{ID: sub}, rememberMe)
	if err != nil {
		logger.WithError(err).Error("failed to generate new tokens")
		return nil, utils.NewError(http.StatusInternalServerError, "failed to generate new tokens", err)
	}

	logger.Info("token refresh successful")
	return tokens, nil
}
