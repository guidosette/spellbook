package content

import (
	"decodica.com/flamel/model"
	"decodica.com/spellbook"
	"encoding/json"
	"fmt"
	"time"
)

const (
	// global name for attachments without parents
	AttachmentGlobalParent = "GLOBAL"

	// supported attachments
	AttachmentTypeGallery    = "gallery"
	AttachmentTypeAttachment = "attachments"
	AttachmentTypeVideo      = "video"
)

type Attachment struct {
	spellbook.GormModel `model:"-"`
	model.Model      `json:"-"`
	Name             string    `json:"name"`
	Description      string    `json:"description";model:"noindex"`
	ResourceUrl      string    `json:"resourceUrl";model:"noindex"`
	ResourceThumbUrl string    `json:"resourceThumbUrl";model:"noindex"`
	Group            string    `json:"group"`
	Type             string    `json:"type"`
	ParentKey        string    `json:"parentKey"` // encode key of content
	Created          time.Time `json:"created"`
	Updated          time.Time `json:"updated"`
	Uploader         string    `json:"uploader"`
	AltText          string    `json:"altText"`
	Seo              int64     `json:"seo"`
}

func (attachment *Attachment) UnmarshalJSON(data []byte) error {

	alias := struct {
		Name             string    `json:"name"`
		Description      string    `json:"description"`
		ResourceUrl      string    `json:"resourceUrl"`
		ResourceThumbUrl string    `json:"resourceThumbUrl"`
		Group            string    `json:"group"`
		Type             string    `json:"type"`
		ParentKey        string    `json:"parentKey"`
		Created          time.Time `json:"created"`
		Updated          time.Time `json:"updated"`
		Uploader         string    `json:"uploader"`
		AltText          string    `json:"altText"`
		Seo              int64     `json:"seo"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	attachment.Name = alias.Name
	attachment.Description = alias.Description
	attachment.ResourceUrl = alias.ResourceUrl
	attachment.ResourceThumbUrl = alias.ResourceThumbUrl
	attachment.Group = alias.Group
	attachment.Type = alias.Type
	attachment.ParentKey = alias.ParentKey
	attachment.Created = alias.Created
	attachment.Updated = alias.Updated
	attachment.Uploader = alias.Uploader
	attachment.AltText = alias.AltText
	attachment.Seo = alias.Seo

	return nil
}

func (attachment *Attachment) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Id               string     `json:"id"`
		Name             string    `json:"name"`
		Description      string    `json:"description"`
		ResourceUrl      string    `json:"resourceUrl"`
		ResourceThumbUrl string    `json:"resourceThumbUrl"`
		Group            string    `json:"group"`
		Type             string    `json:"type"`
		ParentKey        string    `json:"parentKey"`
		Created          time.Time `json:"created"`
		Updated          time.Time `json:"updated"`
		Uploader         string    `json:"uploader"`
		AltText          string    `json:"altText"`
		Seo              int64     `json:"seo"`
	}

	return json.Marshal(&struct {
		Alias
	}{
		Alias{
			Id:               attachment.Id(),
			Name:             attachment.Name,
			Description:      attachment.Description,
			ResourceUrl:      attachment.ResourceUrl,
			ResourceThumbUrl: attachment.ResourceThumbUrl,
			Group:            attachment.Group,
			Type:             attachment.Type,
			ParentKey:        attachment.ParentKey,
			Created:          attachment.Created,
			Updated:          attachment.Updated,
			Uploader:         attachment.Uploader,
			AltText:          attachment.AltText,
			Seo:              attachment.Seo,
		},
	})
}

func (attachment *Attachment) Id() string {
	if id  := attachment.EncodedKey(); id != "" {
		return id
	}
	return fmt.Sprintf("%d", attachment.ID)
}

func (attachment *Attachment) FromRepresentation(rtype spellbook.RepresentationType, data []byte) error {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Unmarshal(data, attachment)
	}
	return spellbook.NewUnsupportedError()
}

func (attachment *Attachment) ToRepresentation(rtype spellbook.RepresentationType) ([]byte, error) {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Marshal(attachment)
	}
	return nil, spellbook.NewUnsupportedError()
}
