package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

var (
	yearsRegex   = regexp.MustCompile(`(?i)(\d+)\s*(?:yr?|year)s?`)
	monthsRegex  = regexp.MustCompile(`(?i)(\d+)\s*(?:(?-i:M)|month)s?`)
	weeksRegex   = regexp.MustCompile(`(?i)(\d+)\s*(?:wk?|week)s?`)
	daysRegex    = regexp.MustCompile(`(?i)(\d+)\s*(?:d|day)s?`)
	hoursRegex   = regexp.MustCompile(`(?i)(\d+)\s*(?:hr?|hour)s?`)
	minutesRegex = regexp.MustCompile(`(?i)(\d+)\s*(?:(?-i:m)|min(?:ute)?)s?`)
	secondsRegex = regexp.MustCompile(`(?i)(\d+)\s*(?:s(?:ec(?:ond)?)?)s?`)
)

// TimePeriod represents a period of time.
type TimePeriod struct {
	Years   int
	Months  int
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

// IsValid returns whether the time period is valid.
func (p TimePeriod) IsValid() bool {
	return !p.IsNull()
}

// IsNull returns whether the time period filled with zero values.
func (p TimePeriod) IsNull() bool {
	return p.Years == 0 &&
		p.Months == 0 &&
		p.Days == 0 &&
		p.Hours == 0 &&
		p.Minutes == 0 &&
		p.Seconds == 0
}

// String formats the time period in a human readable format.
func (p TimePeriod) String() string {
	var ts strings.Builder

	if p.Years > 0 {
		y := fmt.Sprintf("%dY", p.Years)
		ts.WriteString(y)
		ts.WriteString(" ")
	}
	if p.Months > 0 {
		m := fmt.Sprintf("%dM", p.Months)
		ts.WriteString(m)
		ts.WriteString(" ")
	}
	if p.Days > 0 {
		d := fmt.Sprintf("%dd", p.Days)
		ts.WriteString(d)
		ts.WriteString(" ")
	}
	if p.Hours > 0 {
		h := fmt.Sprintf("%dh", p.Hours)
		ts.WriteString(h)
		ts.WriteString(" ")
	}
	if p.Minutes > 0 {
		m := fmt.Sprintf("%dm", p.Minutes)
		ts.WriteString(m)
		ts.WriteString(" ")
	}

	s := fmt.Sprintf("%ds", p.Seconds)
	ts.WriteString(s)

	return ts.String()
}

// Duration converts the time period to a Go duration.
func (p TimePeriod) Duration() time.Duration {
	start := time.Now()
	newTime := p.AfterTime(start)

	return newTime.Sub(start)
}

// AfterTime returns the new time after adding this period to the provided
// start time.
func (p TimePeriod) AfterTime(start time.Time) time.Time {
	hours := time.Hour * time.Duration(p.Hours)
	mins := time.Minute * time.Duration(p.Minutes)
	secs := time.Second * time.Duration(p.Seconds)
	duration := hours + mins + secs

	start = start.AddDate(p.Years, p.Months, p.Days)
	start = start.Add(duration)

	return start
}

// ParseTimePeriod parses a human readable string into a time period.
func ParseTimePeriod(periodString string) TimePeriod {
	var (
		period  TimePeriod
		seconds int
		minutes int
		hours   int
		days    int
		months  int
		years   int
	)

	secondsMatch := secondsRegex.FindStringSubmatch(periodString)
	minutesMatch := minutesRegex.FindStringSubmatch(periodString)
	hoursMatch := hoursRegex.FindStringSubmatch(periodString)
	daysMatch := daysRegex.FindStringSubmatch(periodString)
	weeksMatch := weeksRegex.FindStringSubmatch(periodString)
	monthsMatch := monthsRegex.FindStringSubmatch(periodString)
	yearsMatch := yearsRegex.FindStringSubmatch(periodString)

	if len(secondsMatch) >= 2 {
		seconds, _ = strconv.Atoi(secondsMatch[1])
	}
	if len(minutesMatch) >= 2 {
		minutes, _ = strconv.Atoi(minutesMatch[1])
	}
	if len(hoursMatch) >= 2 {
		hours, _ = strconv.Atoi(hoursMatch[1])
	}
	if len(daysMatch) >= 2 {
		days, _ = strconv.Atoi(daysMatch[1])
	}
	if len(weeksMatch) >= 2 {
		w, _ := strconv.Atoi(weeksMatch[1])
		days += w * 7
	}
	if len(monthsMatch) >= 2 {
		months, _ = strconv.Atoi(monthsMatch[1])
	}
	if len(yearsMatch) >= 2 {
		years, _ = strconv.Atoi(yearsMatch[1])
	}

	period.Seconds += seconds % 60
	period.Minutes += seconds / 60

	period.Minutes += minutes % 60
	period.Hours += minutes / 60

	period.Hours += hours % 24
	period.Days += hours / 24

	period.Days += days % 365
	period.Years += days / 365

	period.Months += months % 12
	period.Years += months / 12

	period.Years += years

	return period
}

// TimeDiff returns a time period object representing the time difference
// between the provided times by all time measurements.
func TimeDiff(old, new time.Time) TimePeriod {
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

	return TimePeriod{
		Years:   years,
		Months:  months,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}
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
	} else if timeSince < humanize.Day {
		hours := timeSince / time.Hour
		timeAgoString = fmt.Sprintf("%dh", hours)
	} else if timeSince < humanize.Week {
		days := timeSince / humanize.Day
		timeAgoString = fmt.Sprintf("%dd", days)
	} else {
		weeks := timeSince / humanize.Week
		timeAgoString = fmt.Sprintf("%dwk", weeks)
	}

	return timeAgoString + " ago"
}

func daysIn(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
