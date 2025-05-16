package lastfm

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/gobble-fm/lastfm"
	"github.com/twoscott/haseul-bot-2/router"
)

var timePeriodChoices = []discord.StringChoice{
	{Name: "Past Week", Value: string(lastfm.PeriodWeek)},
	{Name: "Past Month", Value: string(lastfm.PeriodMonth)},
	{Name: "Past 3 Months", Value: string(lastfm.Period3Months)},
	{Name: "Past 6 Months", Value: string(lastfm.Period6Months)},
	{Name: "Past Year", Value: string(lastfm.PeriodYear)},
	{Name: "All Time", Value: string(lastfm.PeriodOverall)},
}

var lastFMTopCommandGroup = &router.SubCommandGroup{
	Name:        "top",
	Description: "Commands for displaying top stats from your Last.fm library",
}
