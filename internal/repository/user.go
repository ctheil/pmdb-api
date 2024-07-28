package repository

import (
	"database/sql"
	"log"

	"github.com/ctheil/pmdb-api/internal/model"
)

type UserRespoitory struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepositoryInterface {
	return &UserRespoitory{DB: db}
}

func (u *UserRespoitory) InsertUser(post model.PostUser) bool {
	stmt, err := u.DB.Prepare(`INSERT INTO users (username, hashed_pw, refresh_token_version) VALUES ($1, $2, $3)`)
	if err != nil {
		log.Println(err)
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(post.Username, post.HashedPW, post.RefreshTokenVersion)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (u *UserRespoitory) SelectUserByID(id string) (model.User, error) {
	user := model.User{}
	row := u.DB.QueryRow("SELECT * FROM users WHERE id = ($1)", id)
	if err := row.Scan(&user); err != nil {
		return user, err
	}
	return user, nil
}

func (u *UserRespoitory) SelectUserByUsername(username string) (model.User, error) {
	user := model.User{}
	row := u.DB.QueryRow("SELECT * FROM users WHERE username = ($1)", username)
	if err := row.Scan(&user); err != nil {
		return user, err
	}
	return user, nil
}
