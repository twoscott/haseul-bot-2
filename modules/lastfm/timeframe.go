package lastfm

type timeframe struct {
	apiPeriod     string
	datePreset    string
	displayPeriod string
}

func getTimeframe(period string) *timeframe {
	switch period {
	case "7", "7day", "week", "wk":
		return &timeframe{
			apiPeriod:     "7day",
			datePreset:    "LAST_7_DAYS",
			displayPeriod: "Last Week",
		}
	case "30", "30day", "month", "m":
		return &timeframe{
			apiPeriod:     "1month",
			datePreset:    "LAST_30_DAYS",
			displayPeriod: "Last Month",
		}
	case "90", "90day", "90days", "3month", "3months", "3m":
		return &timeframe{
			apiPeriod:     "3month",
			datePreset:    "LAST_90_DAYS",
			displayPeriod: "Last 3 Months",
		}
	case "190", "180day", "180days", "6month", "6months", "6m":
		return &timeframe{
			apiPeriod:     "6month",
			datePreset:    "LAST_180_DAYS",
			displayPeriod: "Last 6 Months",
		}
	case "365", "365day", "year", "yr", "y":
		return &timeframe{
			apiPeriod:     "12month",
			datePreset:    "LAST_365_DAYS",
			displayPeriod: "Last Year",
		}
	case "overall", "alltime", "at", "all":
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
