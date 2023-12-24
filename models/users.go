package models

type User struct {
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	Salt         string `db:"salt"`
	Pepper       string `db:"Pepper"`
}
