package post

import (
	"distudio.com/mage/model"
	"encoding/json"
	"time"
)

type Multimedia struct {
	model.Model `json:"-"`
	Name string `json:"name"`
	Description string `json:"description"`
	ResourceUrl string `json:"resource"`
	Group string `json:"group"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Uploader string `json:"uploader"`
}

// add the json interface to multimedia pointers
func (multimedia *Multimedia) UnmarshalJSON(data []byte) error {

	mm := Multimedia{}

	err := json.Unmarshal(data, mm)
	if err != nil {
		return err
	}

	*multimedia = mm

	return nil
}

func (multimedia *Multimedia) MarshalJSON() ([]byte, error) {
	mm := *multimedia
	return json.Marshal(mm)
}
