package handlers

import (
	"net/http"

	"github.com/Stenoliv/didlydoodash_api/internal/config"
	"github.com/Stenoliv/didlydoodash_api/internal/dto"
	"github.com/Stenoliv/didlydoodash_api/internal/middleware"
	"github.com/Stenoliv/didlydoodash_api/internal/services"
	"github.com/Stenoliv/didlydoodash_api/pkg/logging"
	"github.com/Stenoliv/didlydoodash_api/pkg/permissions"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MembershipHandlerServices struct {
	Member       *services.MembershipService
	Organisation *services.OrganisationService
	Checker      *services.Checker
}

type MembershipHandler struct {
	services *MembershipHandlerServices
	cfg      *config.EnvConfig
}

func NewMembershipHandler(services MembershipHandlerServices, cfg *config.EnvConfig) *MembershipHandler {
	return &MembershipHandler{
		services: &services,
		cfg:      cfg,
	}
}

func (h *MembershipHandler) Routes(router *gin.RouterGroup) {
	base := router.Group("/organisations/:id")
	base.Use(middleware.AuthMiddleware(h.cfg))

	membership := base.Group("/members")
	roles := base.Group("")

	// Members
	membership.GET("", middleware.RequirePermission(h.services.Checker, permissions.OrgViewMembers), h.GetMembers)
	membership.POST("", middleware.RequirePermission(h.services.Checker, permissions.OrgInviteMembers), h.CreateMember)

	// Roles
	roles.GET("/permissions", h.EffectivePermissions)
}

// Members
func (h *MembershipHandler) GetMembers(c *gin.Context) {
	ctx := c.Request.Context()
	orgID := c.Param("id")
	userID := utils.GetUserID(c)

	logger := logging.WithLayer(ctx, "handler", "membership").WithFields(logrus.Fields{
		"org_id":  orgID,
		"user_id": userID,
	})

	logger.Info("trying to fetch organisation members")
	// Try to get members
	members, err := h.services.Member.List(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	logger.Infof("fetched %d organisation members", len(members))
	c.JSON(http.StatusOK, members)
}

func (h *MembershipHandler) CreateMember(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logging.WithLayer(ctx, "handler", "membership")

	var body dto.CreateOrganisationMember
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithError(err).Warn("invalid input provided")
		c.Error(utils.NewError(http.StatusBadRequest, "invalid input", err))
		return
	}

	logger.Info("trying to create a new organisation member")

	member, err := h.services.Member.Create(ctx, &body)
	if err != nil {
		logger.WithError(err).Warn("failed to create organisation member")
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.CreateOrganisationMemberResponse{
		Member: *member,
	})
}

// Roles & Permissions
func (h *MembershipHandler) EffectivePermissions(c *gin.Context) {
	ctx := c.Request.Context()

	orgID := c.Param("id")
	userID := utils.GetUserID(c)

	logger := logging.WithLayer(ctx, "handler", "membership").WithFields(logrus.Fields{
		"org_id":  orgID,
		"user_id": userID,
	})

	logger.Info("getting user permissions in organisation")

	// Fetch role and permissions
	role, perms, err := h.services.Member.GetUserPermissions(ctx, userID, orgID)
	if err != nil {
		c.Error(err)
		return
	}

	// Convert to DTO
	permKeys := make([]string, 0, len(perms))
	for _, p := range perms {
		permKeys = append(permKeys, p.PermissionKey)
	}

	output := dto.GetEffectivePermissionsResponse{
		Role: dto.OrganisationRole{
			ID:          role.ID,
			Name:        role.Name,
			Description: utils.PgTextToPtr(role.Description),
		},
		Permissions: permKeys,
	}

	logger.WithField("perm_count", len(permKeys)).Info("successfully fetched permissions")

	c.JSON(http.StatusOK, output)
}

// ------------- HELPERS --------------
