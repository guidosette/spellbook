package content

import (
	"cloud.google.com/go/storage"
	"context"
	"decodica.com/flamel"
	"decodica.com/spellbook"
	"google.golang.org/appengine/log"
	"net/http"
)

type fileHandler struct {
	spellbook.BaseRestHandler
}

func (handler fileHandler) HandlePost(ctx context.Context, out *flamel.ResponseOutput) flamel.HttpResponse {

	res, err := handler.Manager.NewResource(ctx)
	if err != nil {
		return handler.ErrorToStatus(ctx, err, out)
	}

	if err = handler.Manager.Create(ctx, res, nil); err != nil {
		return handler.ErrorToStatus(ctx, err, out)
	}

	renderer := flamel.JSONRenderer{}
	renderer.Data = res
	out.Renderer = &renderer
	return flamel.HttpResponse{Status: http.StatusCreated}
}

// Converts an error to its equivalent HTTP representation
func (handler fileHandler) ErrorToStatus(ctx context.Context, err error, out *flamel.ResponseOutput) flamel.HttpResponse {
	if err == storage.ErrObjectNotExist {
		log.Errorf(ctx, "%s", err.Error())
		return flamel.HttpResponse{Status: http.StatusNotFound}
	}
	return handler.BaseRestHandler.ErrorToStatus(ctx, err, out)
}
