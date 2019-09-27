package identity

import (
	"context"
	"decodica.com/flamel"
	"decodica.com/spellbook"
	"decodica.com/spellbook/sql"
	"google.golang.org/appengine/user"
)

type SqlAuthenticator struct {
	flamel.Authenticator
}

func (authenticator SqlAuthenticator) Authenticate(ctx context.Context) context.Context {
	inputs := flamel.InputsFromContext(ctx)
	if tkn, ok := inputs[spellbook.HeaderToken]; ok {
		token := tkn.Value()
		// grab the last chars after hashLength
		u := User{}
		db := sql.FromContext(ctx)
		err := db.Where("token = ?", token).First(&u).Error
		if err != nil {
			return ctx
		}

		if !u.IsEnabled() {
			return ctx
		}

		return spellbook.ContextWithIdentity(ctx, u)
	}

	return ctx
}

type SqlGSupportAuthenticator struct {
	flamel.Authenticator
}

func (authenticator SqlGSupportAuthenticator) Authenticate(ctx context.Context) context.Context {
	guser := user.Current(ctx)
	if guser == nil {
		// try with the canonical authenticator
		ua := SqlAuthenticator{}
		return ua.Authenticate(ctx)
	}

	// else populate a flamel user with usable fields
	u := User{}
	u.gUser = guser
	u.Email = guser.Email
	// if admin, grant all permissions
	u.GrantAll()
	return spellbook.ContextWithIdentity(ctx, u)
}
