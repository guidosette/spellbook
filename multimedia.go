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
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"reflect"
	"strconv"
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
		return mage.Redirect{Status: http.StatusCreated}
	case http.MethodGet:
		// check if current user has permission
		me := ctx.Value(identity.KeyUser)
		_, ok := me.(identity.User)

		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		params := mage.RoutingParams(ctx)
		// try to get the username.
		// if there is no param then it is a list request
		param, ok := params["id"]
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

			var result interface{}
			l := 0
			// check property
			property, ok := ins["property"]
			if ok {
				// property
				properties, err := controller.HandleResourceProperties(ctx, property.Value(), page, size)
				if err != nil {
					log.Errorf(ctx, "Error retrieving multimedia %+v", err)
					return mage.Redirect{Status: http.StatusInternalServerError}
				}
				l = len(properties)
				result = properties[:controller.GetCorrectCountForPaging(size, l)]
			} else {
				// list posts
				var posts []*post.Multimedia
				q := model.NewQuery(&post.Multimedia{})
				q = q.OffsetBy(page * size)
				// get one more so we know if we are done
				q = q.Limit(size + 1)
				err := q.GetMulti(ctx, &posts)
				if err != nil {
					log.Errorf(ctx, "Error retrieving multimedia %+v", err)
					return mage.Redirect{Status: http.StatusInternalServerError}
				}
				l = len(posts)
				result = posts[:controller.GetCorrectCountForPaging(size, l)]
			}

			// todo: generalize list handling and responses
			response := struct {
				Items interface{} `json:"items"`
				More  bool        `json:"more"`
			}{result, l > size}
			renderer := mage.JSONRenderer{}
			renderer.Data = response
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusOK}
		}

		id := param.Value()
		item := post.Multimedia{}
		err := model.FromStringID(ctx, &item, id, nil)
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}

		if err != nil {
			log.Errorf(ctx, "error retrieving multimedia %s: %s", id, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		response := struct {
			*post.Multimedia
		}{&item}

		renderer := mage.JSONRenderer{}
		renderer.Data = response
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}
	case http.MethodPut:
		me := ctx.Value(identity.KeyUser)
		_, ok := me.(identity.User)
		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		// control if the user has been specified
		params := mage.RoutingParams(ctx)
		param, ok := params["id"]
		if !ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		j, ok := ins[mage.KeyRequestJSON]
		if !ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// handle the json request
		jdata := j.Value()

		jmultimedia := post.Multimedia{}

		err := json.Unmarshal([]byte(jdata), &jmultimedia)
		if err != nil {
			log.Errorf(ctx, "malformed json: %s", err.Error())
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// retrieve the user
		id := param.Value()
		p := post.Multimedia{}
		err = model.FromStringID(ctx, &p, id, nil)
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}

		if err != nil {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		p.Name = jmultimedia.Name
		p.Description = jmultimedia.Description
		p.ResourceUrl = jmultimedia.ResourceUrl
		p.Group = jmultimedia.Group
		p.Updated = time.Now().UTC()

		err = model.Update(ctx, &p)
		if err != nil {
			log.Errorf(ctx, "error updating multimedia %s: %s", id, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = &p
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}

	}
	return mage.Redirect{Status: http.StatusNotImplemented}
}

func (controller *MultimediaController) GetCorrectCountForPaging(size int, l int) int {
	count := size
	if l < size {
		count = l
	}
	return count
}

func (controller *MultimediaController) HandleResourceProperties(ctx context.Context, property string, page int, size int) ([]interface{}, error) {
	// todo: generalize
	name := ""
	switch property {
	case "group":
		name = "Group"
	default:
		return nil, errors.New("No property found")
	}

	var posts []*post.Post
	q := model.NewQuery(&post.Post{})
	q = q.OffsetBy(page * size)
	q = q.Distinct(name)
	// get one more so we know if we are done
	q = q.Limit(size + 1)
	err := q.GetAll(ctx, &posts)
	if err != nil {
		log.Errorf(ctx, "Error retrieving result: %+v", err)
		return nil, err
	}
	var result []interface{}
	for _, p := range posts {
		value := reflect.ValueOf(p).Elem().FieldByName(name).String()
		result = append(result, &value)
	}
	return result, nil

}
