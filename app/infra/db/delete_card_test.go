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

	cards := fx.CreateCard(t, decks[0].ID, make([]fixture.CardInput, 10)...)

	t.Run("正常系", func(t *testing.T) {
		client := NewTestClient(t, db)

		var result int64
		err := client.Tx(context.Background(), func(ctx context.Context) error {
			r, err := client.DeleteCard(ctx, cards[0].ID)
			if err != nil {
				return err
			}

			result = r

			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), result)

		var c model.Card
		err = db.Get(&c, "SELECT * FROM cards WHERE id=$1", cards[0].ID)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("トランザクションの外でメソッドを実行した場合", func(t *testing.T) {
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
