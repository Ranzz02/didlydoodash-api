package dto

import (
	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
)

type UpdateOrganisationInput struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	Website       *string `json:"website"`
	LogoUrl       *string `json:"logoUrl"`
	Location      *string `json:"location"`
	Timezone      *string `json:"timezone"`
	IsActive      *bool   `json:"isActive"`
	DefaultRoleID *string `json:"defaultRoleId"`
}

type UpdateOrganisationResponse struct {
	Organisation repository.Organisation `json:"organisation"`
}
