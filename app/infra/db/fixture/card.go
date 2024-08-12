package fixture

import (
	"fmt"
	"testing"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/samber/lo"
)

type CardInput struct {
	Question  *string
	Answer    *string
	AIAnswer  *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (f *fixture) CreateCard(t *testing.T, deckID int64, cards ...CardInput) []*model.Card {
	var list []*model.Card
	offset := len(f.Cards)
	for i, d := range cards {
		list = append(list, &model.Card{
			DeckID:    deckID,
			Question:  lo.FromPtrOr(d.Question, fmt.Sprintf("question-%d", i+offset)),
			Answer:    lo.FromPtrOr(d.Answer, fmt.Sprintf("answer-%d", i+offset)),
			AIAnswer:  lo.FromPtrOr(d.AIAnswer, fmt.Sprintf("**【日本語訳】** answer-%d", i+offset)),
			CreatedAt: lo.FromPtrOr(d.CreatedAt, time.Now()),
			UpdatedAt: lo.FromPtrOr(d.UpdatedAt, time.Now()),
		})
		// 作成日の作成順を担保するためスリープする
		time.Sleep(time.Millisecond)
	}

	query := `
		INSERT INTO cards (deck_id, question, answer, ai_answer, created_at, updated_at) 
		VALUES (:deck_id, :question, :answer, :ai_answer, :created_at, :updated_at)
		RETURNING *
	`

	tx := f.db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		t.Fatal(err)

		return nil
	}

	var insertedCards []*model.Card

	for _, l := range list {
		var result model.Card

		if err = stmt.QueryRowx(l).StructScan(&result); err != nil {
			fmt.Print(fmt.Errorf("%w", err))
			if inerr := tx.Rollback(); inerr != nil {
				t.Fatal(inerr)

				return nil
			}

			t.Fatal(err)

			return nil
		}

		insertedCards = append(insertedCards, &result)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)

		return nil
	}

	f.Cards = append(f.Cards, insertedCards...)

	return insertedCards
}
