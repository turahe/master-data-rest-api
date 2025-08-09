package seeders

import (
	"strings"
)

// isNotFoundError checks if the error is a "not found" error
func isNotFoundError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}
