package identity

import (
	"context"
	"decodica.com/flamel/model"
	"decodica.com/spellbook"
	"fmt"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func NewTokenController() *spellbook.RestController {
	handler := spellbook.BaseRestHandler{}
	handler.Manager = tokenManager{}
	return spellbook.NewRestController(handler)
}

func NewTokenControllerWithKey(key string) *spellbook.RestController {
	handler := spellbook.BaseRestHandler{Manager: tokenManager{}}
	c := spellbook.NewRestController(handler)
	c.Key = key
	return c
}

type tokenManager struct{}

func (manager tokenManager) NewResource(ctx context.Context) (spellbook.Resource, error) {
	return &Token{}, nil
}

func (manager tokenManager) FromId(ctx context.Context, id string) (spellbook.Resource, error) {

	// todo
	current := spellbook.IdentityFromContext(ctx)
	if current == nil {
		return nil, spellbook.NewPermissionError(spellbook.PermissionName(spellbook.PermissionReadUser))
	}

	user, ok := current.(User)
	if ok && id == user.Id() {
		return &user, nil
	}

	us := User{}
	if err := model.FromStringID(ctx, &us, id, nil); err != nil {
		log.Errorf(ctx, "could not retrieve user %s: %s", id, err.Error())
		return nil, err
	}

	return &us, nil

	//return nil, spellbook.NewUnsupportedError()
}

func (manager tokenManager) ListOf(ctx context.Context, opts spellbook.ListOptions) ([]spellbook.Resource, error) {
	return nil, spellbook.NewUnsupportedError()
}

func (manager tokenManager) ListOfProperties(ctx context.Context, opts spellbook.ListOptions) ([]string, error) {
	return nil, spellbook.NewUnsupportedError()
}

func (manager tokenManager) Create(ctx context.Context, res spellbook.Resource, bundle []byte) error {

	token := res.(*Token)

	// checks the provided credentials. If correct creates a token, saves the user and returns the token
	nick := spellbook.NewRawField("username", true, token.Username)
	if _, err := nick.Value(); err != nil {
		return spellbook.NewFieldError("username", err)
	}

	password := spellbook.NewRawField("password", true, token.Password)
	password.AddValidator(spellbook.LenValidator{MinLen: 8})
	if _, err := password.Value(); err != nil {
		return spellbook.NewFieldError("password", err)
	}

	u := User{}
	err := model.FromStringID(ctx, &u, token.Username, nil)

	if err == datastore.ErrNoSuchEntity {
		return err
	}

	if err != nil {
		return err
	}

	hp := HashPassword(token.Password, salt)
	if u.Password != hp {
		return datastore.ErrNoSuchEntity
	}

	u.Token, err = u.GenerateToken()
	if err != nil {
		return fmt.Errorf("error generating token for user %s: %s", u.StringID(), err.Error())
	}

	err = model.Update(ctx, &u)
	if err != nil {
		return fmt.Errorf("error updating user token: %s", err.Error())
	}

	token.Value = u.Token

	return nil
}

func (manager tokenManager) Update(ctx context.Context, res spellbook.Resource, bundle []byte) error {
	return spellbook.NewUnsupportedError()
}

func (manager tokenManager) Delete(ctx context.Context, res spellbook.Resource) error {

	u := spellbook.IdentityFromContext(ctx)
	user, ok := u.(User)
	if !ok {
		return spellbook.NewPermissionError(spellbook.PermissionName(spellbook.PermissionEnabled))
	}

	user.Token = ""
	err := model.Update(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}
