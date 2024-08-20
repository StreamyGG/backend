package database

func CreateTables(db *DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY,
            name TEXT,
            email TEXT
        )`,
	}

	for _, query := range queries {
		if err := db.Session.Query(query).Exec(); err != nil {
			return err
		}
	}

	return nil
}
