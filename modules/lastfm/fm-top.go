package lastfm

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var timePeriodChoices = []discord.IntegerChoice{
	{Name: "Last Week", Value: int(weekPeriod)},
	{Name: "Last Month", Value: int(monthPeriod)},
	{Name: "Last 3 Months", Value: int(threeMonthPeriod)},
	{Name: "Last 6 Months", Value: int(sixMonthPeriod)},
	{Name: "Last Year", Value: int(yearPeriod)},
	{Name: "All Time", Value: int(allTimePeriod)},
}

var fmTopCommandGroup = &router.SubCommandGroup{
	Name:        "top",
	Description: "Commands for displaying top stats from your Last.fm library",
}
