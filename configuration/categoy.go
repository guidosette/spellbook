package configuration

import (
	"context"
	"distudio.com/page"
)

type Category string


func (category Category) Id() string {
	return ""
}

func (category Category) Create(ctx context.Context) error {
	return page.NewUnsupportedError()
}

func (category Category) Update(ctx context.Context, other page.Resource) error {
	return page.NewUnsupportedError()
}

func NewCategoryController() *page.RestController {
	man := categoryManager{}
	return page.NewRestController(page.BaseRestHandler{Manager: man})
}

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
		resources[i] = page.Resource(category)
	}

	return resources, nil
}

func (manager categoryManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	return nil, page.NewUnsupportedError()
}

func (manager categoryManager) Save(ctx context.Context, res page.Resource) error {
	return page.NewUnsupportedError()
}

func (manager categoryManager) Delete(ctx context.Context, res page.Resource) error {
	return page.NewUnsupportedError()
}
