package navigation

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"errors"
	"fmt"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
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

type seoManager struct{}

func (manager seoManager) NewResource(ctx context.Context) (page.Resource, error) {
	return &Seo{}, nil
}

func (manager seoManager) FromId(ctx context.Context, id string) (page.Resource, error) {
	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionReadSeo) {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadSeo))
	}

	cont := Seo{}
	if err := model.FromStringID(ctx, &cont, id, nil); err != nil {
		log.Errorf(ctx, "could not retrieve seo %s: %s", id, err.Error())
		return nil, err
	}

	return &cont, nil
}

func (manager seoManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {
	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionReadSeo) {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadSeo))
	}

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
	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionReadSeo) {
		return nil, page.NewPermissionError(page.PermissionName(page.PermissionReadSeo))
	}

	ws := page.Application()

	staticPages := ws.Options().StaticPages

	from := opts.Page * opts.Size
	if from > len(staticPages) {
		return make([]string, 0), nil
	}

	to := from + opts.Size
	if to > len(staticPages) {
		to = len(staticPages)
	}

	items := staticPages[from:to]
	codes := make([]string, len(items))

	for i := range items {
		staticPage := page.StaticPageCode(items[i])
		codes[i] = string(staticPage)
	}

	return codes, nil
}

func (manager seoManager) Create(ctx context.Context, res page.Resource, bundle []byte) error {

	current := page.IdentityFromContext(ctx)
	if current == nil || !current.HasPermission(page.PermissionWriteSeo) {
		return page.NewPermissionError(page.PermissionName(page.PermissionWriteSeo))
	}

	seo := res.(*Seo)

	// if the same seo already exists, we must return false
	existing, _ := manager.NewResource(ctx)
	err := model.FromStringID(ctx, existing.(*Seo), PageId(seo.Locale, seo.Url), nil)
	if err == datastore.ErrNoSuchEntity {
		// we can create the new seo element if another one with the same code doesn't exists
		q := model.NewQuery((*Seo)(nil))
		q.WithField("Code =", seo.Code)
		q.WithField("Locale =", seo.Locale)

		err = q.First(ctx, &Seo{})

		// seo already exists for given code, can't create
		if err == nil {
			msg := fmt.Sprintf("a page for %s already exists.", seo.Code)
			return page.NewFieldError("", errors.New(msg))
		}

		if err == datastore.ErrNoSuchEntity {
			seo.IsRoot = seo.Url == rootUrl
			opts := model.NewCreateOptions()
			opts.WithStringId(PageId(seo.Locale, seo.Url))
			err = model.CreateWithOptions(ctx, seo, &opts)
		}

		if err != nil {
			return err
		}

		return nil
	}

	if  err != nil {
		return err
	}

	// the page with the given url has already been allocated. can't create seo
	msg := fmt.Sprintf("a seo for url %s already exists.", seo.Url)
	return page.NewFieldError("", errors.New(msg))
}

func (manager seoManager) Update(ctx context.Context, res page.Resource, bundle []byte) error {

	current := page.IdentityFromContext(ctx)
	if current == nil || !current.HasPermission(page.PermissionWriteSeo) {
		return page.NewPermissionError(page.PermissionName(page.PermissionWriteSeo))
	}

	seo := res.(*Seo)

	or, _ := manager.NewResource(ctx)
	other := or.(*Seo)

	if err := other.FromRepresentation(page.RepresentationTypeJSON, bundle); err != nil {
		return page.NewFieldError("", fmt.Errorf("invalid json for seo %s: %s", seo.StringID(), err.Error()))
	}

	if err := model.FromStringID(ctx, seo, PageId(other.Locale, other.Url), nil); err != nil {
		return page.NewUnsupportedError()
	}

	q := model.NewQuery((*Seo)(nil))
	q.WithField("Code =", other.Code)
	q.WithField("Locale =", other.Locale)

	existing := Seo{}
	err := q.First(ctx, &existing)

	// seo already exists for given code, can't create
	if err == nil && existing.StringID() != PageId(other.Locale, other.Url) {
		msg := fmt.Sprintf("a page for %s already exists.", other.Code )
		return page.NewFieldError("", errors.New(msg))
	}

	seo.Title = other.Title
	seo.MetaDesc = other.MetaDesc
	seo.Code = other.Code

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
