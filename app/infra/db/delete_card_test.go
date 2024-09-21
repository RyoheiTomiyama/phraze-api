package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/test/assertion"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCard(t *testing.T) {
	db := db_test.GetDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fx := fixture.New(db)
	decks := fx.CreateDeck(
		t,
		&fixture.DeckInput{UserID: lo.ToPtr("own")},
	)

	t.Run("正常系", func(t *testing.T) {
		cases := []struct {
			name    string
			arrange func() (cardId int64)
			assert  func(result int64, err error, cardID int64)
		}{
			{
				name: "cardsのみの場合",
				arrange: func() (cardId int64) {
					cards := fx.CreateCard(t, decks[0].ID, fixture.CardInput{})

					return cards[0].ID
				},
				assert: func(result int64, err error, cardID int64) {
					assert.NoError(t, err)
					assert.Equal(t, int64(1), result)

					var c model.Card
					err = db.Get(&c, "SELECT * FROM cards WHERE id=$1", cardID)
					assert.Equal(t, sql.ErrNoRows, err)
				},
			},
			{
				name: "card_reviews, card_schedulesと紐づいているカードの場合",
				arrange: func() (cardId int64) {
					cards := fx.CreateCard(t, decks[0].ID, fixture.CardInput{})
					fx.CreateCardSchedule(t, fixture.CardScheduleInput{CardID: cards[0].ID})
					fx.CreateCardReview(t, fixture.CardReviewInput{CardID: cards[0].ID})

					return cards[0].ID
				},
				assert: func(result int64, err error, cardID int64) {
					assert.NoError(t, err)
					assert.Equal(t, int64(3), result)

					var c model.Card
					err = db.Get(&c, "SELECT * FROM cards WHERE id=$1", cardID)
					assert.Equal(t, sql.ErrNoRows, err)
				},
			},
		}

		client := NewTestClient(t, db)

		for _, cs := range cases {
			id := cs.arrange()
			var result int64
			err := client.Tx(context.Background(), func(ctx context.Context) error {
				r, err := client.DeleteCard(ctx, id)
				if err != nil {
					return err
				}
				result = r

				return nil
			})
			cs.assert(result, err, id)
		}
	})

	t.Run("トランザクションの外でメソッドを実行した場合", func(t *testing.T) {
		cards := fx.CreateCard(t, decks[0].ID, fixture.CardInput{})
		client := NewTestClient(t, db)

		result, err := client.DeleteCard(context.Background(), cards[0].ID)

		assertion.AssertError(t, "transaction内で実行してください", errutil.CodeInternalError, err)
		assert.Equal(t, int64(0), result)
	})

	t.Run("存在しないカードを削除しようとした場合", func(t *testing.T) {
		client := NewTestClient(t, db)

		var result int64
		err := client.Tx(context.Background(), func(ctx context.Context) error {
			r, err := client.DeleteCard(ctx, -1)
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
