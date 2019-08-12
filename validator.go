package snippets

import (
	"errors"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Validate checks if the values of a User struct has met a set of requirements
// and returns an error should it fail any of it
func (u User) Validate() error {
	u.Username = strings.Trim(u.Username, " ")
	u.FirstName = strings.Trim(u.FirstName, " ")
	u.LastName = strings.Trim(u.LastName, " ")
	if len(strings.Trim(u.Password, " ")) < 1 {
		u.Password = ""
	}

	return validation.ValidateStruct(&u,
		validation.Field(&u.ID, validation.Skip, is.UUIDv4),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Username, validation.Required, validation.Length(2, 25)),
		validation.Field(&u.Password, passwordRules...),
		validation.Field(&u.PasswordHash, validation.Skip),
		validation.Field(&u.FirstName, validation.Required, validation.Length(1, 50)),
		validation.Field(&u.LastName, validation.Required, validation.Length(1, 50)),
	)
}

// Regex based custom validation rules that implements the validation.Rule interface
var checkLowercasePresent = createRegexValidator(`(?:.*[a-z].*)`, "at least 1 lowercase character is required")
var checkUppercasePresent = createRegexValidator(`(?:.*[A-Z].*)`, "at least 1 uppercase character is required")
var checkNumberPresent = createRegexValidator(`(?:.*[0-9].*)`, "at least 1 number is required")
var checkSpecialCharPresent = createRegexValidator(`(?:.*[!@#$%^&*].*)`, "at least 1 special character is required")

// A collection of rules for password input validation
var passwordRules = []validation.Rule{
	validation.Required,
	validation.Length(8, 127),
	validation.By(checkLowercasePresent),
	validation.By(checkUppercasePresent),
	validation.By(checkNumberPresent),
	validation.By(checkSpecialCharPresent),
}

// createRegexValidator generates a regex validator that implements the validation.Rule interface
func createRegexValidator(pattern string, err string) func(interface{}) error {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errors.New("only string is allowed")
		}
		if !matchStringToPattern(pattern, s) {
			return errors.New(err)
		}
		return nil
	}
}

// matchStringToPattern normalizes regexp.MatchString's return value as a boolean
func matchStringToPattern(pattern string, value string) bool {
	match, _ := regexp.MatchString(pattern, value)
	return match
}
