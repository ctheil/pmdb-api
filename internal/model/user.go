package model

import "github.com/jmoiron/sqlx"

type User struct {
	Id                  uint   `json:"id" db:"id"`
	Username            string `json:"username" db:"username"`
	Password            string `json:"password" db:"hashed_pw"`
	RefreshTokenVersion int    `json:"refresh_token_version" db:"refresh_token_version"`
}

type PostUser struct {
	Username            string `json:"username" db:"username"`
	Password            string `json:"password" db:"hashed_pw"`
	RefreshTokenVersion int    `json:"refresh_token_version" db:"refresh_token_version"`
}

func (u *User) Save(tx sqlx.Tx) error {
	_, err := tx.Exec("UPDATE users SET refresh_token_version = $1 WHERE id = $2", u.RefreshTokenVersion, u.Id)

	return err
}
