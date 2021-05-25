package async_events

import (
	"fmt"
	"reflect"
)

type ChannelAlreadyRegistered struct {
	Channel string
	Event   Event
}

func (car ChannelAlreadyRegistered) Error() string {
	return fmt.Sprintf("Channel \"%s\" already registered for event \"%s\"", car.Channel, reflect.TypeOf(car.Event).Elem().Name())
}


type ListenerChannelAlreadyRegistered struct {
	Chn string
}

func (lcar ListenerChannelAlreadyRegistered) Error() string {
	return fmt.Sprintf("Channel %s already registere in listener!", lcar.Chn)
}

type NoEventHandlerRegistered struct {
	EventName string
}

func (ner NoEventHandlerRegistered) Error() string {
	return fmt.Sprintf("No event registered for name: %s", ner.EventName)
}


type NoEventRegistered struct {
	EventName string
}

func (ner NoEventRegistered) Error() string {
	return fmt.Sprintf("No event registered for name: %s", ner.EventName)
}


type NoRouteRegistered struct {
	Event Event
}

func (nrr NoRouteRegistered) Error() string {
	return fmt.Sprintf("No route registered for event: %s", reflect.TypeOf(nrr.Event).Elem().Name())
}

type ErrorHandler func(err error)

type HandlesErrors interface {
	HandleError(err error)
}

func LogError(err error) {
	fmt.Println(err)
}
