// Package persistence contains database logic.
package persistence

import "celtra-programming-assigment/cmd/tracker/dto"

// DB is an active database connection
var DB Database

// Database interface represents the connection to the database
// and defines methods that can be implemented by various database providers.
//
// It can also be used to create a mocked implementation for testing purposes.
type Database interface {
	// IsActiveAccount check if a given account ID is active or not
	IsActiveAccount(ID int) (bool, error)
	// CreateAccount creates a new account.
	//
	// - name     - required
	//
	// - isActive - optional (default: false)
	CreateAccount(name string, isActive bool) (*dto.Account, error)
	// GetAccount returns an account record matching the ID.
	GetAccount(ID int) (*dto.Account, error)
}
