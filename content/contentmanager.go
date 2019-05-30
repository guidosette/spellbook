package content

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"distudio.com/page/identity"
	"errors"
	"fmt"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NewContentController() *page.RestController {
	return NewContentControllerWithKey("")
}

func NewContentControllerWithKey(key string) *page.RestController {
	handler := page.BaseRestHandler{Manager: contentManager{}}
	c := page.NewRestController(handler)
	c.Key = key
	return c
}

type contentManager struct{}

func (manager contentManager) NewResource(ctx context.Context) (page.Resource, error) {
	return &Content{}, nil
}

func (manager contentManager) FromId(ctx context.Context, id string) (page.Resource, error) {

	intid, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		err := fmt.Errorf("Invalid ID for content: %s", err.Error())
		return nil, page.NewFieldError("id", err)
	}

	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionReadContent) {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadContent))
	}

	cont := Content{}
	if err := model.FromIntID(ctx, &cont, intid, nil); err != nil {
		log.Errorf(ctx, "could not retrieve content %s: %s", id, err.Error())
		return nil, err
	}

	q := model.NewQuery((*Attachment)(nil))
	q = q.WithField("Parent =", cont.Slug)
	if err := q.GetMulti(ctx, &cont.Attachments); err != nil {
		log.Errorf(ctx, "could not retrieve content %s attachments: %s", id, err.Error())
		return nil, err
	}

	return &cont, nil
}

func (manager contentManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {

	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionReadContent) {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadContent))
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
	for _, filter := range opts.Filters {
		if filter.Field != "" {
			q = q.WithField(filter.Field+" =", filter.Value)
		}
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

func (manager contentManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {

	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionReadContent) {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadContent))
	}

	a := []string{"category", "topic", "name"} // list property accepted
	name := opts.Property

	if name == "" {
		return nil, page.NewFieldError("property", fmt.Errorf("empty property"))
	}

	i := sort.Search(len(a), func(i int) bool { return name <= a[i] })
	if i == len(a) {
		return nil, datastore.ErrNoSuchEntity
	}

	name = fmt.Sprintf("%s%s", strings.ToUpper(string(name[0])), name[1:])
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

func (manager contentManager) Create(ctx context.Context, res page.Resource, bundle []byte) error {

	current := page.IdentityFromContext(ctx)
	if current == nil || !current.HasPermission(page.PermissionCreateContent) {
		return page.NewPermissionError(page.PermissionName(page.PermissionCreateContent))
	}

	content := res.(*Content)

	content.Created = time.Now().UTC()
	content.Revision = 1
	if !content.Published.IsZero() {
		content.Published = time.Now().UTC()
	}

	if content.Title == "" || content.Name == "" {
		return page.NewFieldError("title", errors.New("title and name can't be empty"))
	}

	if content.Slug == "" {
		content.Slug = url.PathEscape(content.Title)
	}

	// if the same slug already exists, we must return
	// otherwise we would overwrite an existing entry, which is not in the spirit of the create method
	q := model.NewQuery((*Content)(nil))
	q = q.WithField("Slug =", content.Slug)
	count, err := q.Count(ctx)
	if err != nil {
		return page.NewFieldError("slug", fmt.Errorf("error verifying slug uniqueness: %s", err.Error()))
	}

	if count > 0 {
		msg := fmt.Sprintf("a content with slug %s already exists. Slug must be unique.", content.Slug)
		return page.NewFieldError("slug", errors.New(msg))
	}

	if user, ok := current.(identity.User); ok {
		content.Author = user.Username()
	}

	// input is valid, create the resource
	opts := model.CreateOptions{}
	//opts.WithStringId(content.Slug)

	// // WARNING: the volatile field Multimedia because Memcache (Gob)
	//	can't ignore field
	tmp := content.Attachments
	content.Attachments = nil

	err = model.CreateWithOptions(ctx, content, &opts)
	if err != nil {
		log.Errorf(ctx, "error creating post %s: %s", content.Slug, err)
		return err
	}

	// return the swapped multimedia value
	content.Attachments = tmp

	return nil
}

func (manager contentManager) Update(ctx context.Context, res page.Resource, bundle []byte) error {

	current := page.IdentityFromContext(ctx)

	if current == nil || !current.HasPermission(page.PermissionEditContent) {
		return page.NewPermissionError(page.PermissionName(page.PermissionEditContent))
	}

	content := res.(*Content)

	other := Content{}
	if err := other.FromRepresentation(page.RepresentationTypeJSON, bundle); err != nil {
		return page.NewFieldError("", fmt.Errorf("invalid json for content %s: %s", content.StringID(), err.Error()))
	}

	if other.Title == "" || other.Name == "" {
		return page.NewFieldError("title", errors.New("title and name can't be empty"))
	}

	content.Name = other.Name
	content.Title = other.Title
	content.Subtitle = other.Subtitle
	content.Category = other.Category
	content.Topic = other.Topic
	content.Locale = other.Locale
	content.Description = other.Description
	content.Body = other.Body
	content.Cover = other.Cover
	content.Revision = other.Revision
	content.Order = other.Order
	content.Updated = time.Now().UTC()
	content.Tags = other.Tags

	if user, ok := current.(identity.User); ok {
		content.Author = user.Username()
	}

	if other.Published.IsZero() {
		// not set
		content.Published = time.Time{}
	} else {
		// set
		// check previous data
		if content.Published.IsZero() {
			content.Published = time.Now().UTC()
		}
	}

	tmp := content.Attachments
	content.Attachments = nil

	if err := model.Update(ctx, content); err != nil {
		return fmt.Errorf("error updating post %s: %s", content.Slug, err)
	}

	// return the swapped multimedia value
	content.Attachments = tmp

	return nil
}

func (manager contentManager) Delete(ctx context.Context, res page.Resource) error {

	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionEditContent) {
		return page.NewPermissionError(page.PermissionName(page.PermissionEditContent))
	}

	content := res.(*Content)
	err := model.Delete(ctx, content, nil)
	if err != nil {
		log.Errorf(ctx, "error deleting content %s: %s", content.Slug, err.Error())
		return err
	}

	// delete attachments with parent = slug
	attachments := make([]*Attachment, 0, 0)
	q := model.NewQuery(&Attachment{})
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
