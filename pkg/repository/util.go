package repository

import "strings"

func IsUniqueError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
