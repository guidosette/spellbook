package page

import (
	"distudio.com/mage"
	"distudio.com/page/identity"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

type LocaleController struct {
	mage.Controller
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
		me := ctx.Value(identity.KeyUser)
		_, ok := me.(identity.User)

		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		// handle query params for page data:
		size := 20
		//page := 0
		//if pin, ok := ins["page"]; ok {
		//	if num, err := strconv.Atoi(pin.Value()); err == nil {
		//		page = num
		//	} else {
		//		return mage.Redirect{Status: http.StatusBadRequest}
		//	}
		//}

		if sin, ok := ins["results"]; ok {
			if num, err := strconv.Atoi(sin.Value()); err == nil {
				size = num
				// cap the size to 100
				if size > 100 {
					size = 100
				}
			} else {
				return mage.Redirect{Status: http.StatusBadRequest}
			}
		}

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

func (controller *LocaleController) GetCorrectCountForPaging(size int, l int) int {
	count := size
	if l < size {
		count = l
	}
	return count
}
