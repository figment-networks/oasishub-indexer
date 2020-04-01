package time

import timePckg "time"

func BeginningOfHour(t timePckg.Time) timePckg.Time {
	y, m, d := t.Date()
	return timePckg.Date(y, m, d, t.Hour(), 0, 0, 0, t.Location())
}

func BeginningOfDay(t timePckg.Time) timePckg.Time {
	y, m, d := t.Date()
	return timePckg.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func EndOfHour(t timePckg.Time) timePckg.Time {
	return BeginningOfHour(t).Add(timePckg.Hour - timePckg.Nanosecond)
}

func EndOfDay(t timePckg.Time) timePckg.Time {
	y, m, d := t.Date()
	return timePckg.Date(y, m, d, 23, 59, 59, int(timePckg.Second - timePckg.Nanosecond), t.Location())
}