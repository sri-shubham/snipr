package util

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("Not Found")

func PresentStorageErrors(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrNotFound
	default:
		return err
	}
}
