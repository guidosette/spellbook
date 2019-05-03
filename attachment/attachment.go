package attachment

import (
	"distudio.com/mage/model"
	"distudio.com/page"
	"distudio.com/page/identity"
	"distudio.com/page/validators"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
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

func (attachment *Attachment) Create(ctx context.Context) error {
	current, _ := ctx.Value(identity.KeyUser).(identity.User)
	// todo permission?
	//if !current.HasPermission(identity.PermissionCreateContent) {
	//	return resource.NewPermissionError(identity.PermissionCreateContent)
	//}

	// attachment parent is required.
	// if not attachment is to be specified the default value must be used
	if attachment.Parent == "" {
		msg := fmt.Sprintf("attachment parent can't be empty. Use %s as a parent for global attachments", AttachmentGlobalParent)
		return validators.NewFieldError("parent", errors.New(msg))
	}

	attachment.Created = time.Now().UTC()
	attachment.Uploader = current.Username()

	return nil
}

func (attachment *Attachment) Update(ctx context.Context, res page.Resource) error {
	// todo permission?
	//current, _ := ctx.Value(identity.KeyUser).(identity.User)
	//if !current.HasPermission(identity.PermissionEditContent) {
	//	return resource.NewPermissionError(identity.PermissionEditContent)
	//}

	other := res.(*Attachment)
	attachment.Name = other.Name
	attachment.Description = other.Description
	attachment.ResourceUrl = other.ResourceUrl
	attachment.Group = other.Group
	attachment.Parent = other.Parent
	attachment.Updated = time.Now().UTC()

	if attachment.Parent == "" {
		msg := fmt.Sprintf("attachment parent can't be empty. Use %s as a parent for global attachments", AttachmentGlobalParent)
		return validators.NewFieldError("parent", errors.New(msg))
	}

	return nil
}

func (attachment *Attachment) Id() string {
	return attachment.StringID()
}
