package content

import (
	"cloud.google.com/go/storage"
	"context"
	"distudio.com/page"
	"distudio.com/page/identity"
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

func (manager fileManager) Save(ctx context.Context, res page.Resource) error {
	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionLoadFiles) {
		return page.NewPermissionError(identity.PermissionName(identity.PermissionLoadFiles))
	}

	if err := res.Create(ctx); err != nil {
		return err
	}

	return nil
}

func (manager fileManager) Delete(ctx context.Context, res page.Resource) error {
	return page.NewUnsupportedError()
}
