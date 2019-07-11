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

	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionReadContent) {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadContent))
	}

	cont := Content{}

	if err := model.FromEncodedKey(ctx, &cont, id); err != nil {
		log.Errorf(ctx, "could not retrieve content %s: %s", id, err.Error())
		return nil, err
	}

	// attachment
	q := model.NewQuery((*Attachment)(nil))
	q = q.WithField("ParentKey =", cont.EncodedKey())
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

	a := []string{"Category", "Locale", "Name", "Topic"} // list property accepted
	name := opts.Property

	if name == "" {
		return nil, page.NewFieldError("property", fmt.Errorf("empty property"))
	}

	i := sort.Search(len(a), func(i int) bool { return name <= a[i] })
	if i == len(a) {
		return nil, datastore.ErrNoSuchEntity
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
	if current == nil || !current.HasPermission(page.PermissionWriteContent) {
		return page.NewPermissionError(page.PermissionName(page.PermissionWriteContent))
	}

	content := res.(*Content)

	content.Created = time.Now().UTC()
	if content.IdTranslate == "" {
		content.IdTranslate = time.Now().Format(time.RFC3339Nano)
	} else {
		// check same idTranslate and Locale
		q := model.NewQuery((*Content)(nil))
		q = q.WithField("IdTranslate =", content.IdTranslate)
		q = q.WithField("Locale = ", content.Locale)
		count, err := q.Count(ctx)
		if err != nil {
			return page.NewFieldError("locale", fmt.Errorf("error verifying locale translate: %s", err.Error()))
		}
		if count > 0 {
			msg := fmt.Sprintf("a content with IdTranslate '%s' already exists. Locale must be unique.", content.Locale)
			return page.NewFieldError("locale", errors.New(msg))
		}
	}
	content.Revision = 1

	if content.IsPublished() {
		content.PublicationState = PublicationStatePublished
		content.Published = time.Now().UTC()
	} else {
		content.PublicationState = PublicationStateUnpublished
	}

	if content.Type == "" {
		return page.NewFieldError("type", errors.New("type can't be 0"))
	}

	if content.Slug == "" {
		content.Slug = url.PathEscape(content.Title)
	}
	// if the same slug already exists, we must return
	// otherwise we would overwrite an existing entry, which is not in the spirit of the create method
	q := model.NewQuery((*Content)(nil))

	if content.Code == "" {
		q = q.WithField("Slug =", content.Slug)
		q = q.WithField("Locale = ", content.Locale)
	} else {
		q = q.WithField("Code =", content.Code)
		q = q.WithField("Locale = ", content.Locale)
	}
	count, err := q.Count(ctx)
	if err != nil {
		return page.NewFieldError("slug", fmt.Errorf("error verifying slug uniqueness: %s", err.Error()))
	}
	if count > 0 {
		msg := ""
		if content.Code == "" {
			msg = fmt.Sprintf("a content with slug  %s already exists. Slug must be unique.", content.Slug)
		} else {
			msg = fmt.Sprintf("a content with code %s already exists. Code must be unique.", content.Code)
		}
		return page.NewFieldError("slug", errors.New(msg))
	}

	switch content.Type {
	case page.KeyTypeContent:
		if content.Title == "" {
			return page.NewFieldError("title", errors.New("title can't be empty"))
		}
	case page.KeyTypeEvent:
		if content.StartDate.IsZero() {
			msg := fmt.Sprintf("start date can't be empty. %v", content.StartDate)
			return page.NewFieldError("startDate", errors.New(msg))
		}
		if content.EndDate.IsZero() {
			msg := fmt.Sprintf("end date can't be empty. %v", content.StartDate)
			return page.NewFieldError("endDate", errors.New(msg))
		}
		if content.EndDate.Before(content.StartDate) {
			msg := fmt.Sprintf("end date %v can't be before start date %v", content.EndDate, content.StartDate)
			return page.NewFieldError("endDate", errors.New(msg))
		}
	default:
		return fmt.Errorf("error no type %v", content)
	}

	if user, ok := current.(identity.User); ok {
		content.Author = user.Username()
	}

	// // WARNING: the volatile field Multimedia because Memcache (Gob)
	//	can't ignore field
	tmp := content.Attachments
	content.Attachments = nil

	err = model.Create(ctx, content)
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

	if current == nil || !current.HasPermission(page.PermissionWriteContent) {
		return page.NewPermissionError(page.PermissionName(page.PermissionWriteContent))
	}

	content := res.(*Content)

	other := &Content{}
	if err := other.FromRepresentation(page.RepresentationTypeJSON, bundle); err != nil {
		return page.NewFieldError("", fmt.Errorf("invalid json for content %s: %s", content.StringID(), err.Error()))
	}

	if other.Title == "" {
		return page.NewFieldError("title", errors.New("title can't be empty"))
	}

	// if the same slug already exists, we must return
	// otherwise we would overwrite an existing entry, which is not in the spirit of the create method
	q := model.NewQuery((*Content)(nil))
	if other.Code == "" {
		q = q.WithField("Slug =", other.Slug)
		q = q.WithField("Locale = ", other.Locale)
	} else {
		q = q.WithField("Code =", other.Code)
		q = q.WithField("Locale = ", other.Locale)
	}

	count, err := q.Count(ctx)
	if err != nil {
		return page.NewFieldError("slug", fmt.Errorf("error verifying slug uniqueness: %s", err.Error()))
	}
	log.Infof(ctx, "count %d", count)
	if count > 0 {
		var contents []*Content
		err := q.GetMulti(ctx, &contents)
		if err != nil {
			log.Errorf(ctx, "could not retrieve contents %s", err.Error())
			return page.NewFieldError("check", errors.New("could not retrieve contents"))
		}
		valid := true
		for _, c := range contents {
			if c.EncodedKey() != content.EncodedKey() {
				valid = false
				break
			}
		}
		if !valid {
			if content.Code == "" {
				msg := fmt.Sprintf("a content with slug  %s already exists. Slug must be unique.", content.Slug)
				return page.NewFieldError("slug", errors.New(msg))
			} else {
				msg := fmt.Sprintf("a content with code %s already exists. Code must be unique.", content.Code)
				return page.NewFieldError("code", errors.New(msg))
			}
		}
	}

	content.Type = other.Type
	content.Title = other.Title
	content.Subtitle = other.Subtitle
	content.Category = other.Category
	content.Topic = other.Topic
	content.Locale = other.Locale
	content.Description = other.Description
	content.Code = other.Code
	content.Body = other.Body
	content.Cover = other.Cover
	content.Revision = other.Revision
	content.Editor = other.Editor
	content.Order = other.Order
	content.Updated = time.Now().UTC()
	content.Tags = other.Tags
	content.Slug = other.Slug
	content.ParentKey = other.ParentKey

	if !other.IsPublished() {
		// not set
		content.Published = time.Time{} // zero
	} else {
		// set
		// check previous data
		if !content.IsPublished() {
			content.Published = time.Now().UTC()
		}
	}

	if content.IsPublished() {
		content.PublicationState = PublicationStatePublished
	} else {
		content.PublicationState = PublicationStateUnpublished
	}

	switch content.Type {
	case page.KeyTypeContent:
	case page.KeyTypeEvent:
		if other.StartDate.IsZero() {
			msg := fmt.Sprintf("start date can't be empty. %v", other.StartDate)
			return page.NewFieldError("startDate", errors.New(msg))
		}
		if other.EndDate.IsZero() {
			msg := fmt.Sprintf("end date can't be empty. %v", other.EndDate)
			return page.NewFieldError("endDate", errors.New(msg))
		}
		if other.EndDate.Before(other.StartDate) {
			msg := fmt.Sprintf("end date %v can't be before start date %v", other.EndDate, other.StartDate)
			return page.NewFieldError("endDate", errors.New(msg))
		}
		content.StartDate = other.StartDate
		content.EndDate = other.EndDate
	default:
		return fmt.Errorf("error no type %v", content)
	}

	if user, ok := current.(identity.User); ok {
		content.Author = user.Username()
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

	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionWriteContent) {
		return page.NewPermissionError(page.PermissionName(page.PermissionWriteContent))
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
	q.WithField("ParentKey =", content.EncodedKey())
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
