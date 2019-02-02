package identity

type Permission int64

const (
	PermissionEnabled = 1 << iota
	PermissionLogIn
	PermissionCreateUser
	PermissionEditUser
	PermissionBlockUser
	PermissionCreatePost
	PermissionEditPost
	PermissionPublishPost
)

func (user User) HasPermission(permission Permission) bool {
	return user.Permission & permission != 0
}

func (user *User) AddPermission(permission Permission) {
	user.Permission |= permission
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

