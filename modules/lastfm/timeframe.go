package lastfm

type lastFmPeriod int

const (
	weekPeriod lastFmPeriod = iota
	monthPeriod
	threeMonthPeriod
	sixMonthPeriod
	yearPeriod
	allTimePeriod
)

type timeframe struct {
	apiPeriod     string
	datePreset    string
	displayPeriod string
}

func (p lastFmPeriod) Timeframe() *timeframe {
	switch p {
	case weekPeriod:
		return &timeframe{
			apiPeriod:     "7day",
			datePreset:    "LAST_7_DAYS",
			displayPeriod: "Last Week",
		}
	case monthPeriod:
		return &timeframe{
			apiPeriod:     "1month",
			datePreset:    "LAST_30_DAYS",
			displayPeriod: "Last Month",
		}
	case threeMonthPeriod:
		return &timeframe{
			apiPeriod:     "3month",
			datePreset:    "LAST_90_DAYS",
			displayPeriod: "Last 3 Months",
		}
	case sixMonthPeriod:
		return &timeframe{
			apiPeriod:     "6month",
			datePreset:    "LAST_180_DAYS",
			displayPeriod: "Last 6 Months",
		}
	case yearPeriod:
		return &timeframe{
			apiPeriod:     "12month",
			datePreset:    "LAST_365_DAYS",
			displayPeriod: "Last Year",
		}
	case allTimePeriod:
		return &timeframe{
			apiPeriod:     "overall",
			datePreset:    "ALL",
			displayPeriod: "All Time",
		}
	default:
		return &timeframe{
			apiPeriod:     "7day",
			datePreset:    "LAST_7_DAYS",
			displayPeriod: "Last Week",
		}
	}
}
