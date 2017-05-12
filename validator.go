package page

import (
	"distudio.com/mage"
)

type Validator interface {
	Validate(value string) error
	ValidateArray(values []string) bool
	ErrorMessage() string
}

type FieldError struct {
	Field        string
	FieldForUser string
	Message      string
}

func (e FieldError) Error() string {
	if (len(e.FieldForUser) > 0 ) {
		return "Field: " + e.FieldForUser + ", Message: " + e.Message;
	} else {
		return "Field: " + e.Field + ", Message: " + e.Message;
	}
}

type Field struct {
	Name         string
	NameForUser  string
	validators   []Validator
	Required     bool
	Valid        bool
	ErrorMissing string
	value        string
	values       []string
	hasValue     bool
	isArray      bool
}

func NewFieldArray(name string, in mage.RequestInputs) *Field {
	f := NewField(name, in);
	f.isArray = true;
	if (f.hasValue) {
		f.values = in[f.Name].Values();
	}

	return f;
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
		if (!field.isArray) {
			if !field.hasValue || field.value == "" {
				return FieldError{Field:field.Name, Message:field.ErrorMissing, FieldForUser:field.NameForUser};
			}
		} else {
			if !field.hasValue || field.values[0] == "" {
				return FieldError{Field:field.Name, Message:field.ErrorMissing, FieldForUser:field.NameForUser};
			}
		}

	}
	if (!field.isArray) {
		for _, v := range field.validators {
			if v.Validate(field.value) != nil {
				return FieldError{Field:field.Name, Message:v.ErrorMessage(), FieldForUser:field.NameForUser};
			}
		}
	} else {
		for _, v := range field.validators {
			if !v.ValidateArray(field.values) {
				return FieldError{Field:field.Name, Message:v.ErrorMessage(), FieldForUser:field.NameForUser};
			}
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