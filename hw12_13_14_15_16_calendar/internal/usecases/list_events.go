package usecases

import (
	"context"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/services/event_service"
)

type ListEventsUsecase struct {
	Service *event_service.CalendarService
}

func (u *ListEventsUsecase) Do(ctx context.Context, _ any) (any, error) {
	return u.Service.ListEvents(ctx)
}
