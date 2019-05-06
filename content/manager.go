package content

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"distudio.com/page/attachment"
	"distudio.com/page/identity"
	"errors"
	"google.golang.org/appengine/log"
	"reflect"
	"sort"
)

type Manager struct{}

func (manager Manager) NewResource(ctx context.Context) (page.Resource, error) {
	return &Content{}, nil
}

func (manager Manager) FromId(ctx context.Context, id string) (page.Resource, error) {
	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionReadContent) {
		return nil, page.NewPermissionError(identity.PermissionName(identity.PermissionReadContent))
	}

	cont := Content{}
	if err := model.FromStringID(ctx, &cont, id, nil); err != nil {
		log.Errorf(ctx, "could not retrieve content %s: %s", id, err.Error())
		return nil, err
	}

	q := model.NewQuery((*attachment.Attachment)(nil))
	q = q.WithField("Parent =", cont.Slug)
	if err := q.GetMulti(ctx, &cont.Attachments); err != nil {
		log.Errorf(ctx, "could not retrieve content %s attachments: %s", id, err.Error())
		return nil, err
	}

	return &cont, nil
}

func (manager Manager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {

	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionReadContent) {
		return nil, page.NewPermissionError(identity.PermissionName(identity.PermissionReadContent))
	}

	var conts []*Content
	q := model.NewQuery(&Content{})
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
	err := q.GetMulti(ctx, &conts)
	if err != nil {
		return nil, err
	}

	resources := make([]page.Resource, len(conts))
	for i := range conts {
		resources[i] = page.Resource(conts[i])
	}

	return resources, nil
}

func (manager Manager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {

	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionReadContent) {
		return nil, page.NewPermissionError(identity.PermissionName(identity.PermissionReadContent))
	}

	a := []string{"Category", "Topic", "Name"} // list property accepted
	name := opts.Property

	i := sort.Search(len(a), func(i int) bool { return name <= a[i] })
	if i < len(a) && a[i] == name {
		// found
	} else {
		return nil, errors.New("no property found")
	}

	var conts []*Content
	q := model.NewQuery(&Content{})
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

	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionEditContent) {
		return page.NewPermissionError(identity.PermissionName(identity.PermissionEditContent))
	}

	content := res.(*Content)
	// input is valid, create the resource
	opts := model.CreateOptions{}
	opts.WithStringId(content.Slug)

	// // WARNING: the volatile field Multimedia because Memcache (Gob)
	//	can't ignore field
	tmp := content.Attachments
	content.Attachments = nil

	err := model.CreateWithOptions(ctx, content, &opts)
	if err != nil {
		log.Errorf(ctx, "error creating post %s: %s", content.Slug, err)
		return err
	}

	// return the swapped multimedia value
	content.Attachments = tmp
	return nil
}

func (manager Manager) Delete(ctx context.Context, res page.Resource) error {

	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionEditContent) {
		return page.NewPermissionError(identity.PermissionName(identity.PermissionEditContent))
	}

	content := res.(*Content)
	err := model.Delete(ctx, content, nil)
	if err != nil {
		log.Errorf(ctx, "error deleting content %s: %s", content.Slug, err.Error())
		return err
	}

	// delete attachments with parent = slug
	attachments := make([]*attachment.Attachment, 0, 0)
	q := model.NewQuery(&attachment.Attachment{})
	q.WithField("Parent =", content.Slug)
	err = q.GetMulti(ctx, &attachments)
	if err != nil {
		log.Errorf(ctx, "error retrieving attachments: %s", err)
		return err
	}

	for _, att := range attachments {
		err = model.Delete(ctx, att, nil)
		if err != nil {
			log.Errorf(ctx, "error deleting attachment %s: %s", att.Name, err.Error())
			return err
		}
	}

	return nil
}
