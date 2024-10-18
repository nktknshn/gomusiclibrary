package database

import (
	"fmt"

	"zombiezen.com/go/sqlite"
)

type Database struct {
	conn *sqlite.Conn
}

func New(conn *sqlite.Conn) *Database {
	return &Database{conn}
}

func NewFromFile(path string) (*Database, error) {
	conn, err := sqlite.OpenConn(path, sqlite.OpenReadWrite)

	if err != nil {
		return nil, fmt.Errorf("database: create connection: %w", err)
	}

	return New(conn), nil
}
