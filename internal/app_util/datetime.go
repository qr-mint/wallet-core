package app_util

import "time"

func TimeToStartOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func TimeToEndOfDay(value time.Time) time.Time {
	return value.AddDate(0, 0, 1).Add(-time.Nanosecond)
}
