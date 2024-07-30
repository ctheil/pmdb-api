package repository

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

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
	_, err := u.TX.NamedExec("INSERT INTO users (username,email,hashed_pw,refresh_token_version) VALUES (:username,:email,:hashed_pw,:refresh_token_version)", &post)
	if err != nil {
		fmt.Printf("error inserting user: %e", err)
		return 0, false
	}
	return 0, true
}

func (u *UserRespoitory) GetById(id string) (model.User, error) {
	user := model.User{}
	if err := u.TX.Get(&user, "SELECT * FROM users WHERE id=$1", id); err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		return user, err
	}
	return user, nil
}

func (u *UserRespoitory) GetByUsername(username string) (model.User, error) {
	user := model.User{}

	if err := u.TX.Get(&user, "SELECT * FROM users WHERE username=$1", username); err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		fmt.Printf("[GetByUsername]: error scanning: %e\n", err)
		return user, err
	}
	return user, nil
}

func (u *UserRespoitory) GetByEmail(email string) (model.User, error) {
	user := model.User{}

	if err := u.TX.Get(&user, "SELECT * FROM users WHERE email=$1", email); err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		fmt.Printf("[GetByEmail]: error scanning: %e\n", err)
		return user, err
	}
	return user, nil
}

func (u *UserRespoitory) Save(user model.User) error {
	v := reflect.ValueOf(user)
	typeOfUser := v.Type()

	var fields []string
	var values []interface{}
	i := 0
	for ; i < v.NumField(); i++ {
		// grab db tag on struct
		key := typeOfUser.Field(i).Tag.Get("db")
		val := v.Field(i).Interface()
		// skip id update
		// id is first value, so subsequent values will be 1..n == $1, $2, ..$n
		if key != "id" {
			// set value to $i
			fields = append(fields, fmt.Sprintf("%s = $%d", strings.ToLower(key), i))
			values = append(values, val)
		}
	}
	// last value should be userid
	values = append(values, user.Id)

	// i = last idx, use to set the final WHERE clause
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(fields, ", "), i)
	fmt.Printf("\n[userRepo.Save()]: Query: %s Values: %v\n", query, values)
	_, err := u.TX.Exec(query, values...)

	return err
}

func (u *UserRespoitory) UpdateField(user model.User, field string, value any) error {
	query := fmt.Sprintf("UPDATE users SET %s = $1 WHERE id = $2", field)
	_, err := u.TX.Exec(query, value, user.Id)
	return err
}
