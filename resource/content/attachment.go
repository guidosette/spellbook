package content

import (
	"distudio.com/mage/model"
	"encoding/json"
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
	model.Model `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description";model:"noindex"`
	ResourceUrl string    `json:"resourceUrl";model:"noindex"`
	Group       string    `json:"group"`
	Type        string    `json:"type"`
	Parent      string    `json:"parent"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Uploader    string    `json:"uploader"`
}

func (attachment *Attachment) UnmarshalJSON(data []byte) error {

	alias := struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		ResourceUrl string    `json:"resourceUrl"`
		Group       string    `json:"group"`
		Type        string    `json:"type"`
		Parent      string    `json:"parent"`
		Created     time.Time `json:"created"`
		Updated     time.Time `json:"updated"`
		Uploader    string    `json:"uploader"`
		//Id          int64     `json:"id"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	attachment.Name = alias.Name
	attachment.Description = alias.Description
	attachment.ResourceUrl = alias.ResourceUrl
	attachment.Group = alias.Group
	attachment.Type = alias.Type
	attachment.Parent = alias.Parent
	attachment.Created = alias.Created
	attachment.Updated = alias.Updated
	attachment.Uploader = alias.Uploader
	//attachment.Id = alias.Id

	return nil
}

func (attachment *Attachment) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		ResourceUrl string    `json:"resourceUrl"`
		Group       string    `json:"group"`
		Type        string    `json:"type"`
		Parent      string    `json:"parent"`
		Created     time.Time `json:"created"`
		Updated     time.Time `json:"updated"`
		Uploader    string    `json:"uploader"`
		Id          int64     `json:"id"`
	}

	return json.Marshal(&struct {
		Alias
	}{
		Alias{
			Name:        attachment.Name,
			Description: attachment.Description,
			ResourceUrl: attachment.ResourceUrl,
			Group:       attachment.Group,
			Type:        attachment.Type,
			Parent:      attachment.Parent,
			Created:     attachment.Created,
			Updated:     attachment.Updated,
			Uploader:    attachment.Uploader,
			Id:          attachment.IntID(),
		},
	})
}
