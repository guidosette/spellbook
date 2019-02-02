package controller

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
	return mage.Redirect{Location:controller.To, Status:http.StatusFound}
}

func (controller *RedirectController) OnDestroy(ctx context.Context) {

}
