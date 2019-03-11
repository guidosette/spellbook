package content

import (
	"distudio.com/mage/model"
	"encoding/json"
	"strings"
	"time"
)

var ZeroTime = time.Time{}

type Content struct {
	model.Model
	Slug        string
	Name        string `model:"search"`
	Title       string `model:"search"`
	Subtitle    string `model:"search"`
	Body        string `model:"search,noindex"`
	Tags        string `model:"search"`
	Category    string `model:"search";page:"gettable,category"`
	Topic       string `model:"search"`
	Locale      string
	Cover       string
	Revision    int
	Attachments []*Attachment `model:"-"`
	// username of the author
	Author    string `model:"search"`
	Created   time.Time
	Updated   time.Time
	Published time.Time
}

func (content *Content) UnmarshalJSON(data []byte) error {

	alias := struct {
		Slug        string        `json:"slug"`
		Name        string        `json:"name"`
		Title       string        `json:"title"`
		Subtitle    string        `json:"subtitle"`
		Body        string        `json:"body"`
		Tags        []string      `json:"tags"`
		Category    string        `json:"category"`
		Topic       string        `json:"topic"`
		Locale      string        `json:"locale"`
		Revision    int           `json:"revision"`
		Attachments []*Attachment `json:"attachments"`
		Author      string        `json:"author"`
		Cover       string        `json:"cover"`
		Created     time.Time     `json:"created"`
		Updated     time.Time     `json:"updated"`
		Published   time.Time     `json:"published"`
		IsPublished bool          `json:"isPublished"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	content.Slug = alias.Slug
	content.Name = alias.Name
	content.Title = alias.Title
	content.Subtitle = alias.Subtitle
	content.Body = alias.Body
	content.Category = alias.Category
	content.Topic = alias.Topic
	content.Locale = alias.Locale
	content.Revision = alias.Revision
	content.Author = alias.Author
	content.Cover = alias.Cover
	content.Attachments = alias.Attachments
	content.Created = alias.Created
	content.Updated = alias.Updated
	if alias.IsPublished {
		content.Published = time.Now().UTC()
	}
	content.Tags = strings.Join(alias.Tags[:], ";")

	return nil
}

func (content *Content) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Slug        string        `json:"slug"`
		Name        string        `json:"name"`
		Title       string        `json:"title"`
		Subtitle    string        `json:"subtitle"`
		Body        string        `json:"body"`
		Tags        []string      `json:"tags"`
		Category    string        `json:"category"`
		Topic       string        `json:"topic"`
		Locale      string        `json:"locale"`
		Revision    int           `json:"revision"`
		Attachments []*Attachment `json:"attachments"`
		Author      string        `json:"author"`
		Cover       string        `json:"cover"`
		Created     time.Time     `json:"created"`
		Updated     time.Time     `json:"updated"`
		Published   time.Time     `json:"published"`
	}

	tags := make([]string, 0, 0)
	if len(content.Tags) > 0 {
		tags = strings.Split(content.Tags, ";")
	}

	isPublished := content.Published != ZeroTime

	return json.Marshal(&struct {
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"isPublished"`
		Alias
	}{
		tags,
		isPublished,
		Alias{
			Slug:        content.Slug,
			Name:        content.Name,
			Title:       content.Title,
			Subtitle:    content.Subtitle,
			Body:        content.Body,
			Category:    content.Category,
			Topic:       content.Topic,
			Locale:      content.Locale,
			Cover:       content.Cover,
			Revision:    content.Revision,
			Attachments: content.Attachments,
			Author:      content.Author,
			Created:     content.Created,
			Updated:     content.Updated,
			Published:   content.Published,
		},
	})
}
