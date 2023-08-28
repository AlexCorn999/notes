package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore() *Store {
	return &Store{
		db: nil,
	}
}

// Open opens a database connection
func (s *Store) Open() error {
	db, err := sql.Open("postgres", "host=localhost user=kode dbname=kode sslmode=disable password=5427")
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db
	return nil
}

// Close closes a database connection
func (s *Store) Close() {
	s.db.Close()
}
