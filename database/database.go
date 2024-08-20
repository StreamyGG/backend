package database

import (
	"time"

	"github.com/gocql/gocql"
)

type DB struct {
	Session *gocql.Session
}

func InitDB() (*DB, error) {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "streamygg"
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = time.Second * 5

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	db := &DB{Session: session}

	if err := CreateTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Close() {
	db.Session.Close()
}
