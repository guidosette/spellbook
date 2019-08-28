package content

import (
	"decodica.com/flamel/model"
	"decodica.com/spellbook"
	"encoding/json"
	"strings"
	"time"
)

type ByOrder []*Content
type ByTitle []*Content
type ByStartDate []*Content
type ByPublished []*Content

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

/**
ByTitle start
*/
func (content ByTitle) Len() int {
	return len(content)
}

func (content ByTitle) Swap(i, j int) {
	content[i], content[j] = content[j], content[i]
}

func (content ByTitle) Less(i, j int) bool {
	return content[i].Title < content[j].Title
}

/**
ByTitle end
*/

/**
ByStartDate start
*/
func (content ByStartDate) Len() int {
	return len(content)
}

func (content ByStartDate) Swap(i, j int) {
	content[i], content[j] = content[j], content[i]
}

func (content ByStartDate) Less(i, j int) bool {
	return content[i].StartDate.Before(content[j].StartDate)
}

/**
ByStartDate end
*/

/**
ByPublished start
*/
func (content ByPublished) Len() int {
	return len(content)
}

func (content ByPublished) Swap(i, j int) {
	content[i], content[j] = content[j], content[i]
}

func (content ByPublished) Less(i, j int) bool {
	return content[i].Published.Before(content[j].Published)
}

/**
ByPublished end
*/

type PublicationState string

const PublicationStatePublished PublicationState = "PUBLISHED"
const PublicationStateUnpublished PublicationState = "UNPUBLISHED"

type Content struct {
	model.Model `json:"-"`
	Type        string `model:"search"`
	IdTranslate string
	Slug        string
	Title       string `model:"search"`
	Subtitle    string `model:"search"`
	Body        string `model:"search,noindex,HTML"`
	Tags        string `model:"search"`
	Category    string `model:"search,atom";page:"gettable,category"`
	Topic       string `model:"search"`
	Locale      string `model:"search,atom"`
	Description string `model:"search"`
	Cover       string
	Revision    int
	Order       int `model:"search"`
	Attachments []*Attachment `model:"-"`
	// username of the author
	Author           string `model:"search"`
	Editor           string `model:"search"`
	Created          time.Time
	Updated          time.Time `model:"search"`
	Published        time.Time        `model:"search"`
	PublicationState PublicationState `model:"search,atom"`
	ParentKey        string `model:"search,atom"`
	Code             string // special

	// KeyTypeEvent
	StartDate time.Time
	EndDate   time.Time
}

func (content Content) IsPublished() bool {
	return !content.Published.IsZero()
}

func (content *Content) UnmarshalJSON(data []byte) error {

	alias := struct {
		Type        string `json:"type"`
		IdTranslate string                `json:"idTranslate"`
		ParentKey   string                `json:"parentKey"`
		Slug        string                `json:"slug"`
		Title       string                `json:"title"`
		Subtitle    string                `json:"subtitle"`
		Body        string                `json:"body"`
		Tags        []string              `json:"tags"`
		Category    string                `json:"category"`
		Topic       string                `json:"topic"`
		Locale      string                `json:"locale"`
		Description string                `json:"description"`
		Revision    int                   `json:"revision"`
		Order       int                   `json:"order"`
		Attachments []*Attachment         `json:"attachments"`
		Author      string                `json:"author"`
		Editor      string                `json:"editor"`
		Cover       string                `json:"cover"`
		Code        string                `json:"code"`
		Created     time.Time             `json:"created"`
		Updated     time.Time             `json:"updated"`
		Published   time.Time             `json:"published"`
		IsPublished bool                  `json:"isPublished"`
		StartDate   time.Time             `json:"startDate"`
		EndDate     time.Time             `json:"endDate"`
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
	content.Editor = alias.Editor
	content.Attachments = alias.Attachments
	content.Created = alias.Created
	content.Updated = alias.Updated
	content.StartDate = alias.StartDate
	content.EndDate = alias.EndDate
	content.Code = alias.Code
	content.IdTranslate = alias.IdTranslate
	content.ParentKey = alias.ParentKey
	if alias.IsPublished {
		content.Published = time.Now().UTC()
	}
	content.Tags = strings.Join(alias.Tags[:], ";")

	return nil
}

func (content *Content) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Type        string `json:"type"`
		IdTranslate string                `json:"idTranslate"`
		Slug        string                `json:"slug"`
		Title       string                `json:"title"`
		Subtitle    string                `json:"subtitle"`
		Body        string                `json:"body"`
		Tags        []string              `json:"tags"`
		Category    string                `json:"category"`
		Topic       string                `json:"topic"`
		Locale      string                `json:"locale"`
		Description string                `json:"description"`
		Revision    int                   `json:"revision"`
		Order       int                   `json:"order"`
		Attachments []*Attachment         `json:"attachments"`
		Author      string                `json:"author"`
		Editor      string                `json:"editor"`
		Cover       string                `json:"cover"`
		Code        string                `json:"code"`
		Created     time.Time             `json:"created"`
		Updated     time.Time             `json:"updated"`
		Published   time.Time             `json:"published"`
		Key         string                `json:"key"`
		ParentKey   string                `json:"parentKey"`
		StartDate   time.Time             `json:"startDate"`
		EndDate     time.Time             `json:"endDate"`
	}

	tags := make([]string, 0, 0)
	if len(content.Tags) > 0 {
		tags = strings.Split(content.Tags, ";")
	}

	isPublished := content.IsPublished()

	return json.Marshal(&struct {
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"isPublished"`
		Alias
	}{
		tags,
		isPublished,
		Alias{
			Type:        content.Type,
			IdTranslate: content.IdTranslate,
			Slug:        content.Slug,
			Title:       content.Title,
			Subtitle:    content.Subtitle,
			Body:        content.Body,
			Category:    content.Category,
			Topic:       content.Topic,
			Locale:      content.Locale,
			Description: content.Description,
			Cover:       content.Cover,
			Editor:      content.Editor,
			Revision:    content.Revision,
			Order:       content.Order,
			Attachments: content.Attachments,
			Author:      content.Author,
			Created:     content.Created,
			Code:        content.Code,
			Updated:     content.Updated,
			Published:   content.Published,
			StartDate:   content.StartDate,
			EndDate:     content.EndDate,
			Key:         content.EncodedKey(),
			ParentKey:   content.ParentKey,
		},
	})
}

/**
* Resource representation
 */

func (content *Content) Id() string {
	return content.StringID()
}

func (content *Content) FromRepresentation(rtype spellbook.RepresentationType, data []byte) error {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Unmarshal(data, content)
	}
	return spellbook.NewUnsupportedError()
}

func (content *Content) ToRepresentation(rtype spellbook.RepresentationType) ([]byte, error) {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Marshal(content)
	}
	return nil, spellbook.NewUnsupportedError()
}
