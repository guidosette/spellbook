package subscription

import (
	"decodica.com/flamel/model"
	"decodica.com/spellbook"
	"encoding/json"
	"time"
)

type Subscription struct {
	model.Model  `json:"-"`
	Email        string    `model:"search,atom"`
	Country      string    `model:"search"`
	FirstName    string    `model:"search"`
	LastName     string    `model:"search"`
	Organization string    `model:"search"`
	Created      time.Time `model:"search"`
	Updated      time.Time `model:"search"`
}

func (subscription *Subscription) UnmarshalJSON(data []byte) error {

	alias := struct {
		Email        string    `json:"email"`
		Country      string    `json:"country"`
		FirstName    string    `json:"firstName"`
		LastName     string    `json:"lastName"`
		Organization string    `json:"organization"`
		Created      time.Time `json:"created"`
		Updated      time.Time `json:"updated"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	subscription.Email = alias.Email
	subscription.Country = alias.Country
	subscription.FirstName = alias.FirstName
	subscription.LastName = alias.LastName
	subscription.Organization = alias.Organization
	subscription.Created = alias.Created
	subscription.Updated = alias.Updated

	return nil
}

func (subscription *Subscription) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Email        string    `json:"email"`
		Country      string    `json:"country"`
		FirstName    string    `json:"firstName"`
		LastName     string    `json:"lastName"`
		Organization string    `json:"organization"`
		Created      time.Time `json:"created"`
		Updated      time.Time `json:"updated"`
		Key          string    `json:"key"`
	}

	return json.Marshal(&struct {
		Alias
	}{
		Alias{
			Email:        subscription.Email,
			Country:      subscription.Country,
			FirstName:    subscription.FirstName,
			LastName:     subscription.LastName,
			Organization: subscription.Organization,
			Created:      subscription.Created,
			Updated:      subscription.Updated,
			Key:          subscription.EncodedKey(),
		},
	})
}

func (subscription *Subscription) Id() string {
	return subscription.StringID()
}

func (subscription *Subscription) FromRepresentation(rtype spellbook.RepresentationType, data []byte) error {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Unmarshal(data, subscription)
	}
	return spellbook.NewUnsupportedError()
}

func (subscription *Subscription) ToRepresentation(rtype spellbook.RepresentationType) ([]byte, error) {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Marshal(subscription)
	}
	return nil, spellbook.NewUnsupportedError()
}
