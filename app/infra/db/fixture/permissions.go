package fixture

import (
	"fmt"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/samber/lo"
)

type PermissionInput struct {
	Key  *string
	Name *string
}

func (f *fixture) CreatePermission(t *testing.T, inputs ...*PermissionInput) []*model.Permission {
	var list []*model.Permission
	offset := len(f.Permissions)
	for i, input := range inputs {
		list = append(list, &model.Permission{
			Key:  lo.FromPtrOr(input.Key, fmt.Sprintf("key-%d", i+offset)),
			Name: lo.FromPtrOr(input.Name, fmt.Sprintf("name-%d", i+offset)),
		})
	}

	query := `
		INSERT INTO permissions (key, name) 
		VALUES (:key, :name)
		RETURNING *
	`

	tx := f.db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		t.Fatal(err)

		return nil
	}

	var insertedPermission []*model.Permission

	for _, l := range list {
		var result model.Permission

		if err = stmt.QueryRowx(l).StructScan(&result); err != nil {
			fmt.Print(fmt.Errorf("%w", err))
			if inerr := tx.Rollback(); inerr != nil {
				t.Fatal(inerr)

				return nil
			}

			t.Fatal(err)

			return nil
		}

		insertedPermission = append(insertedPermission, &result)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)

		return nil
	}

	f.Permissions = append(f.Permissions, insertedPermission...)

	return insertedPermission
}
