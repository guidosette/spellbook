package page

import (
	"context"
	"distudio.com/mage"
	"distudio.com/mage/model"
	"distudio.com/page/identity"
	"distudio.com/page/post"
	"distudio.com/page/validators"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/appengine/log"
	"net/http"
	"time"
)

type MultimediaController struct {
	mage.Controller
}

func (controller *MultimediaController) OnDestroy(ctx context.Context) {}

func (controller *MultimediaController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	ins := mage.InputsFromContext(ctx)
	method := ins[mage.KeyRequestMethod].Value()
	switch method {
	case http.MethodPost:
		u := ctx.Value(identity.KeyUser)
		user, ok := u.(identity.User)
		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		// todo: permissions?
		j, ok := ins[mage.KeyRequestJSON]
		if !ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		errs := validators.Errors{}
		media := post.Multimedia{}
		err := json.Unmarshal([]byte(j.Value()), &media)
		if err != nil {
			msg := fmt.Sprintf("bad json input")
			errs.AddError("", errors.New(msg))
		}

		if errs.HasErrors() {
			log.Errorf(ctx, "wrong input to create multimedia: %s", errs)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		media.Created = time.Now().UTC()
		media.Uploader = user.Username()

		err = model.Create(ctx, &media)
		if err != nil {
			log.Errorf(ctx, "error creating multimedia %s: %s", media.Name, err)
			errs.AddError("", err)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = &media
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}
	}
	return mage.Redirect{Status: http.StatusNotImplemented}
}
