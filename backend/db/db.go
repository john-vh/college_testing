package db

import (
	"errors"
)

var ErrNoRows = errors.New("No matching rows")
var ErrForeignKey = errors.New("Foreign key violation")
var ErrUnique = errors.New("Unique constraint violation")
var ErrDB = errors.New("Internal database error")
