package auth

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
)

type User domain.User

type IUSerUtil interface {
	WithCtx(ctx context.Context) context.Context
}

func New(u *domain.User) IUSerUtil {
	return &User{
		ID:     u.ID,
		Name:   u.Name,
		Avatar: u.Avatar,
	}
}

type userCtxKey struct{}

// 引数で与えられたロガーを context に詰め、新たな context を返す
func (u *User) WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, userCtxKey{}, u)
}

// context からロガーを取り出す。取り出せない場合はデフォルトのロガーを返す
func FromCtx(ctx context.Context) *User {
	u, ok := ctx.Value(userCtxKey{}).(*User)
	if !ok {
		return nil
	}

	return u
}
