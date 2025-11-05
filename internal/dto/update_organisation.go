package dto

import (
	"encoding/json"

	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
)

type UpdateOrganisationInput struct {
	Name        *string          `json:"name"`
	Description *string          `json:"description"`
	Website     *string          `json:"website"`
	LogoUrl     *string          `json:"logoUrl"`
	Location    *string          `json:"location"`
	Timezone    *string          `json:"timezone"`
	IsActive    *bool            `json:"isActive"`
	Settings    *json.RawMessage `json:"settings"`
}

type UpdateOrganisationResponse struct {
	Organisation repository.Organisation `json:"organisation"`
}
