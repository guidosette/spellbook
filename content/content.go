package content

import (
	"distudio.com/mage/model"
	"distudio.com/page"
	"encoding/json"
	"strings"
	"time"
)

var ZeroTime = time.Time{}

type ByOrder []*Content

/**
ByOrder start
*/
func (content ByOrder) Len() int {
	return len(content)
}

func (content ByOrder) Swap(i, j int) {
	content[i], content[j] = content[j], content[i]
}

func (content ByOrder) Less(i, j int) bool {
	return content[i].Order < content[j].Order
}

/**
ByOrder end
*/

type Content struct {
	model.Model `json:"-"`
	Type        page.ContentType
	Slug        string
	Title       string `model:"search"`
	Subtitle    string `model:"search"`
	Body        string `model:"search,noindex,HTML"`
	Tags        string `model:"search"`
	Category    string `model:"search,atom";page:"gettable,category"`
	Topic       string `model:"search"`
	Locale      string
	Description string
	Cover       string
	Revision    int
	Order       int
	Attachments []*Attachment `model:"-"`
	// username of the author
	Author    string `model:"search"`
	Created   time.Time
	Updated   time.Time
	Published time.Time

	// KeyTypeEvent
	Date time.Time
}

func (content *Content) UnmarshalJSON(data []byte) error {

	alias := struct {
		Type        page.ContentType `json:"type"`
		Slug        string           `json:"slug"`
		Title       string           `json:"title"`
		Subtitle    string           `json:"subtitle"`
		Body        string           `json:"body"`
		Tags        []string         `json:"tags"`
		Category    string           `json:"category"`
		Topic       string           `json:"topic"`
		Locale      string           `json:"locale"`
		Description string           `json:"description"`
		Revision    int              `json:"revision"`
		Order       int              `json:"order"`
		Attachments []*Attachment    `json:"attachments"`
		Author      string           `json:"author"`
		Cover       string           `json:"cover"`
		Created     time.Time        `json:"created"`
		Updated     time.Time        `json:"updated"`
		Published   time.Time        `json:"published"`
		IsPublished bool             `json:"isPublished"`
		Date        time.Time        `json:"date"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	content.Type = alias.Type
	content.Slug = alias.Slug
	content.Title = alias.Title
	content.Subtitle = alias.Subtitle
	content.Body = alias.Body
	content.Category = alias.Category
	content.Topic = alias.Topic
	content.Locale = alias.Locale
	content.Description = alias.Description
	content.Revision = alias.Revision
	content.Order = alias.Order
	content.Author = alias.Author
	content.Cover = alias.Cover
	content.Attachments = alias.Attachments
	content.Created = alias.Created
	content.Updated = alias.Updated
	content.Date = alias.Date
	if alias.IsPublished {
		content.Published = time.Now().UTC()
	}
	content.Tags = strings.Join(alias.Tags[:], ";")

	return nil
}

func (content *Content) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Type        page.ContentType `json:"type"`
		Slug        string           `json:"slug"`
		Title       string           `json:"title"`
		Subtitle    string           `json:"subtitle"`
		Body        string           `json:"body"`
		Tags        []string         `json:"tags"`
		Category    string           `json:"category"`
		Topic       string           `json:"topic"`
		Locale      string           `json:"locale"`
		Description string           `json:"description"`
		Revision    int              `json:"revision"`
		Order       int              `json:"order"`
		Attachments []*Attachment    `json:"attachments"`
		Author      string           `json:"author"`
		Cover       string           `json:"cover"`
		Created     time.Time        `json:"created"`
		Updated     time.Time        `json:"updated"`
		Published   time.Time        `json:"published"`
		Id          string           `json:"id"`
		Date        time.Time        `json:"date"`
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
			Type:        content.Type,
			Slug:        content.Slug,
			Title:       content.Title,
			Subtitle:    content.Subtitle,
			Body:        content.Body,
			Category:    content.Category,
			Topic:       content.Topic,
			Locale:      content.Locale,
			Description: content.Description,
			Cover:       content.Cover,
			Revision:    content.Revision,
			Order:       content.Order,
			Attachments: content.Attachments,
			Author:      content.Author,
			Created:     content.Created,
			Updated:     content.Updated,
			Published:   content.Published,
			Date:        content.Date,
			Id:          content.EncodedKey(),
		},
	})
}

/**
* Resource representation
 */

func (content *Content) Id() string {
	return content.StringID()
}

func (content *Content) FromRepresentation(rtype page.RepresentationType, data []byte) error {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Unmarshal(data, content)
	}
	return page.NewUnsupportedError()
}

func (content *Content) ToRepresentation(rtype page.RepresentationType) ([]byte, error) {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Marshal(content)
	}
	return nil, page.NewUnsupportedError()
}
