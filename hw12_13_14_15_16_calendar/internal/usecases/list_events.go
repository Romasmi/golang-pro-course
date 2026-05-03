package usecases

import (
	"context"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/services/eventservice"
)

type ListEventsUsecase struct {
	Service *eventservice.CalendarService
}

func (u *ListEventsUsecase) Do(ctx context.Context, _ any) (any, error) {
	return u.Service.ListEvents(ctx)
}
