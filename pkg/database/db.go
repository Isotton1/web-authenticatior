package database

import (
	"database/sql"
	"errors"

	"github.com/Isotton1/web-authenticatior/internal/common"
	"github.com/Isotton1/web-authenticatior/models"
	_ "modernc.org/sqlite"
)

var Db *sql.DB

func InitDB(url string) error {
	var err error
	Db, err = sql.Open("sqlite", url)
	if err != nil {
		return err
	}
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		salt TEXT NOT NULL,
		pepper TEXT NOT NULL
	);`
	_, err = Db.Exec(usersTable)
	if err != nil {
		return err
	}

	return nil
}

var ErrUserExists = errors.New("User already exists")

func InsertUser(user *models.User) error {
	exist, err := HasUser(user.Username)
	if err != nil {
		return err
	}
	if exist {
		return ErrUserExists
	}
	query := `INSERT INTO users(username, password_hash, salt, pepper) VALUES(?, ?, ?, ?)`
	_, err = Db.Exec(query, user.Username, user.PasswordHash, user.Salt, user.Pepper)
	return err
}

func GetUser(username string) (models.User, error) {
	var passwordHash, salt, pepper string
	err := Db.QueryRow("SELECT password_hash, salt, pepper FROM users WHERE username = ?", username).Scan(&passwordHash, &salt, &pepper)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, common.ErrNoUserFound
		}
		return models.User{}, err
	}
	user := models.User{
		Username:     username,
		PasswordHash: passwordHash,
		Salt:         salt,
		Pepper:       pepper,
	}
	return user, nil
}

func HasUser(username string) (bool, error) {
	var exist bool
	err := Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exist)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exist, nil
}
