package main

import (
	"context"
	"decodica.com/flamel"
	"cloud.google.com/go/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/search"
	"net/http"
)

type CleanController struct{}

func (controller *CleanController) OnDestroy(ctx context.Context) {}

func (controller *CleanController) Process(ctx context.Context, out *flamel.ResponseOutput) flamel.HttpResponse {
	log.Infof(ctx, "CleanController")
	ins := flamel.InputsFromContext(ctx)
	method := ins[flamel.KeyRequestMethod]
	_, ok := ins["entity"]
	if !ok {
		//log.Errorf(ctx, "unable to perform request. Entity not specified")
		//return mage.Redirect{Status: http.StatusBadRequest}
	}

	//entity := ins["entity"].Value()
	entity := "Attachment"

	switch method.Value() {
	case http.MethodGet:
		idx, err := search.Open(entity)
		if err != nil {
			return flamel.HttpResponse{Status: http.StatusNotFound}
		}

		opts := search.ListOptions{}
		opts.IDsOnly = true

		counter := 0

		for it := idx.List(ctx, &opts); ; {
			k, e := it.Next(nil)

			if e == search.Done {
				break
			}

			key, err := datastore.DecodeKey(k)

			if err != nil {
				log.Errorf(ctx, "error decoding key %s: %s", k, err.Error())
				continue
			}

			if err := datastore.Get(ctx, key, struct{}{}); err == datastore.ErrNoSuchEntity {
				if err = idx.Delete(ctx, k); err != nil {
					log.Errorf(ctx, "error deleting entity of type %s with key %s: %s", entity, k, err.Error())
					continue
				}
				log.Debugf(ctx, "index %s of entity of type %s has been removed", k, entity)
				counter++
			}
		}

		log.Infof(ctx, "Removed index of %d products", counter)

		return flamel.HttpResponse{Status: http.StatusOK}
	}
	return flamel.HttpResponse{Status: http.StatusNotImplemented}
}
