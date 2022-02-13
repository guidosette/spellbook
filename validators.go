package spellbook

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

// checks if the given string is an email address
type EmailValidator struct{}

func (validator EmailValidator) Validate(value string) error {
	_, err := mail.ParseAddress(value)
	return err
}

// Validates the len of a string
// It can validate both maximum and minimum len
// If MaxLen or MinLen is a number lesser or equal to zero the constraint is ignored
type LenValidator struct {
	MinLen int
	MaxLen int
}

func (v LenValidator) Validate(value string) error {

	validate := false
	if v.MaxLen <= 0 && v.MinLen <= 0 {
		validate = true
	}

	l := len(value)

	if v.MaxLen <= 0 {
		validate = l >= v.MinLen
		if !validate {
			return fmt.Errorf("field must be at least %d characters", v.MinLen)
		} else {
			return nil
		}
	}

	if v.MinLen <= 0 {
		validate = l <= v.MaxLen
		if !validate {
			return fmt.Errorf("field can't be more than %d characters", v.MaxLen)
		} else {
			return nil
		}
	}

	validate = l >= v.MinLen && l <= v.MaxLen
	if !validate {
		return fmt.Errorf("field length must be between %d and %d characters", v.MinLen, v.MaxLen)
	} else {
		return nil
	}
}

// Checks if a given string is a valid datastore name
type DatastoreKeyNameValidator struct{}

func (v DatastoreKeyNameValidator) Validate(value string) error {
	if value == "" {
		return errors.New("string is empty")
	}

	if len(value) > 2 && value[:2] == "__" {
		return fmt.Errorf("%s can't start with '__'", value)
	}

	return nil
}

// Checks if a given string is a valid datastore name
type FileNameValidator struct {
	AllowEmpty bool
}

func (v FileNameValidator) Validate(value string) error {
	if len(value) > 1024 {
		return errors.New("file name can't be larger than 1024 bytes")
	}

	if value == "." || value == "..." || value == ".." {
		return fmt.Errorf("invalid file name: %s", value)
	}

	if strings.HasPrefix(value, ".well-known/acme-challenge") {
		return errors.New("file name can't start with '.well-known/acme-challenge'")
	}

	// todo: validate against unicode chars
	if strings.Contains(value, "\n") || strings.Contains(value, "\r\n") {
		return errors.New("file name can't contain new lines or line feeds")
	}

	if value == "" && !v.AllowEmpty {
		return errors.New("string is empty")
	}

	for _, s := range value {
		if s == '#' || s == '[' || s == ']' || s == '*' || s == '?' {
			return fmt.Errorf("file name %q contains invalid character %q", value, s)
		}
	}

	return nil
}
