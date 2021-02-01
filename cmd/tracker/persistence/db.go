// Package persistence contains database logic.
package persistence

import "celtra-programming-assigment/cmd/tracker/dto"

// Database interface represents the connection to the database
// and defines methods that can be implemented by various datasources.
//
// It can also be used to create a mocked implementation for testing purposes.
type Database interface {
	// IsActiveAccount check if a given account ID is active or not
	IsActiveAccount(ID int) (bool, error)
	// CreateAccount creates a new account.
	//
	// - ID       - required
	//
	// - name     - required
	//
	// - isActive - optional (default: false)
	CreateAccount(ID int, name string, isActive bool) (*dto.Account, error)
	// DeactivateAccount flags an account as inactive.
	DeactivateAccount(ID int) error
	// GetAccount returns an account record matching the ID.
	GetAccount(ID int) (*dto.Account, error)
}
