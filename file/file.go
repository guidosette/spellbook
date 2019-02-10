package file

import (
	"cloud.google.com/go/storage"
	"context"
	"distudio.com/mage"
	"distudio.com/page/identity"
	"google.golang.org/appengine/log"
	"net/http"
)

type UploadController struct {
	mage.Controller
}

func (controller *UploadController) OnDestroy(ctx context.Context) {}

func (controller *UploadController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	ins := mage.InputsFromContext(ctx)
	method := ins[mage.KeyRequestMethod].Value()
	switch method {
	case http.MethodPost:
		u := ctx.Value(identity.KeyUser)
		user, ok := u.(identity.User)
		if !ok {
			return mage.Redirect{Status:http.StatusUnauthorized}
		}

		if !user.HasPermission(identity.PermissionLoadFiles) {
			return mage.Redirect{Status:http.StatusForbidden}
		}

		// handle the upload to Google Cloud Storage
		bucket, err := mage.GetBucketName(ctx)
		if err != nil {
			log.Errorf(ctx, "can't retrieve bucket name: %s", err.Error())
			return mage.Redirect{Status:http.StatusInternalServerError}
		}

		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Errorf(ctx, "failed to create client: %s", err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}
		defer client.Close()

		handle := client.Bucket(bucket)
		filename := "file"
		writer := handle.Object(filename).NewWriter(ctx)
		writer.ContentType = "text/plain"


		if _, err := writer.Write([]byte(ins["file"].Value())); err != nil {
			log.Errorf(ctx,"upload: unable to write file %s to bucket %s: %s", filename, bucket, err.Error())
			return mage.Redirect{Status:http.StatusInternalServerError}
		}

		if err := writer.Close(); err != nil {
			log.Errorf(ctx, "upload: unable to close bucket %s: %s", bucket, err.Error())
		}

		return mage.Redirect{Status:http.StatusCreated}

	}
	return mage.Redirect{Status:http.StatusNotImplemented}
}
