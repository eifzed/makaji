package databaseerr

import "errors"

var (
	ErrorDataNotFound = errors.New("Data not found")
	ErrorNoUpdate     = errors.New("No data updated")
	ErrorNoInsert     = errors.New("No data inserted")
	ErrorNoDelete     = errors.New("No data deleted")
)
