package identity

import "distudio.com/page"

func (user *User) GrantPermission(permission page.Permission) {
	user.Permission |= permission
}

func (user *User) GrantNamedPermission(name string) {
	permission := page.NamedPermissionToPermission(name)
	user.GrantPermission(permission)
}

func (user *User) GrantNamedPermissions(names []string) {
	for _, name := range names {
		user.GrantNamedPermission(name)
	}
}

func (user *User) GrantAll() {
	for permission := range page.Permissions {
		user.GrantPermission(permission)
	}
}

func (user *User) RemovePermission(permission page.Permission) {
	user.Permission &= ^permission
}

func (user *User) TogglePermission(permission page.Permission) {
	user.Permission ^= permission
}

func (user User) IsEnabled() bool {
	return user.HasPermission(page.PermissionEnabled)
}

func (user *User) Ban() {
	user.RemovePermission(page.PermissionEnabled)
}

/**
comparison without PermissionEnabled
 **/
func (user User) ChangedPermission(oldUser User) bool {
	if oldUser.HasPermission(page.PermissionEnabled) {
		oldUser.RemovePermission(page.PermissionEnabled)
	}
	if user.HasPermission(page.PermissionEnabled) {
		user.RemovePermission(page.PermissionEnabled)
	}
	changed := user.Permission != oldUser.Permission
	return changed
}
