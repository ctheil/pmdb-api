package repository

import "github.com/ctheil/pmdb-api/internal/model"

type UserRepositoryInterface interface {
	GetById(id string) (model.User, error)
	Insert(post model.PostUser) bool
}

type OAuthUserRepositroyInterface interface {
	GetByEmail(email string) (model.OAuthUser, error)
	Insert(post model.OAuthUser) bool
	UpdateRefreshToeknVersion(new_version uint) (ok bool)
}
