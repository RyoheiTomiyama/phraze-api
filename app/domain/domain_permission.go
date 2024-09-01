package domain

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
	Key  string
	Name string
}
