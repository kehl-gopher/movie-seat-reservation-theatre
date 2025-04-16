package utility

import (
	"fmt"
	"net/mail"
	"regexp"
)

var NameRegexPattern = `^[A-Za-z]+$`

func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %v", v.Errors)
}

func (v *ValidationError) CheckError() bool {
	return len(v.Errors) > 0
}

func (v *ValidationError) AddValidationError(key, msg string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}
	v.Errors[key] = msg
}

func (v *ValidationError) CheckValid(check bool, field string, msg string) {
	if !check {
		v.AddValidationError(field, msg)
	}
}

func (v *ValidationError) ValidateEmailAddress(email string) string {
	e, err := mail.ParseAddress(email)

	if err != nil {
		v.AddValidationError("email", "Invalid email address")
		return ""
	}
	return e.Address
}

func MatchRegexPattern(name string, pattern string) bool {
	r, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	return r.MatchString(name)
}
func (v *ValidationError) ValidatePassword(password string) string {
	if len(password) < 3 {
		v.AddValidationError("password", "Password length is short")
		return ""
	}
	return password
}

func (v *ValidationError) ValidateName(name string, tag string) string {
	if name == "" {
		v.AddValidationError(tag, fmt.Sprintf("%s is empty", tag))
		return ""
	} else if len(name) < 3 {
		v.AddValidationError(tag, fmt.Sprintf("%s length is short", tag))
		return ""
	} else if len(name) > 50 {
		v.AddValidationError(tag, fmt.Sprintf("%s length is long", tag))
		return ""
	} else if !MatchRegexPattern(name, NameRegexPattern) {
		v.AddValidationError(tag, fmt.Sprintf("%s cannot be alpha numeric or contain invalid character", tag))
		return ""
	}

	return name
}
