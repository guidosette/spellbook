package post

import (
	"distudio.com/mage/model"
	"time"
)

type Multimedia struct {
	model.Model `json:"-"`
	Name string `json:"name"`
	Description string `json:"description"`
	ResourceUrl string `json:"resourceUrl"`
	Group string `json:"group"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Uploader string `json:"uploader"`
}
