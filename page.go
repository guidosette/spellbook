package page

import (
	"distudio.com/page/internal/page"
	"golang.org/x/text/language"
)

const salt = "AnticmS"

type Options struct {
	Languages []language.Tag
}

func NewWebsite(opts *Options) *page.Website {
	ws := page.Application()
	if opts != nil {
		if opts.Languages != nil {
			// create the language matcher
			ws.Router = NewInternationalRouter()
			ws.Router.matcher = language.NewMatcher(opts.Languages)
			ws.SetOptions(*opts)
		}
	}
	return ws
}
