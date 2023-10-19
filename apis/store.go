package apis

import (
	"database/sql"

	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
)

type Store struct {
	*db.Queries
	db *sql.DB
}

func NewStore(dbs *sql.DB) *Store {
	return &Store{
		db:      dbs,
		Queries: db.New(dbs),
	}
}

func (s *Store) DB_instatnce() *sql.DB {
	return s.db
}
