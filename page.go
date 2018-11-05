package page

import (
	"distudio.com/mage"
	"golang.org/x/net/context"
)

type Website struct {
	mage.Application
}

func (w *Website) AuthenticatorForPath(token string) error {
	return nil
}

func (app Website) OnStart(ctx context.Context) context.Context {
	return ctx
}

func (app Website) AfterResponse(ctx context.Context) {}
