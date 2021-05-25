package map_dispatcher

import (
	"go-async-events/async-events"
	"reflect"
)



type MapEventDispatcher struct {
	EventMap map[reflect.Type]async_events.EventHandler
}

func NewMapEventDispatcher() *MapEventDispatcher {
	return &MapEventDispatcher{EventMap: make(map[reflect.Type]async_events.EventHandler)}
}

func (emp *MapEventDispatcher) RegisterEventHandler(event async_events.Event, handler async_events.EventHandler) {
	emp.EventMap[reflect.TypeOf(event)] = handler
}

func (emp *MapEventDispatcher) HandleEvent(ev async_events.Event) (err error) {
	var handler async_events.EventHandler
	handler, ok := emp.EventMap[reflect.TypeOf(ev)]

	if ok == false {
		return async_events.NoEventHandlerRegistered{EventName: reflect.TypeOf(ev).Name()}
	}

	return handler(ev)
}
