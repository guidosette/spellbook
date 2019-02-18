package page

import (
	"distudio.com/mage"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"sync"
)

const salt = "AnticmS"

var instance *website
var once sync.Once

type website struct {
	mage.Application
	Router  InternationalRouter
	options Options
}

type Options struct {
	Languages []language.Tag
}

//singleton instance
func application() *website {

	once.Do(func() {
		instance = &website{}
	})

	return instance
}

func NewWebsite(opts *Options) *website {
	ws := application()
	if opts != nil {
		if opts.Languages != nil {
			// create the language matcher
			ws.Router = NewInternationalRouter()
			ws.Router.matcher = language.NewMatcher(opts.Languages)
			ws.options = *opts
		}
	}
	return ws
}

func (app website) OnStart(ctx context.Context) context.Context {
	return ctx
}

func (app website) AfterResponse(ctx context.Context) {}
