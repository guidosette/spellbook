package page

import (
	"context"
	"distudio.com/mage"
	"distudio.com/page"
	"sync"
)

var once sync.Once
var instance *Website

type Website struct {
	mage.Application
	Router  page.InternationalRouter
	options page.Options
}

//singleton instance
func Application() *Website {

	once.Do(func() {
		instance = &Website{}
	})

	return instance
}

func (app Website) OnStart(ctx context.Context) context.Context {
	return ctx
}

func (app Website) AfterResponse(ctx context.Context) {}

func (app *Website) SetOptions(opts page.Options) {
	app.options = opts
}

func (app Website) Options() page.Options {
	return app.options
}
