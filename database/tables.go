package database

func CreateTables(db *DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY,
            name TEXT,
            email TEXT,
            password TEXT
        )`,
		`CREATE INDEX IF NOT EXISTS ON users (email)`,
		`CREATE TABLE IF NOT EXISTS sessions (
            "token" TEXT PRIMARY KEY,
            user_id UUID,
            ip TEXT,
            device TEXT,
            trusted BOOLEAN,
            created_at TIMESTAMP,
            expires_at TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS users_by_email (
            email TEXT PRIMARY KEY,
            user_id UUID
        )`,
	}

	for _, query := range queries {
		if err := db.Session.Query(query).Exec(); err != nil {
			return err
		}
	}

	return nil
}
