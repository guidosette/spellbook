package page

import (
	"context"
	"distudio.com/mage"
	"distudio.com/mage/model"
	"distudio.com/page/identity"
	"distudio.com/page/validators"
	"fmt"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
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
	nick := validators.NewField("nickname", true, ins)
	nickname, err := nick.Value()
	if err != nil {
		errs.AddError("nickname", err)
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
	err = model.FromStringID(ctx, &u, nickname, nil)

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
		// creates an user if doesn't exists
		u := ctx.Value(identity.KeyUser)
		user, ok := u.(identity.User)
		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		if !user.HasPermission(identity.PermissionCreateUser) {
			return mage.Redirect{Status: http.StatusForbidden}
		}

		// user can create a new user, check if data are correct
		errs := Errors{}
		nick := validators.NewField("nickname", true, ins)
		nickname, err := nick.Value()
		if err != nil {
			errs.AddError("nickname", err)
		}

		// if we have the nickname, check if the user does not exists
		newuser := identity.User{}

		err = model.FromStringID(ctx, &newuser, nickname, nil)
		if err != nil {
			// clear the errors
			errs.Clear()
			errs.AddError("nickname", fmt.Errorf("user with nickname %s already exists", nickname))
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusConflict}
		}

		// the user does not exist. Go on and create it
		password := validators.NewField("password", true, ins)
		password.AddValidator(validators.LenValidator{MinLen: 8})
		pwd, err := password.Value()
		if err != nil {
			errs.AddError("password", err)
		}

		newuser.Password = pwd

		if name, ok := ins["name"]; ok {
			newuser.Name = name.Value()
		}

		if surname, ok := ins["surname"]; ok {
			newuser.Surname = surname.Value()
		}

		if locale, ok := ins["locale"]; ok {
			// locale will be evaluated later, after language is resolved.
			// if the locale is not valid or is not supported, the default locale will be used
			newuser.Locale = locale.Value()
		}

		if errs.HasErrors() {
			renderer := mage.JSONRenderer{}
			renderer.Data = errs
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// else, create the new user
		opts := model.CreateOptions{}
		opts.WithStringId(nickname)
		err = model.CreateWithOptions(ctx, &newuser, &opts)
		if err != nil {
			errs.Clear()
			errs.AddError("", err)
			return mage.Redirect{Status:http.StatusInternalServerError}
		}

		return mage.Redirect{Status: http.StatusCreated}
	}
	return mage.Redirect{Status: http.StatusMethodNotAllowed}
}