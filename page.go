package page

import (
	"distudio.com/mage"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
)

const salt = "AnticmS"

type Website struct {
	mage.Application
	Router InternationalRouter
}

type Options struct {
	Languages []language.Tag
}

func NewWebsite(opts *Options) *Website {
	ws := Website{}
	if opts != nil {
		if opts.Languages != nil {
			// create the language matcher
			ws.Router = NewInternationalRouter()
			ws.Router.matcher = language.NewMatcher(opts.Languages)
		}
	}
	return &ws
}

func (app Website) OnStart(ctx context.Context) context.Context {
	return ctx
}

func (app Website) AfterResponse(ctx context.Context) {}
