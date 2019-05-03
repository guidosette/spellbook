package page

import (
	"distudio.com/mage"
	"distudio.com/mage/model"
	"distudio.com/page/content"
	"distudio.com/page/identity"
	"distudio.com/page/newsletter"
	"distudio.com/page/validators"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"strings"
)

type NewsletterController struct {
	mage.Controller
	BaseController
}

func (controller *NewsletterController) OnDestroy(ctx context.Context) {}

func (controller *NewsletterController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	ins := mage.InputsFromContext(ctx)
	method := ins[mage.KeyRequestMethod].Value()
	switch method {
	case http.MethodPost:
		//u := ctx.Value(identity.KeyUser)
		//user, ok := u.(identity.User)
		//if !ok {
		//	return mage.Redirect{Status: http.StatusUnauthorized}
		//}
		//
		//if !user.HasPermission(identity.PermissionReadNewsletter) {
		//	return mage.Redirect{Status: http.StatusForbidden}
		//}

		// get the content data
		j, ok := ins[mage.KeyRequestJSON]
		if !ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		errs := validators.Errors{}

		thenewsletter := newsletter.Newsletter{}
		err := json.Unmarshal([]byte(j.Value()), &thenewsletter)
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

		// validate input fields

		if thenewsletter.Email == "" {
			msg := fmt.Sprintf("Email can't be empty")
			errs.AddError("", errors.New(msg))
			log.Errorf(ctx, msg)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}
		if !strings.Contains(thenewsletter.Email, "@") || !strings.Contains(thenewsletter.Email, ".") {
			msg := fmt.Sprintf("Email not valid")
			errs.AddError("", errors.New(msg))
			log.Errorf(ctx, msg)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// list newsletter
		var emails []*newsletter.Newsletter
		q := model.NewQuery(&newsletter.Newsletter{})
		q = q.WithField("Email =", thenewsletter.Email)
		err = q.GetMulti(ctx, &emails)
		if err != nil {
			log.Errorf(ctx, "Error retrieving list newsletter %+v", err)
			return mage.Redirect{Status: http.StatusInternalServerError}
		}
		if len(emails) > 0 {
			msg := fmt.Sprintf("Email already exist")
			errs.AddError("", errors.New(msg))
			log.Errorf(ctx, msg)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// input is valid, create the resource
		// input is valid, create the resource
		opts := model.CreateOptions{}
		opts.WithStringId(thenewsletter.Email)

		err = model.CreateWithOptions(ctx, &thenewsletter, &opts)
		if err != nil {
			log.Errorf(ctx, "error creating newsletter %s: %s", thenewsletter.Email, err)
			errs.AddError("", err)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = &thenewsletter
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusCreated}
	case http.MethodGet:
		// check if current user has permission
		me := ctx.Value(identity.KeyUser)
		current, ok := me.(identity.User)

		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		if !current.HasPermission(identity.PermissionReadNewsletter) {
			return mage.Redirect{Status: http.StatusForbidden}
		}

		params := mage.RoutingParams(ctx)
		// try to get the username.
		// if there is no param then it is a list request
		param, ok := params["email"]
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
			_, ok := ins["csv"]
			if ok {
				// csv
				// list newsletter
				var newsletters []*newsletter.Newsletter
				q := model.NewQuery(&newsletter.Newsletter{})
				err := q.GetMulti(ctx, &newsletters)
				if err != nil {
					log.Errorf(ctx, "Error retrieving list newsletter %+v", err)
					return mage.Redirect{Status: http.StatusInternalServerError}
				}

				files := make([]string, len(newsletters), cap(newsletters))

				csvString := ""
				for i, newsletter := range newsletters {
					if i == 0 {
						csvString = fmt.Sprintf("\"%s\"", newsletter.Email)
					} else {
						csvString = fmt.Sprintf("%s\n\"%s\"", csvString, newsletter.Email)
					}
					files = append(files, newsletter.Email)
				}

				renderer := mage.TextRenderer{}
				renderer.Data = csvString
				out.Renderer = &renderer
				out.AddHeader("Content-type", "text/csv")
				out.AddHeader("Content-Disposition", "attachment;filename=newsletter.csv")
				return mage.Redirect{Status: http.StatusOK}
			} else {
				// list newsletter
				var newsletters []*newsletter.Newsletter
				q := model.NewQuery(&newsletter.Newsletter{})
				q = q.OffsetBy(page * size)
				// get one more so we know if we are done
				q = q.Limit(size + 1)
				err := q.GetMulti(ctx, &newsletters)
				if err != nil {
					log.Errorf(ctx, "Error retrieving list newsletter %+v", err)
					return mage.Redirect{Status: http.StatusInternalServerError}
				}
				l = len(newsletters)
				result = newsletters[:controller.GetCorrectCountForPaging(size, l)]
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

		email := param.Value()

		// get info content
		item := content.Content{}
		err := model.FromStringID(ctx, &item, email, nil)
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}

		if err != nil {
			log.Errorf(ctx, "error retrieving newsletter %s: %s", email, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = &item
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}

	case http.MethodDelete:
		//u := ctx.Value(identity.KeyUser)
		//user, ok := u.(identity.User)
		//if !ok {
		//	return mage.Redirect{Status: http.StatusUnauthorized}
		//}
		//
		//if !user.HasPermission(identity.PermissionEditNewsletter) {
		//	return mage.Redirect{Status: http.StatusForbidden}
		//}

		params := mage.RoutingParams(ctx)
		param, ok := params["email"]
		if !ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		email := param.Value()
		newsletter := newsletter.Newsletter{}
		err := model.FromStringID(ctx, &newsletter, email, nil)
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}
		if err != nil {
			log.Errorf(ctx, "error retrieving content %s: %s", email, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		err = model.Delete(ctx, &newsletter, nil)
		if err != nil {
			log.Errorf(ctx, "error deleting content %s: %s", email, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = nil
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}
	}

	return mage.Redirect{Status: http.StatusNotImplemented}
}
