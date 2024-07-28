package repository

import "github.com/ctheil/pmdb-api/internal/model"

type UserRepositoryInterface interface {
	GetById(id string) (model.User, error)
	Insert(post model.PostUser) bool
}
