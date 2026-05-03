package usecases

import (
	"context"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/services/eventservice"
)

type UpdateEventUsecase struct {
	Service *eventservice.CalendarService
}

func (u *UpdateEventUsecase) Do(ctx context.Context, req any) (any, error) {
	return u.Service.UpdateEvent(ctx, req.(domain.Event))
}
