package rolesdb

import (
	"database/sql"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

type Tier struct {
	ID          int32           `db:"id"`
	GuildID     discord.GuildID `db:"guildid"`
	Name        string          `db:"name"`
	Description sql.NullString  `db:"description"`
}

func (t Tier) Title() string {
	return util.TitleCase(t.Name)
}

const (
	createRoleTiersTableQuery = `
		CREATE TABLE IF NOT EXISTS RoleTiers(
			id          SERIAL,
			guildID     INT8          NOT NULL,
			name        VARCHAR(32)   NOT NULL,
			description VARCHAR(1024),
			PRIMARY KEY(ID),
			UNIQUE(guildID, name)
		)`
	getTierByNameQuery = `
		SELECT * FROM RoleTiers 
		WHERE guildID = $1 AND name ILIKE $2`
	getAllTiersByGuildQuery = `
		SELECT * FROM RoleTiers WHERE guildID = $1`
	addTierQuery = `
		INSERT INTO RoleTiers (guildID, name, description) 
		VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	removeTierQuery = `
		DELETE FROM RoleTiers WHERE guildID = $1 AND name ILIKE $2`
)

// GetTierByName returns a role tier with the provided name.
func (db *DB) GetTierByName(
	guildID discord.GuildID, name string) (*Tier, error) {

	var tier Tier
	err := db.Get(&tier, getTierByNameQuery, guildID, name)

	return &tier, err
}

// GetTiersByGuild returns all tiers added to a given guild.
func (db *DB) GetAllTiersByGuild(guildID discord.GuildID) ([]Tier, error) {
	var tiers []Tier
	err := db.Select(&tiers, getAllTiersByGuildQuery, guildID)

	return tiers, err
}

// AddTier adds a role tier to the provided guild with the provided name.
func (db *DB) AddTier(
	guildID discord.GuildID, name, description string) (bool, error) {

	res, err := db.Exec(addTierQuery, guildID, name, description)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// RemoveTier removes a role tier from the provided guild with the
// provided name.
func (db *DB) RemoveTier(
	guildID discord.GuildID, name string) (bool, error) {

	res, err := db.Exec(removeTierQuery, guildID, name)
	if err != nil {
		return false, err
	}

	removed, err := res.RowsAffected()
	return removed > 0, err
}
