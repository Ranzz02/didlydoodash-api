package dto

import (
	"time"

	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
)

type OrganisationRolePermission struct {
	PermissionKey string `json:"key"`
	Allowed       bool   `json:"allowed"`
}

type OrganisationRole struct {
	ID          string                       `json:"id"`
	Name        string                       `json:"name"`
	Description *string                      `json:"description,omitempty"`
	BaseRoleID  *string                      `json:"base_role_id,omitempty"`
	Permissions []OrganisationRolePermission `json:"permissions,omitempty"`
}

type OrganisationMember struct {
	UserID   string           `json:"user_id"`
	Username string           `json:"username"`
	Email    string           `json:"email,omitempty"`
	JoinedAt time.Time        `json:"joined_at"`
	Role     OrganisationRole `json:"role"`
}

func NewOrganisationMember(user repository.User, member repository.OrganisationMember, role repository.Role) OrganisationMember {
	return OrganisationMember{
		UserID:   user.ID,
		Username: user.Username,
		JoinedAt: member.JoinedAt.Time,
		Role: OrganisationRole{
			ID:          role.ID,
			Name:        role.Name,
			Description: utils.PgTextToPtr(role.Description),
			BaseRoleID:  utils.PgTextToPtr(role.BaseRoleID),
		},
	}
}
