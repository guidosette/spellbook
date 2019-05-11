package content

import (
	"cloud.google.com/go/storage"
	"context"
	"distudio.com/mage"
	"distudio.com/page"
	"distudio.com/page/identity"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
	"strings"
)

func NewFileController() *page.RestController {
	return NewFileControllerWithKey("")
}

func NewFileControllerWithKey(key string) *page.RestController {
	handler := fileHandler{page.BaseRestHandler{Manager: fileManager{}}}
	c := page.NewRestController(handler)
	c.Key = key
	return c
}

type fileManager struct{}

func (manager fileManager) NewResource(ctx context.Context) (page.Resource, error) {
	return &File{}, nil
}

func (manager fileManager) FromId(ctx context.Context, id string) (page.Resource, error) {

	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionReadContent) {
		return nil, page.NewPermissionError(identity.PermissionName(identity.PermissionReadContent))
	}

	//todo: set bucket name in configuration
	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving default bucket %s", err.Error())
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %s", err.Error())
	}
	defer client.Close()

	handle := client.Bucket(bucket)

	reader, err := handle.Object(id).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	res, _ := manager.NewResource(ctx)
	f := res.(*File)
	f.Name = id
	f.ResourceUrl = fmt.Sprintf(publicURL, bucket, id)

	return f, nil
}

func (manager fileManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {

	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionReadContent) {
		return nil, page.NewPermissionError(identity.PermissionName(identity.PermissionReadContent))
	}

	//todo: set bucket name in configuration
	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving default bucket %s", err.Error())
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %s", err.Error())
	}
	defer client.Close()

	var files []*File

	handle := client.Bucket(bucket)

	q := &storage.Query{}
	q.Versions = false

	it := handle.Objects(ctx, q)
	// todo: handle pagination (https://godoc.org/google.golang.org/api/iterator)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf(ctx, "listBucket: unable to list bucket %q: %v", bucket, err)
			break
		}
		name := obj.Name
		s := strings.Split(obj.Name, "/")
		if len(s) > 0 {
			name = s[len(s)-1]
		}
		res, _ := manager.NewResource(ctx)
		f := res.(*File)
		f.Name = name
		f.ResourceUrl = obj.MediaLink
		files = append(files, f)
	}

	resources := make([]page.Resource, len(files))
	for i := range files {
		resources[i] = page.Resource(files[i])
	}

	return resources, nil
}

func (manager fileManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	return nil, page.NewUnsupportedError()
}

func (manager fileManager) Create(ctx context.Context, res page.Resource, bundle []byte) error {
	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionLoadFiles) {
		return page.NewPermissionError(identity.PermissionName(identity.PermissionLoadFiles))
	}

	rfile := res.(*File)

	ins := mage.InputsFromContext(ctx)

	tv := page.NewField("type", true, ins)
	tv.AddValidator(page.FileNameValidator{})
	typ, err := tv.Value()
	if err != nil {
		return page.NewFieldError("type", err)
	}

	// namespace is the sub folder where the file will be loaded
	nsv := page.NewField("namespace", false, ins)
	nsv.AddValidator(page.FileNameValidator{AllowEmpty: true})
	namespace, err := nsv.Value()
	if err != nil {
		return page.NewFieldError("namespace", err)
	}

	// prepend a slash to build the filename
	if namespace != "" {
		namespace = fmt.Sprintf("/%s", namespace)
	}

	nv := page.NewField("name", true, ins)
	nv.AddValidator(page.FileNameValidator{})
	name, err := nv.Value()
	if err != nil {
		return page.NewFieldError("name", err)
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
		return page.NewFieldError("buffer", errors.New(msg))
	}

	_, err = f.Read(buffer)
	if err == nil {
		// reset the buffer
		_, err = f.Seek(0, 0)
		if err != nil {
			msg := fmt.Sprintf("error Seek buffer: %s", err.Error())
			return page.NewFieldError("buffer", errors.New(msg))
		}
	} else {
		return page.NewFieldError("read", err)
	}

	// build the filename
	filename := fmt.Sprintf("%s%s/%s", typ, namespace, name)

	// handle the upload to Google Cloud Storage
	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		msg := fmt.Sprintf("can't retrieve bucket name: %s", err.Error())
		return page.NewFieldError("bucket", errors.New(msg))
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		msg := fmt.Sprintf("failed to create client: %s", err.Error())
		return page.NewFieldError("client", errors.New(msg))
	}
	defer client.Close()

	handle := client.Bucket(bucket)
	writer := handle.Object(filename).NewWriter(ctx)
	writer.ContentType = fh.Header.Get("Content-Type")

	if _, err := writer.Write(buffer); err != nil {
		msg := fmt.Sprintf("upload: unable to write file %s to bucket %s: %s", filename, bucket, err.Error())
		return page.NewFieldError("parent", errors.New(msg))
	}

	if err := writer.Close(); err != nil {
		msg := fmt.Sprintf("upload: unable to close bucket %s: %s", bucket, err.Error())
		return page.NewFieldError("parent", errors.New(msg))
	}

	uri := fmt.Sprintf(publicURL, bucket, filename)

	rfile.ResourceUrl = uri
	rfile.Name = name

	return nil
}

func (manager fileManager) Update(ctx context.Context, res page.Resource, bundle []byte) error {
	return page.NewUnsupportedError()
}

func (manager fileManager) Delete(ctx context.Context, res page.Resource) error {
	return page.NewUnsupportedError()
}
