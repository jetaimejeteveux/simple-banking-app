package helper

import (
	"errors"

	"gorm.io/gorm"
)

// IsRecordNotFound checks if the error is a gorm.ErrRecordNotFound error.
func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
