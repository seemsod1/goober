package forms

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		return false
	}
	return true
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// Maxlength checks for string minimum length
func (f *Form) Maxlength(field string, length int) bool {
	x := f.Get(field)
	if len(x) > length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at most %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Var(f.Get(field), "email"); err != nil {
		f.Errors.Add(field, "Invalid email address")
	}
}

// IsPhone checks for valid phone number
func (f *Form) IsPhone(field string) {
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Var(f.Get(field), "e164"); err != nil {
		f.Errors.Add(field, "Invalid email address")
	}

}

// IsNumber checks for valid number
func (f *Form) IsNumber(field string) {
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Var(f.Get(field), "number"); err != nil {
		f.Errors.Add(field, "This field must be a number")
	}
}

// IsPlate checks for valid plate number
func (f *Form) IsPlate(field string) {
	pattern := "^[ABEIKMHOPCTXYZ]{2}\\d{4}[ABEIKMHOPCTXYZ]{2}$"

	if match, _ := regexp.MatchString(pattern, f.Get(field)); !match {
		f.Errors.Add(field, "Invalid plate number")
	}

}

// IsName checks for valid name
func (f *Form) IsName(field string) {

	pattern := "^[a-zA-Zа-яА-ЯіІґҐєЄїЇ]+$"

	if match, _ := regexp.MatchString(pattern, f.Get(field)); !match {
		f.Errors.Add(field, "Invalid name")
	}

}

// MinNumber checks for minimum number
func (f Form) MinNumber(field string, number int) {
	if f.Has(field) {
		x, _ := strconv.Atoi(f.Get(field))
		if x < number {
			f.Errors.Add(field, fmt.Sprintf("This field must be at least %d", number))
		}
	}
}

// MaxNumber checks for maximum number
func (f Form) MaxNumber(field string, number int) {
	if f.Has(field) {
		x, _ := strconv.Atoi(f.Get(field))
		if x > number {
			f.Errors.Add(field, fmt.Sprintf("This field must be at most %d", number))
		}
	}
}

// IsPasswordMatch checks for matching passwords
func (f *Form) IsPasswordMatch(field1, field2 string) {
	if f.Get(field1) != f.Get(field2) {
		f.Errors.Add(field1, "Passwords do not match")
		f.Errors.Add(field2, "Passwords do not match")
	}
}
