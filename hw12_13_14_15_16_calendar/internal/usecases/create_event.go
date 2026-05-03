package usecases

import (
	"context"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/services/event_service"
)

type CreateEventUsecase struct {
	Service *event_service.CalendarService
}

func (u *CreateEventUsecase) Do(ctx context.Context, req any) (any, error) {
	return u.Service.CreateEvent(ctx, req.(domain.Event))
}
