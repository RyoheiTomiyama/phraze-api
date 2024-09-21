package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestDeleteDeck(t *testing.T) {
	db := db_test.GetDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fx := fixture.New(db)

	t.Run("正常系", func(t *testing.T) {
		cases := []struct {
			name    string
			arrange func() (deckID int64)
			assert  func(result int64, err error, deckID int64)
		}{
			{
				name: "紐づくカードがない場合",
				arrange: func() (deckID int64) {
					decks := fx.CreateDeck(t, &fixture.DeckInput{UserID: lo.ToPtr("own")})

					return decks[0].ID
				},
				assert: func(result int64, err error, deckID int64) {
					assert.NoError(t, err)
					assert.Equal(t, int64(1), result)

					var c model.Deck
					err = db.Get(&c, "SELECT * FROM decks WHERE id=$1", deckID)
					assert.Equal(t, sql.ErrNoRows, err)
				},
			},
			{
				name: "カードを持っている場合",
				arrange: func() (deckID int64) {
					decks := fx.CreateDeck(t, &fixture.DeckInput{UserID: lo.ToPtr("own")})
					fx.CreateCard(t, decks[0].ID, fixture.CardInput{})

					return decks[0].ID
				},
				assert: func(result int64, err error, deckID int64) {
					assert.NoError(t, err)
					assert.Equal(t, int64(2), result)

					var c model.Deck
					err = db.Get(&c, "SELECT * FROM decks WHERE id=$1", deckID)
					assert.Equal(t, sql.ErrNoRows, err)
				},
			},
			{
				name: "学習したことあるデッキの場合",
				arrange: func() (deckID int64) {
					decks := fx.CreateDeck(t, &fixture.DeckInput{UserID: lo.ToPtr("own")})
					cards := fx.CreateCard(t, decks[0].ID, fixture.CardInput{})
					fx.CreateCardSchedule(t, fixture.CardScheduleInput{CardID: cards[0].ID})
					fx.CreateCardReview(t, fixture.CardReviewInput{CardID: cards[0].ID})

					return decks[0].ID
				},
				assert: func(result int64, err error, deckID int64) {
					assert.NoError(t, err)
					assert.Equal(t, int64(4), result)

					var c model.Deck
					err = db.Get(&c, "SELECT * FROM decks WHERE id=$1", deckID)
					assert.Equal(t, sql.ErrNoRows, err)
				},
			},
		}

		client := NewTestClient(t, db)

		for _, cs := range cases {
			id := cs.arrange()
			var result int64
			err := client.Tx(context.Background(), func(ctx context.Context) error {
				r, err := client.DeleteDeck(ctx, id)
				if err != nil {
					return err
				}
				result = r

				return nil
			})
			cs.assert(result, err, id)
		}

	})

	t.Run("存在しないデッキを削除しようとした場合", func(t *testing.T) {
		client := NewTestClient(t, db)
		var result int64
		err := client.Tx(context.Background(), func(ctx context.Context) error {
			r, err := client.DeleteDeck(ctx, -1)
			if err != nil {
				return err
			}
			result = r

			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(0), result)
	})
}
