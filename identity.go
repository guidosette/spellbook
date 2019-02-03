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
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
	"net/http"
)

type TokenController struct {}

func (controller *TokenController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	ins := mage.InputsFromContext(ctx)
	method := ins[mage.KeyRequestMethod].Value()
	switch method {
	case http.MethodPost:
		return controller.HandlePost(ctx, ins, out)
	case http.MethodDelete:
		return controller.HandleDelete(ctx, ins, out)
	}
	return mage.Redirect{Status: http.StatusMethodNotAllowed}
}

func (controller TokenController) HandlePost(ctx context.Context, ins mage.RequestInputs, out *mage.ResponseOutput) mage.Redirect {
	renderer := mage.JSONRenderer{}
	out.Renderer = &renderer

	// checks the provided credentials. If correct creates a token, saves the user and returns the token
	errs := Errors{}
	nick := validators.NewField("username", true, ins)
	username, err := nick.Value()
	if err != nil {
		errs.AddError("username", err)
	}

	password := validators.NewField("password", true, ins)
	password.AddValidator(validators.LenValidator{MinLen: 8})
	pwd, err := password.Value()
	if err != nil {
		errs.AddError("password", err)
	}

	if errs.HasErrors() {
		renderer.Data = errs
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	u := identity.User{}
	err = model.FromStringID(ctx, &u, username, nil)

	if err == datastore.ErrNoSuchEntity {
		return mage.Redirect{Status: http.StatusNotFound}
	}

	if err != nil {
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	if u.Password != identity.HashPassword(pwd, salt) {
		return mage.Redirect{Status: http.StatusNotFound}
	}

	token, err := u.GenerateToken()
	if err != nil {
		log.Errorf(ctx, "error generating token for user %s: %s", u.StringID(), err.Error())
	}

	renderer.Data = token

	return mage.Redirect{Status: http.StatusOK}
}

func (controller *TokenController) HandleDelete(ctx context.Context, ins mage.RequestInputs, out *mage.ResponseOutput) mage.Redirect {
	u := ctx.Value(identity.KeyUser)
	user, ok := u.(identity.User)
	if !ok {
		return mage.Redirect{Status: http.StatusUnauthorized}
	}

	user.Token = ""
	err := model.Update(ctx, &user)
	if err != nil {
		return mage.Redirect{Status: http.StatusInternalServerError}
	}
	return mage.Redirect{Status:http.StatusOK}
}

func (controller *TokenController) OnDestroy(ctx context.Context) {}

// Controller to handle users
type UserController struct {}

func (controller *UserController) OnDestroy(ctx context.Context) {}

func (controller *UserController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	ins := mage.InputsFromContext(ctx)
	method := ins[mage.KeyRequestMethod].Value()
	switch method {
	case http.MethodPost:
		// check if we have an enabled user or a guser admin
		guser := user.Current(ctx)
		u := ctx.Value(identity.KeyUser)

		// we neither have a
		if guser == nil && u == nil {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		if guser != nil && !guser.Admin {
			return mage.Redirect{Status: http.StatusForbidden}
		}

		// we at least have an user.
		if me, ok := u.(identity.User); (guser == nil && !ok) || (ok && !me.HasPermission(identity.PermissionCreateUser)) {
			return mage.Redirect{Status: http.StatusForbidden}
		}

		log.Debugf(ctx, "guser is %v. user is %v", guser, u)

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
		err = model.FromStringID(ctx, &identity.User{}, username, nil)
		if err == nil {
			// user already exists
			msg := fmt.Sprintf("user %s already exists.", username)
			log.Errorf(ctx, msg)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusConflict}
		}

		if err != datastore.ErrNoSuchEntity {
			// user already exists
			msg := fmt.Sprintf("error retrieving user with username %s: %s.", username, err)
			log.Errorf(ctx, msg)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		pf := validators.NewRawField("password", true, meta.Password)
		pf.AddValidator(validators.LenValidator{MinLen:8})

		if err = pf.Validate(); err != nil {
			log.Errorf(ctx, "invalid password %s for username %s", meta.Password, username)
			errs.AddError(pf.Name, err)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		newuser := identity.User{}
		err = json.Unmarshal([]byte(jdata), &newuser)
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

		newuser.AddPermission(identity.PermissionLogIn)
		err = model.CreateWithOptions(ctx, &newuser, &opts)
		if err != nil {
			log.Errorf(ctx, "error creating user %s: %s", username, err)
			errs.AddError("", err)
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = &newuser
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusCreated}
	}
	return mage.Redirect{Status: http.StatusMethodNotAllowed}
}
