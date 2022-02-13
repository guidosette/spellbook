package identity

import (
	"decodica.com/spellbook"
)

func (user *User) GrantPermission(permission spellbook.Permission) {
	user.Permission |= permission
}

func (user *User) GrantNamedPermission(name string) {
	permission := spellbook.NamedPermissionToPermission(name)
	user.GrantPermission(permission)
}

func (user *User) GrantNamedPermissions(names []string) {
	for _, name := range names {
		user.GrantNamedPermission(name)
	}
}

func (user *User) GrantAll() {
	for permission := range spellbook.Permissions {
		user.GrantPermission(permission)
	}
}

func (user *User) RemovePermission(permission spellbook.Permission) {
	user.Permission &= ^permission
}

func (user *User) TogglePermission(permission spellbook.Permission) {
	user.Permission ^= permission
}

func (user User) IsEnabled() bool {
	return user.HasPermission(spellbook.PermissionEnabled)
}

func (user *User) Ban() {
	user.RemovePermission(spellbook.PermissionEnabled)
}

/**
comparison without PermissionEnabled
 **/
func (user User) ChangedPermission(oldUser User) bool {
	if oldUser.HasPermission(spellbook.PermissionEnabled) {
		oldUser.RemovePermission(spellbook.PermissionEnabled)
	}
	if user.HasPermission(spellbook.PermissionEnabled) {
		user.RemovePermission(spellbook.PermissionEnabled)
	}
	changed := user.Permission != oldUser.Permission
	return changed
}
