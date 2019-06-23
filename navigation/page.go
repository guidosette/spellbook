package navigation

import (
	"distudio.com/mage/model"
	"distudio.com/page"
	"encoding/json"
)

const rootUrl = ""

type Page struct {
	model.Model `json:"-"`
	Label string
	Title       string
	MetaDesc    string
	Url         string
	Order int
	IsRoot      bool
	Code        page.StaticPageCode
	Locale      string
}

func (p Page) LocalizedUrl() string {
	return "/" + p.Locale + "/" + p.Url
}

func (p *Page) UnmarshalJSON(data []byte) error {
	alias := struct {
		Label string `json:"label"`
		Order int `json:"order"`
		Title    string              `json:"title"`
		MetaDesc string              `json:"metadesc"`
		Url      string              `json:"url"`
		Locale   string              `json:"locale"`
		Code     page.StaticPageCode `json:"code"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	p.Label = alias.Label
	p.Order = alias.Order
	p.Title = alias.Title
	p.MetaDesc = alias.MetaDesc
	p.Url = alias.Url
	p.Code = alias.Code
	p.Locale = alias.Locale

	return nil
}

func (p *Page) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Label string `json:"label"`
		Title    string              `json:"title"`
		MetaDesc string              `json:"metadesc"`
		Url      string              `json:"url"`
		Order int `json:"order"`
		Locale   string              `json:"locale"`
		Code     page.StaticPageCode `json:"code"`
	}

	return json.Marshal(&struct {
		Id string `json:"id"`
		Alias
	}{
		p.StringID(),
		Alias{
			Label: p.Label,
			Title:    p.Title,
			MetaDesc: p.MetaDesc,
			Url:      p.Url,
			Order: p.Order,
			Code:     p.Code,
			Locale:   p.Locale,
		},
	})
}

func (p *Page) Id() string {
	return p.StringID()
}

func (p *Page) FromRepresentation(rtype page.RepresentationType, data []byte) error {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Unmarshal(data, p)
	}
	return page.NewUnsupportedError()
}

func (p *Page) ToRepresentation(rtype page.RepresentationType) ([]byte, error) {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Marshal(p)
	}
	return nil, page.NewUnsupportedError()
}
