package db

import (
	"context"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/stretchr/testify/assert"
)

func TestCreateDeck(t *testing.T) {
	ctx := context.Background()

	db := db_test.GetDB(t)
	defer db.Close()

	client := NewTestClient(t, db)

	t.Run("正常系", func(t *testing.T) {
		newDeck := &domain.Deck{
			Name:   "test",
			UserID: "user-1",
		}
		deck, err := client.CreateDeck(ctx, newDeck)

		assert.NoError(t, err)
		assert.Equal(t, "test", deck.Name)
	})
}
