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

func NewAttachmentController() *page.RestController {
	return NewContentControllerWithKey("")
}

func NewAttachmentControllerWithKey(key string) *page.RestController {
	man := attachmentManager{}
	handler := page.BaseRestHandler{Manager:man}
	c := page.NewRestController(handler)
	c.Key = key
	return c
}

type attachmentManager struct{}

func (manager attachmentManager) NewResource(ctx context.Context) (page.Resource, error) {
	return &Attachment{}, nil
}

func (manager attachmentManager) FromId(ctx context.Context, id string) (page.Resource, error) {

	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionReadContent) {
		return nil, page.NewPermissionError(identity.PermissionName(identity.PermissionReadContent))
	}

	att := Attachment{}
	if err := model.FromStringID(ctx, &att, id, nil); err != nil {
		log.Errorf(ctx, "could not retrieve attachment %s: %s", id, err.Error())
		return nil, err
	}

	return &att, nil
}

func (manager attachmentManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {
	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionReadContent) {
		return nil, page.NewPermissionError(identity.PermissionName(identity.PermissionReadContent))
	}

	var attachments []*Attachment
	q := model.NewQuery(&Attachment{})
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
	err := q.GetMulti(ctx, &attachments)
	if err != nil {
		return nil, err
	}

	resources := make([]page.Resource, len(attachments))
	for i := range attachments {
		resources[i] = page.Resource(attachments[i])
	}

	return resources, nil
}

func (manager attachmentManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionReadContent) {
		return nil, page.NewPermissionError(identity.PermissionName(identity.PermissionReadContent))
	}

	a := []string{"Group"} // list property accepted
	name := opts.Property

	i := sort.Search(len(a), func(i int) bool { return name <= a[i] })
	if i < len(a) && a[i] == name {
		// found
	} else {
		return nil, errors.New("no property found")
	}

	var conts []*Attachment
	q := model.NewQuery(&Attachment{})
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

func (manager attachmentManager) Save(ctx context.Context, res page.Resource) error {
	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionEditContent) {
		return page.NewPermissionError(identity.PermissionName(identity.PermissionEditContent))
	}

	attachment := res.(*Attachment)

	err := model.Create(ctx, attachment)
	if err != nil {
		log.Errorf(ctx, "error creating post %s: %s", attachment.Name, err)
		return err
	}

	return nil
}

func (manager attachmentManager) Delete(ctx context.Context, res page.Resource) error {
	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionEditContent) {
		return page.NewPermissionError(identity.PermissionName(identity.PermissionEditContent))
	}

	attachment := res.(*Attachment)
	err := model.Delete(ctx, attachment, nil)
	if err != nil {
		log.Errorf(ctx, "error deleting attachment %s: %s", attachment.Name, err.Error())
		return err
	}

	return nil
}
