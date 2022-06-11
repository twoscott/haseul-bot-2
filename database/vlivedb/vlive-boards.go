package vlivedb

import "github.com/diamondburned/arikawa/v3/discord"

// Board represents a VLIVE board database entry.
type Board struct {
	ID                int64  `db:"id"`
	LastPostTimestamp int64  `db:"lastposttimestamp"`
	lastPostID        string `db:"lastpostid"`
}

const (
	createVLIVEBoardsTableQuery = `
		CREATE TABLE IF NOT EXISTS VLIVEBoards(
			id                INT8        NOT NULL,
			lastPostTimestamp INT8        DEFAULT  0,
			lastPostID		  VARCHAR(64) DEFAULT  0,
			PRIMARY KEY(id)
		)
	`

	addBoardQuery = `
		INSERT INTO VLIVEBoards VALUES($1, $2, $3) ON CONFLICT DO NOTHING`
	setLastDataQuery = `
		UPDATE VLIVEBoards 
		SET lastPostTimestamp = $2, lastPostID = $3 
		WHERE id = $1`

	getBoardQuery     = `SELECT * FROM VLIVEBoards WHERE id = $1`
	getAllBoardsQuery = `SELECT * FROM VLIVEBoards`
	removeBoardQuery  = `DELETE FROM VLIVEBoards WHERE id = $1`

	getGuildBoardCountQuery = `
		SELECT COUNT(*) FROM VLIVEBoards WHERE id IN (
			SELECT vliveBoardID FROM VLIVEFeeds WHERE guildID = $1
		)
	`
	getBoardByGuildQuery = `
		SELECT * FROM VLIVEBoards WHERE id = $2 AND id IN (
			SELECT vliveBoardID FROM VLIVEFeeds WHERE guildID = $1
		)
	`
	getBoardsByGuildQuery = `
		SELECT * FROM VLIVEBoards WHERE id IN (
			SELECT vliveBoardID FROM VLIVEFeeds WHERE guildID = $1
		)
	`
)

// AddBoard adds a VLIVE board to the database.
func (db *DB) AddBoard(
	vliveBoardID string, lastTweetID int64) (bool, error) {

	res, err := db.Exec(addBoardQuery, vliveBoardID, lastTweetID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// GetBoard returns a VLIVE board from the database.
func (db *DB) GetBoard(vliveBoardID string) (*Board, error) {
	var vliveBoard Board
	err := db.Get(&vliveBoard, getBoardQuery, vliveBoardID)

	return &vliveBoard, err
}

// GetAllBoards returns all VLIVE boards from the database.
func (db *DB) GetAllBoards() ([]Board, error) {
	var VLIVEBoards []Board
	err := db.Select(&VLIVEBoards, getAllBoardsQuery)

	return VLIVEBoards, err
}

// GetBoardByGuild returns a VLIVE board if the guild ID has a VLIVE feed
// set up for the provided VLIVE board ID.
func (db *DB) GetBoardByGuild(
	guildID discord.GuildID, vliveBoardID string) (Board, error) {

	var vliveBoard Board
	err := db.Get(
		&vliveBoard, getBoardByGuildQuery, guildID, vliveBoardID,
	)

	return vliveBoard, err
}

// GetBoardsByGuild returns all VLIVE boards that are included in one or
// more VLIVE feeds in the given guild ID.
func (db *DB) GetBoardsByGuild(guildID discord.GuildID) ([]Board, error) {
	var VLIVEBoards []Board
	err := db.Select(&VLIVEBoards, getBoardsByGuildQuery, guildID)

	return VLIVEBoards, err
}

// GetGuildBoardCount returns how many unique VLIVE boards a guild has
// VLIVE feeds set up for.
func (db *DB) GetGuildBoardCount(guildID discord.GuildID) (int, error) {
	var vliveCount int
	err := db.Get(&vliveCount, getGuildBoardCountQuery, guildID)

	return vliveCount, err
}

// RemoveBoard removes a VLIVE board from the database.
func (db *DB) RemoveBoard(vliveBoardID string) (bool, error) {
	res, err := db.Exec(removeBoardQuery, vliveBoardID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// SetLastTimestamp updates the last timestamp for a given VLIVE board ID.
func (db *DB) SetLastTimestamp(
	vliveBoardID string, timestamp int64, postID string) (int64, error) {

	res, err := db.Exec(setLastDataQuery, vliveBoardID, timestamp, postID)
	if err != nil {
		return 0, err
	}

	cleared, err := res.RowsAffected()
	return cleared, err
}
