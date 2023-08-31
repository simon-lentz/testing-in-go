package main

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
)

func TestUserStore(t *testing.T) {
	const (
		createDB        = `CREATE DATABASE test_user_store;`
		dropDB          = `DROP DATABASE IF EXISTS test_user_store;`
		createUserTable = `CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT UNIQUE NOT NULL,
		);`
	)
	psql, err := sql.Open("postgres",
		"host=localhost port=5432 user=simon sslmode=disable")
	if err != nil {
		t.Fatalf("sql.Open() err = %s", err)
	}
	defer psql.Close()

	_, err = psql.Exec(dropDB)
	if err != nil {
		t.Fatalf("psql.Exec() err = %s", err)
	}
	_, err = psql.Exec(createDB)
	if err != nil {
		t.Fatalf("psql.Exec() err = %s", err)
	}
	defer func() {
		// teardown
		_, err = psql.Exec(dropDB)
		if err != nil {
			t.Fatalf("psql.Exec() err = %s", err)
		}
	}()

	db, err := sql.Open("postgres",
		"host=localhost port=5432 user=simon sslmode=disable dbname=test_user_store")
	if err != nil {
		t.Fatalf("sql.Open() err = %s", err)
	}
	defer db.Close()
	_, err = db.Exec(createUserTable)
	if err != nil {
		t.Fatalf("db.Exec() err = %s", err)
	}

	us := &UserStore{
		sql: db,
	}

	t.Run("Find()", testUserStore_Find(us))
	/*
		t.Run("Create()", testUserStore_Create(us))
		t.Run("Delete()", testUserStore_Delete(us))
	*/
}

func testUserStore_Find(us *UserStore) func(t *testing.T) {
	return func(t *testing.T) {
		simon := &User{
			Name:  "Simon",
			Email: "simon@test.com",
		}
		err := us.Create(simon)
		if err != nil {
			t.Errorf("us.Create() err = %s", err)
		}
		defer func() {
			err := us.Delete(simon.ID)
			if err != nil {
				t.Errorf("us.Delete() err = %s", err)
			}
		}()

		tests := []struct {
			name    string
			id      int
			want    *User
			wantErr error
		}{
			{"Found", simon.ID, simon, nil},
			{"Not Found", -1, nil, ErrNotFound},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := us.Find(tt.id)
				if err != tt.wantErr {
					t.Errorf("us.Find() err = %s", err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("us.Find() = %+v; want %+v", got, tt.want)
				}
			})
		}
	}
}
