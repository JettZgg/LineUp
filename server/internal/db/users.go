// File: internal/db/users.go
package db

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateUser(user *User) error {
	_, err := DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)",
		user.Username, user.Password)
	return err
}

func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := DB.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
