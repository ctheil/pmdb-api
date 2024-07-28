package model

type User struct {
	Id                  uint   `json:"id"`
	Username            string `json:"username"`
	HashedPW            string `json:"hashed_pw"`
	RefreshTokenVersion int    `json:"refresh_token_version"`
}

type PostUser struct {
	Username            string `json:"username"`
	HashedPW            string `json:"hashed_pw"`
	RefreshTokenVersion int    `json:"refresh_token_version"`
}
