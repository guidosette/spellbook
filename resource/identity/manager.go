package identity

//import (
//	"context"
//	"distudio.com/mage/model"
//	"distudio.com/page/resource"
//	"errors"
//	"google.golang.org/appengine/log"
//	"reflect"
//	"sort"
//)
//
//type Manager struct{}
//
//func (manager Manager) NewResource(ctx context.Context) (resource.Resource, error) {
//	return &User{}, nil
//}
//
//func (manager Manager) FromId(ctx context.Context, id string) (resource.Resource, error) {
//	current, _ := ctx.Value(KeyUser).(User)
//	if !current.HasPermission(PermissionReadUser) {
//		return nil, resource.NewPermissionError(PermissionReadUser)
//	}
//
//	att := User{}
//	if err := model.FromStringID(ctx, &att, id, nil); err != nil {
//		log.Errorf(ctx, "could not retrieve user %s: %s", id, err.Error())
//		return nil, err
//	}
//
//	return &att, nil
//}
//
//func (manager Manager) ListOf(ctx context.Context, opts resource.ListOptions) ([]resource.Resource, error) {
//	current, _ := ctx.Value(KeyUser).(User)
//	if !current.HasPermission(PermissionReadUser) {
//		return nil, resource.NewPermissionError(PermissionReadUser)
//	}
//
//	var users []*User
//	q := model.NewQuery(&User{})
//	q = q.OffsetBy(opts.Page * opts.Size)
//
//	if opts.Order != "" {
//		dir := model.ASC
//		if opts.Descending {
//			dir = model.DESC
//		}
//		q = q.OrderBy(opts.Order, dir)
//	}
//
//	if opts.FilterField != "" {
//		q = q.WithField(opts.FilterField+" =", opts.FilterValue)
//	}
//
//	// get one more so we know if we are done
//	q = q.Limit(opts.Size + 1)
//	err := q.GetMulti(ctx, &users)
//	if err != nil {
//		return nil, err
//	}
//
//	resources := make([]resource.Resource, len(users))
//	for i := range users {
//		resources[i] = resource.Resource(users[i])
//	}
//
//	return resources, nil
//}
//
//func (manager Manager) ListOfProperties(ctx context.Context, opts resource.ListOptions) ([]string, error) {
//	current, _ := ctx.Value(KeyUser).(User)
//	if !current.HasPermission(PermissionReadUser) {
//		return nil, resource.NewPermissionError(PermissionReadUser)
//	}
//
//	a := []string{"group"} // list property accepted
//	name := opts.Property
//
//	i := sort.Search(len(a), func(i int) bool { return name <= a[i] })
//	if i < len(a) && a[i] == name {
//		// found
//	} else {
//		return nil, errors.New("no property found")
//	}
//
//	var conts []*User
//	q := model.NewQuery(&User{})
//	q = q.OffsetBy(opts.Page * opts.Size)
//
//	if opts.Order != "" {
//		dir := model.ASC
//		if opts.Descending {
//			dir = model.DESC
//		}
//		q = q.OrderBy(opts.Order, dir)
//	}
//
//	if opts.FilterField != "" {
//		q = q.WithField(opts.FilterField+" =", opts.FilterValue)
//	}
//
//	q = q.Distinct(name)
//	q = q.Limit(opts.Size + 1)
//	err := q.GetAll(ctx, &conts)
//	if err != nil {
//		log.Errorf(ctx, "Error retrieving result: %+v", err)
//		return nil, err
//	}
//	var result []string
//	for _, c := range conts {
//		value := reflect.ValueOf(c).Elem().FieldByName(name).String()
//		if len(value) > 0 {
//			result = append(result, value)
//		}
//	}
//	return result, nil
//}
//
//func (manager Manager) Save(ctx context.Context, res resource.Resource) error {
//	current, _ := ctx.Value(KeyUser).(User)
//	if !current.HasPermission(PermissionEditUser) {
//		return resource.NewPermissionError(PermissionEditUser)
//	}
//
//	user := res.(*User)
//	opts := model.CreateOptions{}
//	opts.WithStringId(user.Username())
//
//	err := model.CreateWithOptions(ctx, user, &opts)
//	if err != nil {
//		log.Errorf(ctx, "error creating post %s: %s", user.Name, err)
//		return err
//	}
//
//	return nil
//}
//
//func (manager Manager) Delete(ctx context.Context, res resource.Resource) error {
//	current, _ := ctx.Value(KeyUser).(User)
//	if !current.HasPermission(PermissionEditUser) {
//		return resource.NewPermissionError(PermissionEditUser)
//	}
//
//	user := res.(*User)
//	err := model.Delete(ctx, user, nil)
//	if err != nil {
//		log.Errorf(ctx, "error deleting user %s: %s", user.Name, err.Error())
//		return err
//	}
//
//	return nil
//}
