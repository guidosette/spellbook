package configuration

import (
	"context"
	"distudio.com/page"
	"encoding/json"
	"golang.org/x/text/language"
)

type Locale language.Tag

func (locale Locale) MarshalJSON() ([]byte, error) {
	tag := language.Tag(locale)
	return json.Marshal(tag.String())
}

func (locale Locale) Id() string {
	return ""
}

func (locale Locale) Create(ctx context.Context) error {
	return page.NewUnsupportedError()
}

func (locale Locale) Update(ctx context.Context, other page.Resource) error {
	return page.NewUnsupportedError()
}


type LocaleManager struct{}

func (manager LocaleManager) NewResource(ctx context.Context) (page.Resource, error) {
	return nil, page.NewUnsupportedError()
}

func (manager LocaleManager) FromId(ctx context.Context, id string) (page.Resource, error) {
	return nil, page.NewUnsupportedError()
}

func (manager LocaleManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {

	ws := page.Application()

	langs := ws.Options().Languages

	from := opts.Page * opts.Size
	if from > len(langs) {
		return make([]page.Resource, 0), nil
	}

	to := from + opts.Size
	if to > len(langs) {
		to = len(langs)
	}


	items := langs[from:to]
	resources := make([]page.Resource, len(items))

	for i := range items {
		locale := Locale(items[i])
		resources[i] = page.Resource(locale)
	}

	return resources, nil
}

func (manager LocaleManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	return nil, page.NewUnsupportedError()
}

func (manager LocaleManager) Save(ctx context.Context, res page.Resource) error {
	return page.NewUnsupportedError()
}

func (manager LocaleManager) Delete(ctx context.Context, res page.Resource) error {
	return page.NewUnsupportedError()
}
