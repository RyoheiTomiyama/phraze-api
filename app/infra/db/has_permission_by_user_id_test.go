package db

import (
	"context"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestHasPermissionByUserID(t *testing.T) {
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
		fixture.RoleInput{PermissionIDs: []int64{permissions[1].ID}},
	)
	// user0 has role0, permission0
	users := fx.CreateUser(t, fixture.UserInput{})
	fx.CreateUsersRole(t, fixture.UsersRoleInput{
		UserID: users[0].ID,
		RoleID: roles[0].ID,
	})

	t.Run("所有するロールの場合", func(t *testing.T) {
		client := NewTestClient(t, db)

		result, err := client.HasPermissionByUserID(context.Background(),
			users[0].ID, domain.PermissionUnlimitedCardCreation,
		)
		assert.NoError(t, err)
		assert.True(t, result)
	})
	t.Run("所有しないロールの場合", func(t *testing.T) {
		client := NewTestClient(t, db)

		result, err := client.HasPermissionByUserID(context.Background(),
			users[0].ID, domain.PermissionUnlimitedAIAnswerGeneration,
		)
		assert.NoError(t, err)
		assert.False(t, result)
	})

}
