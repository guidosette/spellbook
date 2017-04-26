package page

import (
	"distudio.com/mage"
	"golang.org/x/net/context"
)

type Website struct {
	mage.Application
}

type user struct {
	mage.Authenticable;
}

func (u *user) Authenticate(ctx context.Context, token string) error {
	return nil;
}

func (app Website) OnCreate(ctx context.Context) context.Context {
	return ctx;
}

func (app Website) NewUser(ctx context.Context) mage.Authenticable {
	return nil;
}

func (app Website) AuthenticatorForPath(path string) mage.Authenticator {
	return nil;
}

func (app Website) OnDestroy(ctx context.Context) {

}

