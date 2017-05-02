package page

import (
	"distudio.com/mage"
)

type Validator interface {
	Validate(value string) error
	ErrorMessage() string
}

type FieldError struct {
	Field string
	Message string
}

func (e FieldError) Error() string {
	return e.Message;
}

type Field struct {
	Name string
	validators []Validator
	Required bool
	Valid bool
	ErrorMissing string
	value string
	hasValue bool
}

func NewField(name string, in mage.RequestInputs) *Field {
	vs := make([]Validator, 0, 0);
	f := &Field{Name:name, validators:vs};
	f.Valid = false;
	_, hasValue := in[f.Name];
	f.hasValue = hasValue;
	if hasValue {
		f.value = in[f.Name].Value();
	}
	return f;
}

func (field *Field) AddValidator(v Validator) {
	field.validators = append(field.validators, v);
}

func (field *Field) Validate() error {

	field.Valid = false;
	if field.Required {
		if !field.hasValue || field.value == "" {
			return FieldError{Field:field.Name, Message:field.ErrorMissing};
		}

	}


	for _, v := range field.validators {
		if v.Validate(field.value) != nil {
			return FieldError{Field:field.Name, Message:v.ErrorMessage()};
		}
	}

	field.Valid = true;
	return nil;
}

//safe to call after validate
func (field *Field) Value() string {
	if !field.Valid {
		return "";
	}
	return field.value;
}