package convert

import (
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/usecases"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/pkg/api"
)

func IntervalFromProto(p api.Interval) usecases.Interval {
	switch p {
	case api.Interval_DAY:
		return usecases.Day
	case api.Interval_WEEK:
		return usecases.Week
	case api.Interval_MONTH:
		return usecases.Month
	case api.Interval_YEAR:
		return usecases.Year
	default:
		return usecases.InvalidInterval
	}
}
