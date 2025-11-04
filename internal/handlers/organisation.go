package handlers

import (
	"net/http"

	"github.com/Stenoliv/didlydoodash_api/internal/config"
	"github.com/Stenoliv/didlydoodash_api/internal/dto"
	"github.com/Stenoliv/didlydoodash_api/internal/middleware"
	"github.com/Stenoliv/didlydoodash_api/internal/services"
	"github.com/Stenoliv/didlydoodash_api/pkg/logging"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/gin-gonic/gin"
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
	logger := logging.WithLayer(ctx, "handler")

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

	orgs, err := h.service.List(c.Request.Context(), userID, search, pagination, ownerOnly)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.GetOrganisationsResponse{
		Organsations: orgs,
		Page:         page,
		Limit:        limit,
	})
}

func (h *OrganisationHandler) Get(c *gin.Context) {

}

func (h *OrganisationHandler) Update(c *gin.Context) {

}

func (h *OrganisationHandler) Delete(c *gin.Context) {

}
