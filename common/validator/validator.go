package validator

import "regexp"

var (
	validator map[string]func(any) bool

	emailMatcher *regexp.Regexp
	phoneMatcher *regexp.Regexp
)

const (
	emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	phoneRegex = `^1[3456789]\d{9}$`
)

func init() {
	validator = make(map[string]func(x any) bool)
	emailMatcher = regexp.MustCompile(emailRegex)
	phoneMatcher = regexp.MustCompile(phoneRegex)

	validator["username"] = func(x any) bool {
		username := x.(string)
		if len(username) < 6 || len(username) > 20 {
			return false
		}
		return true
	}

	validator["password"] = func(x any) bool {
		pwd := x.(string)
		if len(pwd) < 6 || len(pwd) > 20 {
			return false
		}
		return true
	}

	validator["email"] = func(x any) bool {
		email := x.(string)
		return emailMatcher.Match([]byte(email))
	}

	validator["phone"] = func(x any) bool {
		phone := x.(string)
		return phoneMatcher.Match([]byte(phone))
	}
}

func ValidateUsername(username string) bool {
	return validator["username"](username)
}

func ValidatePassword(password string) bool {
	return validator["password"](password)
}

func ValidateEmail(email string) bool {
	return validator["email"](email)
}

func ValidatePhone(phone string) bool {
	return validator["phone"](phone)
}
