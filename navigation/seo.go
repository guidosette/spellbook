package navigation

import (
	"distudio.com/mage/model"
	"distudio.com/page"
	"encoding/json"
)

const rootUrl = ""

type Seo struct {
	model.Model `json:"-"`
	Title       string
	MetaDesc    string
	Url         string
	IsRoot bool
	Code        page.StaticPageCode
	Locale      string
}

func (seo Seo) LocalizedUrl() string {
	return "/" + seo.Locale + "/" + seo.Url
}

func PageId(locale string, url string) string {
	return locale + "-" + url
}

func (seo *Seo) UnmarshalJSON(data []byte) error {
	alias := struct {
		Title    string              `json:"title"`
		MetaDesc string              `json:"metadesc"`
		Url      string              `json:"url"`
		Locale   string `json:"locale"`
		Code     page.StaticPageCode `json:"code"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	seo.Title = alias.Title
	seo.MetaDesc = alias.MetaDesc
	seo.Url = alias.Url
	seo.Code = alias.Code
	seo.Locale = alias.Locale

	return nil
}

func (seo *Seo) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Title    string              `json:"title"`
		MetaDesc string              `json:"metadesc"`
		Url      string              `json:"url"`
		Locale   string              `json:"locale"`
		Code     page.StaticPageCode `json:"code"`
	}

	return json.Marshal(&struct {
		Id string `json:"id"`
		Alias
	}{
		seo.StringID(),
		Alias{
			Title:    seo.Title,
			MetaDesc: seo.MetaDesc,
			Url:      seo.Url,
			Code:     seo.Code,
			Locale:   seo.Locale,
		},
	})
}

func (seo *Seo) Id() string {
	return seo.StringID()
}

func (seo *Seo) FromRepresentation(rtype page.RepresentationType, data []byte) error {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Unmarshal(data, seo)
	}
	return page.NewUnsupportedError()
}

func (seo *Seo) ToRepresentation(rtype page.RepresentationType) ([]byte, error) {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Marshal(seo)
	}
	return nil, page.NewUnsupportedError()
}
