package dto

import "github.com/Stenoliv/didlydoodash_api/internal/db/repository"

type CreateOrganisationInput struct {
	Name string `json:"name" binding:"required"`
}

type CreateOrganisationResponse struct {
	Organisation repository.Organisation `json:"organisation"`
}
