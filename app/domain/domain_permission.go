package domain

type PermissionKey string

const (
	PermissionUnlimitedCardCreation       PermissionKey = "unlimitedCardCreation"
	PermissionUnlimitedAIAnswerGeneration PermissionKey = "unlimitedAIAnswerGeneration"
	PermissionLimitedAIAnswer100rpd       PermissionKey = "limitedAIAnswer100rpd"
)