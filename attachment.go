package page

import (
	"distudio.com/mage"
	"distudio.com/mage/model"
	"distudio.com/page/resource/attachment"
	"distudio.com/page/resource/content"
	user2 "distudio.com/page/resource/identity"
	"distudio.com/page/validators"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type AttachmentController struct {
	mage.Controller
	BaseController
}

func (controller *AttachmentController) OnDestroy(ctx context.Context) {}

func (controller *AttachmentController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	ins := mage.InputsFromContext(ctx)
	method := ins[mage.KeyRequestMethod].Value()
	switch method {
	case http.MethodPost:
		u := ctx.Value(user2.KeyUser)
		user, ok := u.(user2.User)
		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		// todo: permissions?
		j, ok := ins[mage.KeyRequestJSON]
		if !ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		errs := validators.Errors{}
		att := attachment.Attachment{}
		err := json.Unmarshal([]byte(j.Value()), &att)
		if err != nil {
			msg := fmt.Sprintf("bad json input: %s", err.Error())
			errs.AddError("", errors.New(msg))
		}

		// attachment parent is required.
		// if not attachment is to be specified the default value must be used
		if att.Parent == "" {
			msg := fmt.Sprintf("attachment parent can't be empty. Use %s as a parent for global attachments", attachment.AttachmentGlobalParent)
			errs.AddError("Parent", errors.New(msg))
		}

		if errs.HasErrors() {
			log.Errorf(ctx, "wrong input to create attachment: %s", errs)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		att.Created = time.Now().UTC()
		att.Uploader = user.Username()

		err = model.Create(ctx, &att)
		if err != nil {
			log.Errorf(ctx, "error creating attachment %s: %s", att.Name, err)
			errs.AddError("", err)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = &att
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusCreated}
	case http.MethodGet:
		// check if current user has permission
		me := ctx.Value(user2.KeyUser)
		_, ok := me.(user2.User)

		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		params := mage.RoutingParams(ctx)
		// try to get the username.
		// if there is no param then it is a list request
		param, ok := params["id"]
		if !ok {

			// handle query params for page data:
			paging, err := controller.GetPaging(ins)
			if err != nil {
				return mage.Redirect{Status: http.StatusBadRequest}
			}
			page := paging.page
			size := paging.size

			var result interface{}
			l := 0
			// check property
			property, ok := ins["property"]
			if ok {
				// property
				properties, err := controller.HandleResourceProperties(ctx, property.Value(), page, size)
				if err != nil {
					log.Errorf(ctx, "Error retrieving attachment %+v", err)
					return mage.Redirect{Status: http.StatusInternalServerError}
				}
				l = len(properties)
				result = properties[:controller.GetCorrectCountForPaging(size, l)]
			} else {
				var attachments []*attachment.Attachment
				q := model.NewQuery(&attachment.Attachment{})
				q = q.OffsetBy(page * size)
				// get one more so we know if we are done
				q = q.Limit(size + 1)
				err := q.GetMulti(ctx, &attachments)
				if err != nil {
					log.Errorf(ctx, "Error retrieving attachment %+v", err)
					return mage.Redirect{Status: http.StatusInternalServerError}
				}
				l = len(attachments)
				result = attachments[:controller.GetCorrectCountForPaging(size, l)]
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
		item := attachment.Attachment{}
		err := model.FromStringID(ctx, &item, id, nil)
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}

		if err != nil {
			log.Errorf(ctx, "error retrieving attachment %s: %s", id, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		response := struct {
			*attachment.Attachment
		}{&item}

		renderer := mage.JSONRenderer{}
		renderer.Data = response
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}
	case http.MethodPut:
		me := ctx.Value(user2.KeyUser)
		_, ok := me.(user2.User)
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

		jatt := attachment.Attachment{}
		err := json.Unmarshal([]byte(jdata), &jatt)
		if err != nil {
			log.Errorf(ctx, "malformed json: %s", err.Error())
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// retrieve the attachment
		id := param.Value()
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Errorf(ctx, "error convert id %s: %s", id, err.Error())
			return mage.Redirect{Status: http.StatusBadRequest}
		}
		att := attachment.Attachment{}
		err = model.FromIntID(ctx, &att, idInt, nil)
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}

		errs := validators.Errors{}
		if att.Parent == "" {
			msg := fmt.Sprintf("attachment parent can't be empty. Use %s as a parent for global attachments", attachment.AttachmentGlobalParent)
			errs.AddError("parent", errors.New(msg))
		}

		if err != nil {
			errs.AddError("", err)
		}

		if errs.HasErrors() {
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		att.Name = jatt.Name
		att.Description = jatt.Description
		att.ResourceUrl = jatt.ResourceUrl
		att.Group = jatt.Group
		att.Updated = time.Now().UTC()

		err = model.Update(ctx, &att)
		if err != nil {
			log.Errorf(ctx, "error updating multimedia %s: %s", id, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = &att
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}

	case http.MethodDelete:
		u := ctx.Value(user2.KeyUser)
		user, ok := u.(user2.User)
		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		if !user.HasPermission(user2.PermissionEditContent) {
			return mage.Redirect{Status: http.StatusForbidden}
		}

		params := mage.RoutingParams(ctx)
		param, ok := params["id"]
		if !ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		id := param.Value()
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Errorf(ctx, "error convert id %s: %s", id, err.Error())
			return mage.Redirect{Status: http.StatusBadRequest}
		}
		attachment := attachment.Attachment{}
		err = model.FromIntID(ctx, &attachment, idInt, nil)
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}
		if err != nil {
			log.Errorf(ctx, "error retrieving attachment %s: %s", id, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		// delete attachment
		err = model.Delete(ctx, &attachment, nil)
		if err != nil {
			log.Errorf(ctx, "error deleting attachment %s: %s", id, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = nil
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}

	}
	return mage.Redirect{Status: http.StatusNotImplemented}
}

func (controller *AttachmentController) HandleResourceProperties(ctx context.Context, property string, page int, size int) ([]interface{}, error) {
	// todo: generalize
	name := ""
	switch property {
	case "group":
		name = "Group"
	default:
		return nil, errors.New("No property found")
	}

	var posts []*content.Content
	q := model.NewQuery(&content.Content{})
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
