package repository

import "github.com/ctheil/pmdb-api/internal/model"

type UserRepositoryInterface interface {
	SelectUserByID(id string) (model.User, error)
	InsertUser(post model.PostUser) bool
}
