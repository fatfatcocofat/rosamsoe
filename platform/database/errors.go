package database

import (
	"strings"

	"gorm.io/gorm"
)

func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "duplicate key value violates unique")
}

func IsRecordNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	return err == gorm.ErrRecordNotFound
}
