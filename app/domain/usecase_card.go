package domain

type GetCardsInput struct {
	Where  *CardsWhere
	Limit  *int
	Offset *int
}
type GetPendingCardsInput struct {
	Where  *CardsWhere
	Limit  *int
	Offset *int
}

type CardsWhere struct {
	DeckID *int64
	UserID *string
}

type UpdateCardInput struct {
	Field UpdateCardField
}

type UpdateCardField struct {
	DeckID   *int64
	Question *string
	Answer   *string
	AIAnswer *string
}
