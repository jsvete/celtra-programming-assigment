// Package dto contains implementations of data transfer objects.
package dto

// Account DTO used to represent a single account with account ID, name and if it's active or not.
type Account struct {
	ID       int
	Name     string
	isActive bool
}
