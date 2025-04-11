package dbfuncs

import (
	"database/sql"
	"errors"
	"testing"
)

// create a bunch of fake databases structs that act in different ways for each test I want to do
// or a single fake database with SetQueryRow, SetPrepare, SetQuery which sets the method behaviour per test
// maybe MockDatabase should be lowercase, maybe extract and have it common to all tests.
/*
type Database interface {
	QueryRow(query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
}*/

type MockDatabase struct {
	MockQueryRow func(query string, args ...any) *sql.Row
	MockPrepare  func(query string) (*sql.Stmt, error)
	MockQuery    func(query string, args ...any) (*sql.Rows, error)
}

func (db MockDatabase) QueryRow(query string, args ...any) *sql.Row {
	return db.MockQueryRow(query, args)
}

func (db MockDatabase) Prepare(query string) (*sql.Stmt, error) {
	return db.MockPrepare(query)
}

func (db MockDatabase) Query(query string, args ...any) (*sql.Rows, error) {
	return db.MockQuery(query, args)
}

// test pattern to follow for more tests that are doable without a ton of IO mocking required,
// maybe do IO heavy ones at some point but not priority
func TestAddCommentWillNotExecuteBadPreparedStatement(t *testing.T) {

	//Arrange
	testComment := Comment{}

	mockDb := MockDatabase{}
	expected_error_message := "prepare error"
	mockDb.MockPrepare = func(query string) (*sql.Stmt, error) {
		return nil, errors.New(expected_error_message)
	}
	var nil_imagefolder *string = nil //I *want* a nil panic if this path is actually used when there is no image in the comment
	Configure(mockDb, nil_imagefolder)

	//Act
	_, err := AddComment(&testComment, nil)

	//Assert
	if err.Error() != expected_error_message {
		t.Error(err)
		t.Fatal("should have given a prepare error")
	}
}

// func TestAddComment(t *testing.T) {
// 	// make empty mock
// 	mockDb := MockDatabase{}

// 	// set up methods for upcoming tests
// 	mockDb.MockQueryRow = func(query string, args ...any) *sql.Row {
// 		return nil
// 	}

// 	mockDb.MockPrepare = func(query string) (*sql.Stmt, error) {
// 		return nil, errors.New("prepare error")
// 	}

// 	mockDb.MockQuery = func(query string, args ...any) (*sql.Rows, error) {
// 		return nil, errors.New("query error")
// 	}

// 	// e.g. first test: if preprare fails, does addcomment also fail and return the error
// 	// if it doesn't fail, test if statement.Exec works (maybe need to mock up Exec as well?)
// 	//make a directory
// 	//defer delete directory
// 	//setting directory to nil, if it ever tries to use it it should nil pointer derefence panic
// 	Configure(sdb, nil)

// 	//testing without image file
// 	testComment := Comment{}
// 	_, err := AddComment(&testComment, nil)
// 	if err.Error() != "default prepare error" {
// 		t.Error(err)
// 		t.Fatal("should have given a prepare error")
// 	}

// 	// testing with image file
// 	// testImageDirectory := "./testimages"
// 	// err = os.Mkdir(testImageDirectory, 0)
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// }

// 	// set methods to different functions for upcoming test, e.g. what if prepare gives particular error or no error, or particular outputs
// 	// in the case of addcomment, it only uses prepare so other two are irrelevant apart from needing to exist to meet interface.
// 	mockDb.MockQueryRow = func(query string, args ...any) *sql.Row {
// 		return nil
// 	}
// 	// returning an empty statement
// 	mockDb.MockPrepare = func(query string) (*sql.Stmt, error) {
// 		outputStmt := sql.Stmt{}
// 		return &outputStmt, nil
// 	}

// 	mockDb.MockQuery = func(query string, args ...any) (*sql.Rows, error) {
// 		return nil, errors.New("query error")
// 	}

// }
