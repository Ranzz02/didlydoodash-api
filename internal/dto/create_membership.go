package dto

type CreateOrganisationMember struct {
	UserID string `json:"userId" binding:"required"`
	OrgID  string `json:"orgId" binding:"required"`
	RoleID string `json:"roleId"`
}

type CreateOrganisationMemberResponse struct {
	Member OrganisationMember `json:"member"`
}
