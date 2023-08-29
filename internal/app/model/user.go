package model

import "errors"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates users
func (u *User) Validate() error {
	if len(u.Email) < 5 {
		return errors.New("incorrect email")
	}
	if len(u.Password) < 3 {
		return errors.New("incorrect password")
	}
	return nil
}
