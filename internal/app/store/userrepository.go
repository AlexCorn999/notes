package store

import "github.com/AlexCorn999/notes/internal/app/model"

type UserRepository struct {
	store *Store
}

// Create adds the user to the database and assign an id
func (r *UserRepository) Create(u *model.User) (*model.User, error) {
	if err := r.store.db.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.Password,
	).Scan(&u.ID); err != nil {
		return nil, err
	}
	return u, nil
}

// FindByEmail helps to find a user in the database
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, password FROM users WHERE email = $1",
		email,
	).Scan(&u.ID,
		&u.Email,
		&u.Password,
	); err != nil {
		return nil, err
	}
	return u, nil
}
