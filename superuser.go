package page

import (
	"context"
	"distudio.com/mage"
	"distudio.com/mage/model"
	"distudio.com/page/identity"
	"distudio.com/page/validators"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
	"net/http"
)

// Returns 200 if the user is authenticated within the appengine framework
type IsSuperuserController struct {}

func (controller *IsSuperuserController) OnDestroy(ctx context.Context) {}

func (controller *IsSuperuserController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	// check if the user is a google user.
	// Gusers as admin bypass normal users controls
	guser := user.Current(ctx)
	if guser == nil {
		return mage.Redirect{Status: http.StatusUnauthorized}
	}

	return mage.Redirect{Status: http.StatusOK}
}

// this must be put behind a login endpoint
type SuperuserController struct {}

func (controller *SuperuserController) OnDestroy(ctx context.Context) {}

func (controller *SuperuserController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	guser := user.Current(ctx)
	if guser == nil {
		return mage.Redirect{Status: http.StatusUnauthorized}
	}

	if !guser.Admin {
		return mage.Redirect{Status: http.StatusForbidden}
	}

	ins := mage.InputsFromContext(ctx)
	method := ins[mage.KeyRequestMethod].Value()
	switch method {
	case http.MethodPost:
		j, ok := ins[mage.KeyRequestJSON]
		if !ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}
		errs := Errors{}

		jdata := j.Value()

		meta := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		err := json.Unmarshal([]byte(jdata), &meta)

		// check the username
		username := identity.SanitizeUserName(meta.Username)
		if username == "" {
			msg := fmt.Sprintf("invalid username %s", meta.Username)
			log.Errorf(ctx, msg)
			errs.AddError("username", errors.New(msg))
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// check if the user already exists
		u := identity.User{}
		err = model.FromStringID(ctx, &u, username, nil)
		if err == nil {
			// user already exists
			msg := fmt.Sprintf("user %s already exists.", username)
			log.Errorf(ctx, msg)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusConflict}
		}

		pf := validators.NewRawField("password", true, meta.Password)
		pf.AddValidator(validators.LenValidator{MinLen:8})

		if err = pf.Validate(); err != nil {
			errs.AddError(pf.Name, err)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusConflict}
		}

		err = json.Unmarshal([]byte(jdata), &u)
		if err != nil {
			log.Errorf(ctx, "error unmarshaling data %s : %s", jdata, err)
			errs.AddError("", err)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		opts := model.CreateOptions{}
		opts.WithStringId(username)
		u.AddPermission(identity.PermissionLogIn)
		err = model.CreateWithOptions(ctx, &u, &opts)
		if err != nil {
			log.Errorf(ctx, "error creating user %s: %s", username, err)
			errs.AddError("", err)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		return mage.Redirect{Status: http.StatusCreated}
	}

	return mage.Redirect{Status: http.StatusOK}
}
