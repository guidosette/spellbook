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
	"reflect"
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

		// get the p data
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
		if thepost.Published != post.ZeroTime {
			thepost.Published = time.Now().UTC()
		}
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

		// retrieve the multimedia groups
		var media []post.Multimedia
		err = json.Unmarshal([]byte(j.Value()), media)
		if err != nil {
			msg := fmt.Sprintf("bad input for multimedia: %s", err.Error())
			errs.AddError("multimedia", errors.New(msg))
			log.Errorf(ctx, msg)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		for _, m := range media {
			if !thepost.HasMultimedia(m) {
				thepost.AddMultimedia(m)
			}
		}

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

			var result interface{}
			l := 0
			// check property
			property, ok := ins["property"]
			if ok {
				// property
				properties, err := controller.HandleResourceProperties(ctx, property.Value(), page, size)
				if err != nil {
					log.Errorf(ctx, "Error retrieving posts %+v", err)
					return mage.Redirect{Status: http.StatusInternalServerError}
				}
				l = len(properties)
				result = properties[:controller.GetCorrectCountForPaging(size, l)]
			} else {
				// list posts
				var posts []*post.Post
				q := model.NewQuery(&post.Post{})
				q = q.OffsetBy(page * size)
				// get one more so we know if we are done
				q = q.Limit(size + 1)
				err := q.GetMulti(ctx, &posts)
				if err != nil {
					log.Errorf(ctx, "Error retrieving posts %+v", err)
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

		slug := param.Value()
		item := post.Post{}
		err := model.FromStringID(ctx, &item, slug, nil)
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}

		if err != nil {
			log.Errorf(ctx, "error retrieving p %s: %s", slug, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		// get post related multimedia
		q := model.NewQuery(&post.Multimedia{})
		for _, m := range item.MultimediaGroups {
			var mm []*post.Multimedia
			q.WithField("Group =", m)
			err := q.GetMulti(ctx, &mm)
			if err != nil {
				log.Errorf(ctx, "error retrieving multimedia: %s", err)
				return mage.Redirect{Status: http.StatusInternalServerError}
			}
			for _, val := range mm {
				item.Multimedia = append(item.Multimedia, *val)
			}
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

		if !current.HasPermission(identity.PermissionEditPost) {
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

		jpost := post.Post{}

		err := json.Unmarshal([]byte(jdata), &jpost)
		if err != nil {
			log.Errorf(ctx, "malformed json: %s", err.Error())
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// retrieve the user
		slug := param.Value()
		p := post.Post{}
		err = model.FromStringID(ctx, &p, slug, nil)
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}

		if err != nil {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		p.Name = jpost.Name
		p.Title = jpost.Title
		p.Subtitle = jpost.Subtitle
		p.Category = jpost.Category
		p.Topic = jpost.Topic
		p.Locale = jpost.Locale
		p.Body = jpost.Body
		p.Cover = jpost.Cover
		p.Revision = jpost.Revision
		p.Updated = time.Now().UTC()
		p.Tags = jpost.Tags
		p.Author = current.Username()
		if jpost.Published == post.ZeroTime {
			// not setted
			p.Published = post.ZeroTime
		} else {
			// setted
			// check previous data
			if p.Published == post.ZeroTime {
				p.Published = time.Now().UTC()
			}
		}

		// retrieve the multimedia groups
		var media []post.Multimedia
		err = json.Unmarshal([]byte(j.Value()), media)
		errs := validators.Errors{}
		if err != nil {
			msg := fmt.Sprintf("bad input for multimedia: %s", err.Error())
			errs.AddError("multimedia", errors.New(msg))
			log.Errorf(ctx, msg)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// reset multimedia groups
		p.MultimediaGroups = nil
		for _, m := range media {
			if !p.HasMultimedia(m) {
				p.AddMultimedia(m)
			}
		}

		err = model.Update(ctx, &p)
		if err != nil {
			log.Errorf(ctx, "error updating p %s: %s", slug, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = &p
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}
	}

	return mage.Redirect{Status: http.StatusNotImplemented}
}

func (controller *PostController) GetCorrectCountForPaging(size int, l int) int {
	count := size
	if l < size {
		count = l
	}
	return count
}

func (controller *PostController) HandleResourceProperties(ctx context.Context, property string, page int, size int) ([]interface{}, error) {
	// todo: generalize
	name := ""
	switch property {
	case "category":
		name = "Category"
	case "topic":
		name = "Topic"
	case "name":
		name = "Name"
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
