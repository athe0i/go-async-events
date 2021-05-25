package map_dispatcher

import (
	"go-async-events/async-events"
	"go-async-events/test"
	"reflect"
	"testing"
)

func TestEventMapDispatcher_RegisterEventHandler(t *testing.T) {
	handler := func(te async_events.Event) (err error) {
		return err
	}

	ed := NewMapEventDispatcher()
	ed.RegisterEventHandler(test.NewTestEvent(test.TestPayload{}), handler)

	if ed.EventMap[reflect.TypeOf(test.NewTestEvent(test.TestPayload{}))] == nil {
		t.Errorf("Handler was not registered")
	}
}

func TestEventMapDispatcher_HandleEvent(t *testing.T) {
	var field string

	handler := func(te async_events.Event) (err error) {
		field = te.GetPayload().(test.TestPayload).Text
		return err
	}

	ed := NewMapEventDispatcher()
	ed.RegisterEventHandler(test.NewTestEvent(test.TestPayload{}), handler)

	te := test.NewTestEvent(test.TestPayload{Text: "event payload"})

	ed.HandleEvent(te)

	if field != "event payload" {
		t.Errorf("handler wasn't fired!")
	}
}
