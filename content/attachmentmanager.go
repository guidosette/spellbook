package content

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"distudio.com/page/identity"
	"errors"
	"fmt"
	"google.golang.org/appengine/log"
	"reflect"
	"sort"
	"strconv"
	"time"
)

func NewAttachmentController() *page.RestController {
	return NewAttachmentControllerWithKey("")
}

func NewAttachmentControllerWithKey(key string) *page.RestController {
	man := attachmentManager{}
	handler := page.BaseRestHandler{Manager: man}
	c := page.NewRestController(handler)
	c.Key = key
	return c
}

type attachmentManager struct{}

func (manager attachmentManager) NewResource(ctx context.Context) (page.Resource, error) {
	return &Attachment{}, nil
}

func (manager attachmentManager) FromId(ctx context.Context, strId string) (page.Resource, error) {

	if current := page.IdentityFromContext(ctx); current == nil || (!current.HasPermission(page.PermissionReadContent) && !current.HasPermission(page.PermissionReadMedia)) {
		var p page.Permission
		p = page.PermissionReadContent
		if !current.HasPermission(page.PermissionReadMedia) {
			p = page.PermissionReadMedia
		}
		return nil, page.NewPermissionError(page.PermissionName(p))
	}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return nil, page.NewFieldError(strId, err)
	}

	att := Attachment{}
	if err := model.FromIntID(ctx, &att, id, nil); err != nil {
		log.Errorf(ctx, "could not retrieve attachment %s: %s", id, err.Error())
		return nil, err
	}

	return &att, nil
}

func (manager attachmentManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {
	if current := page.IdentityFromContext(ctx); current == nil || (!current.HasPermission(page.PermissionReadContent) && !current.HasPermission(page.PermissionReadMedia)) {
		var p page.Permission
		p = page.PermissionReadContent
		if !current.HasPermission(page.PermissionReadMedia) {
			p = page.PermissionReadMedia
		}
		return nil, page.NewPermissionError(page.PermissionName(p))
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

	for _, filter := range opts.Filters {
		if filter.Field != "" {
			q = q.WithField(filter.Field+" =", filter.Value)
		}
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
	if current := page.IdentityFromContext(ctx); current == nil || (!current.HasPermission(page.PermissionReadContent) && !current.HasPermission(page.PermissionReadMedia)) {
		var p page.Permission
		p = page.PermissionReadContent
		if !current.HasPermission(page.PermissionReadMedia) {
			p = page.PermissionReadMedia
		}
		return nil, page.NewPermissionError(page.PermissionName(p))
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

func (manager attachmentManager) Create(ctx context.Context, res page.Resource, bundle []byte) error {
	current := page.IdentityFromContext(ctx)
	if current := page.IdentityFromContext(ctx); current == nil || (!current.HasPermission(page.PermissionWriteContent) && !current.HasPermission(page.PermissionWriteMedia)) {
		var p page.Permission
		p = page.PermissionWriteContent
		if !current.HasPermission(page.PermissionWriteMedia) {
			p = page.PermissionWriteMedia
		}
		return page.NewPermissionError(page.PermissionName(p))
	}

	attachment := res.(*Attachment)

	// attachment parent is required.
	// if not attachment is to be specified the default value must be used
	if attachment.ParentKey == "" {
		msg := fmt.Sprintf("attachment parent can't be empty. Use %s as a parent for global attachments", AttachmentGlobalParent)
		return page.NewFieldError("parent", errors.New(msg))
	}

	attachment.Created = time.Now().UTC()
	attachment.Uploader = current.(identity.User).Username()

	err := model.Create(ctx, attachment)
	if err != nil {
		log.Errorf(ctx, "error creating post %s: %s", attachment.Name, err)
		return err
	}

	return nil
}

func (manager attachmentManager) Update(ctx context.Context, res page.Resource, bundle []byte) error {
	if current := page.IdentityFromContext(ctx); current == nil || (!current.HasPermission(page.PermissionWriteContent) && !current.HasPermission(page.PermissionWriteMedia)) {
		var p page.Permission
		p = page.PermissionWriteContent
		if !current.HasPermission(page.PermissionWriteMedia) {
			p = page.PermissionWriteMedia
		}
		return page.NewPermissionError(page.PermissionName(p))
	}

	other := Attachment{}
	if err := other.FromRepresentation(page.RepresentationTypeJSON, bundle); err != nil {
		return page.NewFieldError("", fmt.Errorf("bad json %s", string(bundle)))
	}

	attachment := res.(*Attachment)
	attachment.Name = other.Name
	attachment.Description = other.Description
	attachment.ResourceUrl = other.ResourceUrl
	attachment.Group = other.Group
	attachment.ParentKey = other.ParentKey
	attachment.Updated = time.Now().UTC()
	attachment.AltText = other.AltText
	attachment.Seo = other.Seo

	if attachment.ParentKey == "" {
		msg := fmt.Sprintf("attachment parent can't be empty. Use %s as a parent for global attachments", AttachmentGlobalParent)
		return page.NewFieldError("parent", errors.New(msg))
	}

	return model.Update(ctx, attachment)
}

func (manager attachmentManager) Delete(ctx context.Context, res page.Resource) error {
	if current := page.IdentityFromContext(ctx); current == nil || (!current.HasPermission(page.PermissionWriteContent) && !current.HasPermission(page.PermissionWriteMedia)) {
		var p page.Permission
		p = page.PermissionWriteContent
		if !current.HasPermission(page.PermissionWriteMedia) {
			p = page.PermissionWriteMedia
		}
		return page.NewPermissionError(page.PermissionName(p))
	}

	attachment := res.(*Attachment)
	err := model.Delete(ctx, attachment, nil)
	if err != nil {
		log.Errorf(ctx, "error deleting attachment %s: %s", attachment.Name, err.Error())
		return err
	}

	return nil
}
