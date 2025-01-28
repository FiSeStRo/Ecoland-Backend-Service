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
		if i == len(email)-1 && v == '@' {
			return false
		}
	}
	return hasAt
}
