package rolesdb

import (
	"database/sql"

	"github.com/diamondburned/arikawa/v3/discord"
)

type Role struct {
	ID          discord.RoleID `db:"roleid"`
	TierID      int32          `db:"tierid"`
	Description sql.NullString `db:"description"`
}

const (
	createRolePickerTableQuery = `
		CREATE TABLE IF NOT EXISTS RolePicker(
			roleID INT8              NOT NULL,
			tierID INT4              NOT NULL,
			description VARCHAR(100),
			PRIMARY KEY(roleID, tierID),
			FOREIGN KEY(tierID) REFERENCES RoleTiers(id) ON DELETE CASCADE
		)`
	addRoleQuery = `
		INSERT INTO RolePicker VALUES($1, $2, $3) ON CONFLICT DO NOTHING`
	removeRoleQuery = `
		DELETE FROM RolePicker WHERE roleID = $1 AND tierID = $2`
	getRolesByTierQuery = `SELECT * FROM RolePicker WHERE tierID = $1`
)

// AddRole adds a role to a role picker tier.
func (db *DB) AddRole(
	roleID discord.RoleID, tierID int32, description string) (bool, error) {

	res, err := db.Exec(addRoleQuery, roleID, tierID, description)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// RemoveRole removes a role from a role picker tier.
func (db *DB) RemoveRole(roleID discord.RoleID, tierID int32) (bool, error) {
	res, err := db.Exec(removeRoleQuery, roleID, tierID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// GetAllRolesByTier returns all the roles for a role picker tier.
func (db *DB) GetAllRolesByTier(tierID int32) ([]Role, error) {
	var roles []Role
	err := db.Select(&roles, getRolesByTierQuery, tierID)

	return roles, err
}
