package page

import (
	"context"
	"encoding/json"
)

func (staticpage *StaticPageCode) MarshalJSON() ([]byte, error) {
	staticPageCode := StaticPageCode(*staticpage)
	return json.Marshal(staticPageCode)
}

func (staticpage *StaticPageCode) Id() string {
	return string(*staticpage)
}

func (staticpage *StaticPageCode) FromRepresentation(rtype RepresentationType, data []byte) error {
	return NewUnsupportedError()
}

func (staticpage *StaticPageCode) ToRepresentation(rtype RepresentationType) ([]byte, error) {
	switch rtype {
	case RepresentationTypeJSON:
		return json.Marshal(staticpage)
	}
	return nil, NewUnsupportedError()
}

func NewStaticPageCodeController() *RestController {
	man := staticPageCodeManager{}
	return NewRestController(BaseRestHandler{Manager: man})
}

type staticPageCodeManager struct{}

func (manager staticPageCodeManager) NewResource(ctx context.Context) (Resource, error) {
	return nil, NewUnsupportedError()
}

func (manager staticPageCodeManager) FromId(ctx context.Context, id string) (Resource, error) {
	return nil, NewUnsupportedError()
}

func (manager staticPageCodeManager) ListOf(ctx context.Context, opts ListOptions) ([]Resource, error) {

	ws := Application()

	staticPageCodes := ws.Options().StaticPages

	from := opts.Page * opts.Size
	if from > len(staticPageCodes) {
		return make([]Resource, 0), nil
	}

	to := from + opts.Size
	if to > len(staticPageCodes) {
		to = len(staticPageCodes)
	}

	items := staticPageCodes[from:to]
	resources := make([]Resource, len(items))

	for i := range items {
		staticPage := StaticPageCode(items[i])
		resources[i] = Resource(&staticPage)
	}

	return resources, nil
}

func (manager staticPageCodeManager) ListOfProperties(ctx context.Context, opts ListOptions) ([]string, error) {
	return nil, NewUnsupportedError()
}

func (manager staticPageCodeManager) Create(ctx context.Context, res Resource, bundle []byte) error {
	return NewUnsupportedError()
}

func (manager staticPageCodeManager) Update(ctx context.Context, res Resource, bundle []byte) error {
	return NewUnsupportedError()
}

func (manager staticPageCodeManager) Delete(ctx context.Context, res Resource) error {
	return NewUnsupportedError()
}
