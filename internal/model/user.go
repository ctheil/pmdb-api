package model

type User struct {
	Id                  uint   `json:"id" db:"id"`
	Username            string `json:"username" db:"username"`
	Password            string `json:"password" db:"hashed_pw"`
	Email               string `json:"email" db:"email"`
	RefreshTokenVersion int    `json:"refresh_token_version" db:"refresh_token_version"`
}

type PostUser struct {
	Username            string `json:"username" db:"username"`
	Password            string `json:"password" db:"hashed_pw"`
	Email               string `json:"email" db:"email"`
	RefreshTokenVersion int    `json:"refresh_token_version" db:"refresh_token_version"`
}

type OAuthUser struct {
	Id                  uint   `json:"id" db:"id"`
	Name                string `json:"name" db:"name"`
	Email               string `json:"email" db:"email"`
	Picture             string `json:"picture" db:"picture"`
	RefreshTokenVersion uint   `json:"refresh_token_version" db:"refresh_token_version"`
}

// type PostOAuthUser struct {
// 	Name                string `json:"name" db:"name"`
// 	Email               string `json:"email" db:"email"`
// 	Picture             string `json:"picture" db:"picture"`
// 	RefreshTokenVersion int    `json:"refresh_token_version" db:"refresh_token_version"`
// }
