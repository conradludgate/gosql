package gosql

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

type statements struct {
	ior *sql.Stmt
}

type DB struct {
	DB         *sql.DB
	DriverName string

	stmts map[string]statements
}

func (db DB) Create(t interface{}) error {
	if db.DB == nil {
		return errors.New("Uninitialised DB")
	}

	typ := reflect.TypeOf(t)
	if typ.Kind() != reflect.Struct {
		return errors.New("Parameter must be a struct type")
	}

	typs := types[db.DriverName]
	if typs == nil {
		return errors.New("Unknown driver")
	}

	n := typ.NumField()
	fields := make([]string, 0, n)

	for i := 0; i < n; i++ {
		sf := typ.Field(i)

		ts := typs[sf.Type]

		if ts == "" {
			return errors.New("Invalid type in struct")
		}

		fields = append(fields, sf.Name+" "+sf.Type.String()+" "+sf.Tag.Get("gosql"))
	}

	_, err := db.DB.Exec("CREATE TABLE " + typ.Name() + "(" + strings.Join(fields, ",") + ");")
	return err
}

func (db DB) InsertOrReplace(t interface{}) error {
	if db.DB == nil {
		return errors.New("Uninitialised DB")
	}

	typ := reflect.TypeOf(t)
	if typ.Kind() != reflect.Struct {
		return errors.New("Parameter must be a struct type")
	}

	typs := types[db.DriverName]
	if typs == nil {
		return errors.New("Unknown driver")
	}

	if db.stmts == nil {
		db.stmts = make(map[string]statements)
	}

	n := typ.NumField()

	values := make([]interface{}, 0, n)

	v := reflect.ValueOf(t)
	stmt, ok := db.stmts[typ.Name()]
	if !ok || stmt.ior == nil {
		query := "INSERT OR REPLACE INTO " + typ.Name() + "("
		end := ") VALUES ("
		first := true
		for i := 0; i < n; i++ {
			if first {
				first = false
			} else {
				query += ","
				end += ","
			}

			query += typ.Field(i).Name
			end += "?"

			values = append(values, v.Field(i).Interface())
		}

		var err error
		stmt.ior, err = db.DB.Prepare(query + end + ")")
		if err != nil {
			return err
		}
	}

	if len(values) == 0 {
		for i := 0; i < n; i++ {
			values = append(values, v.Field(i).Interface())
		}
	}

	_, err := stmt.ior.Exec(values...)
	return err
}
