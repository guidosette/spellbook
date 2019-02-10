package page

import (
	"context"
	"distudio.com/mage"
	"encoding/json"
	"net/http"
)

type fieldError struct {
	error
	field string
}

func (err fieldError) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Field string
		Error string
	}
	return json.Marshal(Alias{err.field, err.Error()})
}

type Errors struct {
	errors []fieldError
}

func (errs Errors) HasErrors() bool {
	return errs.errors != nil
}

func (errs *Errors) AddError(name string, value error) {
	errs.errors = append(errs.errors, fieldError{error:value, field:name})
}

func (errs *Errors) Clear() {
	errs.errors = nil
}

func (errs *Errors) MarshalJSON() ([]byte, error) {
	return json.Marshal(errs.errors)
}

type RedirectController struct {
	mage.Controller
	To string
}

func (controller *RedirectController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	return mage.Redirect{Location:controller.To, Status:http.StatusFound}
}

func (controller *RedirectController) OnDestroy(ctx context.Context) {}
