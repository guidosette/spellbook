package newsletter

import (
	"distudio.com/mage/model"
	"distudio.com/page"
	"encoding/json"
)

type Newsletter struct {
	model.Model `json:"-"`
	Email string `json:"email"`
}

func (newsletter *Newsletter) Id() string {
	return newsletter.StringID()
}

func (newsletter *Newsletter) FromRepresentation(rtype page.RepresentationType, data []byte) error {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Unmarshal(data, newsletter)
	}
	return page.NewUnsupportedError()
}

func (newsletter *Newsletter) ToRepresentation(rtype page.RepresentationType) ([]byte, error) {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Marshal(newsletter)
	}
	return nil, page.NewUnsupportedError()
}
