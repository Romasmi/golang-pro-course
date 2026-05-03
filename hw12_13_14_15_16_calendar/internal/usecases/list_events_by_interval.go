package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/services/eventservice"
)

type Interval string

const (
	InvalidInterval Interval = ""
	Day             Interval = "day"
	Week            Interval = "week"
	Month           Interval = "month"
	Year            Interval = "year"
)

type ListEventsByIntervalRequest struct {
	Date     time.Time
	Interval Interval
}

type ListEventsByIntervalUsecase struct {
	Service *eventservice.CalendarService
}

func (i Interval) IsValid() bool {
	return i != InvalidInterval
}

func (u *ListEventsByIntervalUsecase) Do(ctx context.Context, req any) (any, error) {
	r := req.(ListEventsByIntervalRequest)
	var from, to time.Time
	switch r.Interval {
	case Day:
		from = time.Date(r.Date.Year(), r.Date.Month(), r.Date.Day(), 0, 0, 0, 0, r.Date.Location())
		to = from.AddDate(0, 0, 1)
	case Week:
		// Assuming r.Date is the start of the week
		from = r.Date
		to = from.AddDate(0, 0, 7)
	case Month:
		from = time.Date(r.Date.Year(), r.Date.Month(), 1, 0, 0, 0, 0, r.Date.Location())
		to = from.AddDate(0, 1, 0)
	case Year:
		from = time.Date(r.Date.Year(), 1, 1, 0, 0, 0, 0, r.Date.Location())
		to = from.AddDate(1, 0, 0)
	case InvalidInterval:
		return nil, fmt.Errorf("invalid interval: %s", r.Interval)
	default:
		return nil, fmt.Errorf("unknown interval: %s", r.Interval)
	}
	return u.Service.ListEventsByFilter(ctx, from, to)
}
