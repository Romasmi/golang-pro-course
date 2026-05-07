package usecases

type Type int

const (
	CreateEvent Type = iota
	UpdateEvent
	GetEvent
	ListEvents
	ListEventsByInterval
	Hello
)
