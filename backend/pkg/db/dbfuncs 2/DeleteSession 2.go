package dbfuncs

import (
	"fmt"
)

func DeleteSessionColumn(column string, value interface{}) error {
	stmt, err := database.Prepare(fmt.Sprintf("DELETE FROM Sessions WHERE %s = ?", column))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(value)
	if err != nil {
		return err
	}
	return nil

}
