package content

import (
	"context"
	"distudio.com/page"
	"encoding/json"
)

type Category page.SupportedCategory

func (category *Category) UnmarshalJSON(data []byte) error {

	alias := struct {
		Name  string           `json:"name"`
		Label string           `json:"label"`
		Type  page.ContentType `json:"type"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	category.Type = alias.Type
	category.Name = alias.Name
	category.Label = alias.Label

	return nil
}

func (category *Category) MarshalJSON() ([]byte, error) {
	alias := struct {
		Name                   string                        `json:"name"`
		Label                  string                        `json:"label"`
		Type                   page.ContentType              `json:"type"`
		DefaultAttachmentGroup []page.DefaultAttachmentGroup `json:"defaultAttachmentGroups"`
	}{category.Name, category.Label, category.Type, category.DefaultAttachmentGroups}

	return json.Marshal(&alias)
}

/**
* Resource representation
 */

func (category *Category) Id() string {
	return category.Name
}

func (category *Category) FromRepresentation(rtype page.RepresentationType, data []byte) error {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Unmarshal(data, category)
	}
	return page.NewUnsupportedError()
}

func (category *Category) ToRepresentation(rtype page.RepresentationType) ([]byte, error) {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Marshal(category)
	}
	return nil, page.NewUnsupportedError()
}

func NewCategoryController() *page.RestController {
	man := categoryManager{}
	return page.NewRestController(page.BaseRestHandler{Manager: man})
}

/*
* Category manager
 */

type categoryManager struct{}

func (manager categoryManager) NewResource(ctx context.Context) (page.Resource, error) {
	return nil, page.NewUnsupportedError()
}

func (manager categoryManager) FromId(ctx context.Context, id string) (page.Resource, error) {
	return nil, page.NewUnsupportedError()
}

func (manager categoryManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {

	ws := page.Application()

	categories := ws.Options().Categories

	from := opts.Page * opts.Size
	if from > len(categories) {
		return make([]page.Resource, 0), nil
	}

	to := from + opts.Size
	if to > len(categories) {
		to = len(categories)
	}

	items := categories[from:to]
	resources := make([]page.Resource, len(items))

	for i := range items {
		category := Category(items[i])
		resources[i] = page.Resource(&category)
	}

	return resources, nil
}

func (manager categoryManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	return nil, page.NewUnsupportedError()
}

func (manager categoryManager) Create(ctx context.Context, res page.Resource, bundle []byte) error {
	return page.NewUnsupportedError()
}

func (manager categoryManager) Update(ctx context.Context, res page.Resource, bundle []byte) error {
	return page.NewUnsupportedError()
}

func (manager categoryManager) Delete(ctx context.Context, res page.Resource) error {
	return page.NewUnsupportedError()
}
