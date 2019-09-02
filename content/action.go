package content

import (
	"context"
	"decodica.com/spellbook"

	"encoding/json"
)

type Action spellbook.SupportedAction

func (action *Action) UnmarshalJSON(data []byte) error {

	alias := struct {
		Name     string               `json:"name"`
		Endpoint string               `json:"endpoint"`
		Type     spellbook.ActionType `json:"type"`
		Method   string               `json:"method"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	action.Type = alias.Type
	action.Name = alias.Name
	action.Endpoint = alias.Endpoint
	action.Method = alias.Method

	return nil
}

func (action *Action) MarshalJSON() ([]byte, error) {
	alias := struct {
		Name     string               `json:"name"`
		Endpoint string               `json:"endpoint"`
		Type     spellbook.ActionType `json:"type"`
		Method   string               `json:"method"`
	}{action.Name, action.Endpoint, action.Type, action.Method}

	return json.Marshal(&alias)
}

/**
* Resource representation
 */
func (action *Action) Id() string {
	return action.Name
}

func (action *Action) FromRepresentation(rtype spellbook.RepresentationType, data []byte) error {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Unmarshal(data, action)
	}
	return spellbook.NewUnsupportedError()
}

func (action *Action) ToRepresentation(rtype spellbook.RepresentationType) ([]byte, error) {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Marshal(action)
	}
	return nil, spellbook.NewUnsupportedError()
}

func NewActionController() *spellbook.RestController {
	man := actionManager{}
	return spellbook.NewRestController(spellbook.BaseRestHandler{Manager: man})
}

/*
* Action manager
 */

type actionManager struct{}

func (manager actionManager) NewResource(ctx context.Context) (spellbook.Resource, error) {
	return nil, spellbook.NewUnsupportedError()
}

func (manager actionManager) FromId(ctx context.Context, id string) (spellbook.Resource, error) {
	return nil, spellbook.NewUnsupportedError()
}

func (manager actionManager) ListOf(ctx context.Context, opts spellbook.ListOptions) ([]spellbook.Resource, error) {

	ws := spellbook.Application()

	actions := ws.Options().Actions

	from := opts.Page * opts.Size
	if from > len(actions) {
		return make([]spellbook.Resource, 0), nil
	}

	to := from + opts.Size
	if to > len(actions) {
		to = len(actions)
	}

	items := actions[from:to]
	resources := make([]spellbook.Resource, len(items))

	for i := range items {
		action := Action(items[i])
		resources[i] = &action
	}

	return resources, nil
}

func (manager actionManager) ListOfProperties(ctx context.Context, opts spellbook.ListOptions) ([]string, error) {
	return nil, spellbook.NewUnsupportedError()
}

func (manager actionManager) Create(ctx context.Context, res spellbook.Resource, bundle []byte) error {
	return spellbook.NewUnsupportedError()
}

func (manager actionManager) Update(ctx context.Context, res spellbook.Resource, bundle []byte) error {
	return spellbook.NewUnsupportedError()
}

func (manager actionManager) Delete(ctx context.Context, res spellbook.Resource) error {
	return spellbook.NewUnsupportedError()
}
