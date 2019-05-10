package showcase

import (
	"context"
	"distudio.com/mage"
	"distudio.com/page"
	"distudio.com/page/configuration"
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
	opts.Categories = []page.Category{
		"services",
		"news",
	}

	instance := page.NewWebsite(&opts)

	// superuser endpoints
	instance.Router.SetUniversalRoute("/api/me", func(ctx context.Context) mage.Controller {
		var key string
		if u := ctx.Value(identity.KeyUser); u != nil {
			user := u.(identity.User)
			key = user.Id()
		}
		return identity.NewUserControllerWithKey(key)
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/file", func(ctx context.Context) mage.Controller {
		return content.NewFileController()
	}, &identity.GSupportAuthenticator{})

	// superuser endpoints
	instance.Router.SetUniversalRoute("/api/users", func(ctx context.Context) mage.Controller {
		return identity.NewUserController()
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/users/:username", func(ctx context.Context) mage.Controller {
		params := mage.RoutingParams(ctx)
		key := params["username"].Value()
		return identity.NewUserControllerWithKey(key)
	}, &identity.GSupportAuthenticator{})
	// backend
	instance.Router.SetUniversalRoute("/api/tokens", func(ctx context.Context) mage.Controller {
		return identity.NewTokenController()
	}, nil)

	instance.Router.SetUniversalRoute("/api/content", func(ctx context.Context) mage.Controller {
		return content.NewContentController()
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/content/:slug", func(ctx context.Context) mage.Controller {
		params := mage.RoutingParams(ctx)
		key := params["slug"].Value()
		return content.NewContentControllerWithKey(key)
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/languages", func(ctx context.Context) mage.Controller {
		return configuration.NewLocaleController()
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/categories", func(ctx context.Context) mage.Controller {
		return configuration.NewCategoryController()
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/attachment", func(ctx context.Context) mage.Controller {
		return content.NewAttachmentController()
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/attachment/:id", func(ctx context.Context) mage.Controller {
		params := mage.RoutingParams(ctx)
		key := params["id"].Value()
		return content.NewAttachmentControllerWithKey(key)
	}, &identity.GSupportAuthenticator{})

	m.Router = &instance.Router
	m.LaunchApp(instance)
	http.HandleFunc("/", m.Run)
}
