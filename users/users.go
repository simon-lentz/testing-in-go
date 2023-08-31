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
	ID      int
	Name    string
	Email   string
	Balance int
}

type UserStore interface {
	Tx(func(UserStore) error) error
	Find(id int) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id int) error
}

// PsqlUserStore is used to interact with the user store.
type PsqlUserStore struct {
	tx interface {
		Begin() (*sql.Tx, error)
	}
	sql interface {
		Exec(query string, args ...interface{}) (sql.Result, error)
		QueryRow(query string, args ...interface{}) *sql.Row
	}
}

// implements the transaction method for the UserStore interface
func (us *PsqlUserStore) Tx(fn func(us UserStore) error) error {
	tx, err := us.tx.Begin()
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "race: failed to begin transaction")
	}
	txStore := &PsqlUserStore{
		// using specific transaction rather than database connection
		sql: tx,
	}
	err = fn(txStore)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "race: failed to commit transaction")
	}
	return nil
}

// Find will retrieve a user via a user id,
// returning either a user, ErrNotFound, or a wrapped error.
func (us *PsqlUserStore) Find(id int) (*User, error) {
	const query = `SELECT id, name, email, balance FROM users WHERE id=$1;`
	row := us.sql.QueryRow(query, id)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Balance)
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
func (us *PsqlUserStore) Create(user *User) error {
	const query = `INSERT INTO users (name, email, balance) VALUES ($1, $2, $3) RETURNING id;`
	err := us.sql.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return errors.Wrap(err, "psql: error creating new user")
	}
	return nil
}

// Update will update a user in the DB.
func (us *PsqlUserStore) Update(user *User) error {
	const query = `UPDATE users SET name=$2 email=$3 balance=$4 WHERE id=$1;`
	_, err := us.sql.Exec(query, user.ID, user.Name, user.Email, user.Balance)
	if err != nil {
		return errors.Wrap(err, "race: error updating user")
	}
	return nil
}

// Delete removes a user from the database,
// or returns a wrapped error.
func (us *PsqlUserStore) Delete(id int) error {
	const query = `DELETE FROM users WHERE id=$1;`
	_, err := us.sql.Exec(query, id)
	if err != nil {
		return errors.Wrap(err, "psql: error deleted user")
	}
	return nil
}

type Transaction interface {
	Tx(func(UserStore) error) error
}

func Spend(tx Transaction, userID int, amount int) error {
	return tx.Tx(func(us UserStore) error {
		user, err := us.Find(userID)
		if err != nil {
			return err
		}
		user.Balance -= amount
		return us.Update(user)
	})
}
