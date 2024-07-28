package repository

import (
	"fmt"

	"github.com/ctheil/pmdb-api/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRespoitory struct {
	TX *sqlx.Tx
}

func NewUserRepository(tx *sqlx.Tx) *UserRespoitory {
	return &UserRespoitory{TX: tx}
}

func (u *UserRespoitory) Insert(post model.PostUser) (id int64, ok bool) {
	_, err := u.TX.NamedExec("INSERT INTO users (username,hashed_pw,refresh_token_version) VALUES (:username,:hashed_pw,:refresh_token_version)", &post)
	if err != nil {
		fmt.Printf("error inserting user: %e", err)
		return 0, false
	}
	//  u.DB.MustExec(`INSERT INTO users (username, hashed_pw, refresh_token_version) VALUES ($1, $2, $3)`, )
	// stmt, err := u.DB.Prepare(`INSERT INTO users (username, hashed_pw, refresh_token_version) VALUES ($1, $2, $3)`)
	// if err != nil {
	// 	log.Println(err)
	// 	return 0, false
	// }
	// defer stmt.Close()
	// result, err := stmt.Exec(post.Username, post.HashedPW, post.RefreshTokenVersion)
	// if err != nil {
	// 	log.Println(err)
	// 	return 0, false
	// }
	// id, _ = result.LastInsertId()
	// return id, true
	return 0, true
}

func (u *UserRespoitory) GetById(id string) (model.User, error) {
	user := model.User{}
	if err := u.TX.Get(&user, "SELECT * FROM users WHERE id=$1", id); err != nil {
		return user, err
	}
	return user, nil
}

func (u *UserRespoitory) GetByUsername(username string) (model.User, error) {
	user := model.User{}

	if err := u.TX.Get(&user, "SELECT * FROM users WHERE username=$1", username); err != nil {
		fmt.Printf("[GetByUsername]: error scanning: %e\n", err)
		return user, err
	}
	fmt.Printf("Found user: %v", user)
	return user, nil
}
