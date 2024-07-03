package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDeck(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		client, err := NewClient(testDataSourceOption())
		if err != nil {
			t.Fatalf("%+v", err)
		}

		deck, err := client.GetDeck(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, "sample", deck.Name)
	})
}
