package vlivedb

import "github.com/diamondburned/arikawa/v3/discord"

// Board represents a VLIVE board database entry.
type Board struct {
	ID                int64  `db:"id"`
	ChannelCode       string `db:"channelcode"`
	LastPostTimestamp int64  `db:"lastposttimestamp"`
	LastPostID        string `db:"lastpostid"`
}

const (
	createVLIVEBoardsTableQuery = `
		CREATE TABLE IF NOT EXISTS VLIVEBoards(
			id                INT8        NOT NULL,
			channelCode       VARCHAR(64) NOT NULL,
			lastPostTimestamp INT8        DEFAULT  0,
			lastPostID		  VARCHAR(64) DEFAULT  '0-0',
			PRIMARY KEY(id)
		)
	`

	addBoardQuery = `
		INSERT INTO VLIVEBoards VALUES($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	setLastDataQuery = `
		UPDATE VLIVEBoards 
		SET lastPostTimestamp = $2, lastPostID = $3 
		WHERE id = $1`

	getBoardQuery     = `SELECT * FROM VLIVEBoards WHERE id = $1`
	getAllBoardsQuery = `SELECT * FROM VLIVEBoards`
	removeBoardQuery  = `DELETE FROM VLIVEBoards WHERE id = $1`

	getBoardsByVLIVEChannelQuery = `
		SELECT * FROM VLIVEBoards WHERE channelCode = $1
	`
	getGuildBoardCountQuery = `
		SELECT COUNT(*) FROM VLIVEBoards WHERE id IN (
			SELECT boardID FROM VLIVEFeeds WHERE guildID = $1
		)
	`
	getBoardByGuildQuery = `
		SELECT * FROM VLIVEBoards WHERE id = $2 AND id IN (
			SELECT boardID FROM VLIVEFeeds WHERE guildID = $1
		)
	`
	getBoardsByGuildQuery = `
		SELECT * FROM VLIVEBoards WHERE id IN (
			SELECT boardID FROM VLIVEFeeds WHERE guildID = $1
		)
	`
	getChannelCodesByGuildQuery = `
		SELECT DISTINCT channelCode FROM VLIVEBoards WHERE  id IN (
			SELECT boardID FROM VLIVEFeeds WHERE guildID = $1
		)	
	`
)

// AddBoard adds a VLIVE board to the database.
func (db *DB) AddBoard(
	boardID int64,
	channelCode string,
	postTimestamp int64,
	postID string) (bool, error) {

	res, err := db.Exec(
		addBoardQuery, boardID, channelCode, postTimestamp, postID,
	)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// GetBoard returns a VLIVE board from the database.
func (db *DB) GetBoard(boardID int64) (*Board, error) {
	var vliveBoard Board
	err := db.Get(&vliveBoard, getBoardQuery, boardID)

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
	guildID discord.GuildID, boardID int64) (Board, error) {

	var vliveBoard Board
	err := db.Get(
		&vliveBoard, getBoardByGuildQuery, guildID, boardID,
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

// GetBoardsByVLIVEChannel returns all VLIVE boards that belong to a given
// VLIVE channel code.
func (db *DB) GetBoardsByVLIVEChannel(channelCode string) ([]Board, error) {
	var VLIVEBoards []Board
	err := db.Select(&VLIVEBoards, getBoardsByVLIVEChannelQuery, channelCode)

	return VLIVEBoards, err
}

// GetChannelCodesByGuild returns all VLIVE channels' codes that belong to a
// given guild.
func (db *DB) GetChannelCodesByGuild(guildID discord.GuildID) ([]string, error) {
	var channelCodes []string
	err := db.Select(&channelCodes, getChannelCodesByGuildQuery, guildID)

	return channelCodes, err
}

// GetGuildBoardCount returns how many unique VLIVE boards a guild has
// VLIVE feeds set up for.
func (db *DB) GetGuildBoardCount(guildID discord.GuildID) (uint64, error) {
	var vliveCount uint64
	err := db.Get(&vliveCount, getGuildBoardCountQuery, guildID)

	return vliveCount, err
}

// RemoveBoard removes a VLIVE board from the database.
func (db *DB) RemoveBoard(boardID int64) (bool, error) {
	res, err := db.Exec(removeBoardQuery, boardID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// SetLastTimestamp updates the last timestamp for a given VLIVE board ID.
func (db *DB) SetLastTimestamp(
	boardID int64, timestamp int64, postID string) (int64, error) {

	res, err := db.Exec(setLastDataQuery, boardID, timestamp, postID)
	if err != nil {
		return 0, err
	}

	cleared, err := res.RowsAffected()
	return cleared, err
}
