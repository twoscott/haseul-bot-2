package guilddb

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

const (
	defaultWelcomeTitle   string      = "New Member!"
	defaultWelcomeMessage welcomeText = "Welcome {mention} to {server}!"
)

type welcomeText string

// String returns the welcome text as a string
func (t welcomeText) String() string {
	return string(t)
}

// Format returns the welcome text formatted according to the following rules:
// - {mention} is replaced with the user mention of the new member.
// - {username} is replaced with the username of the new member.
// - {tag} is replaced with the user tag of the new member.
// - {server} is relpaced with the name of the server the new member joined.
// - {member-number} is relpaced new member count of the server.
func (t welcomeText) Format(member discord.Member, guild discord.Guild) string {
	text := t.String()
	if t == "" {
		return text
	}

	memberNumber := strconv.FormatInt(int64(guild.ApproximateMembers), 10)

	text = strings.ReplaceAll(text, "{mention}", member.Mention())
	text = strings.ReplaceAll(text, "{username}", member.User.Username)
	text = strings.ReplaceAll(text, "{tag}", member.User.Tag())
	text = strings.ReplaceAll(text, "{server}", guild.Name)
	text = strings.ReplaceAll(text, "{member-number}", memberNumber)

	return text
}

type welcomeConfig struct {
	RawTitle   string        `db:"welcometitle"`
	RawMessage welcomeText   `db:"welcomemessage"`
	RawColour  sql.NullInt32 `db:"welcomecolour"`
}

// Welcome represents the welcome configuration fields of the guild config.
type Welcome struct {
	ChannelID discord.ChannelID `db:"welcomechannelid"`
	welcomeConfig
}

// Title returns a welcome config's message title.
func (w Welcome) Title() string {
	if w.RawTitle == "" {
		return defaultWelcomeTitle
	}

	return w.RawTitle
}

// Message returns a welcome config's message.
func (w Welcome) Message() string {
	if w.RawMessage == "" {
		return defaultWelcomeMessage.String()
	}

	return w.RawMessage.String()
}

// FormattedMessage returns a welcome config's message, formatted with the details of the
// provided new member and the guild they joined.
func (w Welcome) FormattedMessage(member discord.Member, guild discord.Guild) string {
	if w.RawMessage == "" {
		return defaultWelcomeMessage.Format(member, guild)
	}

	return w.RawMessage.Format(member, guild)
}

// Colour returns the Discord colour of the welcome config. If the colour is
// not set, a default embed background colour is returned.
func (w Welcome) Colour() discord.Color {
	if !w.RawColour.Valid {
		return dctools.BlurpleColour
	}

	return discord.Color(w.RawColour.Int32)
}

const (
	setWelcomeChannelQuery = `
		UPDATE GuildConfigs SET welcomeChannelID = $1
		WHERE guildID = $2`
	setWelcomeChannelNullQuery = `
		UPDATE GuildConfigs SET welcomeChannelID = 0 # 0
		WHERE guildID = $1`
	setWelcomeMessageQuery = `
		UPDATE GuildConfigs SET	welcomeMessage = $1
		WHERE guildID = $2`
	setWelcomeTitleQuery = `
		UPDATE GuildConfigs SET	welcomeTitle = $1
		WHERE guildID = $2`
	setWelcomeColourQuery = `
		UPDATE GuildConfigs SET	welcomeColour = $1
		WHERE guildID = $2`
	getWelcomeConfigQuery = `
		SELECT welcomeChannelID, welcomeTitle, welcomeMessage, welcomeColour 
		FROM GuildConfigs WHERE guildID = $1`
)

// SetWelcomeChannel updates the guild config of the given guild ID and sets
// the welcome logs channel ID to the provided channel ID.
func (db *DB) SetWelcomeChannel(
	guildID discord.GuildID, channelID discord.ChannelID) (bool, error) {

	res, err := db.Exec(setWelcomeChannelQuery, channelID, guildID)
	if err != nil {
		return false, err
	}

	updated, err := res.RowsAffected()
	return updated > 0, err
}

// DisableWelcomeLogs updates the guild config of the given guild ID and sets
// the welcome logs channel ID to a value that can be interpreted as null.
func (db *DB) DisableWelcomeLogs(guildID discord.GuildID) (bool, error) {
	res, err := db.Exec(setWelcomeChannelNullQuery, guildID)
	if err != nil {
		return false, err
	}

	updated, err := res.RowsAffected()
	return updated > 0, err
}

// SetWelcomeMessage updates the guild config of the given guild ID and sets
// the welcome message.
func (db *DB) SetWelcomeMessage(
	guildID discord.GuildID, message string) (bool, error) {

	res, err := db.Exec(setWelcomeMessageQuery, message, guildID)
	if err != nil {
		return false, err
	}

	updated, err := res.RowsAffected()
	return updated > 0, err
}

// SetWelcomeTitle updates the guild config of the given guild ID and sets
// the welcome title.
func (db *DB) SetWelcomeTitle(
	guildID discord.GuildID, title string) (bool, error) {

	res, err := db.Exec(setWelcomeTitleQuery, title, guildID)
	if err != nil {
		return false, err
	}

	updated, err := res.RowsAffected()
	return updated > 0, err
}

// SetWelcomeColour updates the guild config of the given guild ID and sets
// the welcome colour.
func (db *DB) SetWelcomeColour(
	guildID discord.GuildID, colour discord.Color) (bool, error) {

	res, err := db.Exec(setWelcomeColourQuery, colour, guildID)
	if err != nil {
		return false, err
	}

	updated, err := res.RowsAffected()
	return updated > 0, err
}

// WelcomeConfig returns the welcome logs config from the guild config
// corresponding to the provided guild ID.
func (db *DB) WelcomeConfig(guildID discord.GuildID) (*Welcome, error) {

	var config Welcome
	err := db.Get(&config, getWelcomeConfigQuery, guildID)

	return &config, err
}
