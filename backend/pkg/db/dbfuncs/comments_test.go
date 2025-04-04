package dbfuncs

import (
	"database/sql"
	"errors"
	"testing"
)

type sickDatabase struct {
}

var sdb Database = sickDatabase{}

/*
type Database interface {
	QueryRow(query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
}*/

func (db sickDatabase) QueryRow(query string, args ...any) *sql.Row {
	return nil
}

func (db sickDatabase) Prepare(query string) (*sql.Stmt, error) {
	return nil, errors.New("default prepare error")
}

func (db sickDatabase) Query(query string, args ...any) (*sql.Rows, error) {
	return nil, errors.New("default query error")
}

func TestAddComment(t *testing.T) {
	//make a directory
	//defer delete directory
	//setting directory to nil, if it ever tries to use it it should nil pointer derefence panic
	Configure(sdb, nil)

	//testing without image file
	testComment := Comment{}
	_, err := AddComment(&testComment, nil)
	if err.Error() != "default prepare error" {
		t.Error(err)
		t.Fatal("should have given a prepare error")
	}

	// testing with image file
	// testImageDirectory := "./testimages"
	// err = os.Mkdir(testImageDirectory, 0)
	// if err != nil {
	// 	t.Error(err)
	// }

}
