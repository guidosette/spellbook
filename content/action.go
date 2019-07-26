package content

import (
	"context"
	"distudio.com/page"
	"encoding/json"
)

type Action page.SupportedAction

func (action *Action) UnmarshalJSON(data []byte) error {

	alias := struct {
		Name     string          `json:"name"`
		Endpoint string          `json:"endpoint"`
		Type     page.ActionType `json:"type"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	action.Type = alias.Type
	action.Name = alias.Name
	action.Endpoint = alias.Endpoint

	return nil
}

func (action *Action) MarshalJSON() ([]byte, error) {
	alias := struct {
		Name     string          `json:"name"`
		Endpoint string          `json:"endpoint"`
		Type     page.ActionType `json:"type"`
	}{action.Name, action.Endpoint, action.Type}

	return json.Marshal(&alias)
}

/**
* Resource representation
 */

func (action *Action) Id() string {
	return action.Name
}

func (action *Action) FromRepresentation(rtype page.RepresentationType, data []byte) error {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Unmarshal(data, action)
	}
	return page.NewUnsupportedError()
}

func (action *Action) ToRepresentation(rtype page.RepresentationType) ([]byte, error) {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Marshal(action)
	}
	return nil, page.NewUnsupportedError()
}

func NewActionController() *page.RestController {
	man := actionManager{}
	return page.NewRestController(page.BaseRestHandler{Manager: man})
}

/*
* Action manager
 */

type actionManager struct{}

func (manager actionManager) NewResource(ctx context.Context) (page.Resource, error) {
	return nil, page.NewUnsupportedError()
}

func (manager actionManager) FromId(ctx context.Context, id string) (page.Resource, error) {
	return nil, page.NewUnsupportedError()
}

func (manager actionManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {

	ws := page.Application()

	actions := ws.Options().Actions

	from := opts.Page * opts.Size
	if from > len(actions) {
		return make([]page.Resource, 0), nil
	}

	to := from + opts.Size
	if to > len(actions) {
		to = len(actions)
	}

	items := actions[from:to]
	resources := make([]page.Resource, len(items))

	for i := range items {
		action := Action(items[i])
		resources[i] = page.Resource(&action)
	}

	return resources, nil
}

func (manager actionManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	return nil, page.NewUnsupportedError()
}

func (manager actionManager) Create(ctx context.Context, res page.Resource, bundle []byte) error {
	return page.NewUnsupportedError()
}

func (manager actionManager) Update(ctx context.Context, res page.Resource, bundle []byte) error {
	return page.NewUnsupportedError()
}

func (manager actionManager) Delete(ctx context.Context, res page.Resource) error {
	return page.NewUnsupportedError()
}
