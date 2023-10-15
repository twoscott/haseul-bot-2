package user

import "time"

func getNextRepReset(lastRepTime time.Time) time.Time {
	tomorrow := time.Now().Add(24 * time.Hour)
	return time.Date(
		tomorrow.Year(),
		tomorrow.Month(),
		tomorrow.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

func getNextRepResetFromNow() time.Time {
	return getNextRepReset(time.Now())
}

func getRepCutoff() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}
