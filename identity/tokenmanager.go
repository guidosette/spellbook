package identity

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"fmt"
	"google.golang.org/appengine/datastore"
)

type TokenManager struct{}

func (manager TokenManager) NewResource(ctx context.Context) (page.Resource, error) {
	return &Token{}, nil
}

func (manager TokenManager) FromId(ctx context.Context, id string) (page.Resource, error) {
	return nil, page.NewUnsupportedError()
}

func (manager TokenManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {
	return nil, page.NewUnsupportedError()
}

func (manager TokenManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	return nil, page.NewUnsupportedError()
}

func (manager TokenManager) Save(ctx context.Context, res page.Resource) error {

	token := res.(Token)
	u := User{}
	err := model.FromStringID(ctx, &u, token.Username, nil)

	if err == datastore.ErrNoSuchEntity {
		return err
	}

	if err != nil {
		return err
	}

	if u.Password != HashPassword(token.Password, salt) {
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

	return nil
}

func (manager TokenManager) Delete(ctx context.Context, res page.Resource) error {

	u := ctx.Value(KeyUser)
	user, ok := u.(User)
	if !ok {
		return page.NewPermissionError(PermissionName(PermissionEnabled))
	}

	user.Token = ""
	err := model.Update(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}
