package identity

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"fmt"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func NewTokenController() *page.RestController {
	handler := page.BaseRestHandler{}
	handler.Manager = tokenManager{}
	return page.NewRestController(handler)
}

func NewTokenControllerWithKey(key string) *page.RestController {
	handler := page.BaseRestHandler{Manager: tokenManager{}}
	c := page.NewRestController(handler)
	c.Key = key
	return c
}

type tokenManager struct{}

func (manager tokenManager) NewResource(ctx context.Context) (page.Resource, error) {
	return &Token{}, nil
}

func (manager tokenManager) FromId(ctx context.Context, id string) (page.Resource, error) {

	// todo
	current := page.IdentityFromContext(ctx)
	if current == nil {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadUser))
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

	//return nil, page.NewUnsupportedError()
}

func (manager tokenManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {
	return nil, page.NewUnsupportedError()
}

func (manager tokenManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	return nil, page.NewUnsupportedError()
}

func (manager tokenManager) Create(ctx context.Context, res page.Resource, bundle []byte) error {

	token := res.(*Token)

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

func (manager tokenManager) Update(ctx context.Context, res page.Resource, bundle []byte) error {
	return page.NewUnsupportedError()
}

func (manager tokenManager) Delete(ctx context.Context, res page.Resource) error {

	u := page.IdentityFromContext(ctx)
	user, ok := u.(User)
	if !ok {
		return page.NewPermissionError(page.PermissionName(page.PermissionEnabled))
	}

	user.Token = ""
	err := model.Update(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}
