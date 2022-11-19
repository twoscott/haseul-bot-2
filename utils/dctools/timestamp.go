package dctools

import (
	"fmt"
	"time"
)

type TimestampStyle string

const (
	ShortTime     TimestampStyle = "t"
	LongTime      TimestampStyle = "T"
	ShortDate     TimestampStyle = "d"
	LongDate      TimestampStyle = "D"
	ShortDateTime TimestampStyle = "f"
	LongDateTime  TimestampStyle = "F"
	RelativeTime  TimestampStyle = "R"
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
