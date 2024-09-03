package models

import (
	"backend/database"

	"github.com/gocql/gocql"
)

type User struct {
	ID       gocql.UUID `json:"id"`
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Password string     `json:"password,omitempty"`
}

func (u *User) Save(db *database.DB) error {
	return db.Session.Query(`INSERT INTO users (id, name, email, password) 
		VALUES (?, ?, ?, ?) IF NOT EXISTS`,
		u.ID, u.Name, u.Email, u.Password).Exec()
}

func GetUserByID(db *database.DB, id gocql.UUID) (*User, error) {
	var user User
	if err := db.Session.Query(`SELECT id, name, email, password FROM users WHERE id = ?`, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) CreateSession(db *database.DB, ip string, device string) (*Session, error) {
	return CreateSession(db, u.ID, ip, device)
}

func (u *User) GetSessions(db *database.DB) ([]*Session, error) {
	return GetSessionsByUserID(db, u.ID)
}

func (u *User) DeleteAllSessions(db *database.DB) error {
	return DeleteSessionsByUserID(db, u.ID)
}

func (u *User) LogoutSession(db *database.DB, token string) error {
	return LogoutSession(db, token)
}
