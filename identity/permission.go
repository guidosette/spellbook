package identity

type Permission int64

const (
	PermissionEnabled = 1 << iota
	PermissionReadUser
	PermissionCreateUser
	PermissionEditUser
	PermissionBlockUser
	PermissionReadPost
	PermissionCreatePost
	PermissionEditPost
	PermissionPublishPost
)

var Permissions = map[Permission]string{
	PermissionEnabled: "PERMISSION_ENABLED",
	PermissionReadUser: "PERMISSION_READ_USERS",
	PermissionCreateUser: "PERMISSION_CREATE_USERS",
	PermissionEditUser: "PERMISSION_UPDATE_USERS",
	PermissionBlockUser: "PERMISSION_BLOCK_USERS",
	PermissionReadPost: "PERMISSION_READ_POSTS",
	PermissionCreatePost: "PERMISSION_CREATE_POSTS",
	PermissionEditPost: "PERMISSION_UPDATE_POSTS",
	PermissionPublishPost: "PERMISSION_PUBLISH_POSTS",
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

