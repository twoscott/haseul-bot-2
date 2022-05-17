package dctools

import "github.com/diamondburned/arikawa/v3/discord"

// HasAnyPerm returns whether the target permissions includes any one of
// the required permissions.
func HasAnyPerm(
	targetPermissions discord.Permissions,
	requiredPermissions discord.Permissions) bool {

	return targetPermissions&requiredPermissions > 0x0
}

// HasAnyPermOrAdmin returns whether the target permissions includes any one of
// the required permissions or the Administrator permission.
func HasAnyPermOrAdmin(
	targetPermissions discord.Permissions,
	requiredPermissions discord.Permissions) bool {

	requiredPermissions.Add(discord.PermissionAdministrator)
	return HasAnyPerm(targetPermissions, requiredPermissions)
}
