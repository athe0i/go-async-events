package map_router

import (
	"go-async-events/async-events"
	"reflect"
)

type MapRouter struct {
	EventMap map[reflect.Type][]string
}

func NewRouter() *MapRouter {
	return &MapRouter{EventMap: make(map[reflect.Type][]string)}
}

func (mr *MapRouter) RegisterEventRoute(ev async_events.Event, channelName string) (err error) {
	typ := reflect.TypeOf(ev)

	if _, ok := mr.EventMap[typ]; ok == false {
		mr.EventMap[typ] = []string{channelName}
	} else {
		_, err := searchStringInArray(channelName, mr.EventMap[typ])
		if err != nil {
			mr.EventMap[typ] = append(mr.EventMap[typ], channelName)
		} else {
			return async_events.ChannelAlreadyRegistered{Channel: channelName, Event: ev}
		}
	}

	return
}

func (mr *MapRouter) GetEventRoutes(ev async_events.Event) (channels []string, err error) {
	typ := reflect.TypeOf(ev)

	if channels, ok := mr.EventMap[typ]; ok {
		return channels, err
	} else {
		return []string{}, async_events.NoRouteRegistered{Event: ev}
	}
}

type NeedleNotFound struct{}

func (nnf NeedleNotFound) Error() string { return "String was not found in array" }

func searchStringInArray(needle string, arr []string) (idx int, err error) {
	for idx, el := range arr {
		if el == needle {
			return idx, err
		}
	}

	return 0, NeedleNotFound{}
}
