package content

import (
	"distudio.com/mage/model"
	"time"
)

const AttachmentGlobalParent = "GLOBAL"

type Attachment struct {
	model.Model `json:"-"`
	Name string `json:"name"`
	Description string `json:"description"`
	ResourceUrl string `json:"resourceUrl"`
	Group string `json:"group"`
	Type string `json:"type"`
	Parent string `json:"parent"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Uploader string `json:"uploader"`
}
