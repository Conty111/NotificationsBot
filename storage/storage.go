package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgreDB struct {
	host     string
	user     string
	password string
	dbname   string
	port     int
	DB       *sql.DB
}

type NewSeria struct {
	ID        int
	AnimeName string
	Text      string
	Href      string
}

// Return the *PostgreDB with params
func New(host, user, password, dbname string, port int) (*PostgreDB, error) {
	p := PostgreDB{
		host:     host,
		user:     user,
		password: password,
		dbname:   dbname,
		port:     port,
	}
	err := p.Connect()
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Connect to the database with params
func (s *PostgreDB) Connect() error {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.host, s.port, s.user, s.password, s.dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return err
	}
	s.DB = db
	log.Printf("Connected to database %s", s.dbname)
	return nil
}

// Count rows in sql.Rows
func countRows(res *sql.Rows) int {
	var i int
	for res.Next() {
		i += 1
	}
	return i
}
