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
	"net/http"
)

// This controller is responsible for dispensing new tokens
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


	j, ok := ins[mage.KeyRequestJSON]

	if !ok {
		return mage.Redirect{Status:http.StatusBadRequest}
	}

	credentials := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} {}

	err := json.Unmarshal([]byte(j.Value()), &credentials)
	if err != nil {
		return mage.Redirect{Status:http.StatusBadRequest}
	}

	// checks the provided credentials. If correct creates a token, saves the user and returns the token
	errs := Errors{}
	nick := validators.NewRawField("username", true, credentials.Username)
	username, err := nick.Value()
	if err != nil {
		errs.AddError("username", err)
	}

	password := validators.NewRawField("password", true, credentials.Password)
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

	u.Token, err = u.GenerateToken()
	if err != nil {
		log.Errorf(ctx, "error generating token for user %s: %s", u.StringID(), err.Error())
	}

	err = model.Update(ctx, &u)
	if err != nil {
		log.Errorf(ctx, "error updating user token: %s", err.Error())
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	renderer.Data = u.Token

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

// identity controller, used to operate on the current user
type IdentityController struct {}

func (controller *IdentityController) OnDestroy(ctx context.Context) {}

func (controller *IdentityController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	u := ctx.Value(identity.KeyUser)
	me, ok := u.(identity.User)
	if !ok {
		log.Debugf(ctx, "user is not an user: +v", u)
		return mage.Redirect{Status: http.StatusUnauthorized}
	}

	ins := mage.InputsFromContext(ctx)
	switch ins[mage.KeyRequestMethod].Value() {
	case http.MethodGet:
		renderer := mage.JSONRenderer{}
		renderer.Data = &me
		out.Renderer = &renderer
		return mage.Redirect{Status:http.StatusOK}
	}

	return mage.Redirect{Status:http.StatusNotImplemented}
}

// This controller handles user's CRUD operations
// Before performing each operation test if the user requesting the method
// holds permissions relating to user manipulations
type UserController struct {}

func (controller *UserController) OnDestroy(ctx context.Context) {}

func (controller *UserController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	ins := mage.InputsFromContext(ctx)
	method := ins[mage.KeyRequestMethod].Value()
	switch method {
	case http.MethodPost:

		u := ctx.Value(identity.KeyUser)
		// we at least have an user.
		me, ok := u.(identity.User)
		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		if !me.HasPermission(identity.PermissionCreateUser) {
			return mage.Redirect{Status: http.StatusForbidden}
		}

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

		// retrieve the other fields
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

		newuser.Password = identity.HashPassword(meta.Password, salt)
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
	case http.MethodGet:
		// check if current user has permission
		me := ctx.Value(identity.KeyUser)
		current, ok := me.(identity.User)

		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		if !current.HasPermission(identity.PermissionReadUser) {
			return mage.Redirect{Status: http.StatusForbidden}
		}

		params := mage.RoutingParams(ctx)
		// try to get the username.
		// if there is no param then it is a list request
		param, ok := params["username"]
		if !ok {
			// handle pagination of users. todo
			return mage.Redirect{Status: http.StatusNotImplemented}
		}

		username := param.Value()
		item := identity.User{}
		err := model.FromStringID(ctx, &item, username, nil)
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status: http.StatusNotFound}
		}

		if err != nil {
			log.Errorf(ctx, "error retrieving user %s: %s", username, err.Error())
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
		param, ok := params["username"]
		if !ok {
			return mage.Redirect{Status:http.StatusBadRequest}
		}

		j, ok := ins[mage.KeyRequestJSON]
		if !ok {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// handle the json request
		jdata := j.Value()

		juser := struct {
			*identity.User
			Password string `json:"password"`
		}{User:&identity.User{}}

		err := json.Unmarshal([]byte(jdata), &juser)
		if err != nil {
			log.Errorf(ctx, "malformed json: %s", err.Error())
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// retrieve the user
		username := param.Value()
		target := identity.User{}
		err = model.FromStringID(ctx, &target, username, nil)
		if err == datastore.ErrNoSuchEntity {
			return mage.Redirect{Status:http.StatusNotFound}
		}

		if err != nil {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// check if the password is correct in case the user supplied it
		errs := Errors{}

		if juser.Password != "" {
			pf := validators.NewRawField("password", true, juser.Password)
			pf.AddValidator(validators.LenValidator{MinLen:8})

			if err = pf.Validate(); err != nil {
				err = fmt.Errorf("invalid password %s for username %s", juser.Password, username)
				errs.AddError(pf.Name, err)
				renderer := mage.JSONRenderer{}
				renderer.Data = errs
				out.Renderer = &renderer
				return mage.Redirect{Status: http.StatusBadRequest}
			}
			target.Password = juser.Password
		}

		if juser.Email != "" {
			ef := validators.NewRawField("email", true, juser.Email)
			if err = ef.Validate(); err != nil {
				err := fmt.Errorf("invalid email address: %s", juser.Email)
				errs.AddError(ef.Name, err)
				renderer := mage.JSONRenderer{}
				renderer.Data = errs
				out.Renderer = &renderer
				return mage.Redirect{Status: http.StatusBadRequest}
			}
			target.Email = juser.Email
		}

		target.Name = juser.Name
		target.Surname = juser.Surname
		target.Permission = juser.Permission

		err = model.Update(ctx, &target)
		if err != nil {
			log.Errorf(ctx, "error updating user %s: %s", username, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = &target
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}
	}
	return mage.Redirect{Status: http.StatusMethodNotAllowed}
}
