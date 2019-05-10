package content

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"distudio.com/page/identity"
	"errors"
	"google.golang.org/appengine/log"
	"reflect"
	"sort"
)

func NewFileController() *page.RestController {
	handler := page.BaseRestHandler{Manager: fileManager{}}
	return page.NewRestController(handler)
}

type fileManager struct{}

func (manager fileManager) NewResource(ctx context.Context) (page.Resource, error) {
	return &File{}, nil
}

func (manager fileManager) FromId(ctx context.Context, id string) (page.Resource, error) {

	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionReadContent) {
		return nil, page.NewPermissionError(identity.PermissionName(identity.PermissionReadContent))
	}

	att := File{}
	if err := model.FromStringID(ctx, &att, id, nil); err != nil {
		log.Errorf(ctx, "could not retrieve file %s: %s", id, err.Error())
		return nil, err
	}

	return &att, nil
}

func (manager fileManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {

	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionReadContent) {
		return nil, page.NewPermissionError(identity.PermissionName(identity.PermissionReadContent))
	}

	var files []*File
	q := model.NewQuery(&File{})
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
	err := q.GetMulti(ctx, &files)
	if err != nil {
		return nil, err
	}

	resources := make([]page.Resource, len(files))
	for i := range files {
		resources[i] = page.Resource(files[i])
	}

	return resources, nil
}

func (manager fileManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionReadContent) {
		return nil, page.NewPermissionError(identity.PermissionName(identity.PermissionReadContent))
	}

	a := []string{"Name"} // list property accepted
	name := opts.Property

	i := sort.Search(len(a), func(i int) bool { return name <= a[i] })
	if i < len(a) && a[i] == name {
		// found
	} else {
		return nil, errors.New("no property found")
	}

	var conts []*File
	q := model.NewQuery(&File{})
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

func (manager fileManager) Save(ctx context.Context, res page.Resource) error {
	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionEditContent) {
		return page.NewPermissionError(identity.PermissionName(identity.PermissionEditContent))
	}

	file := res.(*File)

	err := model.Create(ctx, file)
	if err != nil {
		log.Errorf(ctx, "error creating post %s: %s", file.Name, err)
		return err
	}

	return nil
}

func (manager fileManager) Delete(ctx context.Context, res page.Resource) error {

	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionEditContent) {
		return page.NewPermissionError(identity.PermissionName(identity.PermissionEditContent))
	}

	file := res.(*File)
	err := model.Delete(ctx, file, nil)
	if err != nil {
		log.Errorf(ctx, "error deleting file %s: %s", file.Name, err.Error())
		return err
	}

	return nil
}
