package page

import (
	"context"
	"distudio.com/mage"
	"google.golang.org/appengine/user"
	"net/http"
)

// Returns 200 if the user is authenticated within the appengine framework
type IsSuperuserController struct{}

func (controller *IsSuperuserController) OnDestroy(ctx context.Context) {}

func (controller *IsSuperuserController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	// check if the user is a google user.
	// Gusers as admin bypass normal users controls
	guser := user.Current(ctx)
	if guser == nil {
		return mage.Redirect{Status: http.StatusUnauthorized}
	}

	return mage.Redirect{Status: http.StatusOK}
}
