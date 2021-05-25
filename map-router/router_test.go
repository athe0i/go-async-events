package map_router

import (
	"errors"
	async_events "go-async-events/async-events"
	"go-async-events/test"
	"reflect"
	"testing"
)

func TestMapRouter_RegisterEventRoute(t *testing.T) {
	router := NewRouter()
	testEvent := test.NewTestEvent(test.TestPayload{})
	err := router.RegisterEventRoute(testEvent, "test-channel")
	if err != nil {
		t.Errorf("Failed to create route: %s", err.Error())
		return
	}

	if router.EventMap[reflect.TypeOf(testEvent)][0] != "test-channel" {
		t.Errorf("Channel was not registered!")
		return
	}

	err = router.RegisterEventRoute(testEvent, "test-channel-two")
	if err != nil {
		t.Errorf("Failed to create route: %s", err.Error())
		return
	}

	if router.EventMap[reflect.TypeOf(testEvent)][1] != "test-channel-two" {
		t.Errorf("It should've register new route!")
		return
	}

	err = router.RegisterEventRoute(testEvent, "test-channel-two")
	if errors.As(err, &async_events.ChannelAlreadyRegistered{}) == false {
		t.Errorf("It should throw an error if channel is already registered")
		return
	}
}

func TestMapRouter_GetEventRoutes(t *testing.T) {
	router := NewRouter()
	testEvent := test.NewTestEvent(test.TestPayload{})

	_, err := router.GetEventRoutes(testEvent)
	if errors.As(err, &async_events.NoRouteRegistered{}) == false {
		t.Errorf("It should throw an error if no routes available for the event.")
		return
	}

	err = router.RegisterEventRoute(testEvent, "test-channel")
	if err != nil {
		t.Errorf("Failed to create route: %s", err.Error())
		return
	}

	routes, err := router.GetEventRoutes(testEvent)
	if err != nil {
		t.Errorf("Failed to retrive routes!")
		return
	}

	if routes[0] != "test-channel" || len(routes) != 1 {
		t.Errorf("Invalid routes list returned")
		return
	}
}
