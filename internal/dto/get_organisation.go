package dto

import "github.com/Stenoliv/didlydoodash_api/internal/db/repository"

type GetOrganisationsResponse struct {
	Organisations []repository.Organisation `json:"organisations"`
	Page          int                       `json:"page"`
	Limit         int                       `json:"limit"`
}

type GetOrganisationResponse struct {
	Organisation repository.Organisation `json:"organisation"`
}
