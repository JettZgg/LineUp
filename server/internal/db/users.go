// File: internal/db/users.go

package db

import (
	"time"
	"database/sql"
	"errors"
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

func GetUserByID(userID int64) (*User, error) {
	var user User
	err := DB.QueryRow("SELECT uid, username FROM users WHERE uid = $1", userID).Scan(&user.UID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}