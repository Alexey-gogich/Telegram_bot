package storage

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type Storage struct {
	db        *sql.DB
	reminders *Reminders
}

func NewStorage() (*Storage, error) {
	db_connection, _ := sql.Open("postgres", "postgres://postgres:postgres@0.0.0.0:8888/"+os.Getenv("DATABASE_NAME")+"?sslmode=disable")
	stor := Storage{
		db:        db_connection,
		reminders: nil,
	}
	if err := stor.db.Ping(); err != nil {
		return nil, err
	}
	return &stor, nil
}

func (st *Storage) Reminder() *Reminders {
	if st.reminders == nil {
		st.reminders = &Reminders{
			storage: st,
		}
	}
	return st.reminders
}
