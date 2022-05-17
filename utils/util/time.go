package util

import (
	"fmt"
	"strings"
	"time"
)

const (
	// Day is a time constant for 24 hours.
	Day = 24 * time.Hour
	// Week is a time constant for 7 days.
	Week = 7 * Day
)

type TimePeriod struct {
	Years   int
	Months  int
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

func (t TimePeriod) String() string {
	var ts strings.Builder

	if t.Years > 0 {
		y := fmt.Sprintf("%dY", t.Years)
		ts.WriteString(y)
		ts.WriteString(" ")
	}
	if t.Months > 0 {
		m := fmt.Sprintf("%dM", t.Months)
		ts.WriteString(m)
		ts.WriteString(" ")
	}
	if t.Days > 0 {
		d := fmt.Sprintf("%dd", t.Days)
		ts.WriteString(d)
		ts.WriteString(" ")
	}
	if t.Hours > 0 {
		h := fmt.Sprintf("%dh", t.Hours)
		ts.WriteString(h)
		ts.WriteString(" ")
	}
	if t.Minutes > 0 {
		m := fmt.Sprintf("%dm", t.Minutes)
		ts.WriteString(m)
		ts.WriteString(" ")
	}

	s := fmt.Sprintf("%ds", t.Seconds)
	ts.WriteString(s)

	return ts.String()
}

func daysIn(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// MaxTimeAgoString returns a shorthand string representing how long ago
// an epoch was, measured by the highest possible time increment up to
// a week.
func MaxTimeAgoString(epoch time.Time) string {
	timeSince := time.Since(epoch)
	var timeAgoString string

	if timeSince < time.Minute {
		seconds := timeSince / time.Second
		timeAgoString = fmt.Sprintf("%ds", seconds)
	} else if timeSince < time.Hour {
		minutes := timeSince / time.Minute
		timeAgoString = fmt.Sprintf("%dm", minutes)
	} else if timeSince < Day {
		hours := timeSince / time.Hour
		timeAgoString = fmt.Sprintf("%dh", hours)
	} else if timeSince < Week {
		days := timeSince / Day
		timeAgoString = fmt.Sprintf("%dd", days)
	} else {
		weeks := timeSince / Week
		timeAgoString = fmt.Sprintf("%dwk", weeks)
	}

	return timeAgoString + " ago"
}

// TimeDiff returns a time period object representing the time difference
// between the provided times by all time measurements.
func TimeDiff(old, new time.Time) *TimePeriod {
	if old.Location() != new.Location() {
		new = new.In(old.Location())
	}
	if old.After(new) {
		old, new = new, old
	}

	y1, M1, d1 := old.Date()
	y2, M2, d2 := new.Date()
	years := y2 - y1
	months := int(M2 - M1)
	days := d2 - d1

	h1, m1, s1 := old.Clock()
	h2, m2, s2 := new.Clock()
	hours := h2 - h1
	minutes := m2 - m1
	seconds := s2 - s1

	if seconds < 0 {
		seconds += 60
		minutes--
	}
	if minutes < 0 {
		minutes += 60
		hours--
	}
	if hours < 0 {
		hours += 24
		days--
	}
	if days < 0 {
		days += daysIn(y2, M2-1)
		months--
	}
	if months < 0 {
		months += 12
		years--
	}

	return &TimePeriod{
		Years:   years,
		Months:  months,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}
}
