package identity

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"errors"
	"fmt"
	"google.golang.org/appengine/log"
	"reflect"
	"sort"
)

type UserManager struct{}

func (manager UserManager) NewResource(ctx context.Context) (page.Resource, error) {
	return &User{}, nil
}

func (manager UserManager) FromId(ctx context.Context, id string) (page.Resource, error) {

	current, _ := ctx.Value(KeyUser).(User)
	if !current.HasPermission(PermissionReadUser) {
		return nil, page.NewPermissionError(PermissionName(PermissionReadUser))
	}

	att := User{}
	if err := model.FromStringID(ctx, &att, id, nil); err != nil {
		log.Errorf(ctx, "could not retrieve user %s: %s", id, err.Error())
		return nil, err
	}

	return &att, nil
}

func (manager UserManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {

	current, _ := ctx.Value(KeyUser).(User)
	if !current.HasPermission(PermissionReadUser) {
		return nil, page.NewPermissionError(PermissionName(PermissionReadUser))
	}

	var users []*User
	q := model.NewQuery(&User{})
	q = q.OffsetBy(opts.Page * opts.Size)

	if opts.Order != "" {
		dir := model.ASC
		if opts.Descending {
			dir = model.DESC
		}
		q = q.OrderBy(opts.Order, dir)
	}

	if opts.FilterField != "" {
		q = q.WithField(fmt.Sprintf("%s =", opts.FilterField), opts.FilterValue)
	}

	// get one more so we know if we are done
	q = q.Limit(opts.Size + 1)
	err := q.GetMulti(ctx, &users)
	if err != nil {
		return nil, err
	}

	resources := make([]page.Resource, len(users))
	for i := range users {
		resources[i] = page.Resource(users[i])
	}

	return resources, nil
}

func (manager UserManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	current, _ := ctx.Value(KeyUser).(User)
	if !current.HasPermission(PermissionReadUser) {
		return nil, page.NewPermissionError(PermissionName(PermissionReadUser))
	}

	a := []string{"group"}
	name := opts.Property

	i := sort.Search(len(a), func(i int) bool { return name <= a[i] })
	if i < len(a) && a[i] == name {
		// found
	} else {
		return nil, errors.New("no property found")
	}

	var conts []*User
	q := model.NewQuery(&User{})
	q = q.OffsetBy(opts.Page * opts.Size)

	if opts.Order != "" {
		dir := model.ASC
		if opts.Descending {
			dir = model.DESC
		}
		q = q.OrderBy(opts.Order, dir)
	}

	if opts.FilterField != "" {
		q = q.WithField(fmt.Sprintf("%s =", opts.FilterField), opts.FilterValue)
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

func (manager UserManager) Save(ctx context.Context, res page.Resource) error {

	current, _ := ctx.Value(KeyUser).(User)
	if !current.HasPermission(PermissionEditUser) {
		return page.NewPermissionError(PermissionName(PermissionEditUser))
	}

	user := res.(*User)
	opts := model.CreateOptions{}
	opts.WithStringId(user.Username())

	err := model.CreateWithOptions(ctx, user, &opts)
	if err != nil {
		log.Errorf(ctx, "error creating post %s: %s", user.Name, err)
		return err
	}

	return nil
}

func (manager UserManager) Delete(ctx context.Context, res page.Resource) error {
	current, _ := ctx.Value(KeyUser).(User)
	if !current.HasPermission(PermissionEditUser) {
		return page.NewPermissionError(PermissionName(PermissionEditUser))
	}

	user := res.(*User)
	err := model.Delete(ctx, user, nil)
	if err != nil {
		log.Errorf(ctx, "error deleting user %s: %s", user.Name, err.Error())
		return err
	}

	return nil
}
