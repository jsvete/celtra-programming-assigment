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
	dbName string // DB_NAME
	dbUser string // DB_USER
	dbPwd  string // DB_PWD
	dbAddr string // DB_ADDR
	dbPort string // DB_PORT
)

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
// - name     - required
//
// - isActive - optional (default: false)
func (pg *Postgres) CreateAccount(name string, isActive bool) (*dto.Account, error) {

	return nil, nil
}

// GetAccount returns an account record matching the ID.
func (pg *Postgres) GetAccount(ID int) (*dto.Account, error) {

	return nil, nil
}

// NewPostgres creates a new instance of Postgres.
func NewPostgres() error {
	pg := &Postgres{}
	dataSource := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s", dbUser, dbName, dbPwd, dbAddr, dbPort)

	db, err := sql.Open("posgres", dataSource)
	if err != nil {
		return err
	}

	pg.db = db

	rows, err := pg.db.Query("SELECT * FROM account LIMIT 1")
	if err != nil {
		// database isn't set up yet, lets create an account table and populate it with data
		_, err = pg.db.Exec(`
		CREATE TABLE IF NOT EXISTS account (
			id       SERIAL        PRIMARY KEY,
			name     VARCHAR (255) NOT NULL,
			isActive BOOLEAN       DEFAULT TRUE
		)

		INSERT INTO account (name) SELECT md5(RANDOM()::TEXT) FROM generate_series(1, 1000);
		`)
		if err != nil {
			return err
		}
	} else {
		rows.Close()
	}

	DB = pg

	return err
}

func init() {
	dbName = os.Getenv("DB_NAME")
	dbUser = os.Getenv("DB_USER")
	dbPwd = os.Getenv("DB_PWD")
	dbAddr = os.Getenv("DB_ADDR")
	dbPort = os.Getenv("DB_PORT")
}
