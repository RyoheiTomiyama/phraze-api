package domain

type GetCardsInput struct {
	Where  *CardsWhere
	Limit  *int
	Offset *int
}

type CardsWhere struct {
	DeckID int64
}
