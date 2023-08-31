package users

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	// Can do flag.Parse() if needed.
	exitCode := run(m)
	os.Exit(exitCode)
}

func run(m *testing.M) int {
	const (
		createDB = `CREATE DATABASE test_user_store;`
		dropDB   = `DROP DATABASE IF EXISTS test_user_store;`
	)
	psql, err := sql.Open("postgres",
		"host=localhost port=5432 user=simon sslmode=disable")
	if err != nil {
		panic(fmt.Errorf("sql.Open() err = %s", err))
	}
	defer psql.Close()

	_, err = psql.Exec(dropDB)
	if err != nil {
		panic(fmt.Errorf("psql.Exec() err = %s", err))
	}
	_, err = psql.Exec(createDB)
	if err != nil {
		panic(fmt.Errorf("psql.Exec() err = %s", err))
	}
	defer func() {
		// teardown
		_, err = psql.Exec(dropDB)
		if err != nil {
			panic(fmt.Errorf("psql.Exec() err = %s", err))
		}
	}()

	return m.Run()
}
func TestUserStore(t *testing.T) {
	const createUserTable = `CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT UNIQUE NOT NULL,
	);`
	db, err := sql.Open("postgres",
		"host=localhost port=5432 user=simon sslmode=disable dbname=test_user_store")
	if err != nil {
		panic(fmt.Errorf("sql.Open() err = %s", err))
	}
	defer db.Close()
	_, err = db.Exec(createUserTable)
	if err != nil {
		panic(fmt.Errorf("db.Exec() err = %s", err))
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

type unsafeUserStore struct {
	*UserStore
	wg *sync.WaitGroup
}

func (unsafe *unsafeUserStore) Find(id int) (*User, error) {
	user, err := unsafe.UserStore.Find(id)
	if err != nil {
		return nil, err
	}
	unsafe.wg.Done()
	unsafe.wg.Wait()
	return user, err
}

func TestSpend(t *testing.T) {
	db, err := sql.Open("postgres",
		"host=localhost port=5432 user=simon sslmode=disable dbname=test_user_store")
	if err != nil {
		panic(fmt.Errorf("sql.Open() err = %s", err))
	}
	defer db.Close()

	us := &UserStore{
		sql: db,
	}

	simon := &User{
		Name:    "Simon",
		Email:   "simon@test.com",
		Balance: 100,
	}
	err = us.Create(simon)
	if err != nil {
		t.Errorf("us.Create() err = %s", err)
	}
	defer func() {
		err := us.Delete(simon.ID)
		if err != nil {
			t.Errorf("us.Delete() err = %s", err)
		}
	}()

	unsafe := &unsafeUserStore{
		UserStore: us,
		wg:        &sync.WaitGroup{},
	}
	unsafe.wg.Add(2)
	var spendWg sync.WaitGroup
	// done := make(chan bool)
	for i := 0; i < 2; i++ {
		spendWg.Add(1)
		go func() {
			err := Spend(unsafe, simon.ID, 20)
			if err != nil {
				panic(fmt.Errorf("Spend() err = %s", err))
			}
			spendWg.Done()
			// done <- true
		}()
	}
	// Wait until both goroutines testing the Spend() function
	// are done before testing the outcome. This could also
	// be done using a channel (commented out above).
	spendWg.Wait()
	// Wait until two values have been received from the done channel
	// before moving on.
	// <-done
	// <-done
	got, err := us.Find(simon.ID)
	if err != nil {
		t.Fatalf("us.Find() err = %s", err)
	}
	if got.Balance != 60 {
		t.Fatalf("user.Balance = %d; want %d", got.Balance, 60)
	}

}
