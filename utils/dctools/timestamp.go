package dctools

import (
	"fmt"
	"time"
)

type TimestampStyle string

const (
	// 	16:20
	ShortTime TimestampStyle = "t"
	// 16:20:30
	LongTime TimestampStyle = "T"
	// 20/04/2021
	ShortDate TimestampStyle = "d"
	// 20 April 2021
	LongDate TimestampStyle = "D"
	// 20 April 2021 16:20 - DEFAULT
	ShortDateTime TimestampStyle = "f"
	// Tuesday, 20 April 2021 16:20
	LongDateTime TimestampStyle = "F"
	// 	2 months ago
	RelativeTime TimestampStyle = "R"
)

// UnixTimestampStyled returns markdown for displaying a time in Discord,
// with the provided style.
func UnixTimestampStyled(timestamp time.Time, style TimestampStyle) string {
	if style == "" {
		return fmt.Sprintf("<t:%d>", timestamp.Unix())
	}

	return fmt.Sprintf("<t:%d:%s>", timestamp.Unix(), style)
}

// UnixTimestamp returns markdown for displaying a time in Discord.
func UnixTimestamp(timestamp time.Time) string {
	return UnixTimestampStyled(timestamp, "")
}
