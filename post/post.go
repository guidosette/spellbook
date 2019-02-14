package post

import (
	"distudio.com/mage/model"
	"encoding/json"
	"strings"
	"time"
)

type Post struct {
	model.Model
	Slug     string `json:"slug"`
	Title    string `model:"search";json:"title"`
	Subtitle string `model:"search";json:"subtitle"`
	Body     string `model:"search";json:"body"`
	Tags     string `model:"search";json:"tags"`
	Category string `model:"search";json:"category";page:"gettable,category"`
	Topic    string `model:"search";json:"topic"`
	Locale   string `json:"locale"`
	Revision int    `json:"revision"`
	// username of the author
	Author    string    `model:"search";json:"author"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Published time.Time `json:"published"`
}

func (post *Post) UnmarshalJSON(data []byte) error {

	alias := struct {
		Slug      string    `json:"slug"`
		Title     string    `json:"title"`
		Subtitle  string    `json:"subtitle"`
		Body      string    `json:"body"`
		Tags      []string  `json:"tags"`
		Category  string    `json:"category"`
		Topic     string    `json:"topic"`
		Locale    string    `json:"locale"`
		Revision  int       `json:"revision"`
		Author    string    `json:"author"`
		Created   time.Time `json:"created"`
		Updated   time.Time `json:"updated"`
		Published time.Time `json:"published"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	post.Slug = alias.Slug
	post.Title = alias.Title
	post.Subtitle = alias.Subtitle
	post.Body = alias.Body
	post.Category = alias.Category
	post.Topic = alias.Topic
	post.Locale = alias.Locale
	post.Revision = alias.Revision
	post.Author = alias.Author
	post.Created = alias.Created
	post.Updated = alias.Updated
	post.Published = alias.Published
	post.Tags = strings.Join(alias.Tags[:], ";")

	return nil
}

func (post *Post) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Slug      string    `json:"slug"`
		Title     string    `json:"title"`
		Subtitle  string    `json:"subtitle"`
		Body      string    `json:"body"`
		Tags      []string  `json:"tags"`
		Category  string    `json:"category"`
		Topic     string    `json:"topic"`
		Locale    string    `json:"locale"`
		Revision  int       `json:"revision"`
		Author    string    `json:"author"`
		Created   time.Time `json:"created"`
		Updated   time.Time `json:"updated"`
		Published time.Time `json:"published"`
	}

	tags := strings.Split(post.Tags, ";")
	return json.Marshal(&struct {
		Tags []string `json:"username"`
		Alias
	}{
		tags,
		Alias{
			Slug:     post.Slug,
			Title:    post.Title,
			Subtitle: post.Subtitle,
			Body:     post.Body,
			Category: post.Category,
			Topic:    post.Topic,
			Locale:   post.Locale,
			Revision: post.Revision,
			Author:   post.Author,
			Created:  post.Created,
			Updated:  post.Updated,
			Published:  post.Published,
		},
	})
}
