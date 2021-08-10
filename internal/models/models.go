// Package models contains data structures common through the program.
package models

type User struct {
	ID       uint64 `json:"id"`
	Account  string `json:"account"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Data struct {
	Users []User `json:"users"`
}
