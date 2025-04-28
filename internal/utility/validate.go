package utility

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"
)

var NameRegexPattern = `^[A-Za-z]+$`
var HallNameRegexPattern = `^[A-Za-z0-9\s]+$`

func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation error")
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
		v.AddValidationError(tag, fmt.Sprintf("%s is required", tag))
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

// movie validation

func (v *ValidationError) ValidateMovieName(name string, tag string) string {
	if name == "" {
		v.AddValidationError(tag, fmt.Sprintf("%s is required", tag))
		return ""
	} else if len(name) < 3 {
		v.AddValidationError(tag, fmt.Sprintf("%s length is short", tag))
		return ""
	} else if len(name) > 50 {
		v.AddValidationError(tag, fmt.Sprintf("%s length is long", tag))
		return ""
	}
	return name
}

func (v *ValidationError) ValidateMovieSynopsis(synopsis string) string {
	if synopsis == "" {
		v.AddValidationError("synopsis", "synopsis is required")
		return ""
	} else if utf8.RuneCountInString(synopsis) < 3 {
		v.AddValidationError("synopsis", "synopsis length is short")
		return ""
	} else if utf8.RuneCountInString(synopsis) > 500 {
		v.AddValidationError("synopsis", "synopsis length is long")
		return ""
	}
	return synopsis
}

func (v *ValidationError) ValidateMovieGenreID(genreID []string) []string {
	if len(genreID) == 0 {
		v.AddValidationError("genre_ids", "genre_ids is required")
		return nil
	}
	return genreID
}

func (v *ValidationError) ValidateMovieDuration(duration uint8) uint8 {
	if duration == 0 {
		v.AddValidationError("duration", "duration is required")
		return 0
	}
	return duration
}

func (v *ValidationError) ValidateImage(base64Image string, tag string) ([]byte, string) {
	const maxImageSize = 5 * 1024 * 1024 // 5MB
	var ext string

	if base64Image == "" {
		return nil, ""
	}

	if !ChecksupportedImageFormat(base64Image) {
		v.AddValidationError(tag, "expected base64 image string, got URL instead")
		return nil, ""
	}
	splitBaseImage := strings.Split(base64Image, ",")
	switch splitBaseImage[0] {
	case "data:image/jpeg;base64":
		ext = "jpeg"
	case "data:image/png;base64":
		ext = "png"
	case "data:image/jpg;base64":
		ext = "jpg"
	default:
		v.AddValidationError(tag, "invalid content type: only PNG, JPG, JPEG are supported")
	}

	byt, err := base64.StdEncoding.DecodeString(splitBaseImage[1])

	if err != nil {
		fmt.Println(err, byt)
		v.AddValidationError("image", "Invalid image format")
		return nil, ""
	}

	// if ValidateMimeType(string(byt))
	if len(byt) > maxImageSize {
		v.AddValidationError("image", "Image size is too large")
		return nil, ""
	}

	return byt, ext
}

func ChecksupportedImageFormat(base64Image string) bool {
	imagePrefixes := []string{"data:image/jpeg;base64", "data:image/png;base64", "data:image/jpg;base64"}
	for _, imagePrefix := range imagePrefixes {
		if strings.HasPrefix(base64Image, imagePrefix) {
			return true
		}
		continue
	}
	return false
}

func (v *ValidationError) ValidateHallName(hallName string) bool {
	if !MatchRegexPattern(hallName, HallNameRegexPattern) {
		v.AddValidationError("hall_name", "hall name cannot contain invalid character")
		return false
	}
	return true
}
