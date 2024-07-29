// File: internal/auth/auth.go
package auth

import (
	"errors"
	"log"
	"time"

	"github.com/JettZgg/LineUp/internal/db"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("467731") // Replace with a secure key from config

func RegisterUser(user *db.User) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Store the hashed password
	user.Password = string(hashedPassword)

	// Create the user in the database
	return db.CreateUser(user)
}

func LoginUser(username, password string) (string, error) {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		log.Printf("Error retrieving user %s: %v", username, err)
		return "", errors.New("invalid credentials")
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password mismatch for user %s: %v", username, err)
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating token for user %s: %v", username, err)
		return "", errors.New("could not generate token")
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	}

	return "", errors.New("invalid token")
}
