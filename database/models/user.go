package models

import (
	"backend/database"

	"github.com/gocql/gocql"
)

type User struct {
	ID    gocql.UUID `json:"id"`
	Name  string     `json:"name"`
	Email string     `json:"email"`
}

func (u *User) Save(db *database.DB) error {
	return db.Session.Query(`INSERT INTO users (id, name, email) VALUES (?, ?, ?)`,
		u.ID, u.Name, u.Email).Exec()
}

func GetUserByID(db *database.DB, id gocql.UUID) (*User, error) {
	var user User
	if err := db.Session.Query(`SELECT id, name, email FROM users WHERE id = ?`, id).Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}
