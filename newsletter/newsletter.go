package newsletter

import (
	"distudio.com/mage/model"
	"distudio.com/page"
	"distudio.com/page/validators"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"strings"
	"time"
)

var ZeroTime = time.Time{}

type Newsletter struct {
	model.Model `json:"-"`
	page.Resource
	Email string `json:"email"`
}

func (newsletter *Newsletter) Create(ctx context.Context) error {
	// todo permission?
	//current, _ := ctx.Value(identity.KeyUser).(identity.User)
	//if !current.HasPermission(identity.PermissionCreateContent) {
	//	return resource.NewPermissionError(identity.PermissionCreateContent)
	//}

	if newsletter.Email == "" {
		msg := fmt.Sprintf("Email can't be empty")
		return validators.NewFieldError("Email", errors.New(msg))
	}
	if !strings.Contains(newsletter.Email, "@") || !strings.Contains(newsletter.Email, ".") {
		msg := fmt.Sprintf("Email not valid")
		return validators.NewFieldError("Email", errors.New(msg))
	}

	// list newsletter
	var emails []*Newsletter
	q := model.NewQuery(&Newsletter{})
	q = q.WithField("Email =", newsletter.Email)
	err := q.GetMulti(ctx, &emails)
	if err != nil {
		msg := fmt.Sprintf("Error retrieving list newsletter %+v", err)
		return validators.NewFieldError("Email", errors.New(msg))
	}
	if len(emails) > 0 {
		msg := fmt.Sprintf("Email already exist")
		return validators.NewFieldError("Email", errors.New(msg))
	}

	return nil
}

func (newsletter *Newsletter) Update(ctx context.Context, res page.Resource) error {
	//// todo permission?
	//current, _ := ctx.Value(identity.KeyUser).(identity.User)
	//if !current.HasPermission(identity.PermissionEditContent) {
	//	return resource.NewPermissionError(identity.PermissionEditContent)
	//}

	other := res.(*Newsletter)
	newsletter.Email = other.Email

	return nil
}

func (newsletter *Newsletter) Id() string {
	return newsletter.StringID()
}
