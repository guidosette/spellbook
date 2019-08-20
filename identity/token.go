package identity

import (
	"decodica.com/spellbook"
	"encoding/json"
)

type Token struct {
	Value    string
	Username string
	Password string
}

func (token *Token) UnmarshalJSON(data []byte) error {
	alias := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

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

func (token *Token) FromRepresentation(rtype spellbook.RepresentationType, data []byte) error {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Unmarshal(data, token)
	}
	return spellbook.NewUnsupportedError()
}

func (token *Token) ToRepresentation(rtype spellbook.RepresentationType) ([]byte, error) {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Marshal(token)
	}
	return nil, spellbook.NewUnsupportedError()
}
