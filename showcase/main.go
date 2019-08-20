package main

import (
	"context"
	"decodica.com/flamel"
	"decodica.com/spellbook"
	"decodica.com/spellbook/configuration"
	"decodica.com/spellbook/content"
	"decodica.com/spellbook/identity"
	"decodica.com/spellbook/mailmessage"
	"decodica.com/spellbook/navigation"
	"golang.org/x/text/language"
	"net/http"
)

const (
	HOME     spellbook.StaticPageCode = "HOME"
	PRODUCTS spellbook.StaticPageCode = "PRODUCTS"
)

const (
	HomeBanner spellbook.SpecialCode = "HOME_BANNER"
	FB         spellbook.SpecialCode = "FB"
)

func main() {
	m := flamel.Instance()

	opts := spellbook.Options{}
	opts.Languages = []language.Tag{
		language.Italian,
		language.English,
	}
	opts.Categories = []spellbook.SupportedCategory{
		{Type: spellbook.KeyTypeContent, Name: "services", Label: "Services"},
		{Type: spellbook.KeyTypeContent, Name: "news", Label: "News", DefaultAttachmentGroups: []spellbook.DefaultAttachmentGroup{
			{"Gallery", content.AttachmentTypeGallery, 0, "Prova descr"},
		}},
		{Type: spellbook.KeyTypeEvent, Name: "events", Label: "Events"},
	}
	opts.StaticPages = []spellbook.StaticPageCode{
		HOME,
		PRODUCTS,
	}
	opts.SpecialCodes = []spellbook.SpecialCode{
		HomeBanner,
		FB,
	}
	opts.Actions = []spellbook.SupportedAction{
		{Type: spellbook.ActionTypeNormal, Name: "cleanindextest", Endpoint: "/api/cleanindextest", Method: http.MethodGet},
		{Type: spellbook.ActionTypeUpload, Name: "places", Endpoint: "/api/places", Method: http.MethodGet},
	}

	instance := spellbook.NewWebsite(&opts)

	// superuser endpoints
	instance.Router.SetUniversalRoute("/api/me", func(ctx context.Context) flamel.Controller {
		var key string
		if u := spellbook.IdentityFromContext(ctx); u != nil {
			user := u.(identity.User)
			key = user.Id()
		}
		c := identity.NewUserControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/file", func(ctx context.Context) flamel.Controller {
		c := content.NewFileController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	// superuser endpoints
	instance.Router.SetUniversalRoute("/api/users", func(ctx context.Context) flamel.Controller {
		c := identity.NewUserController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/users/:username", func(ctx context.Context) flamel.Controller {
		params := flamel.RoutingParams(ctx)
		key := params["username"].Value()
		c := identity.NewUserControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})
	// backend
	instance.Router.SetUniversalRoute("/api/tokens", func(ctx context.Context) flamel.Controller {
		c := identity.NewTokenController()
		return c
	}, nil)

	instance.Router.SetUniversalRoute("/api/tokens/:username", func(ctx context.Context) flamel.Controller {
		// todo
		params := flamel.RoutingParams(ctx)
		key := params["username"].Value()
		c := identity.NewTokenControllerWithKey(key)
		return c
	}, nil)

	instance.Router.SetUniversalRoute("/api/content", func(ctx context.Context) flamel.Controller {
		c := content.NewContentController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/content/:id", func(ctx context.Context) flamel.Controller {
		params := flamel.RoutingParams(ctx)
		key := params["id"].Value()
		c := content.NewContentControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/languages", func(ctx context.Context) flamel.Controller {
		c := configuration.NewLocaleController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/categories", func(ctx context.Context) flamel.Controller {
		c := content.NewCategoryController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/actions", func(ctx context.Context) flamel.Controller {
		c := content.NewActionController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/task", func(ctx context.Context) flamel.Controller {
		c := content.NewTaskController("", "", "")
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/task/:id", func(ctx context.Context) flamel.Controller {
		params := flamel.RoutingParams(ctx)
		key := params["id"].Value()
		c := content.NewTaskControllerWithKey(key, "", "", "")
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/attachment", func(ctx context.Context) flamel.Controller {
		c := content.NewAttachmentController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/attachment/:id", func(ctx context.Context) flamel.Controller {
		params := flamel.RoutingParams(ctx)
		key := params["id"].Value()
		c := content.NewAttachmentControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/place", func(ctx context.Context) flamel.Controller {
		c := content.NewPlaceController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/place/:id", func(ctx context.Context) flamel.Controller {
		params := flamel.RoutingParams(ctx)
		key := params["id"].Value()
		c := content.NewPlaceControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/seo", func(ctx context.Context) flamel.Controller {
		c := navigation.NewPageController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/seo/:id", func(ctx context.Context) flamel.Controller {
		params := flamel.RoutingParams(ctx)
		key := params["id"].Value()
		c := navigation.NewPageControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/mailmessage", func(ctx context.Context) flamel.Controller {
		c := mailmessage.NewMailMessageController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/mailmessage/:id", func(ctx context.Context) flamel.Controller {
		params := flamel.RoutingParams(ctx)
		key := params["id"].Value()
		c := mailmessage.NewMailMessageControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/page", func(ctx context.Context) flamel.Controller {
		c := navigation.NewPageController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/page/:id", func(ctx context.Context) flamel.Controller {
		params := flamel.RoutingParams(ctx)
		key := params["id"].Value()
		c := navigation.NewPageControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/staticpage", func(ctx context.Context) flamel.Controller {
		c := spellbook.NewStaticPageCodeController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/specialcode", func(ctx context.Context) flamel.Controller {
		c := spellbook.NewSpecialCodeController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/cleanindextest", func(ctx context.Context) flamel.Controller {
		c := CleanController{}
		return &c
	}, nil)

	m.Router = &instance.Router
	m.LaunchApp(instance)
	m.Run()
}
