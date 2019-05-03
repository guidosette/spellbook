package identity

import (
	"context"
	"distudio.com/mage"
	"distudio.com/page"
	"net/http"
)

type Token string

func (token Token) Id() string {
	return string(token)
}

func (token Token) Create(ctx context.Context) error {
	j, ok := ins[mage.KeyRequestJSON]

	if !ok {
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	credentials := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
}

func (token Token) Update(ctx context.Context, res page.Resource) error {

}

