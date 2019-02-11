package page

import (
	"distudio.com/mage"
	"distudio.com/mage/model"
	"distudio.com/page/identity"
	"distudio.com/page/post"
	"distudio.com/page/validators"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type PostController struct {
	mage.Controller
}

func (controller *PostController) OnDestroy(ctx context.Context) {}

func (controller *PostController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
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
		j, ok := ins[mage.KeyRequestJSON]
		if !ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		errs := validators.Errors{}

		thepost := post.Post{}
		err := json.Unmarshal([]byte(j.Value()), &thepost)
		if err != nil {
			msg := fmt.Sprintf("bad json: %s", err.Error())
			errs.AddError("", errors.New(msg))
			log.Errorf(ctx, msg)
		}

		// check for client input erros
		if errs.HasErrors() {
			log.Errorf(ctx, "error HasErrors %+v", errs)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		thepost.Created = time.Now().UTC()
		thepost.Revision = 1

		// validate input fields

		if thepost.Title == "" || thepost.Body == "" {
			msg := fmt.Sprintf("the body or the title can't be both empty")
			errs.AddError("", errors.New(msg))
			log.Errorf(ctx, msg)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		if thepost.Slug == "" {
			thepost.Slug = url.PathEscape(thepost.Title)
		}
		thepost.Author = user.Username()

		// input is valid, create the resource
		opts := model.CreateOptions{}
		opts.WithStringId(thepost.Slug)

		err = model.CreateWithOptions(ctx, &thepost, &opts)
		if err != nil {
			log.Errorf(ctx, "error creating post %s: %s", thepost.Slug, err)
			errs.AddError("", err)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = &thepost
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusCreated}
	case http.MethodGet:
		// check if current user has permission
		me := ctx.Value(identity.KeyUser)
		current, ok := me.(identity.User)

		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		if !current.HasPermission(identity.PermissionReadPost) {
			return mage.Redirect{Status: http.StatusForbidden}
		}

		params := mage.RoutingParams(ctx)
		// try to get the username.
		// if there is no param then it is a list request
		param, ok := params["slug"]
		if !ok {
			// handle query params for page data:
			page := 0
			size := 20
			if pin, ok := ins["page"]; ok {
				if num, err := strconv.Atoi(pin.Value()); err == nil {
					page = num
				} else {
					return mage.Redirect{Status: http.StatusBadRequest}
				}
			}

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

			var posts []*post.Post
			q := model.NewQuery(&post.Post{})
			q = q.OffsetBy(page * size)
			// get one more so we know if we are done
			q = q.Limit(size + 1)
			err := q.GetMulti(ctx, &posts)
			if err != nil {
				return mage.Redirect{Status: http.StatusInternalServerError}
			}

			// todo: generalize list handling and responses
			l := len(posts)
			count := size
			if l < size {
				count = l
			}
			response := struct {
				Items []*post.Post `json:"items"`
				More  bool         `json:"more"`
			}{posts[:count], l > size}
			renderer := mage.JSONRenderer{}
			renderer.Data = response
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusOK}
		}

		slug := param.Value()
		item := post.Post{}
		err := model.FromStringID(ctx, &item, slug, nil)
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}

		if err != nil {
			log.Errorf(ctx, "error retrieving post %s: %s", slug, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = &item
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}
	case http.MethodPut:
		me := ctx.Value(identity.KeyUser)
		current, ok := me.(identity.User)
		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		if !current.HasPermission(identity.PermissionEditUser) {
			return mage.Redirect{Status: http.StatusForbidden}
		}

		// control if the user has been specified
		params := mage.RoutingParams(ctx)
		param, ok := params["slug"]
		if !ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		j, ok := ins[mage.KeyRequestJSON]
		if !ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// handle the json request
		jdata := j.Value()

		jpost := struct {
			*post.Post
		}{Post: &post.Post{}}

		err := json.Unmarshal([]byte(jdata), &jpost)
		if err != nil {
			log.Errorf(ctx, "malformed json: %s", err.Error())
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// retrieve the user
		slug := param.Value()
		post := post.Post{}
		err = model.FromStringID(ctx, &post, slug, nil)
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}

		if err != nil {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		post.Title = jpost.Title
		post.Subtitle = jpost.Subtitle
		post.Category = jpost.Category
		post.Topic = jpost.Topic
		post.Body = jpost.Body
		//target.Locale = juser.Locale
		post.Revision = jpost.Revision
		post.Updated = time.Now().UTC()
		post.Tags = jpost.Tags
		post.Author = current.Username()

		err = model.Update(ctx, &post)
		if err != nil {
			log.Errorf(ctx, "error updating post %s: %s", slug, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = &post
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}
	}

	return mage.Redirect{Status: http.StatusMethodNotAllowed}
}
