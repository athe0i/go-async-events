package async_events

import (
	"reflect"
)

type EventRegistry interface {
	RegisterEvent(name string, c EventConstructor)
	CreateEvent(name string, payload map[string]interface{}, callbackChannel string) (e Event, err error)
	GetName(ev Event) (name string, err error)
}

type EventConstructor func() Event

type EventConstructorRegistry struct {
	constructorMap map[string]EventConstructor
	typeMap        map[reflect.Type]string
}

func NewEventConstructorRegistry() (er EventRegistry) {
	return &EventConstructorRegistry{
		constructorMap: make(map[string]EventConstructor),
		typeMap:        make(map[reflect.Type]string),
	}
}

func (etr *EventConstructorRegistry) RegisterEvent(name string, constructor EventConstructor) {
	etr.constructorMap[name] = constructor
	etr.typeMap[reflect.TypeOf(constructor())] = name
}

func (etr *EventConstructorRegistry) CreateEvent(name string, payload map[string]interface{}, callbackChannel string) (e Event, err error) {
	var constructor EventConstructor
	constructor, ok := etr.constructorMap[name]
	if ok == false {
		return nil, NoEventRegistered{EventName: name}
	}

	ev := constructor()

	ev.SetName(name)
	ev.SetCallbackChannel(callbackChannel)
	err = ev.ParsePayload(payload)

	return ev, err
}

func (etr *EventConstructorRegistry) GetName(ev Event) (name string, err error) {
	name, ok := etr.typeMap[reflect.TypeOf(ev)]

	if ok == false {
		return name, NoEventRegistered{EventName: reflect.TypeOf(ev).Name()}
	}

	return
}
