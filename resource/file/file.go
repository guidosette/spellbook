package file

import (
	"cloud.google.com/go/storage"
	"distudio.com/mage"
	"distudio.com/mage/model"
	"distudio.com/page/resource"
	"distudio.com/page/validators"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	appengineFile "google.golang.org/appengine/file"
)

type File struct {
	model.Model `json:"-"`
	resource.Resource
	Name        string `json:"name"`
	ResourceUrl string `json:"resourceUrl"`
}

func (file *File) Create(ctx context.Context) error {
	// todo permission?
	//current, _ := ctx.Value(identity.KeyUser).(identity.User)
	//if !current.HasPermission(identity.PermissionCreateContent) {
	//	return resource.NewPermissionError(identity.PermissionCreateContent)
	//}

	ins := mage.InputsFromContext(ctx)

	tv := validators.NewField("type", true, ins)
	tv.AddValidator(validators.FileNameValidator{})
	typ, err := tv.Value()
	if err != nil {
		return validators.NewFieldError("type", err)
	}

	// namespace is the sub folder where the file will be loaded
	nsv := validators.NewField("namespace", false, ins)
	nsv.AddValidator(validators.FileNameValidator{AllowEmpty: true})
	namespace, err := nsv.Value()
	if err != nil {
		return validators.NewFieldError("namespace", err)
	}

	// prepend a slash to build the firename
	if namespace != "" {
		namespace = fmt.Sprintf("/%s", namespace)
	}

	nv := validators.NewField("name", true, ins)
	nv.AddValidator(validators.FileNameValidator{})
	name, err := nv.Value()
	if err != nil {
		return validators.NewFieldError("name", err)
	}

	// get the file headers
	fhs := ins["file"].Files()
	// todo: handle multiple files
	fh := fhs[0]
	f, err := fh.Open()
	defer f.Close()

	buffer := make([]byte, fh.Size)
	if err != nil {
		msg := fmt.Sprintf("error buffer: %s", err.Error())
		return validators.NewFieldError("buffer", errors.New(msg))
	}

	_, err = f.Read(buffer)
	if err == nil {
		// reset the buffer
		_, err = f.Seek(0, 0)
		if err != nil {
			msg := fmt.Sprintf("error Seek buffer: %s", err.Error())
			return validators.NewFieldError("buffer", errors.New(msg))
		}
	} else {
		return validators.NewFieldError("read", err)
	}

	// build the filename
	filename := fmt.Sprintf("%s%s/%s", typ, namespace, name)

	// handle the upload to Google Cloud Storage
	bucket, err := appengineFile.DefaultBucketName(ctx)
	if err != nil {
		msg := fmt.Sprintf("can't retrieve bucket name: %s", err.Error())
		return validators.NewFieldError("bucket", errors.New(msg))
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		msg := fmt.Sprintf("failed to create client: %s", err.Error())
		return validators.NewFieldError("client", errors.New(msg))
	}
	defer client.Close()

	handle := client.Bucket(bucket)
	writer := handle.Object(filename).NewWriter(ctx)
	writer.ContentType = fh.Header.Get("Content-Type")

	if _, err := writer.Write(buffer); err != nil {
		msg := fmt.Sprintf("upload: unable to write file %s to bucket %s: %s", filename, bucket, err.Error())
		return validators.NewFieldError("parent", errors.New(msg))
	}

	if err := writer.Close(); err != nil {
		msg := fmt.Sprintf("upload: unable to close bucket %s: %s", bucket, err.Error())
		return validators.NewFieldError("parent", errors.New(msg))
	}

	// return the file data
	const publicURL = "https://storage.googleapis.com/%s/%s"
	uri := fmt.Sprintf(publicURL, bucket, filename)

	file.ResourceUrl = uri
	file.Name = name

	return nil
}

func (file *File) Update(ctx context.Context, res resource.Resource) error {
	// todo permission?
	//current, _ := ctx.Value(identity.KeyUser).(identity.User)
	//if !current.HasPermission(identity.PermissionEditContent) {
	//	return resource.NewPermissionError(identity.PermissionEditContent)
	//}

	other := res.(*File)
	file.Name = other.Name

	return nil
}

func (file *File) Id() string {
	return file.StringID()
}
