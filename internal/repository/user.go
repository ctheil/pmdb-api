package repository

import (
	"database/sql"
	"log"

	model "github.com/ctheil/pmdb-api/internal/model"
)

type UserRespoitory struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepositoryInterface {
	return &UserRespoitory{DB: db}
}

func (u *UserRespoitory) InsertUser(post model.PostUser) bool {
	stmt, err := u.DB.Prepare("INSERT INTO users (username, hashed_pw) VALUES ($1, $2)")
	if err != nil {
		log.Println(err)
		return false
	}
	defer stmt.Close()
	_, err2 := stmt.Exec(post.Username, post.HashedPW)
	if err2 != nil {
		log.Println(err2)
		return false
	}

	return true
}

func (u *UserRespoitory) SelectUser(n string) []model.User {
	var results []model.User

	rows, err := u.DB.Query("SELECT * FROM users")
	if err != nil {
		log.Println(err)
		return nil
	}

	for rows.Next() {
		var (
			id        uint
			username  string
			hashed_pw string
		)
		err := rows.Scan(&id, &username, &hashed_pw)
		if err != nil {
			log.Println(err)
		} else {
			user := model.User{Id: id, Username: username, HashedPW: hashed_pw}
			results = append(results, user)
		}
	}
	return results
}
