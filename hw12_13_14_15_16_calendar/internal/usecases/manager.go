package usecases

import (
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/services/event_service"
)

func NewUsecases(service *event_service.CalendarService) map[Type]Usecase {
	return map[Type]Usecase{
		CreateEvent:          &CreateEventUsecase{Service: service},
		UpdateEvent:          &UpdateEventUsecase{Service: service},
		GetEvent:             &GetEventUsecase{Service: service},
		ListEvents:           &ListEventsUsecase{Service: service},
		ListEventsByInterval: &ListEventsByIntervalUsecase{Service: service},
		Hello:                &HelloUsecase{},
	}
}
