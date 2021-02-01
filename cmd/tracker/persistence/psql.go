// Package persistence contains database logic.
package persistence

import (
	"celtra-programming-assigment/cmd/tracker/dto"
	"database/sql"
	"fmt"
	"os"
)

// variables defining information to successfully connect to the database
var (
	dbName    string // DB_NAME
	dbUser    string // DB_USER
	dbPwd     string // DB_PWD
	dbAddr    string // DB_ADDR
	dbPort    string // DB_PORT
	redisAddr string // REDIS_ADDR
)

func init() {
	dbName = os.Getenv("DB_NAME")
	dbUser = os.Getenv("DB_USER")
	dbPwd = os.Getenv("DB_PWD")
	dbAddr = os.Getenv("DB_ADDR")
	dbPort = os.Getenv("DB_PORT")
	redisAddr = os.Getenv("REDIS_ADDR")
}

// Postgres implements Database interface and represents a connection to the PostgreSQL database.
type Postgres struct {
	db *sql.DB
}

// IsActiveAccount check if a given account ID is active or not
func (pg *Postgres) IsActiveAccount(ID int) (bool, error) {

	return true, nil
}

// CreateAccount creates a new account.
//
// - ID       - required
//
// - name     - required
//
// - isActive - optional (default: false)
func (pg *Postgres) CreateAccount(ID int, name string, isActive bool) (*dto.Account, error) {

	return nil, nil
}

// DeactivateAccount flags an account as inactive.
func (pg *Postgres) DeactivateAccount(ID int) error {

	return nil
}

// GetAccount returns an account record matching the ID.
func (pg *Postgres) GetAccount(ID int) (*dto.Account, error) {

	return nil, nil
}

// NewPostgres creates a new instance of Postgres.
func NewPostgres() (*Postgres, error) {
	pg := &Postgres{}
	dataSource := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s", dbUser, dbName, dbPwd, dbAddr, dbPort)

	db, err := sql.Open("posgres", dataSource)
	if err != nil {
		return nil, err
	}

	pg.db = db

	// attempt to create table and random data, if it doesnt exist already
	_, err = pg.db.Exec(`
		CREATE TABLE IF NOT EXISTS account (
			id INTEGER PRIMARY KEY,
			name VARCHAR (255) NOT NULL,
			isActive BOOLEAN DEFAULT TRUE
		)

		INSERT INTO account (id, name) SELECT *, md5(RANDOM()::TEXT) FROM generate_series(1, 1000);
	`)
	if err != nil {
		return nil, err
	}

	// // create Redis client for caching layer
	// pg.redis = redis.NewClient(&redis.Options{
	// 	Addr: redisAddr,
	// 	DB:   0,
	// })

	// status := pg.redis.Ping(context.Background())
	// if status.Err() != nil {
	// 	return nil, status.Err()
	// }

	return pg, err
}
