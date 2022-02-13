package spellbook

import (
	"context"
	"encoding/json"
)

func (specialcode *SpecialCode) MarshalJSON() ([]byte, error) {
	staticPageCode := SpecialCode(*specialcode)
	return json.Marshal(staticPageCode)
}

func (specialcode *SpecialCode) Id() string {
	return string(*specialcode)
}

func (specialcode *SpecialCode) FromRepresentation(rtype RepresentationType, data []byte) error {
	return NewUnsupportedError()
}

func (specialcode *SpecialCode) ToRepresentation(rtype RepresentationType) ([]byte, error) {
	switch rtype {
	case RepresentationTypeJSON:
		return json.Marshal(specialcode)
	}
	return nil, NewUnsupportedError()
}

func NewSpecialCodeController() *RestController {
	man := SpecialCodeManager{}
	return NewRestController(BaseRestHandler{Manager: man})
}

type SpecialCodeManager struct{}

func (manager SpecialCodeManager) NewResource(ctx context.Context) (Resource, error) {
	return nil, NewUnsupportedError()
}

func (manager SpecialCodeManager) FromId(ctx context.Context, id string) (Resource, error) {
	return nil, NewUnsupportedError()
}

func (manager SpecialCodeManager) ListOf(ctx context.Context, opts ListOptions) ([]Resource, error) {

	ws := Application()

	specialCodes := ws.Options().SpecialCodes

	from := opts.Page * opts.Size
	if from > len(specialCodes) {
		return make([]Resource, 0), nil
	}

	to := from + opts.Size
	if to > len(specialCodes) {
		to = len(specialCodes)
	}

	items := specialCodes[from:to]
	resources := make([]Resource, len(items))

	for i := range items {
		specialCode := SpecialCode(items[i])
		resources[i] = Resource(&specialCode)
	}

	return resources, nil
}

func (manager SpecialCodeManager) ListOfProperties(ctx context.Context, opts ListOptions) ([]string, error) {
	return nil, NewUnsupportedError()
}

func (manager SpecialCodeManager) Create(ctx context.Context, res Resource, bundle []byte) error {
	return NewUnsupportedError()
}

func (manager SpecialCodeManager) Update(ctx context.Context, res Resource, bundle []byte) error {
	return NewUnsupportedError()
}

func (manager SpecialCodeManager) Delete(ctx context.Context, res Resource) error {
	return NewUnsupportedError()
}
