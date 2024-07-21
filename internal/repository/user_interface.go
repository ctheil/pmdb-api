package repository

import "github.com/ctheil/pmdb-api/internal/model"

type UserRepositoryInterface interface {
	SelectUser() []model.User
	InsertUser(post model.PostUser) bool
}
