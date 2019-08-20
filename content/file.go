package content

import (
	"decodica.com/spellbook"
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

func (file *File) FromRepresentation(rtype spellbook.RepresentationType, data []byte) error {
	return spellbook.NewUnsupportedError()
}

func (file *File) ToRepresentation(rtype spellbook.RepresentationType) ([]byte, error) {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Marshal(file)
	}
	return nil, spellbook.NewUnsupportedError()
}
