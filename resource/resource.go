package resource

import (
	"context"
	"distudio.com/mage"
	"distudio.com/page/resource/identity"
	"distudio.com/page/validators"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"strconv"
)

type ListOptions struct {
	Size       int
	Page       int
	Order      string
	Descending bool
	Property   string
}

type Manager interface {
	NewResource(ctx context.Context) (Resource, error)
	FromId(ctx context.Context, id string) (Resource, error)
	ListOf(ctx context.Context, opts ListOptions) ([]Resource, error)
	ListOfProperties(ctx context.Context, opts ListOptions) ([]string, error)
	Save(ctx context.Context, resource Resource) error
	Delete(ctx context.Context, resource Resource) error
}

type Resource interface {
	Id() string
	Create(ctx context.Context) error
	// todo: other should be a general serializable object
	Update(ctx context.Context, other Resource) error
}

type PermissionError struct {
	error
	message string
}

func NewPermissionError(permission identity.Permission) PermissionError {
	msg := fmt.Sprintf("%s", identity.Permissions[permission])
	return PermissionError{errors.New(msg), msg}
}

type Controller struct {
	mage.Controller
	Manager Manager
}

func (controller *Controller) IsPublicMethod(method string) bool {
	return true
}

func (controller *Controller) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {

	ins := mage.InputsFromContext(ctx)

	method := ins[mage.KeyRequestMethod].Value()
	if !controller.IsPublicMethod(method) {
		_, ok := ctx.Value(identity.KeyUser).(identity.User)
		if !ok {
			log.Errorf(ctx, "non public controller requires authenticated user")
			return mage.Redirect{Status: http.StatusUnauthorized}
		}
	}

	params := mage.RoutingParams(ctx)
	key, hasKey := params["key"]
	prop, hasProperty := ins["property"]

	switch method {
	case http.MethodPost:
		return controller.HandlePost(ctx, out)
	case http.MethodGet:
		if !hasKey {
			if hasProperty {
				return controller.HandlePropertyValues(ctx, out, prop.Value())
			}
			return controller.HandleList(ctx, out)
		}
		return controller.HandleGet(ctx, key.Value(), out)
	case http.MethodPut:
		if !hasKey {
			log.Errorf(ctx, "no item was specify for put method")
			return mage.Redirect{Status: http.StatusBadRequest}
		}
		return controller.HandlePut(ctx, key.Value(), out)
	case http.MethodDelete:
		if !hasKey {
			log.Errorf(ctx, "no item was specify for delete method")
			return mage.Redirect{Status: http.StatusBadRequest}
		}
		return controller.HandleDelete(ctx, key.Value(), out)
	}

	return mage.Redirect{Status: http.StatusNotImplemented}
}

// REST Method handlers
func (controller *Controller) HandleGet(ctx context.Context, key string, out *mage.ResponseOutput) mage.Redirect {
	renderer := mage.JSONRenderer{}
	out.Renderer = &renderer

	resource, err := controller.Manager.FromId(ctx, key)
	if err != nil {
		if _, ok := err.(PermissionError); ok {
			return mage.Redirect{Status: http.StatusForbidden}
		}
		if _, ok := err.(validators.FieldError); ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	renderer.Data = resource
	return mage.Redirect{Status: http.StatusOK}
}

func (controller *Controller) HandlePropertyValues(ctx context.Context, out *mage.ResponseOutput, prop string) mage.Redirect {
	opts := &ListOptions{}
	opts.Property = prop
	opts, err := controller.BuildOptions(ctx, out, opts)
	if err != nil {
		return mage.Redirect{Status: http.StatusBadRequest}
	}
	results, err := controller.Manager.ListOfProperties(ctx, *opts)
	if err != nil {
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	// output
	l := len(results)
	count := opts.Size
	if l < opts.Size {
		count = l
	}

	renderer := mage.JSONRenderer{}
	out.Renderer = &renderer
	renderer.Data = struct {
		Items interface{} `json:"items"`
		More  bool        `json:"more"`
	}{results[:count], l > opts.Size}
	return mage.Redirect{Status: http.StatusOK}
}

func (controller *Controller) HandleList(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	opts := &ListOptions{}
	opts, err := controller.BuildOptions(ctx, out, opts)
	if err != nil {
		return mage.Redirect{Status: http.StatusBadRequest}
	}
	results, err := controller.Manager.ListOf(ctx, *opts)
	if err != nil {
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	// output
	l := len(results)
	count := opts.Size
	if l < opts.Size {
		count = l
	}

	renderer := mage.JSONRenderer{}
	out.Renderer = &renderer
	renderer.Data = struct {
		Items interface{} `json:"items"`
		More  bool        `json:"more"`
	}{results[:count], l > opts.Size}
	return mage.Redirect{Status: http.StatusOK}
}

func (controller *Controller) BuildOptions(ctx context.Context, out *mage.ResponseOutput, opts *ListOptions) (*ListOptions, error) {
	// build paging
	opts.Size = 20
	opts.Page = 0

	ins := mage.InputsFromContext(ctx)
	if pin, ok := ins["page"]; ok {
		if num, err := strconv.Atoi(pin.Value()); err == nil {
			opts.Page = num
		} else {
			msg := fmt.Sprintf("invalid page value : %s. page must be an integer", pin)
			return nil, errors.New(msg)
		}
	}

	if sin, ok := ins["results"]; ok {
		if num, err := strconv.Atoi(sin.Value()); err == nil {
			opts.Size = num
		} else {
			msg := fmt.Sprintf("invalid result size value : %s. results must be an integer", sin)
			return nil, errors.New(msg)
		}
	}

	// order is not mandatory
	if oin, ok := ins["order"]; ok {
		oins := oin.Value()
		// descendig has the format "-fieldname"
		opts.Descending = oins[:1] == "-"
		if opts.Descending {
			opts.Order = oins[1:]
		} else {
			opts.Order = oins
		}
	}

	return opts, nil
}

// handles a post request, ensuring the creation of the resourse
func (controller *Controller) HandlePost(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	renderer := mage.JSONRenderer{}
	out.Renderer = &renderer

	resource, err := controller.Manager.NewResource(ctx)
	if err != nil {
		if _, ok := err.(validators.FieldError); ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	errs := validators.Errors{}
	// get the content data
	ins := mage.InputsFromContext(ctx)
	j, ok := ins[mage.KeyRequestJSON]
	if !ok {
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	err = json.Unmarshal([]byte(j.Value()), resource)
	if err != nil {
		msg := fmt.Sprintf("bad json: %s", err.Error())
		errs.AddError("", errors.New(msg))
		log.Errorf(ctx, msg)
		renderer.Data = errs
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	if err = resource.Create(ctx); err != nil {
		if _, ok := err.(PermissionError); ok {
			return mage.Redirect{Status: http.StatusForbidden}
		}
		errs.AddFieldError(err.(validators.FieldError))
	}

	// check for client input erros
	if errs.HasErrors() {
		log.Errorf(ctx, "invalid request")
		renderer.Data = errs
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	if err = controller.Manager.Save(ctx, resource); err != nil {
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	renderer.Data = resource
	return mage.Redirect{Status: http.StatusOK}
}

// Handles put requests, ensuring the update of the requested resource
func (controller *Controller) HandlePut(ctx context.Context, key string, out *mage.ResponseOutput) mage.Redirect {
	renderer := mage.JSONRenderer{}
	out.Renderer = &renderer

	ins := mage.InputsFromContext(ctx)
	j, ok := ins[mage.KeyRequestJSON]
	if !ok {
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	resource, err := controller.Manager.FromId(ctx, key)
	if err != nil {
		if _, ok := err.(validators.FieldError); ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	errs := validators.Errors{}
	jresource, err := controller.Manager.NewResource(ctx)
	if err != nil {
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	err = json.Unmarshal([]byte(j.Value()), &jresource)
	if err != nil {
		log.Errorf(ctx, "malformed json: %s", err.Error())
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	if err = resource.Update(ctx, jresource); err != nil {
		errs.AddFieldError(err.(validators.FieldError))
	}

	if errs.HasErrors() {
		log.Errorf(ctx, "invalid request")
		renderer.Data = errs
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	if err = controller.Manager.Save(ctx, resource); err != nil {
		if _, ok := err.(PermissionError); ok {
			return mage.Redirect{Status: http.StatusForbidden}
		}
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	renderer.Data = resource
	return mage.Redirect{Status: http.StatusOK}
}

// Handles delete requests
func (controller *Controller) HandleDelete(ctx context.Context, key string, out *mage.ResponseOutput) mage.Redirect {
	renderer := mage.JSONRenderer{}
	out.Renderer = &renderer

	resource, err := controller.Manager.NewResource(ctx)
	if err != nil {
		if _, ok := err.(validators.FieldError); ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	if err = controller.Manager.Delete(ctx, resource); err != nil {
		if _, ok := err.(PermissionError); ok {
			return mage.Redirect{Status: http.StatusForbidden}
		}
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	return mage.Redirect{Status: http.StatusOK}
}

func (controller *Controller) OnDestroy(ctx context.Context) {

}
