package db

import (
	"context"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGetPermissionsByUserID(t *testing.T) {
	db := db_test.GetDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fx := fixture.New(db)
	permissions := fx.CreatePermission(t,
		fixture.PermissionInput{Key: lo.ToPtr(domain.PermissionUnlimitedCardCreation.String())},
		fixture.PermissionInput{Key: lo.ToPtr(domain.PermissionUnlimitedAIAnswerGeneration.String())},
	)
	roles := fx.CreateRole(t,
		fixture.RoleInput{PermissionIDs: []int64{permissions[0].ID}},
		fixture.RoleInput{PermissionIDs: []int64{permissions[0].ID, permissions[1].ID}},
	)
	// user0 has role0, permission0
	users := fx.CreateUser(t, fixture.UserInput{})
	fx.CreateUsersRole(t,
		fixture.UsersRoleInput{UserID: users[0].ID, RoleID: roles[0].ID},
		fixture.UsersRoleInput{UserID: users[0].ID, RoleID: roles[1].ID},
	)

	t.Run("正常系", func(t *testing.T) {
		client := NewTestClient(t, db)

		result, err := client.GetPermissionsByUserID(context.Background(), users[0].ID)
		assert.NoError(t, err)
		assert.Len(t, result, 2)

		expectedKeys := lo.Map(permissions, func(item *model.Permission, _ int) string {
			return item.Key
		})
		actualKeys := lo.Map(result, func(item *domain.Permission, _ int) string {
			return item.Key
		})

		assert.True(t, lo.Every(expectedKeys, actualKeys), "expected: %v, actual: %v", expectedKeys, actualKeys)
	})
}
