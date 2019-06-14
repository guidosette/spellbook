package mailmessage

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"errors"
	"fmt"
	"google.golang.org/appengine/log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NewMailMessageController() *page.RestController {
	return NewMailMessageControllerWithKey("")
}

func NewMailMessageControllerWithKey(key string) *page.RestController {
	man := mailMessageManager{}
	handler := page.BaseRestHandler{Manager: man}
	c := page.NewRestController(handler)
	c.Key = key
	return c
}

type mailMessageManager struct{}

func (manager mailMessageManager) NewResource(ctx context.Context) (page.Resource, error) {
	return &MailMessage{}, nil
}

func (manager mailMessageManager) FromId(ctx context.Context, strId string) (page.Resource, error) {
	// todo permission?
	//current, _ := ctx.Value(identity.KeyUser).(identity.User)
	//if !current.HasPermission(identity.PermissionReadContent) {
	//	return nil, resource.NewPermissionError(identity.PermissionReadContent)
	//}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return nil, page.NewFieldError(strId, err)
	}

	att := MailMessage{}
	if err := model.FromIntID(ctx, &att, id, nil); err != nil {
		log.Errorf(ctx, "could not retrieve mailMessage %s: %s", id, err.Error())
		return nil, err
	}

	return &att, nil
}

func (manager mailMessageManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {
	// todo permission?
	//current, _ := ctx.Value(identity.KeyUser).(identity.User)
	//if !current.HasPermission(identity.PermissionReadContent) {
	//	return nil, resource.NewPermissionError(identity.PermissionReadContent)
	//}

	var mailMessages []*MailMessage
	q := model.NewQuery(&MailMessage{})
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
	err := q.GetMulti(ctx, &mailMessages)
	if err != nil {
		return nil, err
	}

	resources := make([]page.Resource, len(mailMessages))
	for i := range mailMessages {
		resources[i] = page.Resource(mailMessages[i])
	}

	return resources, nil
}

func (manager mailMessageManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	// todo permission?
	//current, _ := ctx.Value(identity.KeyUser).(identity.User)
	//if !current.HasPermission(identity.PermissionReadContent) {
	//	return nil, resource.NewPermissionError(identity.PermissionReadContent)
	//}

	a := []string{"Recipient"} // list property accepted
	name := opts.Property

	i := sort.Search(len(a), func(i int) bool { return name <= a[i] })
	if i < len(a) && a[i] == name {
		// found
	} else {
		return nil, errors.New("no property found")
	}

	var conts []*MailMessage
	q := model.NewQuery(&MailMessage{})
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

func (manager mailMessageManager) Create(ctx context.Context, res page.Resource, bundle []byte) error {

	current := page.IdentityFromContext(ctx)
	if current == nil || !current.HasPermission(page.PermissionWriteMailMessage) {
		return page.NewPermissionError(page.PermissionName(page.PermissionWriteMailMessage))
	}

	mailMessage := res.(*MailMessage)
	mailMessage.Created = time.Now().UTC()

	if mailMessage.Recipient == "" {
		msg := fmt.Sprintf("Recipient can't be empty")
		return page.NewFieldError("Recipient", errors.New(msg))
	}
	if !strings.Contains(mailMessage.Recipient, "@") || !strings.Contains(mailMessage.Recipient, ".") {
		msg := fmt.Sprintf("Recipient not valid")
		return page.NewFieldError("Recipient", errors.New(msg))
	}

	// list mailMessage
	var emails []*MailMessage
	q := model.NewQuery(&MailMessage{})
	q = q.WithField("Recipient =", mailMessage.Recipient)
	err := q.GetMulti(ctx, &emails)
	if err != nil {
		msg := fmt.Sprintf("Error retrieving list mailMessage %+v", err)
		return page.NewFieldError("Recipient", errors.New(msg))
	}
	if len(emails) > 0 {
		msg := fmt.Sprintf("Recipient already exist")
		return page.NewFieldError("Recipient", errors.New(msg))
	}

	err = model.Create(ctx, mailMessage)
	if err != nil {
		log.Errorf(ctx, "error creating mailMessage %s: %s", mailMessage.Name, err)
		return err
	}

	return nil
}

func (manager mailMessageManager) Update(ctx context.Context, res page.Resource, bundle []byte) error {
	current := page.IdentityFromContext(ctx)
	if current == nil || !current.HasPermission(page.PermissionWriteMailMessage) {
		return page.NewPermissionError(page.PermissionName(page.PermissionWriteMailMessage))
	}

	other := MailMessage{}
	if err := other.FromRepresentation(page.RepresentationTypeJSON, bundle); err != nil {
		return page.NewFieldError("", fmt.Errorf("invalid json %s: %s", string(bundle), err.Error()))
	}

	mailMessage := res.(*MailMessage)
	mailMessage.Recipient = other.Recipient
	return model.Update(ctx, mailMessage)
}

func (manager mailMessageManager) Delete(ctx context.Context, res page.Resource) error {
	if current := page.IdentityFromContext(ctx); current == nil || !current.HasPermission(page.PermissionWriteMailMessage) {
		return page.NewPermissionError(page.PermissionName(page.PermissionWriteMailMessage))
	}

	mailMessage := res.(*MailMessage)
	err := model.Delete(ctx, mailMessage, nil)
	if err != nil {
		log.Errorf(ctx, "error deleting mailMessage %s: %s", mailMessage.Name, err.Error())
		return err
	}

	return nil
}
