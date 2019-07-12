package content

import (
	"distudio.com/page"
	"encoding/json"
)

// return the file data
const publicURL = "https://storage.googleapis.com/%s/%s"

type File struct {
	Name             string `json:"name"`
	ResourceUrl      string `json:"resourceUrl"`
	ResourceThumbUrl string `json:"resourceThumbUrl"`
	ContentType      string `json:"contentType"`
}

func (file *File) Id() string {
	return file.Name
}

func (file *File) FromRepresentation(rtype page.RepresentationType, data []byte) error {
	return page.NewUnsupportedError()
}

func (file *File) ToRepresentation(rtype page.RepresentationType) ([]byte, error) {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Marshal(file)
	}
	return nil, page.NewUnsupportedError()
}
