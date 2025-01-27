package utils

func IsEmail(email string) bool {
	hasAt := false
	for i, v := range email {
		if i == 0 && v == '@' {
			return false
		}
		if v == '@' {
			hasAt = true
		}
	}
	return hasAt
}
