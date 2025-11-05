package handlers

import (
	"errors"
	"net/http"

	"github.com/Stenoliv/didlydoodash_api/internal/config"
	"github.com/Stenoliv/didlydoodash_api/internal/dto"
	"github.com/Stenoliv/didlydoodash_api/internal/middleware"
	"github.com/Stenoliv/didlydoodash_api/internal/services"
	"github.com/Stenoliv/didlydoodash_api/pkg/logging"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OrganisationHandler struct {
	service *services.OrganisationService
	cfg     *config.EnvConfig
}

// Create a new organisation handler
func NewOrganisationHandler(service *services.OrganisationService, cfg *config.EnvConfig) *OrganisationHandler {
	return &OrganisationHandler{
		service: service,
		cfg:     cfg,
	}
}

func (h *OrganisationHandler) Routes(rg *gin.RouterGroup) {
	org := rg.Group("/organisations")
	org.Use(middleware.AuthMiddleware(h.cfg))

	org.POST("", h.Create)
	org.GET("", h.GetAll)
	org.GET("/:id", h.Get)
	org.PUT("/:id", h.Update)
	org.DELETE("/:id", h.Delete)
}

// POST /organisations
func (h *OrganisationHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logging.WithLayer(ctx, "handler", "organisation")

	userID := utils.GetUserID(c)
	logger.Infof("user with id: %s is trying to create a organisation", userID)

	var body dto.CreateOrganisationInput
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithError(err).Warn("failed to bind input params")
		c.Error(utils.NewError(http.StatusBadRequest, "invalid input", err))
		return
	}

	org, err := h.service.Create(ctx, userID, body)
	if err != nil {
		logger.WithError(err).Warn("failed to create organisation")
		c.Error(err)
		return
	}

	// Organisation
	logger.Info("organisation successfully returned")
	c.JSON(http.StatusCreated, dto.CreateOrganisationResponse{
		Organisation: *org,
	})
}

// GET /organisations
func (h *OrganisationHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logging.WithLayer(ctx, "handler", "organisation")

	userID := utils.GetUserID(c)

	search := c.Query("search")
	page := utils.ParseIntDefault(c.Query("page"), 1)
	limit := utils.ParseIntDefault(c.Query("limit"), 10)

	ownerOnly := utils.ParseBoolDefault(c.Query("ownerOnly"), false)

	offset := (page - 1) * limit

	pagination := services.Pagination{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	logger.Info("trying to fetch organisations")

	orgs, err := h.service.List(c.Request.Context(), userID, search, pagination, ownerOnly)
	if err != nil {
		logger.WithError(err).Warn("failed to get organisations")
		c.Error(err)
		return
	}

	logger.Info("organisations successfully fetched")
	c.JSON(http.StatusOK, dto.GetOrganisationsResponse{
		Organisations: orgs,
		Page:          page,
		Limit:         limit,
	})
}

// GET /organisation/{id}
func (h *OrganisationHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logging.WithLayer(ctx, "handler", "organisation")

	// Parse organisation ID from URL
	orgID := c.Param("id")
	if orgID == "" {
		logger.Warn("organisation id not provided in path")
		c.Error(utils.NewError(http.StatusBadRequest, "organisation id required", errors.New("missing organisation id")))
		return
	}

	// Parse userId from request
	userID := utils.GetUserID(c)
	if userID == "" {
		logger.Warn("user id missing in request context")
		c.Error(utils.NewError(http.StatusUnauthorized, "unauthorized", errors.New("missing user id in context")))
		return
	}

	// Set logger context
	logger = logger.WithFields(logrus.Fields{
		"org_id":  orgID,
		"user_id": userID,
	})

	logger.Info("attempting to fetch organisation")

	organisation, err := h.service.Get(ctx, orgID, userID)
	if err != nil {
		logger.WithError(err).Warn("failed to get organisation")
		c.Error(err)
		return
	}

	logger.Info("organisation successfully fetched")
	c.JSON(http.StatusOK, dto.GetOrganisationResponse{
		Organisation: *organisation,
	})
}

// PUT /organisation/{id}
func (h *OrganisationHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logging.WithLayer(ctx, "handler", "organisation")

	// Parse organisation ID from url
	orgID := c.Param("id")
	if orgID == "" {
		logger.Warn("organisation id not provided in path")
		c.Error(utils.NewError(http.StatusBadRequest, "organisation id required", errors.New("missing organisation id")))
		return
	}

	// Parse userId from request
	userID := utils.GetUserID(c)
	// Set logger context
	logger = logger.WithFields(logrus.Fields{
		"org_id":  orgID,
		"user_id": userID,
	})

	var body dto.UpdateOrganisationInput
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithError(err).Warn("failed to parse json")
		c.Error(utils.NewError(http.StatusBadRequest, "invalid input", err))
		return
	}

	logger.Info("attempting to update organisation")

	updated, err := h.service.Update(ctx, orgID, userID, body)
	if err != nil {
		logger.WithError(err).Warn("failed to update organisation")
		c.Error(err)
		return
	}

	logger.Info("organisation updated successfully")
	c.JSON(http.StatusOK, dto.UpdateOrganisationResponse{
		Organisation: *updated,
	})
}

func (h *OrganisationHandler) Delete(c *gin.Context) {

}
