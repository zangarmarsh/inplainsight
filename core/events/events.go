package events

import (
	"time"
)

type EventType uint8

const (
	AppInit EventType = iota
	AppLogout
	UserLoggedIn

	SecretDiscovered
	SecretAdded
	SecretUpdated

	UserPreferenceInit
	UserPreferenceChanged
)

type Event struct {
	CreatedAt time.Time
	EventType EventType
	Data      map[string]interface{}
}

type EventeableStruct struct {
	listenersQueue  map[EventType][]func(event Event)
	triggeredEvents []Event
}

func (e *EventeableStruct) AddEventsListener(eventsType []EventType, callback func(event Event)) {
	if e.listenersQueue == nil {
		e.listenersQueue = make(map[EventType][]func(event Event))
	}

	for _, eventType := range eventsType {
		e.listenersQueue[eventType] = append(e.listenersQueue[eventType], callback)
	}
}

func (e *EventeableStruct) Trigger(event Event) {
	for _, callback := range e.listenersQueue[event.EventType] {
		callback(event)
	}
}
