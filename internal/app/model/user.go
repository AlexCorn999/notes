package model

import "errors"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) Validate() error {
	if len(u.Email) < 1 {
		return errors.New("no email")
	}
	if len(u.Password) < 1 {
		return errors.New("no password")
	}
	return nil
}
