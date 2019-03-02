package post

import (
	"distudio.com/mage/model"
	"encoding/json"
	"strings"
	"time"
)

var ZeroTime = time.Time{}

type Post struct {
	model.Model
	Slug             string
	Name             string `model:"search"`
	Title            string `model:"search"`
	Subtitle         string `model:"search"`
	Body             string `model:"search,noindex"`
	Tags             string `model:"search"`
	Category         string `model:"search";page:"gettable,category"`
	Topic            string `model:"search"`
	Locale           string
	Cover            string
	Revision         int
	Attachments      []*Attachment `model:"-"`
	// username of the author
	Author    string `model:"search"`
	Created   time.Time
	Updated   time.Time
	Published time.Time
}

func (post *Post) UnmarshalJSON(data []byte) error {

	alias := struct {
		Slug        string       `json:"slug"`
		Name        string       `json:"name"`
		Title       string       `json:"title"`
		Subtitle    string       `json:"subtitle"`
		Body        string       `json:"body"`
		Tags        []string     `json:"tags"`
		Category    string       `json:"category"`
		Topic       string       `json:"topic"`
		Locale      string       `json:"locale"`
		Revision    int          `json:"revision"`
		Attachments []*Attachment `json:"attachments"`
		Author      string       `json:"author"`
		Cover       string       `json:"cover"`
		Created     time.Time    `json:"created"`
		Updated     time.Time    `json:"updated"`
		Published   time.Time    `json:"published"`
		IsPublished bool         `json:"isPublished"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	post.Slug = alias.Slug
	post.Name = alias.Name
	post.Title = alias.Title
	post.Subtitle = alias.Subtitle
	post.Body = alias.Body
	post.Category = alias.Category
	post.Topic = alias.Topic
	post.Locale = alias.Locale
	post.Revision = alias.Revision
	post.Author = alias.Author
	post.Cover = alias.Cover
	post.Attachments = alias.Attachments
	post.Created = alias.Created
	post.Updated = alias.Updated
	post.Published = alias.Published
	if alias.IsPublished {
		post.Published = time.Now().UTC()
	}
	post.Tags = strings.Join(alias.Tags[:], ";")

	return nil
}

func (post *Post) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Slug        string       `json:"slug"`
		Name        string       `json:"name"`
		Title       string       `json:"title"`
		Subtitle    string       `json:"subtitle"`
		Body        string       `json:"body"`
		Tags        []string     `json:"tags"`
		Category    string       `json:"category"`
		Topic       string       `json:"topic"`
		Locale      string       `json:"locale"`
		Revision    int          `json:"revision"`
		Attachments []*Attachment `json:"attachments"`
		Author      string       `json:"author"`
		Cover       string       `json:"cover"`
		Created     time.Time    `json:"created"`
		Updated     time.Time    `json:"updated"`
		Published   time.Time    `json:"published"`
	}

	tags := strings.Split(post.Tags, ";")
	isPublished := post.Published != ZeroTime

	return json.Marshal(&struct {
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"isPublished"`
		Alias
	}{
		tags,
		isPublished,
		Alias{
			Slug:        post.Slug,
			Name:        post.Name,
			Title:       post.Title,
			Subtitle:    post.Subtitle,
			Body:        post.Body,
			Category:    post.Category,
			Topic:       post.Topic,
			Locale:      post.Locale,
			Cover:       post.Cover,
			Revision:    post.Revision,
			Attachments: post.Attachments,
			Author:      post.Author,
			Created:     post.Created,
			Updated:     post.Updated,
			Published:   post.Published,
		},
	})
}
