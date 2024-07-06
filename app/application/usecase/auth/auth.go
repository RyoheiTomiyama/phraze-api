package auth

import "github.com/RyoheiTomiyama/phraze-api/domain"

type IAuthUsecase interface {
	ParseToken(idToken string) (*domain.User, error)
}

type usecase struct{}

func New() IAuthUsecase {
	return &usecase{}
}

func (u *usecase) ParseToken(idToken string) (*domain.User, error) {
	return nil, nil
}
