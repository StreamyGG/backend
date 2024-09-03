package database

import (
	"fmt"
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
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	db := &DB{Session: session}

	if err := db.Session.Query("SELECT release_version FROM system.local").Exec(); err != nil {
		return nil, fmt.Errorf("failed to execute test query: %v", err)
	}

	if err := CreateTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return db, nil
}

func (db *DB) Close() {
	db.Session.Close()
}
