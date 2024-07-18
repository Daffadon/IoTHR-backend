package utils

import "strings"

func IsUniqueConstraintError(err error) bool {
	// Check for MySQL duplicate entry error
	if strings.Contains(err.Error(), "Error 1062") && strings.Contains(err.Error(), "for key 'uni_users_email'") {
		return true
	}
	return false
}
