package page

import (
	"context"
	"distudio.com/mage"
	"distudio.com/page/identity"
	"distudio.com/page/post"
	"encoding/json"
	"google.golang.org/appengine/log"
	"net/http"
	"net/url"
	"time"
)

type PostController struct {
	mage.Controller
}

func (controller *PostController) OnDestroy(ctx context.Context) {}

func (controller *PostController) Process(ctx context.Context, out mage.ResponseOutput) mage.Redirect {
	ins := mage.InputsFromContext(ctx)
	method := ins[mage.KeyRequestMethod].Value()
	switch method {
	case http.MethodPost:
		u := ctx.Value(identity.KeyUser)
		user, ok := u.(identity.User)
		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		if !user.HasPermission(identity.PermissionCreatePost) {
			return mage.Redirect{Status: http.StatusForbidden}
		}

		// get the post data
		j, ok:= ins[mage.KeyRequestJSON]
		if !ok {
			return mage.Redirect{Status:http.StatusBadRequest}
		}

		thepost := post.Post{}
		err := json.Unmarshal([]byte(j.Value()), &thepost)
		if err != nil {
			log.Errorf(ctx, "bad json: %s", err.Error())
			return mage.Redirect{Status:http.StatusBadRequest}
		}

		thepost.Created = time.Now().UTC()
		thepost.Revision = 1

		if thepost.Title == "" || thepost.Body == "" {
			log.Errorf(ctx, "the body or the title can't be both empty")
			return mage.Redirect{Status:http.StatusBadRequest}
		}

		if thepost.Slug == "" {
			thepost.Slug = url.PathEscape(thepost.Title)
		}


	}
	return mage.Redirect{Status:http.StatusNotImplemented}
}
