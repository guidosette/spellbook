package content

import (
	"distudio.com/mage/model"
	"distudio.com/page"
	"encoding/json"
	"google.golang.org/appengine"
	"time"
)

type Place struct {
	model.Model `json:"-"`
	Name        string             `model:"search";json:"name"`
	Address     string             `model:"search";json:"address"`
	City        string             `model:"search";json:"city"`
	PostalCode  string             `model:"search";json:"postalCode"`
	Country     string             `model:"search";json:"country"`
	Phone       string             `json:"phone";model:"noindex"`
	Description string             `json:"description";model:"noindex"`
	Position    appengine.GeoPoint `model:"search"`
	Website     string             `json:"website";model:"noindex"`
	Created     time.Time          `json:"created"`
	Updated     time.Time          `json:"updated"`
}

func (place *Place) UnmarshalJSON(data []byte) error {

	alias := struct {
		Name        string    `json:"name"`
		Address     string    `json:"address"`
		City        string    `json:"city"`
		PostalCode  string    `json:"postalCode"`
		Country     string    `json:"country"`
		Phone       string    `json:"phone"`
		Description string    `json:"description"`
		Lat         float64   `json:"lat"`
		Lng         float64   `json:"lng"`
		Website     string    `json:"website"`
		Created     time.Time `json:"created"`
		Updated     time.Time `json:"updated"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	place.Name = alias.Name
	place.Address = alias.Address
	place.City = alias.City
	place.PostalCode = alias.PostalCode
	place.Country = alias.Country
	place.Phone = alias.Phone
	place.Description = alias.Description
	place.Website = alias.Website
	place.Created = alias.Created
	place.Updated = alias.Updated
	place.Position = appengine.GeoPoint{Lat: alias.Lat, Lng: alias.Lng}

	return nil
}

func (place *Place) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Name        string    `json:"name"`
		Address     string    `json:"address"`
		City        string    `json:"city"`
		PostalCode  string    `json:"postalCode"`
		Country     string    `json:"country"`
		Phone       string    `json:"phone"`
		Description string    `json:"description"`
		Lat         float64   `json:"lat"`
		Lng         float64   `json:"lng"`
		Website     string    `json:"website"`
		Created     time.Time `json:"created"`
		Updated     time.Time `json:"updated"`
		Id          int64     `json:"id"`
	}

	return json.Marshal(&struct {
		Alias
	}{
		Alias{
			Name:        place.Name,
			Address:     place.Address,
			City:        place.City,
			PostalCode:  place.PostalCode,
			Country:     place.Country,
			Phone:       place.Phone,
			Description: place.Description,
			Lat:         place.Position.Lat,
			Lng:         place.Position.Lng,
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
