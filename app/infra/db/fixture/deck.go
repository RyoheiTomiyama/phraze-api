package fixture

import (
	"fmt"
	"testing"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/samber/lo"
)

type DeckInput struct {
	UserID    string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (f *fixture) CreateDeck(t *testing.T, decks ...*DeckInput) []*model.Deck {
	var list []*model.Deck
	offset := len(f.Decks)
	for i, d := range decks {
		list = append(list, &model.Deck{
			UserID:    lo.Ternary(d.UserID == "", "", d.UserID),
			Name:      lo.Ternary(d.Name == "", fmt.Sprintf("deck-%d", i+offset), d.Name),
			CreatedAt: lo.Ternary(d.CreatedAt.IsZero(), time.Now(), d.CreatedAt),
			UpdatedAt: lo.Ternary(d.CreatedAt.IsZero(), time.Now(), d.UpdatedAt),
		})
	}

	query := `
		INSERT INTO decks (user_id, name, created_at, updated_at) 
		VALUES (:user_id, :name, :created_at, :updated_at)
		RETURNING *
	`

	tx := f.db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		t.Fatal(err)

		return nil
	}

	var insertedDeck []*model.Deck

	for _, l := range list {
		var result model.Deck

		if err = stmt.QueryRowx(l).StructScan(&result); err != nil {
			fmt.Print(fmt.Errorf("%w", err))
			if inerr := tx.Rollback(); inerr != nil {
				t.Fatal(inerr)

				return nil
			}

			t.Fatal(err)

			return nil
		}

		insertedDeck = append(insertedDeck, &result)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)

		return nil
	}

	f.Decks = append(f.Decks, insertedDeck...)

	return insertedDeck
}
