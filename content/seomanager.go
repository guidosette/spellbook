package content

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"errors"
	"fmt"
	"google.golang.org/appengine/log"
	"strconv"
)

func NewSeoController() *page.RestController {
	return NewSeoControllerWithKey("")
}

func NewSeoControllerWithKey(key string) *page.RestController {
	handler := page.BaseRestHandler{Manager: seoManager{}}
	c := page.NewRestController(handler)
	c.Key = key
	return c
}

type seoManager struct {}

func (manager seoManager) NewResource(ctx context.Context) (page.Resource, error) {
	return &Seo{}, nil
}

func (manager seoManager) FromId(ctx context.Context, id string) (page.Resource, error) {

	intid, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		err := fmt.Errorf("Invalid ID for seo: %s", err.Error())
		return nil, page.NewFieldError("id", err)
	}

	cont := Seo{}
	if err := model.FromIntID(ctx, &cont, intid, nil); err != nil {
		log.Errorf(ctx, "could not retrieve seo %s: %s", id, err.Error())
		return nil, err
	}

	return &cont, nil
}

func (manager seoManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {

	var conts []*Seo
	q := model.NewQuery(&Seo{})
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

func (manager seoManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	var result []string
	return result, nil
}


func (manager seoManager) Create(ctx context.Context, res page.Resource, bundle []byte) error {

	current := page.IdentityFromContext(ctx)
	if current == nil || !current.HasPermission(page.PermissionCreateSeo) {
		return page.NewPermissionError(page.PermissionName(page.PermissionCreateSeo))
	}

	seo := res.(*Seo)

	if seo.Title == "" || seo.Url == "" {
		return page.NewFieldError("title", errors.New("title and url can't be empty"))
	}

	// if the same seo already exists, we must return
	q := model.NewQuery((*Seo)(nil))
	q = q.WithField("Url =", seo.Url)
	count, err := q.Count(ctx)
	if err != nil {
		return page.NewFieldError("url", fmt.Errorf("error verifying url uniqueness: %s", err.Error()))
	}

	if count > 0 {
		msg := fmt.Sprintf("a seo for url %s already exists.", seo.Url)
		return page.NewFieldError("url", errors.New(msg))
	}

	err = model.Create(ctx, seo)
	if err != nil {
		log.Errorf(ctx, "error creating seo for url %s: %s", seo.Url, err)
		return err
	}

	return nil
}

func (manager seoManager) Update(ctx context.Context, res page.Resource, bundle []byte) error {

	current := page.IdentityFromContext(ctx)
	if current == nil || !current.HasPermission(page.PermissionEditSeo) {
		return page.NewPermissionError(page.PermissionName(page.PermissionEditSeo))
	}

	seo := res.(*Seo)

	other := Seo{}
	if err := other.FromRepresentation(page.RepresentationTypeJSON, bundle); err != nil {
		return page.NewFieldError("", fmt.Errorf("invalid json for seo %s: %s", seo.StringID(), err.Error()))
	}

	if seo.Title == "" || seo.Url == "" {
		return page.NewFieldError("title", errors.New("title and url can't be empty"))
	}

	if len(seo.MetaDesc) > 160 {
		return page.NewFieldError("metadesc", errors.New("metadesc can be at most 160 characters long"))
	}

	seo.Title = other.Title
	seo.MetaDesc = other.MetaDesc
	seo.Url = other.Url

	if err := model.Update(ctx, seo); err != nil {
		return fmt.Errorf("error updating seo with url %s: %s", seo.Url, err)
	}

	return nil
}

func (manager seoManager) Delete(ctx context.Context, res page.Resource) error {

	seo := res.(*Seo)
	err := model.Delete(ctx, seo, nil)
	if err != nil {
		log.Errorf(ctx, "error deleting seo with url %s: %s", seo.Url, err.Error())
		return err
	}

	return nil
}

