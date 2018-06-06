package gosql

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type TestTable struct {
	ID    int    `gosql:"primary key not null"`
	Value string `gosql:"not null"`
}

func TestDatabase(t *testing.T) {
	fmt.Println(t.Name())

	db := DB{DriverName: "sqlite3"}
	var err error
	db.DB, err = sql.Open("sqlite3", "test.db")
	if err != nil {
		t.Error(err)
		return
	}
	defer db.DB.Close()

	err = db.Create(TestTable{})
	if err != nil {
		t.Error(err)
		return
	}

	err = db.InsertOrReplace(TestTable{5, "Foo"})
	if err != nil {
		t.Error(err)
	}
}
