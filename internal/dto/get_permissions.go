package dto

type GetEffectivePermissionsResponse struct {
	Role        OrganisationRole `json:"role"`
	Permissions []string         `json:"permissions"`
}
