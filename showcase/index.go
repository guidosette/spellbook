package showcase

import (
	"context"
	"distudio.com/mage"
	"distudio.com/page"
	"distudio.com/page/configuration"
	"distudio.com/page/content"
	"distudio.com/page/identity"
	"distudio.com/page/mailmessage"
	"distudio.com/page/navigation"
	"golang.org/x/text/language"
	"net/http"
)

const (
	HOME     page.StaticPageCode = "HOME"
	PRODUCTS page.StaticPageCode = "PRODUCTS"
)

const (
	HOME_BANNER page.SpecialCode = "HOME_BANNER"
	FB          page.SpecialCode = "FB"
)

func init() {
	m := mage.Instance()

	opts := page.Options{}
	opts.Languages = []language.Tag{
		language.Italian,
		language.English,
	}
	opts.Categories = []page.SupportedCategory{
		page.SupportedCategory{Type: page.KeyTypeContent, Name: "services", Label: "Services"},
		page.SupportedCategory{Type: page.KeyTypeContent, Name: "news", Label: "News"},
		page.SupportedCategory{Type: page.KeyTypeEvent, Name: "events", Label: "Events"},
	}
	opts.StaticPages = []page.StaticPageCode{
		HOME,
		PRODUCTS,
	}
	opts.SpecialCodes = []page.SpecialCode{
		HOME_BANNER,
		FB,
	}
	opts.Actions = []page.SupportedAction{
		page.SupportedAction{Type: page.ActionTypeNormal, Name: "Clean index", Endpoint: "/normal"},
		page.SupportedAction{Type: page.ActionTypeUpload, Name: "Import places", Endpoint: "/upload"},
	}

	instance := page.NewWebsite(&opts)

	// superuser endpoints
	instance.Router.SetUniversalRoute("/api/me", func(ctx context.Context) mage.Controller {
		var key string
		if u := page.IdentityFromContext(ctx); u != nil {
			user := u.(identity.User)
			key = user.Id()
		}
		c := identity.NewUserControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/file", func(ctx context.Context) mage.Controller {
		c := content.NewFileController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	// superuser endpoints
	instance.Router.SetUniversalRoute("/api/users", func(ctx context.Context) mage.Controller {
		c := identity.NewUserController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/users/:username", func(ctx context.Context) mage.Controller {
		params := mage.RoutingParams(ctx)
		key := params["username"].Value()
		c := identity.NewUserControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})
	// backend
	instance.Router.SetUniversalRoute("/api/tokens", func(ctx context.Context) mage.Controller {
		c := identity.NewTokenController()
		return c
	}, nil)

	instance.Router.SetUniversalRoute("/api/tokens/:username", func(ctx context.Context) mage.Controller {
		// todo
		params := mage.RoutingParams(ctx)
		key := params["username"].Value()
		c := identity.NewTokenControllerWithKey(key)
		return c
	}, nil)

	instance.Router.SetUniversalRoute("/api/content", func(ctx context.Context) mage.Controller {
		c := content.NewContentController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/content/:id", func(ctx context.Context) mage.Controller {
		params := mage.RoutingParams(ctx)
		key := params["id"].Value()
		c := content.NewContentControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/languages", func(ctx context.Context) mage.Controller {
		c := configuration.NewLocaleController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/categories", func(ctx context.Context) mage.Controller {
		c := content.NewCategoryController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/actions", func(ctx context.Context) mage.Controller {
		c := content.NewActionController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/tasks", func(ctx context.Context) mage.Controller {
		c := content.NewTaskController("", "", "")
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/attachment", func(ctx context.Context) mage.Controller {
		c := content.NewAttachmentController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/attachment/:id", func(ctx context.Context) mage.Controller {
		params := mage.RoutingParams(ctx)
		key := params["id"].Value()
		c := content.NewAttachmentControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/place", func(ctx context.Context) mage.Controller {
		c := content.NewPlaceController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/place/:id", func(ctx context.Context) mage.Controller {
		params := mage.RoutingParams(ctx)
		key := params["id"].Value()
		c := content.NewPlaceControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/seo", func(ctx context.Context) mage.Controller {
		c := navigation.NewPageController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/seo/:id", func(ctx context.Context) mage.Controller {
		params := mage.RoutingParams(ctx)
		key := params["id"].Value()
		c := navigation.NewPageControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/mailmessage", func(ctx context.Context) mage.Controller {
		c := mailmessage.NewMailMessageController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/mailmessage/:id", func(ctx context.Context) mage.Controller {
		params := mage.RoutingParams(ctx)
		key := params["id"].Value()
		c := mailmessage.NewMailMessageControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/page", func(ctx context.Context) mage.Controller {
		c := navigation.NewPageController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/page/:id", func(ctx context.Context) mage.Controller {
		params := mage.RoutingParams(ctx)
		key := params["id"].Value()
		c := navigation.NewPageControllerWithKey(key)
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/staticpage", func(ctx context.Context) mage.Controller {
		c := page.NewStaticPageCodeController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	instance.Router.SetUniversalRoute("/api/specialcode", func(ctx context.Context) mage.Controller {
		c := page.NewSpecialCodeController()
		c.Private = true
		return c
	}, &identity.GSupportAuthenticator{})

	m.Router = &instance.Router
	m.LaunchApp(instance)
	http.HandleFunc("/", m.Run)
}
