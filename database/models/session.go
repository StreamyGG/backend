package models

import (
	"backend/database"
	"time"

	"github.com/gocql/gocql"
)

type Session struct {
	Token     string     `json:"token"`
	UserID    gocql.UUID `json:"user_id"`
	IP        string     `json:"ip"`
	Device    string     `json:"device"`
	Trusted   bool       `json:"trusted"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt time.Time  `json:"expires_at"`
}

func CreateSession(db *database.DB, userID gocql.UUID, ip string, device string) (*Session, error) {
	session := &Session{
		Token:     generateToken(),
		UserID:    userID,
		IP:        ip,
		Device:    device,
		Trusted:   false,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(48 * time.Hour), // Set expiration to 48 hours from now
	}

	err := db.Session.Query(`INSERT INTO sessions ("token", user_id, ip, device, trusted, created_at, expires_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		session.Token, session.UserID, session.IP, session.Device, session.Trusted, session.CreatedAt, session.ExpiresAt).Exec()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func generateToken() string {
	return gocql.TimeUUID().String()
}

func GetSessionByToken(db *database.DB, token string) (*Session, error) {
	var session Session
	err := db.Session.Query(`SELECT "token", user_id, ip, device, trusted, created_at, expires_at FROM sessions WHERE "token" = ?`, token).
		Scan(&session.Token, &session.UserID, &session.IP, &session.Device, &session.Trusted, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func GetSessionsByUserID(db *database.DB, userID gocql.UUID) ([]*Session, error) {
	var sessions []*Session
	iter := db.Session.Query(`SELECT "token", user_id, ip, device, trusted, created_at, expires_at FROM sessions WHERE user_id = ?`, userID).Iter()
	var session Session
	for iter.Scan(&session.Token, &session.UserID, &session.IP, &session.Device, &session.Trusted, &session.CreatedAt, &session.ExpiresAt) {
		sessions = append(sessions, &session)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return sessions, nil
}

func DeleteSessionsByUserID(db *database.DB, userID gocql.UUID) error {
	return db.Session.Query(`DELETE FROM sessions WHERE user_id = ?`, userID).Exec()
}

func DeleteExpiredSessions(db *database.DB) error {
	return db.Session.Query(`DELETE FROM sessions WHERE expires_at < ?`, time.Now()).Exec()
}

func TrustDevice(db *database.DB, token string) error {
	return db.Session.Query(`UPDATE sessions SET trusted = true WHERE "token" = ?`, token).Exec()
}

func GetTrustedDevices(db *database.DB, userID gocql.UUID) ([]*Session, error) {
	var sessions []*Session
	iter := db.Session.Query(`SELECT "token", user_id, ip, device, trusted, created_at, expires_at FROM sessions WHERE user_id = ? AND trusted = true ALLOW FILTERING`, userID).Iter()
	var session Session
	for iter.Scan(&session.Token, &session.UserID, &session.IP, &session.Device, &session.Trusted, &session.CreatedAt, &session.ExpiresAt) {
		sessions = append(sessions, &session)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return sessions, nil
}

func LogoutSession(db *database.DB, token string) error {
	return db.Session.Query(`DELETE FROM sessions WHERE "token" = ?`, token).Exec()
}
