package domain

type GetCardsInput struct {
	Where  *CardsWhere
	Limit  *int
	Offset *int
}

type CardsWhere struct {
	DeckID int64
}

type UpdateCardInput struct {
	Field UpdateCardField
}

type UpdateCardField struct {
	DeckID   *int64
	Question *string
	Answer   *string
}
