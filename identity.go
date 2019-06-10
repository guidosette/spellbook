package page

import "context"

type Permission int64

const (
	HeaderToken string = "X-Authentication"
	keyUser     string = "__pUser__"
)

const (
	PermissionEnabled = 1 << iota
	PermissionReadUser
	PermissionCreateUser
	PermissionEditUser
	PermissionEditPermissions
	PermissionBlockUser
	PermissionReadContent
	PermissionCreateContent
	PermissionEditContent
	PermissionPublishContent
	PermissionReadMailMessage
	PermissionEditMailMessage
	PermissionReadPlace
	PermissionCreatePlace
	PermissionEditPlace
	PermissionReadMedia
	PermissionCreateMedia
	PermissionEditMedia
	PermissionReadSeo
	PermissionCreateSeo
	PermissionEditSeo
)

var Permissions = map[Permission]string{
	PermissionEnabled:         "PERMISSION_ENABLED",
	PermissionReadUser:        "PERMISSION_READ_USERS",
	PermissionCreateUser:      "PERMISSION_CREATE_USERS",
	PermissionEditUser:        "PERMISSION_UPDATE_USERS",
	PermissionEditPermissions: "PERMISSION_EDIT_PERMISSIONS",
	PermissionBlockUser:       "PERMISSION_BLOCK_USERS",
	PermissionReadContent:     "PERMISSION_READ_CONTENT",
	PermissionCreateContent:   "PERMISSION_CREATE_CONTENT",
	PermissionEditContent:     "PERMISSION_UPDATE_CONTENT",
	PermissionPublishContent:  "PERMISSION_PUBLISH_CONTENT",
	PermissionReadMailMessage: "PERMISSION_READ_MAILMESSAGE",
	PermissionEditMailMessage: "PERMISSION_EDIT_MAILMESSAGE",
	PermissionReadPlace:       "PERMISSION_READ_PLACE",
	PermissionCreatePlace:     "PERMISSION_CREATE_PLACE",
	PermissionEditPlace:       "PERMISSION_EDIT_PLACE",
	PermissionReadMedia:       "PERMISSION_READ_MEDIA",
	PermissionCreateMedia:     "PERMISSION_CREATE_MEDIA",
	PermissionEditMedia:       "PERMISSION_EDIT_MEDIA",
	PermissionReadSeo:         "PERMISSION_READ_SEO",
	PermissionCreateSeo:       "PERMISSION_CREATE_SEO",
	PermissionEditSeo:         "PERMISSION_EDIT_SEO",
}

func PermissionName(permission Permission) string {
	return Permissions[permission]
}

func NamedPermissionToPermission(name string) Permission {
	for permission, n := range Permissions {
		if n == name {
			return permission
		}
	}
	return Permission(0)
}

type Identity interface {
	HasPermission(permission Permission) bool
}

func IdentityFromContext(ctx context.Context) Identity {
	if id := ctx.Value(keyUser); id != nil {
		return id.(Identity)
	}
	return nil
}

func ContextWithIdentity(ctx context.Context, id Identity) context.Context {
	return context.WithValue(ctx, keyUser, id)
}
