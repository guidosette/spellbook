package validators

import (
	"distudio.com/mage"
	"strings"
)


// Validator operates on a single field.
// Returns a validation error if the field does not satisfy a set of constraint.
// Each validator acts upon a given constraint
type Validator interface {
	Validate(value string) error
}

type Field struct {
	Name         string
	validators   []Validator
	Required     bool
	IsValidated bool
	value       string
}

func NewField(name string, required bool, in mage.RequestInputs) *Field {
	vs := make([]Validator, 0, 0)
	f := &Field{Name: name, Required: required, validators: vs}
	if val, ok := in[f.Name]; ok {
		f.value = strings.TrimSpace(val.Value())
	}

	return f
}

func NewRawField(name string, required bool, value string) *Field {
	vs := make([]Validator, 0, 0)
	f := &Field{Name: name, Required:required, validators:vs}
	f.value = value
	return f
}

func (field *Field) AddValidator(v Validator) {
	field.validators = append(field.validators, v)
}

func (field *Field) Validate() error {
	if field.Required && field.value == "" {
		return ErrMissingField
	}

	for _, v := range field.validators {
		if err := v.Validate(field.value); err != nil {
			return err
		}
	}

	field.IsValidated = true

	return nil
}

func (field *Field) Value() (string, error) {
	if !field.IsValidated {
		err := field.Validate()
		if err != nil {
			return "", err
		}
	}
	return field.value, nil
}

func (field Field) MustValue() string {
	val, err := field.Value()
	if err != nil {
		panic(err)
	}
	return val
}
