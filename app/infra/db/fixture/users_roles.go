package fixture

import (
	"fmt"
	"testing"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/samber/lo"
)

type UsersRoleInput struct {
	UserID    string
	RoleID    int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (f *fixture) CreateUsersRole(t *testing.T, inputs ...UsersRoleInput) []*model.UsersRole {
	t.Helper()

	var list []*model.UsersRole
	for _, input := range inputs {
		t.Log("input")
		t.Log(input)
		list = append(list, &model.UsersRole{
			UserID:    input.UserID,
			RoleID:    input.RoleID,
			CreatedAt: lo.FromPtrOr(input.CreatedAt, time.Now()),
			UpdatedAt: lo.FromPtrOr(input.UpdatedAt, time.Now()),
		})
	}

	query := `
		INSERT INTO users_roles (user_id, role_id, created_at, updated_at) 
		VALUES (:user_id, :role_id, :created_at, :updated_at)
		RETURNING *
	`

	tx := f.db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		t.Fatal(err)

		return nil
	}

	var insertedModel []*model.UsersRole

	for _, l := range list {
		var result model.UsersRole

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

	f.UsersRoles = append(f.UsersRoles, insertedModel...)

	return insertedModel
}
