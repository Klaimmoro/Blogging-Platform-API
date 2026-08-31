package dberrors

import (
	"errors"

	"github.com/lib/pq"
)

func IsUnique(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "2305"
	}
	return false
}
