package page

import (
	"cloud.google.com/go/storage"
	"distudio.com/mage"
	"distudio.com/page/content"
	"distudio.com/page/identity"
	"distudio.com/page/validators"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		if !user.HasPermission(identity.PermissionLoadFiles) {
			return mage.Redirect{Status: http.StatusForbidden}
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
		nsv.AddValidator(validators.FileNameValidator{AllowEmpty: true})
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
		filename := fmt.Sprintf("%s%s/%s", typ, namespace, name)

		// handle the upload to Google Cloud Storage
		bucket, err := file.DefaultBucketName(ctx)
		if err != nil {
			log.Errorf(ctx, "can't retrieve bucket name: %s", err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
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
			log.Errorf(ctx, "upload: unable to write file %s to bucket %s: %s", filename, bucket, err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		if err := writer.Close(); err != nil {
			log.Errorf(ctx, "upload: unable to close bucket %s: %s", bucket, err.Error())
		}

		// return the file data
		const publicURL = "https://storage.googleapis.com/%s/%s"
		uri := fmt.Sprintf(publicURL, bucket, filename)
		response := struct {
			URI string
		}{uri}

		renderer := mage.JSONRenderer{}
		renderer.Data = response
		out.Renderer = &renderer

		return mage.Redirect{Status: http.StatusCreated}
	case http.MethodGet:
		// check if current user has permission
		me := ctx.Value(identity.KeyUser)
		current, ok := me.(identity.User)

		if !ok {
			return mage.Redirect{Status: http.StatusUnauthorized}
		}

		if !current.HasPermission(identity.PermissionReadContent) {
			return mage.Redirect{Status: http.StatusForbidden}
		}

		// handle the upload to Google Cloud Storage
		bucket, err := file.DefaultBucketName(ctx)
		if err != nil {
			log.Errorf(ctx, "can't retrieve bucket name: %s", err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Errorf(ctx, "failed to create client: %s", err.Error())
			return mage.Redirect{Status: http.StatusInternalServerError}
		}
		defer client.Close()

		handle := client.Bucket(bucket)

		params := mage.RoutingParams(ctx)
		// if there is no param then it is a list request
		param, ok := params["name"]
		if !ok {
			// list
			// handle query params for page data:
			page := 0
			size := 20
			if pin, ok := ins["page"]; ok {
				if num, err := strconv.Atoi(pin.Value()); err == nil {
					page = num
				} else {
					return mage.Redirect{Status: http.StatusBadRequest}
				}
			}

			if sin, ok := ins["results"]; ok {
				if num, err := strconv.Atoi(sin.Value()); err == nil {
					size = num
					// cap the size to 100
					if size > 100 {
						size = 100
					}
				} else {
					return mage.Redirect{Status: http.StatusBadRequest}
				}
			}

			log.Infof(ctx, "page", page) //todo
			var result interface{}
			l := 0

			files := make([]content.File, 0, 0)
			query := &storage.Query{}
			it := handle.Objects(ctx, query)
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
				file := content.File{Name: name, ResourceUrl: obj.MediaLink}
				files = append(files, file)
			}

			l = len(files)
			result = files[:controller.GetCorrectCountForPaging(size, l)]

			response := struct {
				Items interface{} `json:"items"`
				More  bool        `json:"more"`
			}{result, l > size}
			renderer := mage.JSONRenderer{}
			renderer.Data = response
			out.Renderer = &renderer
			return mage.Redirect{Status: http.StatusOK}
		}

		// get single file by name
		name := param.Value()
		reader, err := handle.Object(name).NewReader(ctx)
		if err != nil {
			log.Errorf(ctx, "readFile: unable to open file from bucket %q, file %q: %v", bucket, name, err)
			return mage.Redirect{Status: http.StatusInternalServerError}
		}
		defer reader.Close()
		slurp, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Errorf(ctx, "readFile: unable to read data from bucket %q, file %q: %v", bucket, name, err)
			return mage.Redirect{Status: http.StatusInternalServerError}
		}

		if len(slurp) > 1024 {
			//fmt.Fprintf(d.w, "...%s\n", slurp[len(slurp)-1024:])
		} else {
			//fmt.Fprintf(d.w, "%s\n", slurp)
		}

		renderer := mage.JSONRenderer{}
		renderer.Data = &slurp
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusOK}

	}
	return mage.Redirect{Status: http.StatusNotImplemented}
}

func (controller *FileController) GetCorrectCountForPaging(size int, l int) int {
	count := size
	if l < size {
		count = l
	}
	return count
}
