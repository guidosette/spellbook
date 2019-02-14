package page

import (
	"cloud.google.com/go/storage"
	"context"
	"distudio.com/mage"
	"distudio.com/page/identity"
	"distudio.com/page/validators"
	"fmt"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
	"net/http"
)

type FileController struct {
	mage.Controller
}

func (controller *FileController) OnDestroy(ctx context.Context) {}

func (controller *FileController) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
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

		errs := validators.Errors{}

		// we want our path to be as follows:
		// '/type/namespace/filename.filetype
		// type (not file type) of file upload. Ex. post, profile, user etc
		// it serves as a first directory
		tv := validators.NewField("type", true, ins)
		tv.AddValidator(validators.FileNameValidator{})
		typ, err := tv.Value()
		if err != nil {
			errs.AddError("type", err)
		}

		// namespace is the sub folder where the file will be loaded
		nsv := validators.NewField("namespace", false, ins)
		nsv.AddValidator(validators.FileNameValidator{AllowEmpty:true})
		namespace, err := nsv.Value()
		if err != nil {
			errs.AddError("namespace", err)
		}

		// prepend a slash to build the firename
		if namespace != "" {
			namespace = fmt.Sprintf("/%s", namespace)
		}

		nv := validators.NewField("name", true, ins)
		nv.AddValidator(validators.FileNameValidator{})
		name, err := nv.Value()
		if err != nil {
			errs.AddError("name", err)
		}

		// get the file headers
		fhs := ins["file"].Files()
		// todo: handle multiple files
		fh := fhs[0]
		f, err := fh.Open()
		defer f.Close()

		buffer := make([]byte, fh.Size)
		if err != nil {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		_, err = f.Read(buffer)
		if err == nil {
			// reset the buffer
			_, err = f.Seek(0, 0)
			if err != nil {
				return mage.Redirect{Status: http.StatusInternalServerError}
			}
		} else {
			return mage.Redirect{Status: http.StatusBadRequest}
		}

		// build the filename
		filename := fmt.Sprintf("/%s%s/%s", typ, namespace, name)

		// handle the upload to Google Cloud Storage
		bucket, err := file.DefaultBucketName(ctx)
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
		writer := handle.Object(filename).NewWriter(ctx)
		writer.ContentType = fh.Header.Get("Content-Type")

		if _, err := writer.Write(buffer); err != nil {
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
