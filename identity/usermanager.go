package identity

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"reflect"
	"sort"
)

type userManager struct{}

func NewUserController() *page.RestController {
	return NewUserControllerWithKey("")
}

func NewUserControllerWithKey(key string) *page.RestController {
	manager := userManager{}
	handler := page.BaseRestHandler{Manager: manager}
	c := page.NewRestController(handler)
	c.Key = key
	return c
}

func (manager userManager) NewResource(ctx context.Context) (page.Resource, error) {
	return &User{}, nil
}

func (manager userManager) FromId(ctx context.Context, id string) (page.Resource, error) {
	current := page.IdentityFromContext(ctx)
	if current == nil || !current.HasPermission(page.PermissionReadUser) {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadUser))
	}

	user, ok := current.(User)
	if ok && id == user.Id() {
		return &user, nil
	}

	att := User{}
	if err := model.FromStringID(ctx, &att, id, nil); err != nil {
		log.Errorf(ctx, "could not retrieve user %s: %s", id, err.Error())
		return nil, err
	}

	return &att, nil
}

func (manager userManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {

	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionReadUser) {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadUser))
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

	for _, filter := range opts.Filters {
		if filter.Field != "" {
			q = q.WithField(filter.Field+" =", filter.Value)
		}
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

func (manager userManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionReadUser) {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadUser))
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

	for _, filter := range opts.Filters {
		if filter.Field != "" {
			q = q.WithField(filter.Field+" =", filter.Value)
		}
	}

	q = q.Distinct(name)
	q = q.Limit(opts.Size + 1)
	err := q.GetAll(ctx, &conts)
	if err != nil {
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

func (manager userManager) Create(ctx context.Context, res page.Resource, bundle []byte) error {

	current := page.IdentityFromContext(ctx)
	if current == nil || !current.HasPermission(page.PermissionCreateUser) {
		return page.NewPermissionError(page.PermissionName(page.PermissionCreateUser))
	}

	user := res.(*User)

	if user.StringID() != "" {
		return page.NewFieldError("username", fmt.Errorf("user already exists"))
	}

	meta := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := json.Unmarshal(bundle, &meta)
	if err != nil {
		return page.NewFieldError("json", fmt.Errorf("invalid json: %s", string(bundle)))
	}

	username := meta.Username
	username = SanitizeUserName(username)

	uf := page.NewRawField("username", true, username)
	uf.AddValidator(page.DatastoreKeyNameValidator{})

	// validate the username. Accepted values for the username are implementation dependent
	if err := uf.Validate(); err != nil {
		msg := fmt.Sprintf("invalid username %s", user.Username())
		return page.NewFieldError("username", errors.New(msg))
	}

	pf := page.NewRawField("password", true, meta.Password)
	pf.AddValidator(page.LenValidator{MinLen: 8})

	if err := pf.Validate(); err != nil {
		msg := fmt.Sprintf("invalid password %s for username %s", user.Password, username)
		return page.NewFieldError("password", errors.New(msg))
	}

	if !current.HasPermission(page.PermissionEditPermissions) {
		// user without the EditPermission perm can only enable or disable a user
		if !((len(user.Permissions()) == 1 && user.IsEnabled()) || (len(user.Permissions()) == 0 && !user.IsEnabled())) {
			return page.NewPermissionError(page.PermissionName(page.PermissionEditPermissions))
		}
	}

	// check for user existence
	err = model.FromStringID(ctx, &User{}, username, nil)

	if err == nil {
		// user already exists
		msg := fmt.Sprintf("user %s already exists.", username)
		return page.NewFieldError("user", errors.New(msg))
	}

	if err != datastore.ErrNoSuchEntity {
		// generic datastore error
		msg := fmt.Sprintf("error retrieving user with username %s: %s", username, err.Error())
		return page.NewFieldError("user", errors.New(msg))
	}

	user.Password = HashPassword(user.Password, salt)


	opts := model.CreateOptions{}
	opts.WithStringId(username)

	err = model.CreateWithOptions(ctx, user, &opts)
	if err != nil {
		return fmt.Errorf("error creating post %s: %s", user.Name, err)
	}

	return nil
}

func (manager userManager) Update(ctx context.Context, res page.Resource, bundle []byte) error {
	current := page.IdentityFromContext(ctx)
	if !current.HasPermission(page.PermissionEditUser) {
		return page.NewPermissionError(page.PermissionName(page.PermissionEditUser))
	}

	o, _ := manager.NewResource(ctx)
	if err := o.FromRepresentation(page.RepresentationTypeJSON, bundle); err != nil {
		return page.NewFieldError("", fmt.Errorf("invalid json %s: %s", string(bundle), err.Error()))
	}

	other := o.(*User)
	user := res.(*User)
	user.Name = other.Name

	tkn, _ := tokenManager{}.NewResource(ctx)
	token := tkn.(*Token)
	if err := token.FromRepresentation(page.RepresentationTypeJSON, bundle); err == nil {
		if token.Password != "" {
			pf := page.NewRawField("password", true, token.Password)
			pf.AddValidator(page.LenValidator{MinLen: 8})

			if err := pf.Validate(); err != nil {
				msg := fmt.Sprintf("invalid password %s for username %s", token.Password, other.Username())
				return page.NewFieldError("user", errors.New(msg))
			}
			user.Password = HashPassword(token.Password, salt)
		}
	}

	if other.Email != "" {
		ef := page.NewRawField("email", true, other.Email)
		if err := ef.Validate(); err != nil {
			msg := fmt.Sprintf("invalid email address: %s", other.Email)
			return page.NewFieldError("user", errors.New(msg))
		}
		user.Email = other.Email
	}

	if !current.HasPermission(page.PermissionEditPermissions) && other.ChangedPermission(*user) {
		return page.NewPermissionError(page.PermissionName(page.PermissionEditPermissions))
	}

	user.Name = other.Name
	user.Surname = other.Surname
	user.Permission = other.Permission

	return model.Update(ctx, user)
}

func (manager userManager) Delete(ctx context.Context, res page.Resource) error {
	current := page.IdentityFromContext(ctx)
	if !current.HasPermission(page.PermissionEditUser) {
		return page.NewPermissionError(page.PermissionName(page.PermissionEditUser))
	}

	user := res.(*User)
	err := model.Delete(ctx, user, nil)
	if err != nil {
		return fmt.Errorf("error deleting user %s: %s", user.Name, err.Error())
	}

	return nil
}
