package identity

import (
	"context"
	"distudio.com/mage"
	"distudio.com/mage/model"
)

const (
	HeaderToken string = "X-Authorization"
	KeyUser string = "__pUser__"
)

type UserAuthenticator struct{
	mage.Authenticator
}

func (authenticator UserAuthenticator) Authenticate(ctx context.Context) context.Context {
	inputs := mage.InputsFromContext(ctx)
	if tkn, ok := inputs[HeaderToken]; ok {
		token := tkn.Value()
		// grab the last chars after hashLenght
		encoded := token[hashLen:]
		u := User{}
		err := model.FromEncodedKey(ctx, &u, encoded)
		if err != nil || u.Token != token {
			return ctx
		}

		return context.WithValue(ctx, KeyUser, u)
	}

	return ctx
}


