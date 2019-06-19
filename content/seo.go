package content

import (
	"distudio.com/mage/model"
	"distudio.com/page"
	"encoding/json"
)

type Seo struct {
	model.Model `json:"-"`
	Title       string
	MetaDesc    string
	Url         string
	Code        page.StaticPageCode
}

func (seo *Seo) UnmarshalJSON(data []byte) error {
	alias := struct {
		Title    string              `json:"title"`
		MetaDesc string              `json:"metadesc"`
		Url      string              `json:"url"`
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

	return nil
}

func (seo *Seo) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Title    string              `json:"title"`
		MetaDesc string              `json:"metadesc"`
		Url      string              `json:"url"`
		Code     page.StaticPageCode `json:"code"`
	}

	return json.Marshal(&struct {
		Id int64 `json:"id"`
		Alias
	}{
		seo.IntID(),
		Alias{
			Title:    seo.Title,
			MetaDesc: seo.MetaDesc,
			Url:      seo.Url,
			Code:     seo.Code,
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
