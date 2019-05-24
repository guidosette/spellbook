package content

import (
	"cloud.google.com/go/datastore"
	"distudio.com/mage/model"
	"distudio.com/page"
	"encoding/json"
	"time"
)

type Place struct {
	model.Model `json:"-"`
	Address     string             `json:"address"`
	Phone       string             `json:"phone";model:"noindex"`
	Description string             `json:"description";model:"noindex"`
	Position    datastore.GeoPoint `json:"position"`
	Website     *Attachment        `model:"-"`
	Created     time.Time          `json:"created"`
	Updated     time.Time          `json:"updated"`
}

func (place *Place) UnmarshalJSON(data []byte) error {

	alias := struct {
		Address     string             `json:"address"`
		Phone       string             `json:"phone"`
		Description string             `json:"description"`
		Position    datastore.GeoPoint `json:"position"`
		Website     *Attachment        `json:"website"`
		Created     time.Time          `json:"created"`
		Updated     time.Time          `json:"updated"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	place.Address = alias.Address
	place.Phone = alias.Phone
	place.Description = alias.Description
	place.Position = alias.Position
	place.Website = alias.Website
	place.Created = alias.Created
	place.Updated = alias.Updated

	return nil
}

func (place *Place) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Address     string             `json:"address"`
		Phone       string             `json:"phone"`
		Description string             `json:"description"`
		Position    datastore.GeoPoint `json:"position"`
		Website     *Attachment        `json:"website"`
		Created     time.Time          `json:"created"`
		Updated     time.Time          `json:"updated"`
		Id          int64              `json:"id"`
	}

	return json.Marshal(&struct {
		Alias
	}{
		Alias{
			Address:     place.Address,
			Phone:       place.Phone,
			Description: place.Description,
			Position:    place.Position,
			Website:     place.Website,
			Created:     place.Created,
			Updated:     place.Updated,
			Id:          place.IntID(),
		},
	})
}

func (place *Place) Id() string {
	return place.StringID()
}

func (place *Place) FromRepresentation(rtype page.RepresentationType, data []byte) error {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Unmarshal(data, place)
	}
	return page.NewUnsupportedError()
}

func (place *Place) ToRepresentation(rtype page.RepresentationType) ([]byte, error) {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Marshal(place)
	}
	return nil, page.NewUnsupportedError()
}
