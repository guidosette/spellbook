package page

import (
	"context"
	"distudio.com/mage"
	"net/http"
)

type RedirectController struct {
	mage.Controller
	To string
}

func (controller *RedirectController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	return mage.Redirect{Location: controller.To, Status: http.StatusFound}
}

func (controller *RedirectController) OnDestroy(ctx context.Context) {}

// Returns a 307 redirect
// see https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/307
type TemporaryRedirectController struct {
	mage.Controller
	To string
}

func (controller *TemporaryRedirectController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	return mage.Redirect{Location: controller.To, Status: http.StatusTemporaryRedirect}
}

func (controller *TemporaryRedirectController) OnDestroy(ctx context.Context) {}
