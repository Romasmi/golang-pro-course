package usecases

import (
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/services/eventservice"
)

func NewUsecases(service *eventservice.CalendarService) map[Type]Usecase {
	return map[Type]Usecase{
		CreateEvent:          &CreateEventUsecase{Service: service},
		UpdateEvent:          &UpdateEventUsecase{Service: service},
		GetEvent:             &GetEventUsecase{Service: service},
		ListEvents:           &ListEventsUsecase{Service: service},
		ListEventsByInterval: &ListEventsByIntervalUsecase{Service: service},
		Hello:                &HelloUsecase{},
	}
}
