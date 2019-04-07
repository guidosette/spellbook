package page

import (
	"distudio.com/mage"
	"distudio.com/page/resource"
	"golang.org/x/net/context"
	"net/http"
)

type LocaleController struct {
	mage.Controller
	BaseController
}

func (controller *LocaleController) OnDestroy(ctx context.Context) {}

func (controller *LocaleController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	ins := mage.InputsFromContext(ctx)
	method := ins[mage.KeyRequestMethod].Value()
	switch method {
	case http.MethodPost:
		return mage.Redirect{Status: http.StatusNotImplemented}
	case http.MethodGet:
		// check if current user has permission
		me := ctx.Value(resource.KeyUser)
		_, ok := me.(resource.User)

		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		// handle query params for page data:
		paging, err := controller.GetPaging(ins)
		if err != nil {
			return mage.Redirect{Status: http.StatusBadRequest}
		}
		//page := paging.page
		size := paging.size

		var result interface{}
		l := 0
		// list language

		ws := application()
		languages := ws.options.Languages
		l = len(languages)
		result = languages[:controller.GetCorrectCountForPaging(size, l)]

		// todo: generalize list handling and responses
		response := struct {
			Items interface{} `json:"items"`
			More  bool        `json:"more"`
		}{result, l > size}
		renderer := mage.JSONRenderer{}
		renderer.Data = response
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}

	case http.MethodPut:
		return mage.Redirect{Status: http.StatusNotImplemented}
	}

	return mage.Redirect{Status: http.StatusMethodNotAllowed}
}
