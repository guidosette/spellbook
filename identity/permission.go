package identity

type Permission int64

const (
	PermissionEnabled = 1 << iota
	PermissionLoadFiles
	PermissionReadUser
	PermissionCreateUser
	PermissionEditUser
	PermissionBlockUser
	PermissionReadContent
	PermissionCreateContent
	PermissionEditContent
	PermissionPublishContent
	PermissionReadNewsletter
	PermissionEditNewsletter
)

var Permissions = map[Permission]string{
	PermissionEnabled: "PERMISSION_ENABLED",
	PermissionLoadFiles: "PERMISSION_LOAD_FILES",
	PermissionReadUser: "PERMISSION_READ_USERS",
	PermissionCreateUser: "PERMISSION_CREATE_USERS",
	PermissionEditUser: "PERMISSION_UPDATE_USERS",
	PermissionBlockUser: "PERMISSION_BLOCK_USERS",
	PermissionReadContent: "PERMISSION_READ_CONTENT",
	PermissionCreateContent: "PERMISSION_CREATE_CONTENT",
	PermissionEditContent: "PERMISSION_UPDATE_CONTENT",
	PermissionPublishContent: "PERMISSION_PUBLISH_CONTENT",
	PermissionReadNewsletter: "PERMISSION_READ_NEWSLETTER",
	PermissionEditNewsletter: "PERMISSION_EDIT_NEWSLETTER",
}

func NamedPermissionToPermission(name string) Permission {
	for permission, n := range Permissions {
		if n == name {
			return permission
		}
	}
	return Permission(0)
}

func (user User) HasPermission(permission Permission) bool {
	return user.Permission & permission != 0
}

func (user *User) GrantPermission(permission Permission) {
	user.Permission |= permission
}

func (user *User) GrantNamedPermission(name string) {
	permission := NamedPermissionToPermission(name)
	user.GrantPermission(permission)
}

func (user *User) GrantNamedPermissions(names []string) {
	for _, name := range names {
		user.GrantNamedPermission(name)
	}
}

func (user *User) GrantAll() {
	for permission := range Permissions {
		user.GrantPermission(permission)
	}
}

func (user *User) RemovePermission(permission Permission) {
	user.Permission &= ^permission
}

func (user *User) TogglePermission(permission Permission) {
	user.Permission ^= permission
}

func (user User) IsEnabled() bool {
	return user.HasPermission(PermissionEnabled)
}

func (user *User) Ban() {
	user.RemovePermission(PermissionEnabled)
}

