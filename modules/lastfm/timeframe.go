package lastfm

import "github.com/twoscott/gobble-fm/lastfm"

type timeframe struct {
	apiPeriod     lastfm.Period
	datePreset    string
	displayPeriod string
}

func newTimeframe(period lastfm.Period) *timeframe {
	switch period {
	case lastfm.PeriodWeek:
		return &timeframe{
			apiPeriod:     lastfm.PeriodWeek,
			datePreset:    "LAST_7_DAYS",
			displayPeriod: "Past Week",
		}
	case lastfm.PeriodMonth:
		return &timeframe{
			apiPeriod:     lastfm.PeriodMonth,
			datePreset:    "LAST_30_DAYS",
			displayPeriod: "Past Month",
		}
	case lastfm.Period3Months:
		return &timeframe{
			apiPeriod:     lastfm.Period3Months,
			datePreset:    "LAST_90_DAYS",
			displayPeriod: "Past 3 Months",
		}
	case lastfm.Period6Months:
		return &timeframe{
			apiPeriod:     lastfm.Period6Months,
			datePreset:    "LAST_180_DAYS",
			displayPeriod: "Past 6 Months",
		}
	case lastfm.PeriodYear:
		return &timeframe{
			apiPeriod:     lastfm.PeriodYear,
			datePreset:    "LAST_365_DAYS",
			displayPeriod: "Past Year",
		}
	case lastfm.PeriodOverall:
		return &timeframe{
			apiPeriod:     lastfm.PeriodOverall,
			datePreset:    "ALL",
			displayPeriod: "All Time",
		}
	default:
		return &timeframe{
			apiPeriod:     lastfm.PeriodWeek,
			datePreset:    "LAST_7_DAYS",
			displayPeriod: "Past Week",
		}
	}
}
