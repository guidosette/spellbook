package post

import (
	"distudio.com/mage/model"
	"time"
)

type Post struct {
	model.Model
	Slug string `json:"slug"`
	Title string `model:"search";json:"title"`
	Subtitle string `model:"search";json:"subtitle"`
	Body string `model:"search";json:"body"`
	Tags string `model:"search";json:"tags"`
	Category string `model:"search";json:"category"`
	Topic string `model:"search";json:"topic"`
	Locale string `json:"locale"`
	Revision int `json:"revision"`
	// username of the author
	Author string `model:"search";json:"author"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
