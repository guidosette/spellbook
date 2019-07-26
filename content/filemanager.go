package content

import (
	"cloud.google.com/go/storage"
	"context"
	"distudio.com/mage"
	"distudio.com/page"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
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
	if current := page.IdentityFromContext(ctx); current == nil || (!current.HasPermission(page.PermissionReadContent) && !current.HasPermission(page.PermissionReadMedia)) {
		var p page.Permission
		p = page.PermissionReadContent
		if !current.HasPermission(page.PermissionReadMedia) {
			p = page.PermissionReadMedia
		}
		return nil, page.NewPermissionError(page.PermissionName(p))
	}

	//todo: set bucket name in configuration
	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving default bucket %s", err.Error())
	}

	res, _ := manager.NewResource(ctx)
	f := res.(*File)
	f.Name = id
	f.ResourceUrl = fmt.Sprintf(publicURL, bucket, id)

	return f, nil
}

func (manager fileManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {

	if current := page.IdentityFromContext(ctx); current == nil || (!current.HasPermission(page.PermissionReadContent) && !current.HasPermission(page.PermissionReadMedia)) {
		var p page.Permission
		p = page.PermissionReadContent
		if !current.HasPermission(page.PermissionReadMedia) {
			p = page.PermissionReadMedia
		}
		return nil, page.NewPermissionError(page.PermissionName(p))
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
	objs, err := manager.listPagination(ctx, it, opts)
	if err != nil {
		log.Errorf(ctx, "listBucket: unable to list bucket %q: %v", bucket, err)
		return nil, err
	}
	for _, obj := range objs {
		name := obj.Name
		s := strings.Split(obj.Name, "/")
		if len(s) > 0 {
			name = s[len(s)-1]
		}
		// create file resource
		res, _ := manager.NewResource(ctx)
		f := res.(*File)
		f.Name = name
		f.ResourceUrl = obj.MediaLink
		f.ContentType = obj.ContentType
		// append file
		files = append(files, f)
	}

	resources := make([]page.Resource, len(files))
	for i := range files {
		resources[i] = page.Resource(files[i])
	}

	return resources, nil
}

func (manager fileManager) listPagination(ctx context.Context, it *storage.ObjectIterator, opts page.ListOptions) ([]*storage.ObjectAttrs, error) {
	p := iterator.NewPager(it, opts.Size+1, "")
	var objs []*storage.ObjectAttrs
	for i := 0; i < opts.Page+1; i++ {
		objs = make([]*storage.ObjectAttrs, 0, 0)
		nextPageToken, err := p.NextPage(&objs)
		if err != nil {
			return nil, err
		}
		if nextPageToken == "" {
			// end pagination
			if i != opts.Page {
				// page requested is out of bound
				objs = make([]*storage.ObjectAttrs, 0, 0)
			}
			break
		}
	}
	return objs, nil
}

func (manager fileManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	return nil, page.NewUnsupportedError()
}

func (manager fileManager) Create(ctx context.Context, res page.Resource, bundle []byte) error {

	if current := page.IdentityFromContext(ctx); current == nil || (!current.HasPermission(page.PermissionWriteContent) && !current.HasPermission(page.PermissionWriteMedia)) {
		var p page.Permission
		p = page.PermissionWriteContent
		if !current.HasPermission(page.PermissionWriteMedia) {
			p = page.PermissionWriteMedia
		}
		return page.NewPermissionError(page.PermissionName(p))
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
	if len(fhs) == 0 {
		msg := fmt.Sprintf("error no file: %v", fhs)
		return page.NewFieldError("no file", errors.New(msg))
	}
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
		return page.NewFieldError("bucket", errors.New(msg))
	}

	if err := writer.Close(); err != nil {
		msg := fmt.Sprintf("upload: unable to close bucket %s: %s", bucket, err.Error())
		return page.NewFieldError("bucket", errors.New(msg))
	}

	uri := fmt.Sprintf(publicURL, bucket, filename)

	rfile.ResourceUrl = uri
	rfile.Name = name

	// -----------------------------thumbnail
	if strings.Contains(writer.ContentType, "image/") {
		// get image
		image, err := imaging.Decode(f)
		if err != nil {
			msg := fmt.Sprintf("error in opening image %s", err)
			return page.NewFieldError("bucket", errors.New(msg))
		}
		// create thumbnail
		fileNameThumbnail := fmt.Sprintf("%s%s/thumb/%s", typ, namespace, name)
		afterImage := imaging.Thumbnail(image, 100, 100, imaging.Linear)
		// Save thumbnail
		wc := handle.Object(fileNameThumbnail).NewWriter(ctx)
		wc.ContentType = fh.Header.Get("Content-Type")
		if imaging.Encode(wc, afterImage, imaging.JPEG); err != nil {
			msg := fmt.Sprintf("%s in saving image thumbnail", err.Error())
			return page.NewFieldError("bucket", errors.New(msg))
		}
		if err = wc.Close(); err != nil {
			msg := fmt.Sprintf("CreateFileThumbnail: unable to close bucket %q, file %q: %v", bucket, fileNameThumbnail, err)
			return page.NewFieldError("bucket", errors.New(msg))
		}

		uriThumb := fmt.Sprintf(publicURL, bucket, fileNameThumbnail)
		rfile.ResourceThumbUrl = uriThumb
	}

	return nil
}

func (manager fileManager) Update(ctx context.Context, res page.Resource, bundle []byte) error {
	return page.NewUnsupportedError()
}

func (manager fileManager) Delete(ctx context.Context, res page.Resource) error {
	return page.NewUnsupportedError()
}
