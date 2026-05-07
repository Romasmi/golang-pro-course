package usecases

import (
	"context"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/services/eventservice"
)

type CreateEventUsecase struct {
	Service *eventservice.CalendarService
}

func (u *CreateEventUsecase) Do(ctx context.Context, req any) (any, error) {
	return u.Service.CreateEvent(ctx, req.(domain.Event))
}
