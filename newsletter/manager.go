package newsletter

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"errors"
	"google.golang.org/appengine/log"
	"reflect"
	"sort"
)

type Manager struct{}

func (manager Manager) NewResource(ctx context.Context) (page.Resource, error) {
	return &Newsletter{}, nil
}

func (manager Manager) FromId(ctx context.Context, id string) (page.Resource, error) {
	// todo permission?
	//current, _ := ctx.Value(identity.KeyUser).(identity.User)
	//if !current.HasPermission(identity.PermissionReadContent) {
	//	return nil, resource.NewPermissionError(identity.PermissionReadContent)
	//}

	att := Newsletter{}
	if err := model.FromStringID(ctx, &att, id, nil); err != nil {
		log.Errorf(ctx, "could not retrieve newsletter %s: %s", id, err.Error())
		return nil, err
	}

	return &att, nil
}

func (manager Manager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {
	// todo permission?
	//current, _ := ctx.Value(identity.KeyUser).(identity.User)
	//if !current.HasPermission(identity.PermissionReadContent) {
	//	return nil, resource.NewPermissionError(identity.PermissionReadContent)
	//}

	var newsletters []*Newsletter
	q := model.NewQuery(&Newsletter{})
	q = q.OffsetBy(opts.Page * opts.Size)

	if opts.Order != "" {
		dir := model.ASC
		if opts.Descending {
			dir = model.DESC
		}
		q = q.OrderBy(opts.Order, dir)
	}

	if opts.FilterField != "" {
		q = q.WithField(opts.FilterField+" =", opts.FilterValue)
	}

	// get one more so we know if we are done
	q = q.Limit(opts.Size + 1)
	err := q.GetMulti(ctx, &newsletters)
	if err != nil {
		return nil, err
	}

	resources := make([]page.Resource, len(newsletters))
	for i := range newsletters {
		resources[i] = page.Resource(newsletters[i])
	}

	return resources, nil
}

func (manager Manager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	// todo permission?
	//current, _ := ctx.Value(identity.KeyUser).(identity.User)
	//if !current.HasPermission(identity.PermissionReadContent) {
	//	return nil, resource.NewPermissionError(identity.PermissionReadContent)
	//}

	a := []string{"Email"} // list property accepted
	name := opts.Property

	i := sort.Search(len(a), func(i int) bool { return name <= a[i] })
	if i < len(a) && a[i] == name {
		// found
	} else {
		return nil, errors.New("no property found")
	}

	var conts []*Newsletter
	q := model.NewQuery(&Newsletter{})
	q = q.OffsetBy(opts.Page * opts.Size)

	if opts.Order != "" {
		dir := model.ASC
		if opts.Descending {
			dir = model.DESC
		}
		q = q.OrderBy(opts.Order, dir)
	}

	if opts.FilterField != "" {
		q = q.WithField(opts.FilterField+" =", opts.FilterValue)
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

func (manager Manager) Save(ctx context.Context, res page.Resource) error {
	// todo permission?
	//current, _ := ctx.Value(identity.KeyUser).(identity.User)
	//if !current.HasPermission(identity.PermissionEditContent) {
	//	return resource.NewPermissionError(identity.PermissionEditContent)
	//}

	newsletter := res.(*Newsletter)
	opts := model.CreateOptions{}
	opts.WithStringId(newsletter.Email)

	err := model.CreateWithOptions(ctx, newsletter, &opts)
	if err != nil {
		log.Errorf(ctx, "error creating post %s: %s", newsletter.Name, err)
		return err
	}

	return nil
}

func (manager Manager) Delete(ctx context.Context, res page.Resource) error {
	// todo permission?
	//current, _ := ctx.Value(identity.KeyUser).(identity.User)
	//if !current.HasPermission(identity.PermissionEditContent) {
	//	return resource.NewPermissionError(identity.PermissionEditContent)
	//}

	newsletter := res.(*Newsletter)
	err := model.Delete(ctx, newsletter, nil)
	if err != nil {
		log.Errorf(ctx, "error deleting newsletter %s: %s", newsletter.Name, err.Error())
		return err
	}

	return nil
}
