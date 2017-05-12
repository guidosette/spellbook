package page

import (
	"net/mail"
	"fmt"
	"strings"
	"errors"
)

//structs implementing "Validator"

//checks if the given string is an email address
type EmailValidator struct {
	Message string
}

func (validator EmailValidator) Validate(value string) error {
	_, err := mail.ParseAddress(value);
	return err;
}

func (validator EmailValidator) ValidateArray(value []string) bool {
	return true;
}

func (validator EmailValidator) ErrorMessage() string {
	return validator.Message;
}


//Validates the len of an array
type ArrayLenValidator struct {
	MinLenArray int
	MaxLenArray int
}

func (v ArrayLenValidator) Validate(value string) error {
	return nil;
}

func (v ArrayLenValidator) ValidateArray(values []string) bool {

	if v.MaxLenArray < 0 && v.MinLenArray < 0 {
		return true;
	}

	l := len(values);

	if v.MaxLenArray < 0 {
		return l >= v.MinLenArray
	}

	if v.MinLenArray < 0 {
		return l <= v.MaxLenArray;
	}

	return l >= v.MinLenArray && l <= v.MaxLenArray;

}

func (v ArrayLenValidator) ErrorMessage() string {
	if v.MaxLenArray < 0 && v.MinLenArray < 0 {
		return "";
	}

	if v.MaxLenArray < 0 {
		return fmt.Sprintf("Seleziona almeno %d opzione.", v.MinLenArray)
	}

	if v.MinLenArray < 0 {
		return fmt.Sprintf("Seleziona al massimo %d opzioni.", v.MaxLenArray)
	}

	return fmt.Sprintf("Seleziona da %d a %d opzioni.", v.MinLenArray, v.MaxLenArray)
}





//Validates the len of a string
type LenValidator struct {
	MinLen int
	MaxLen int
}

func (v *LenValidator) Validate(value string) bool {

	if v.MaxLen < 0 && v.MinLen < 0 {
		return true;
	}

	l := len(value);

	if v.MaxLen < 0 {
		return l >= v.MinLen
	}

	if v.MinLen < 0 {
		return l <= v.MaxLen;
	}

	return l >= v.MinLen && l <= v.MaxLen;

}

func (v *LenValidator) ErrorMessage() string {
	if v.MaxLen < 0 && v.MinLen < 0 {
		return "";
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
	Position int
	IgnoreCase bool
	Mode SubstrMode
}

type SubstrMode int;

const (
	SUBSTR_MODE_OR SubstrMode = 0
	SUBSTR_MODE_AND SubstrMode = 1
)

func (v *SubstrValidator) Validate(value string) bool {
	if len(value) - 1 < v.Position {
		panic(errors.New("Index specified is larger than value length!"));
	}

	if v.Position < 0 {
		for _, v := range v.Against {
			if !strings.Contains(value, v) {
				return false;
			}
		}
		return true;
	}

	for _, against := range v.Against {
		l := len(against);
		sub := value[v.Position:l];

		if v.IgnoreCase {
			sub = strings.ToUpper(sub);
			against = strings.ToUpper(against);
		}

		if v.Mode == SUBSTR_MODE_AND && sub != against {
			return false;
		}

		if v.Mode == SUBSTR_MODE_OR && sub == against {
			return true;
		}
	}

	if v.Mode == SUBSTR_MODE_OR {
		return false;
	}

	return true;

}

func (v *SubstrValidator) ErrorMessage() string {
	s := "";
	for _, v := range v.Against {
		s = fmt.Sprintf("%s %s", s, v);
	}

	return fmt.Sprintf("Il valore deve contenere i caratteri: %s", s);
}


