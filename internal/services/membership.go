package services

import (
	"context"
)

type MembershipService struct {
}

func NewMembershipService(ctx context.Context) *MembershipService {
	return &MembershipService{}
}
