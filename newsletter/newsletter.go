package newsletter

import (
	"distudio.com/mage/model"
	"distudio.com/page"
	"distudio.com/page/identity"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"strings"
)

type Newsletter struct {
	model.Model `json:"-"`
	page.Resource
	Email string `json:"email"`
}

func (newsletter *Newsletter) Create(ctx context.Context) error {

	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionEditNewsletter) {
		return page.NewPermissionError(identity.PermissionName(identity.PermissionEditNewsletter))
	}

	if newsletter.Email == "" {
		msg := fmt.Sprintf("Email can't be empty")
		return page.NewFieldError("Email", errors.New(msg))
	}
	if !strings.Contains(newsletter.Email, "@") || !strings.Contains(newsletter.Email, ".") {
		msg := fmt.Sprintf("Email not valid")
		return page.NewFieldError("Email", errors.New(msg))
	}

	// list newsletter
	var emails []*Newsletter
	q := model.NewQuery(&Newsletter{})
	q = q.WithField("Email =", newsletter.Email)
	err := q.GetMulti(ctx, &emails)
	if err != nil {
		msg := fmt.Sprintf("Error retrieving list newsletter %+v", err)
		return page.NewFieldError("Email", errors.New(msg))
	}
	if len(emails) > 0 {
		msg := fmt.Sprintf("Email already exist")
		return page.NewFieldError("Email", errors.New(msg))
	}

	return nil
}

func (newsletter *Newsletter) Update(ctx context.Context, res page.Resource) error {

	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionEditNewsletter) {
		return page.NewPermissionError(identity.PermissionName(identity.PermissionEditNewsletter))
	}

	other := res.(*Newsletter)
	newsletter.Email = other.Email

	return nil
}

func (newsletter *Newsletter) Id() string {
	return newsletter.StringID()
}
