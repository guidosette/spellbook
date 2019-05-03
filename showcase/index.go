package showcase

import (
	"context"
	"distudio.com/mage"
	"distudio.com/page"
	"distudio.com/page/content"
	"distudio.com/page/identity"
	"golang.org/x/text/language"
	"net/http"
)

func init() {
	m := mage.Instance()
	//m.EnforceHostnameRedirect =  // todo

	opts := page.Options{}
	opts.Languages = []language.Tag{
		language.Italian,
		language.English,
	}

	instance := page.NewWebsite(&opts)

	// superuser endpoints
	instance.Router.SetUniversalRoute("/api/me", func(ctx context.Context) mage.Controller {
		return &page.IdentityController{}
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/file", func(ctx context.Context) mage.Controller {
		return &page.FileController{}
	}, &identity.GSupportAuthenticator{})

	// superuser endpoints
	instance.Router.SetUniversalRoute("/api/users", func(ctx context.Context) mage.Controller {
		return &page.UserController{}
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/users/:username", func(ctx context.Context) mage.Controller {
		return &page.UserController{}
	}, &identity.GSupportAuthenticator{})
	// backend
	instance.Router.SetUniversalRoute("/api/tokens", func(ctx context.Context) mage.Controller {
		return &page.TokenController{}
	}, nil)

	instance.Router.SetUniversalRoute("/api/content", func(ctx context.Context) mage.Controller {
		c := &page.Controller{}
		m := content.Manager{}
		c.Manager = m
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/content/:slug", func(ctx context.Context) mage.Controller {
		c := &page.Controller{}
		m := content.Manager{}
		c.Manager = m
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/languages", func(ctx context.Context) mage.Controller {
		return &page.LocaleController{}
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/attachment", func(ctx context.Context) mage.Controller {
		return &page.AttachmentController{}
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/attachment/:id", func(ctx context.Context) mage.Controller {
		return &page.AttachmentController{}
	}, &identity.GSupportAuthenticator{})

	m.Router = &instance.Router
	m.LaunchApp(instance)
	http.HandleFunc("/", m.Run)
}
