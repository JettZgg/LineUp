// File: internal/db/users.go

package db

import (
	"time"
)

type User struct {
	UID       int64     `json:"uid"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateUser(user *User) error {
	_, err := DB.Exec("INSERT INTO users (uid, username, password) VALUES ($1, $2, $3)",
		user.UID, user.Username, user.Password)
	return err
}

func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := DB.QueryRow("SELECT uid, username, password, created_at, updated_at FROM users WHERE username = $1", username).
		Scan(&user.UID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
