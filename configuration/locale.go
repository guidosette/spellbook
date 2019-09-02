package configuration

import (
	"context"
	"decodica.com/spellbook"
	"encoding/json"
	"golang.org/x/text/language"
)

type Locale language.Tag

func (locale *Locale) MarshalJSON() ([]byte, error) {
	tag := language.Tag(*locale)
	return json.Marshal(tag.String())
}

func (locale *Locale) Id() string {
	return language.Tag(*locale).String()
}

func (locale *Locale) FromRepresentation(rtype spellbook.RepresentationType, data []byte) error {
	return spellbook.NewUnsupportedError()
}

func (locale *Locale) ToRepresentation(rtype spellbook.RepresentationType) ([]byte, error) {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Marshal(locale)
	}
	return nil, spellbook.NewUnsupportedError()
}

func NewLocaleController() *spellbook.RestController {
	man := localeManager{}
	return spellbook.NewRestController(spellbook.BaseRestHandler{Manager: man})
}

type localeManager struct{}

func (manager localeManager) NewResource(ctx context.Context) (spellbook.Resource, error) {
	return nil, spellbook.NewUnsupportedError()
}

func (manager localeManager) FromId(ctx context.Context, id string) (spellbook.Resource, error) {
	return nil, spellbook.NewUnsupportedError()
}

func (manager localeManager) ListOf(ctx context.Context, opts spellbook.ListOptions) ([]spellbook.Resource, error) {

	ws := spellbook.Application()

	langs := ws.Options().Languages

	from := opts.Page * opts.Size
	if from > len(langs) {
		return make([]spellbook.Resource, 0), nil
	}

	to := from + opts.Size
	if to > len(langs) {
		to = len(langs)
	}

	items := langs[from:to]
	resources := make([]spellbook.Resource, len(items))

	for i := range items {
		locale := Locale(items[i])
		resources[i] = &locale
	}

	return resources, nil
}

func (manager localeManager) ListOfProperties(ctx context.Context, opts spellbook.ListOptions) ([]string, error) {
	return nil, spellbook.NewUnsupportedError()
}

func (manager localeManager) Create(ctx context.Context, res spellbook.Resource, bundle []byte) error {
	return spellbook.NewUnsupportedError()
}

func (manager localeManager) Update(ctx context.Context, res spellbook.Resource, bundle []byte) error {
	return spellbook.NewUnsupportedError()
}

func (manager localeManager) Delete(ctx context.Context, res spellbook.Resource) error {
	return spellbook.NewUnsupportedError()
}
