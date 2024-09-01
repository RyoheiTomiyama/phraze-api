package domain

import "context"

type PermissionKey string

const (
	PermissionUnlimitedCardCreation       PermissionKey = "unlimitedCardCreation"
	PermissionUnlimitedAIAnswerGeneration PermissionKey = "unlimitedAIAnswerGeneration"
	PermissionLimitedAIAnswer100rpd       PermissionKey = "limitedAIAnswer100rpd"
)

func (pk PermissionKey) String() string {
	return string(pk)
}

type Permission struct {
	ID   int64
	Key  string
	Name string
}

type Permissions []*Permission

func (pp Permissions) HasKey(ctx context.Context, k PermissionKey) bool {
	for _, p := range pp {
		if p.Key == k.String() {
			return true
		}
	}

	return false
}
