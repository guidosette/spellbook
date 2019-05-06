package identity

import (
	"context"
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
	return []byte(token.Value), nil
}

func (token Token) Id() string {
	return token.Value
}

func (token Token) Create(ctx context.Context) error {
	// checks the provided credentials. If correct creates a token, saves the user and returns the token
	nick := page.NewRawField("username", true, token.Username)
	if _, err := nick.Value(); err != nil {
		return page.NewFieldError("username", err)
	}

	password := page.NewRawField("password", true, token.Password)
	password.AddValidator(page.LenValidator{MinLen: 8})
	if _, err := password.Value(); err != nil {
		return page.NewFieldError("password", err)
	}

	return nil
}

func (token Token) Update(ctx context.Context, res page.Resource) error {
	return page.NewUnsupportedError()
}

