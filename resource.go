package page

import (
	"context"
	"distudio.com/mage"
	"google.golang.org/appengine/log"
	"net/http"
)

type ListOptions struct {
	Size       int
	Page       int
	Order      string // field
	Descending bool   // if -Order = desc
	Property   string
	Filters    []Filter // example url: &filter=Locale=it^Category=services
}

type Filter struct {
	Field string
	Value string
}

type ListResponse struct {
	Items interface{} `json:"items"`
	More  bool        `json:"more"`
}

type Manager interface {
	NewResource(ctx context.Context) (Resource, error)
	FromId(ctx context.Context, id string) (Resource, error)
	ListOf(ctx context.Context, opts ListOptions) ([]Resource, error)
	ListOfProperties(ctx context.Context, opts ListOptions) ([]string, error)
	Create(ctx context.Context, resource Resource, bundle []byte) error
	Update(ctx context.Context, resource Resource, bundle []byte) error
	Delete(ctx context.Context, resource Resource) error
}

type RepresentationType int

const (
	RepresentationTypeJSON = iota
	RepresentationTypeUrlencoded
)

type Resource interface {
	Id() string
	ToRepresentation(rtype RepresentationType) ([]byte, error)
	FromRepresentation(rtype RepresentationType, data []byte) error
}

/**
* Base rest controller
 */

type RestController struct {
	Key     string
	Private bool
	RestHandler
	extenders map[string][]Extender
}

func NewBaseRestController() *RestController {
	return NewRestController(BaseRestHandler{})
}

func NewRestController(handler RestHandler) *RestController {
	return &RestController{RestHandler: handler}
}

func (controller *RestController) AddExtender(hook string, extender Extender) {
	e, ok := controller.extenders[hook]
	if !ok {
		e = make([]Extender, 0)
	}
	e = append(e, extender)
	controller.extenders[hook] = e
}

func (controller *RestController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {

	u := IdentityFromContext(ctx)

	if controller.Private && u == nil {
		return mage.Redirect{Status: http.StatusUnauthorized}
	}

	ins := mage.InputsFromContext(ctx)

	method := ins[mage.KeyRequestMethod].Value()
	hasKey := controller.Key != ""
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
		return controller.HandleGet(ctx, controller.Key, out)
	case http.MethodPut:
		if !hasKey {
			log.Errorf(ctx, "no item was specify for put method")
			return mage.Redirect{Status: http.StatusBadRequest}
		}
		return controller.HandlePut(ctx, controller.Key, out)
	case http.MethodDelete:
		if !hasKey {
			log.Errorf(ctx, "no item was specify for delete method")
			return mage.Redirect{Status: http.StatusBadRequest}
		}
		return controller.HandleDelete(ctx, controller.Key, out)
	}

	return mage.Redirect{Status: http.StatusNotImplemented}
}

func (controller *RestController) OnDestroy(ctx context.Context) {}
