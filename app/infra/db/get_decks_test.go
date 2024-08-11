package db

import (
	"context"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/stretchr/testify/assert"
)

func TestGetDecks(t *testing.T) {
	db := db_test.GetDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fx := fixture.New(db)
	decks := fx.CreateDeck(
		t,
		&fixture.DeckInput{UserID: "own"},
		&fixture.DeckInput{UserID: "own"},
		&fixture.DeckInput{UserID: "another user"},
	)

	deck := decks[0]
	if deck == nil {
		t.Fatalf("not found created deck")
	}

	t.Run("test", func(t *testing.T) {
		client := NewTestClient(t, db)

		decks, err := client.GetDecks(context.Background(), "own")
		assert.NoError(t, err)
		assert.Len(t, decks, 2)
	})
}
