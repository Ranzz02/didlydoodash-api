package dto

import "github.com/Stenoliv/didlydoodash_api/internal/db/repository"

type GetOrganisationsResponse struct {
	Organsations []repository.Organisation `json:"organisations"`
	Page         int                       `json:"page"`
	Limit        int                       `json:"limit"`
}
