package content

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"errors"
	"fmt"
	"google.golang.org/appengine/log"
	"reflect"
	"sort"
	"strconv"
	"time"
)

func NewPlaceController() *page.RestController {
	return NewPlaceControllerWithKey("")
}

func NewPlaceControllerWithKey(key string) *page.RestController {
	man := placeManager{}
	handler := page.BaseRestHandler{Manager: man}
	c := page.NewRestController(handler)
	c.Key = key
	return c
}

type placeManager struct{}

func (manager placeManager) NewResource(ctx context.Context) (page.Resource, error) {
	return &Place{}, nil
}

func (manager placeManager) FromId(ctx context.Context, strId string) (page.Resource, error) {

	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionReadPlace) {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadPlace))
	}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return nil, page.NewFieldError(strId, err)
	}

	att := Place{}
	if err := model.FromIntID(ctx, &att, id, nil); err != nil {
		log.Errorf(ctx, "could not retrieve place %s: %s", id, err.Error())
		return nil, err
	}

	return &att, nil
}

func (manager placeManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {
	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionReadPlace) {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadPlace))
	}

	var places []*Place
	q := model.NewQuery(&Place{})
	q = q.OffsetBy(opts.Page * opts.Size)

	if opts.Order != "" {
		dir := model.ASC
		if opts.Descending {
			dir = model.DESC
		}
		q = q.OrderBy(opts.Order, dir)
	}

	for _, filter := range opts.Filters {
		if filter.Field != "" {
			q = q.WithField(filter.Field+" =", filter.Value)
		}
	}

	// get one more so we know if we are done
	q = q.Limit(opts.Size + 1)
	err := q.GetMulti(ctx, &places)
	if err != nil {
		return nil, err
	}

	resources := make([]page.Resource, len(places))
	for i := range places {
		resources[i] = page.Resource(places[i])
	}

	return resources, nil
}

func (manager placeManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionReadPlace) {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadPlace))
	}

	a := []string{"Address"} // list property accepted
	name := opts.Property

	i := sort.Search(len(a), func(i int) bool { return name <= a[i] })
	if i < len(a) && a[i] == name {
		// found
	} else {
		return nil, errors.New("no property found")
	}

	var conts []*Place
	q := model.NewQuery(&Place{})
	q = q.OffsetBy(opts.Page * opts.Size)

	if opts.Order != "" {
		dir := model.ASC
		if opts.Descending {
			dir = model.DESC
		}
		q = q.OrderBy(opts.Order, dir)
	}

	for _, filter := range opts.Filters {
		if filter.Field != "" {
			q = q.WithField(filter.Field+" =", filter.Value)
		}
	}

	q = q.Distinct(name)
	q = q.Limit(opts.Size + 1)
	err := q.GetAll(ctx, &conts)
	if err != nil {
		log.Errorf(ctx, "Error retrieving result: %+v", err)
		return nil, err
	}
	var result []string
	for _, c := range conts {
		value := reflect.ValueOf(c).Elem().FieldByName(name).String()
		if len(value) > 0 {
			result = append(result, value)
		}
	}
	return result, nil
}

func (manager placeManager) Create(ctx context.Context, res page.Resource, bundle []byte) error {
	current := page.IdentityFromContext(ctx)
	if current == nil || !current.HasPermission(page.PermissionCreatePlace) {
		return page.NewPermissionError(page.PermissionName(page.PermissionCreatePlace))
	}

	place := res.(*Place)

	if place.Address == "" || !place.Position.Valid() {
		return page.NewFieldError("address", errors.New("address and position can't be empty"))
	}

	place.Created = time.Now().UTC()

	err := model.Create(ctx, place)
	if err != nil {
		log.Errorf(ctx, "error creating place %s: %s", place.Name, err)
		return err
	}

	return nil
}

func (manager placeManager) Update(ctx context.Context, res page.Resource, bundle []byte) error {
	current := page.IdentityFromContext(ctx)
	if current == nil || !current.HasPermission(page.PermissionEditPlace) {
		return page.NewPermissionError(page.PermissionName(page.PermissionEditPlace))
	}

	other := Place{}
	if err := other.FromRepresentation(page.RepresentationTypeJSON, bundle); err != nil {
		return page.NewFieldError("", fmt.Errorf("bad json %s", string(bundle)))
	}

	place := res.(*Place)
	place.Address = other.Address
	place.Description = other.Description
	place.Position = other.Position
	place.Phone = other.Phone
	place.Website = other.Website
	place.Updated = time.Now().UTC()

	if place.Address == "" || !place.Position.Valid() {
		return page.NewFieldError("address", errors.New("address and position can't be empty"))
	}

	if err := model.Update(ctx, place); err != nil {
		return fmt.Errorf("error updating place %s: %s", place.Address, err)
	}

	return nil
}

func (manager placeManager) Delete(ctx context.Context, res page.Resource) error {
	current := page.IdentityFromContext(ctx)
	if current == nil || !current.HasPermission(page.PermissionEditPlace) {
		return page.NewPermissionError(page.PermissionName(page.PermissionEditPlace))
	}

	place := res.(*Place)
	err := model.Delete(ctx, place, nil)
	if err != nil {
		log.Errorf(ctx, "error deleting place %s: %s", place.Name, err.Error())
		return err
	}

	return nil
}
