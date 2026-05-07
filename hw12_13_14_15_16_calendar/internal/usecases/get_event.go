package usecases

import (
	"context"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/services/eventservice"
)

type GetEventUsecase struct {
	Service *eventservice.CalendarService
}

func (u *GetEventUsecase) Do(ctx context.Context, req any) (any, error) {
	return u.Service.GetEvent(ctx, req.(string))
}
