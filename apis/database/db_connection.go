package database

import "database/sql"

var (
	DB_DRIVER = "postgres"
	DB_SOURCE = "postgresql://postgres:support12@skroman-user.ckwveljlsuux.ap-south-1.rds.amazonaws.com:5432/skroman_client_complaints"
)

func make_connection() (*sql.DB, error) {
	db, err := sql.Open(DB_DRIVER, DB_SOURCE)

	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

var DB_INSTANCE = make_connection

func CloseDBConnection(db *sql.DB) error {
	return db.Close()
}
