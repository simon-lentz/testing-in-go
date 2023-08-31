package users

import (
	"database/sql"

	"github.com/pkg/errors"
)

// Common errors the we want to account for explicitly.
var (
	ErrNotFound = errors.New("psql: resource could not be located")
)

// User would not typically be defined here, but is done for simplicity.
type User struct {
	ID    int
	Name  string
	Email string
}

// UserStore is used to interact with the user store.
type UserStore struct {
	sql interface {
		Exec(query string, args ...interface{}) (sql.Result, error)
		QueryRow(query string, args ...interface{}) *sql.Row
	}
}

// Find will retrieve a user via a user id,
// returning either a user, ErrNotFound, or a wrapped error.
func (us *UserStore) Find(id int) (*User, error) {
	const query = `SELECT id, name, email FROM users WHERE id=$1;`
	row := us.sql.QueryRow(query, id)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	switch err {
	case sql.ErrNoRows:
		return nil, ErrNotFound
	case nil:
		return &user, nil
	default:
		return nil, errors.Wrap(err, "psql: error querying for user by id")
	}
}

// Create will create a new user in the db and
// update the user row with the returned id, or
// return a wrapped error.
func (us *UserStore) Create(user *User) error {
	const query = `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id;`
	err := us.sql.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return errors.Wrap(err, "psql: error creating new user")
	}
	return nil
}

// Delete removes a user from the database,
// or returns a wrapped error.
func (us *UserStore) Delete(id int) error {
	const query = `DELETE FROM users WHERE id=$1;`
	_, err := us.sql.Exec(query, id)
	if err != nil {
		return errors.Wrap(err, "psql: error deleted user")
	}
	return nil
}
