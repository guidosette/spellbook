package content

import (
	"context"
)

type attk string

const attKey attk = "__attachments"

var allowedParents []string

type SupportedAttachments struct{}

var supportedAtt SupportedAttachments

// tests if the parent passed is supported.
// should be used for validation purposes only
func (par SupportedAttachments) IsSupported(attachment *Attachment) bool {
	// we admit empty parent types, if the type is global
	if attachment.ParentType == "" {
		return attachment.ParentKey == AttachmentGlobalParent
	}
	for _, p := range allowedParents {
		if p == attachment.ParentType {
			return true
		}
	}
	return false
}

type AttachmentService struct {
	ParentDeclaration
}

type ParentDeclaration interface {
	// return the supported parents
	DeclareParents() []string
}

const name = "__attachment_service"

// callback methods
func (service *AttachmentService) Name() string {
	return name
}

func (service *AttachmentService) Initialize() {
	// add default supported parents
	allowedParents = append(allowedParents, "content")
	if service.ParentDeclaration != nil {
		parents := service.DeclareParents()
		for _, p := range parents {
			allowedParents = append(allowedParents, p)
		}
	}
}

// adds the appengine client to the context
func (service *AttachmentService) OnStart(ctx context.Context) context.Context {
	return context.WithValue(ctx, attKey, &supportedAtt)
}

func (service *AttachmentService) OnEnd(ctx context.Context) {}

func (service *AttachmentService) Destroy() {}

func SupportedAttachmentsFromContext(ctx context.Context) *SupportedAttachments {
	val := ctx.Value(attKey)
	if val == nil {
		return nil
	}
	return val.(*SupportedAttachments)
}
