package dto

import "encoding/json"

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
