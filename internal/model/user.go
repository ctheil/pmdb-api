package model

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	HashedPW string `json:"hashed_pw"`
}

type PostUser struct {
	Username string `json:"username"`
	HashedPW string `json:"hashed_pw"`
}
