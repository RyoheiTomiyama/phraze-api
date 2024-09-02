package fixture

import (
	"fmt"
	"testing"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/go-faker/faker/v4"
	"github.com/samber/lo"
)

type UserInput struct {
	ID        *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (f *fixture) CreateUser(t *testing.T, inputs ...UserInput) []*model.User {
	t.Helper()

	var list []*model.User
	for _, input := range inputs {
		t.Log("input")
		t.Log(input)
		list = append(list, &model.User{
			ID:        lo.FromPtrOr(input.ID, faker.UUIDDigit()),
			CreatedAt: lo.FromPtrOr(input.CreatedAt, time.Now()),
			UpdatedAt: lo.FromPtrOr(input.UpdatedAt, time.Now()),
		})
	}

	query := `
		INSERT INTO users (id, created_at, updated_at) 
		VALUES (:id, :created_at, :updated_at)
		RETURNING *
	`

	tx := f.db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		t.Fatal(err)

		return nil
	}

	var insertedModel []*model.User

	for _, l := range list {
		var result model.User

		if err = stmt.QueryRowx(l).StructScan(&result); err != nil {
			fmt.Print(fmt.Errorf("%w", err))
			if inerr := tx.Rollback(); inerr != nil {
				t.Fatal(inerr)

				return nil
			}

			t.Fatal(err)

			return nil
		}

		insertedModel = append(insertedModel, &result)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)

		return nil
	}

	f.Users = append(f.Users, insertedModel...)

	return insertedModel
}
