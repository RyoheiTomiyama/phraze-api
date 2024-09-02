package fixture

import (
	"fmt"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/samber/lo"
)

type RoleInput struct {
	Key           *string
	Name          *string
	PermissionIDs []int64
}

func (f *fixture) CreateRole(t *testing.T, inputs ...RoleInput) []*model.Role {
	var list []struct {
		role          *model.Role
		permissionIDs []int64
	}
	offset := len(f.Roles)
	for i, input := range inputs {
		list = append(list, struct {
			role          *model.Role
			permissionIDs []int64
		}{
			role: &model.Role{
				Key:  lo.FromPtrOr(input.Key, fmt.Sprintf("key-%d", i+offset)),
				Name: lo.FromPtrOr(input.Name, fmt.Sprintf("name-%d", i+offset)),
			},
			permissionIDs: input.PermissionIDs,
		})
	}

	query := `
		INSERT INTO roles (key, name) 
		VALUES (:key, :name)
		RETURNING *
	`
	queryRelation := `
		INSERT INTO roles_permissions (role_id, permission_id) 
		VALUES (:role_id, :permission_id)
		RETURNING *
	`

	tx := f.db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		t.Fatal(err)

		return nil
	}
	stmtRelation, err := tx.PrepareNamed(queryRelation)
	if err != nil {
		t.Fatal(err)

		return nil
	}

	var insertedModel []*model.Role

	for _, l := range list {
		var result model.Role

		if err = stmt.QueryRowx(l.role).StructScan(&result); err != nil {
			fmt.Print(fmt.Errorf("%w", err))
			if inerr := tx.Rollback(); inerr != nil {
				t.Fatal(inerr)

				return nil
			}

			t.Fatal(err)

			return nil
		}

		insertedModel = append(insertedModel, &result)

		for _, pid := range l.permissionIDs {
			argsRelation := &model.RolesPermission{
				RoleID:       result.ID,
				PermissionID: pid,
			}

			var result model.RolesPermission

			if err = stmtRelation.QueryRowx(argsRelation).StructScan(&result); err != nil {
				fmt.Print(fmt.Errorf("%w", err))
				if inerr := tx.Rollback(); inerr != nil {
					t.Fatal(inerr)

					return nil
				}
				t.Fatal(err)

				return nil
			}
		}

	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)

		return nil
	}

	f.Roles = append(f.Roles, insertedModel...)

	return insertedModel
}
