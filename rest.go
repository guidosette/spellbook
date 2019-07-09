package page

import (
	"context"
	"distudio.com/mage"
	"errors"
	"fmt"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"strconv"
	"strings"
)

type ReadHandler interface {
	HandleGet(context context.Context, key string, out *mage.ResponseOutput) mage.Redirect
}

type WriteHandler interface {
	HandlePost(context context.Context, out *mage.ResponseOutput) mage.Redirect
	HandlePut(context context.Context, key string, out *mage.ResponseOutput) mage.Redirect
	HandleDelete(context context.Context, key string, out *mage.ResponseOutput) mage.Redirect
}

type ListHandler interface {
	HandleList(context context.Context, out *mage.ResponseOutput) mage.Redirect
	HandlePropertyValues(context context.Context, out *mage.ResponseOutput, property string) mage.Redirect
}

type RestHandler interface {
	ReadHandler
	WriteHandler
	ListHandler
}

type BaseRestHandler struct {
	Manager Manager
}

// Builds the paging options, ordering and standard inputs of a given request
func (handler BaseRestHandler) buildOptions(ctx context.Context, out *mage.ResponseOutput, opts *ListOptions) (*ListOptions, error) {
	// build paging
	opts.Size = 20
	opts.Page = 0

	ins := mage.InputsFromContext(ctx)
	if pin, ok := ins["page"]; ok {
		if num, err := strconv.Atoi(pin.Value()); err == nil {
			if num > 0 {
				opts.Page = num
			}
		} else {
			msg := fmt.Sprintf("invalid page value : %v. page must be an integer", pin)
			return nil, errors.New(msg)
		}
	}

	if sin, ok := ins["results"]; ok {
		if num, err := strconv.Atoi(sin.Value()); err == nil {
			if num > 0 {
				opts.Size = num
			}
		} else {
			msg := fmt.Sprintf("invalid result size value : %v. results must be an integer", sin)
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

	// filter is not mandatory
	if fin, ok := ins["filter"]; ok {
		finv := fin.Value()
		filters := strings.Split(finv, "^")
		opts.Filters = make([]Filter, len(filters), cap(filters))
		for i, filter := range filters {
			farray := strings.Split(filter, "=")
			if len(farray) > 1 {
				opts.Filters[i] = Filter{farray[0], farray[1]}
			}
		}

	}

	return opts, nil
}

// REST Method handlers
func (handler BaseRestHandler) HandleGet(ctx context.Context, key string, out *mage.ResponseOutput) mage.Redirect {
	renderer := mage.JSONRenderer{}
	out.Renderer = &renderer

	resource, err := handler.Manager.FromId(ctx, key)
	if err != nil {
		return handler.ErrorToStatus(ctx, err, out)
	}

	renderer.Data = resource
	return mage.Redirect{Status: http.StatusOK}
}

// Called on GET requests.
// This handler is called when the available values of one property of a resource are requested
// Returns a list of the values that the requested property can assume
func (handler BaseRestHandler) HandlePropertyValues(ctx context.Context, out *mage.ResponseOutput, prop string) mage.Redirect {
	opts := &ListOptions{}
	opts.Property = prop
	opts, err := handler.buildOptions(ctx, out, opts)
	if err != nil {
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	results, err := handler.Manager.ListOfProperties(ctx, *opts)
	if err != nil {
		return handler.ErrorToStatus(ctx, err, out)
	}

	// output
	l := len(results)
	count := opts.Size
	if l < opts.Size {
		count = l
	}

	renderer := mage.JSONRenderer{}
	renderer.Data = ListResponse{results[:count], l > opts.Size}

	out.Renderer = &renderer

	return mage.Redirect{Status: http.StatusOK}
}

// Called on GET requests
// This handler is called when a list of resources is requested.
// Returns a paged result
func (handler BaseRestHandler) HandleList(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	opts := &ListOptions{}
	opts, err := handler.buildOptions(ctx, out, opts)
	if err != nil {
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	results, err := handler.Manager.ListOf(ctx, *opts)
	if err != nil {
		return handler.ErrorToStatus(ctx, err, out)
	}

	// output
	l := len(results)
	count := opts.Size
	if l < opts.Size {
		count = l
	}

	renderer := mage.JSONRenderer{}
	renderer.Data = ListResponse{results[:count], l > opts.Size}

	out.Renderer = &renderer

	return mage.Redirect{Status: http.StatusOK}
}

// handles a POST request, ensuring the creation of the resource.
func (handler BaseRestHandler) HandlePost(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	renderer := mage.JSONRenderer{}
	out.Renderer = &renderer

	resource, err := handler.Manager.NewResource(ctx)
	if err != nil {
		return handler.ErrorToStatus(ctx, err, out)
	}

	errs := Errors{}
	// get the content data
	ins := mage.InputsFromContext(ctx)
	j, ok := ins[mage.KeyRequestJSON]
	if !ok {
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	err = resource.FromRepresentation(RepresentationTypeJSON, []byte(j.Value()))
	if err != nil {
		msg := fmt.Sprintf("bad json: %s", err.Error())
		errs.AddError("", errors.New(msg))
		log.Errorf(ctx, msg)
		renderer.Data = errs
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	if err = handler.Manager.Create(ctx, resource, []byte(j.Value())); err != nil {
		return handler.ErrorToStatus(ctx, err, out)
	}

	renderer.Data = resource
	return mage.Redirect{Status: http.StatusCreated}
}

// Handles put requests, ensuring the update of the requested resource
func (handler BaseRestHandler) HandlePut(ctx context.Context, key string, out *mage.ResponseOutput) mage.Redirect {
	renderer := mage.JSONRenderer{}
	out.Renderer = &renderer

	ins := mage.InputsFromContext(ctx)
	j, ok := ins[mage.KeyRequestJSON]
	if !ok {
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	resource, err := handler.Manager.FromId(ctx, key)
	if err != nil {
		return handler.ErrorToStatus(ctx, err, out)
	}

	if err = handler.Manager.Update(ctx, resource, []byte(j.Value())); err != nil {
		return handler.ErrorToStatus(ctx, err, out)
	}

	renderer.Data = resource
	return mage.Redirect{Status: http.StatusOK}
}

// Handles DELETE requests over a Resource type
func (handler BaseRestHandler) HandleDelete(ctx context.Context, key string, out *mage.ResponseOutput) mage.Redirect {
	renderer := mage.JSONRenderer{}
	out.Renderer = &renderer

	resource, err := handler.Manager.FromId(ctx, key)
	if err != nil {
		return handler.ErrorToStatus(ctx, err, out)
	}

	if err = handler.Manager.Delete(ctx, resource); err != nil {
		return handler.ErrorToStatus(ctx, err, out)
	}
	return mage.Redirect{Status: http.StatusOK}
}

// Converts an error to its equivalent HTTP representation
func (handler BaseRestHandler) ErrorToStatus(ctx context.Context, err error, out *mage.ResponseOutput) mage.Redirect {
	log.Errorf(ctx, "%s", err.Error())
	switch err.(type) {
	case UnsupportedError:
		return mage.Redirect{Status: http.StatusMethodNotAllowed}
	case FieldError:
		renderer := mage.JSONRenderer{}
		renderer.Data = struct {
			Field string
			Error string
		}{
			err.(FieldError).field,
			err.(FieldError).error.Error(),
		}
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusBadRequest}
	case PermissionError:
		renderer := mage.JSONRenderer{}
		renderer.Data = struct {
			Error string
		}{
			err.(PermissionError).Error(),
		}
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusForbidden}
	default:
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}
		return mage.Redirect{Status: http.StatusInternalServerError}
	}
}
