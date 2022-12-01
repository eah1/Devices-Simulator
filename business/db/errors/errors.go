// Package errors provides support for access the database.
package errors

import "fmt"

// PsqlError errors from sql.
type PsqlError struct {
	CodeSQL        string
	TableName      string
	ConstraintName string
	Err            string
}

// Error redefined error from struct PsqlError.
func (err *PsqlError) Error() string {
	return fmt.Sprintf("code sql %s, table %s, constraints %s : %s",
		err.CodeSQL, err.TableName, err.ConstraintName, err.Err)
}
