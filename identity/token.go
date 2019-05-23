package identity

import (
	"distudio.com/page"
	"encoding/json"
)

type Token struct {
	Value string
	Username string
	Password string
}

func (token *Token) UnmarshalJSON(data []byte) error {
	alias := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} {}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	token.Username = alias.Username
	token.Password = alias.Password
	return nil
}

func (token *Token) MarshalJSON() ([]byte, error) {
	return json.Marshal(token.Value)
}

/**
* Resource implementation
 */
func (token *Token) Id() string {
	return token.Value
}

func (token *Token) FromRepresentation(rtype page.RepresentationType, data []byte) error {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Unmarshal(data, token)
	}
	return page.NewUnsupportedError()
}

func (token *Token) ToRepresentation(rtype page.RepresentationType) ([]byte, error) {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Marshal(token)
	}
	return nil, page.NewUnsupportedError()
}

