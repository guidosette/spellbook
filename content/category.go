package content

import (
	"context"
	"decodica.com/spellbook"
	"encoding/json"
)

type Category spellbook.SupportedCategory

func (category *Category) UnmarshalJSON(data []byte) error {

	alias := struct {
		Name  string                `json:"name"`
		Label string                `json:"label"`
		Type  spellbook.ContentType `json:"type"`
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
		Name                   string                             `json:"name"`
		Label                  string                             `json:"label"`
		Type                   spellbook.ContentType              `json:"type"`
		DefaultAttachmentGroup []spellbook.DefaultAttachmentGroup `json:"defaultAttachmentGroups"`
	}{category.Name, category.Label, category.Type, category.DefaultAttachmentGroups}

	return json.Marshal(&alias)
}

/**
* Resource representation
 */

func (category *Category) Id() string {
	return category.Name
}

func (category *Category) FromRepresentation(rtype spellbook.RepresentationType, data []byte) error {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Unmarshal(data, category)
	}
	return spellbook.NewUnsupportedError()
}

func (category *Category) ToRepresentation(rtype spellbook.RepresentationType) ([]byte, error) {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Marshal(category)
	}
	return nil, spellbook.NewUnsupportedError()
}

func NewCategoryController() *spellbook.RestController {
	man := categoryManager{}
	return spellbook.NewRestController(spellbook.BaseRestHandler{Manager: man})
}

/*
* Category manager
 */

type categoryManager struct{}

func (manager categoryManager) NewResource(ctx context.Context) (spellbook.Resource, error) {
	return nil, spellbook.NewUnsupportedError()
}

func (manager categoryManager) FromId(ctx context.Context, id string) (spellbook.Resource, error) {
	return nil, spellbook.NewUnsupportedError()
}

func (manager categoryManager) ListOf(ctx context.Context, opts spellbook.ListOptions) ([]spellbook.Resource, error) {

	ws := spellbook.Application()

	categories := ws.Options().Categories

	from := opts.Page * opts.Size
	if from > len(categories) {
		return make([]spellbook.Resource, 0), nil
	}

	to := from + opts.Size
	if to > len(categories) {
		to = len(categories)
	}

	items := categories[from:to]
	resources := make([]spellbook.Resource, len(items))

	for i := range items {
		category := Category(items[i])
		resources[i] = spellbook.Resource(&category)
	}

	return resources, nil
}

func (manager categoryManager) ListOfProperties(ctx context.Context, opts spellbook.ListOptions) ([]string, error) {
	return nil, spellbook.NewUnsupportedError()
}

func (manager categoryManager) Create(ctx context.Context, res spellbook.Resource, bundle []byte) error {
	return spellbook.NewUnsupportedError()
}

func (manager categoryManager) Update(ctx context.Context, res spellbook.Resource, bundle []byte) error {
	return spellbook.NewUnsupportedError()
}

func (manager categoryManager) Delete(ctx context.Context, res spellbook.Resource) error {
	return spellbook.NewUnsupportedError()
}
