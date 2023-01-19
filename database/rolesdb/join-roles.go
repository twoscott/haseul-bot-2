package rolesdb

import "github.com/diamondburned/arikawa/v3/discord"

const (
	createJoinRolesTableQuery = `
		CREATE TABLE IF NOT EXISTS JoinRoles(
			guildID     INT8 NOT NULL,
			roleID      INT8 NOT NULL,
			PRIMARY KEY(guildID, roleID)
		)`
	addJoinRoleQuery = `
		INSERT INTO JoinRoles VALUES ($1, $2) ON CONFLICT DO NOTHING`
	removeJoinRoleQuery = `
		DELETE FROM JoinRoles WHERE guildID = $1 AND roleID = $2`
	clearGuildJoinRolesQuery = `
		DELETE FROM JoinRoles WHERE guildID = $1`
	getAllGuildJoinRolesQuery = `
		SELECT roleID FROM JoinRoles WHERE guildID = $1`
)

// AddTier adds a role to the list of roles assigned to new users
// who join.
func (db *DB) AddJoinRole(
	guildID discord.GuildID, roleID discord.RoleID) (bool, error) {

	res, err := db.Exec(addJoinRoleQuery, guildID, roleID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// RemoveJoinRole removes a role from the list of roles assigned to new users
// who join.
func (db *DB) RemoveJoinRole(
	guildID discord.GuildID, roleID discord.RoleID) (bool, error) {

	res, err := db.Exec(removeJoinRoleQuery, guildID, roleID)
	if err != nil {
		return false, err
	}

	removed, err := res.RowsAffected()
	return removed > 0, err
}

// ClearGuildJoinRoles removes all roles from the list of roles assigned to new
// users who join.
func (db *DB) ClearGuildJoinRoles(guildID discord.GuildID) (int64, error) {
	res, err := db.Exec(clearGuildJoinRolesQuery, guildID)
	if err != nil {
		return 0, err
	}

	cleared, err := res.RowsAffected()
	return cleared, err
}

// GetAllGuildJoinRoles returns all join role IDs added to a guild.
func (db *DB) GetAllGuildJoinRoles(
	guildID discord.GuildID) ([]discord.RoleID, error) {

	var roleIDs []discord.RoleID
	err := db.Select(&roleIDs, getAllGuildJoinRolesQuery, guildID)

	return roleIDs, err
}
