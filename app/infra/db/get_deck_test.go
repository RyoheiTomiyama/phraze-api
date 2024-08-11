package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/stretchr/testify/assert"
)

func TestGetDeck(t *testing.T) {
	db := db_test.GetDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fx := fixture.New(db)
	decks := fx.CreateDeck(t, &fixture.DeckInput{Name: "test"})

	deck := decks[0]
	if deck == nil {
		t.Fatalf("not found created deck")
	}

	t.Run("test", func(t *testing.T) {
		client := NewTestClient(t, db)

		fmt.Printf("%d", deck.ID)

		deck, err := client.GetDeck(context.Background(), deck.ID)
		fmt.Printf("%v", err)
		assert.NoError(t, err)
		assert.Equal(t, "test", deck.Name)
	})
}
