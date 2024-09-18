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
	decks := fx.CreateDeck(
		t,
		&fixture.DeckInput{UserID: lo.ToPtr("own")},
	)

	t.Run("正常系", func(t *testing.T) {
		client := NewTestClient(t, db)

		result, err := client.DeleteDeck(context.Background(), decks[0].ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), result)

		var c model.Deck
		err = db.Get(&c, "SELECT * FROM decks WHERE id=$1", decks[0].ID)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("存在しないデッキを削除しようとした場合", func(t *testing.T) {
		client := NewTestClient(t, db)

		result, err := client.DeleteDeck(context.Background(), -1)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), result)
	})
}
