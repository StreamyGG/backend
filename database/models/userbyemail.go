package models

import (
	"github.com/gocql/gocql"
)

type UserByEmail struct {
	Email  string     `cql:"email"`
	UserID gocql.UUID `cql:"user_id"`
}

func GetUserIDByEmail(session *gocql.Session, email string) (gocql.UUID, error) {
	var userID gocql.UUID
	err := session.Query("SELECT user_id FROM users_by_email WHERE email = ?", email).Scan(&userID)
	return userID, err
}
