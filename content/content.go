package content

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page"
	"distudio.com/page/attachment"
	"distudio.com/page/identity"
	"distudio.com/page/validators"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
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
	Slug        string
	Name        string `model:"search"`
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
	Attachments []*attachment.Attachment `model:"-"`
	// username of the author
	Author    string `model:"search"`
	Created   time.Time
	Updated   time.Time
	Published time.Time
}

func (content *Content) UnmarshalJSON(data []byte) error {

	alias := struct {
		Slug        string                   `json:"slug"`
		Name        string                   `json:"name"`
		Title       string                   `json:"title"`
		Subtitle    string                   `json:"subtitle"`
		Body        string                   `json:"body"`
		Tags        []string                 `json:"tags"`
		Category    string                   `json:"category"`
		Topic       string                   `json:"topic"`
		Locale      string                   `json:"locale"`
		Description string                   `json:"description"`
		Revision    int                      `json:"revision"`
		Order       int                      `json:"order"`
		Attachments []*attachment.Attachment `json:"attachments"`
		Author      string                   `json:"author"`
		Cover       string                   `json:"cover"`
		Created     time.Time                `json:"created"`
		Updated     time.Time                `json:"updated"`
		Published   time.Time                `json:"published"`
		IsPublished bool                     `json:"isPublished"`
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
	content.Description = alias.Description
	content.Revision = alias.Revision
	content.Order = alias.Order
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
		Slug        string                   `json:"slug"`
		Name        string                   `json:"name"`
		Title       string                   `json:"title"`
		Subtitle    string                   `json:"subtitle"`
		Body        string                   `json:"body"`
		Tags        []string                 `json:"tags"`
		Category    string                   `json:"category"`
		Topic       string                   `json:"topic"`
		Locale      string                   `json:"locale"`
		Description string                   `json:"description"`
		Revision    int                      `json:"revision"`
		Order       int                      `json:"order"`
		Attachments []*attachment.Attachment `json:"attachments"`
		Author      string                   `json:"author"`
		Cover       string                   `json:"cover"`
		Created     time.Time                `json:"created"`
		Updated     time.Time                `json:"updated"`
		Published   time.Time                `json:"published"`
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
			Description: content.Description,
			Cover:       content.Cover,
			Revision:    content.Revision,
			Order:       content.Order,
			Attachments: content.Attachments,
			Author:      content.Author,
			Created:     content.Created,
			Updated:     content.Updated,
			Published:   content.Published,
		},
	})
}

func (content *Content) Create(ctx context.Context) error {
	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionCreateContent) {
		return validators.NewPermissionError(identity.PermissionCreateContent)
	}

	content.Created = time.Now().UTC()
	content.Revision = 1
	if !content.Published.IsZero() {
		content.Published = time.Now().UTC()
	}

	if content.Title == "" || content.Name == "" {
		msg := fmt.Sprintf(" title and name can't be empty")
		return validators.NewFieldError("title", errors.New(msg))
	}

	if content.Slug == "" {
		content.Slug = url.PathEscape(content.Title)
	}

	// if the same slug already exists, we must return
	// otherwise we would overwrite an existing entry, which is not in the spirit of the create method
	q := model.NewQuery((*Content)(nil))
	q = q.WithField("Slug =", content.Slug)
	count, err := q.Count(ctx)
	if err != nil {
		msg := fmt.Sprintf("error verifying slug uniqueness: %s", err.Error())
		return validators.NewFieldError("slug", errors.New(msg))
	}

	if count > 0 {
		msg := fmt.Sprintf("a content with slug %s already exists. Slug must be unique.", content.Slug)
		return validators.NewFieldError("slug", errors.New(msg))
	}

	if user, ok := ctx.Value(identity.KeyUser).(identity.User); ok {
		content.Author = user.Username()
	}

	return nil
}

func (content *Content) Update(ctx context.Context, res page.Resource) error {
	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	if !current.HasPermission(identity.PermissionEditContent) {
		return validators.NewPermissionError(identity.PermissionEditContent)
	}

	other := res.(*Content)
	content.Name = other.Name
	content.Title = other.Title
	content.Subtitle = other.Subtitle
	content.Category = other.Category
	content.Topic = other.Topic
	content.Locale = other.Locale
	content.Description = other.Description
	content.Body = other.Body
	content.Cover = other.Cover
	content.Revision = other.Revision
	content.Order = other.Order
	content.Updated = time.Now().UTC()
	content.Tags = other.Tags

	if user, ok := ctx.Value(identity.KeyUser).(identity.User); ok {
		content.Author = user.Username()
	}

	if other.Published.IsZero() {
		// not setted
		content.Published = time.Time{}
	} else {
		// setted
		// check previous data
		if content.Published.IsZero() {
			content.Published = time.Now().UTC()
		}
	}

	return nil
}

func (content *Content) Id() string {
	return content.StringID()
}
