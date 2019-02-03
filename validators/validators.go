package validators

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

//structs implementing "Validator"

//checks if the given string is an email address
type EmailValidator struct {}

func (validator EmailValidator) Validate(value string) error {
	_, err := mail.ParseAddress(value)
	return err
}

//Validates the len of a string
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
			return errors.New(v.ErrorMessage())
		} else {
			return nil
		}
	}

	if v.MinLen <= 0 {
		validate = l <= v.MaxLen
		if !validate {
			return errors.New(v.ErrorMessage())
		} else {
			return nil
		}
	}

	validate = l >= v.MinLen && l <= v.MaxLen
	if !validate {
		return errors.New(v.ErrorMessage())
	} else {
		return nil
	}
}

func (v LenValidator) ErrorMessage() string {
	if v.MaxLen < 0 && v.MinLen < 0 {
		return ""
	}

	if v.MaxLen < 0 {
		return fmt.Sprintf("Il campo deve essere almeno %d lettere.", v.MinLen)
	}

	if v.MinLen < 0 {
		return fmt.Sprintf("Il campo deve essere massimo %d lettere.", v.MaxLen)
	}

	return fmt.Sprintf("Il campo deve essere compreso fra %d e %d lettere.", v.MinLen, v.MaxLen)

}

//validates against the content of a substring
type SubstrValidator struct {
	Against []string
	//negative values ignore the position
	//positive values Against must start at the position specified
	Position   int
	IgnoreCase bool
	Mode       SubstrMode
}

type SubstrMode int

const (
	ModeSubstringOr  SubstrMode = iota
	ModeSubstringAnd
	)

func (v *SubstrValidator) Validate(value string) error {
	if len(value)-1 < v.Position {
		return errors.New("specified index  is larger than value length!")
	}

	if v.Position < 0 {
		for _, v := range v.Against {
			if !strings.Contains(value, v) {
				return fmt.Errorf("string %s does not contain substring %s", value, v)
			}
		}
		return nil
	}

	for _, against := range v.Against {
		l := len(against)
		sub := value[v.Position:l]

		if v.IgnoreCase {
			sub = strings.ToUpper(sub)
			against = strings.ToUpper(against)
		}

		if v.Mode == ModeSubstringAnd && sub != against {
			return fmt.Errorf("string %s does not contain substring %s", value, against)
		}

		if v.Mode == ModeSubstringOr && sub == against {
			return nil
		}
	}

	if v.Mode == ModeSubstringOr {
		return fmt.Errorf("string %s is not valid", value)
	}

	return nil

}
