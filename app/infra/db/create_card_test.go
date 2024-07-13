package db

import (
	"context"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/stretchr/testify/assert"
)

func TestCreateCard(t *testing.T) {
	ctx := context.Background()

	db := db_test.GetDB(t)
	defer db.Close()

	client := NewTestClient(t, db)

	fx := fixture.New(db)
	decks := fx.CreateDeck(t, &fixture.DeckInput{})

	t.Run("正常系", func(t *testing.T) {
		newCard := &domain.Card{
			DeckID:   decks[0].ID,
			Question: "question",
			Answer:   "answer",
		}
		card, err := client.CreateCard(ctx, newCard)

		assert.NoError(t, err)
		assert.Equal(t, newCard.DeckID, card.DeckID)
		assert.Equal(t, newCard.Question, card.Question)
		assert.Equal(t, newCard.Answer, card.Answer)
	})

	t.Run("存在しないDeckIDで登録しようとした場合", func(t *testing.T) {
		newCard := &domain.Card{
			DeckID:   -1,
			Question: "question",
			Answer:   "answer",
		}
		card, err := client.CreateCard(ctx, newCard)
		assert.Error(t, err)
		assert.Nil(t, card)
	})
}
