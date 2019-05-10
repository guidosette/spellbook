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
	Save(ctx context.Context, resource Resource) error
	Delete(ctx context.Context, resource Resource) error
}

type Resource interface {
	Id() string
	Create(ctx context.Context) error
	// todo: other should be a general serializable object, like a bundle
	Update(ctx context.Context, other Resource) error
}

type RestController struct {
	Key string
	RestHandler
}

func NewBaseRestController() *RestController {
	return NewRestController(BaseRestHandler{})
}

func NewRestController(handler RestHandler) *RestController {
	return &RestController{RestHandler: handler}
}

func (controller *RestController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {

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
