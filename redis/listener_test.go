package redis

import (
	"context"
	"go-async-events/async-events"
	"go-async-events/json-serializer"
	map_dispatcher "go-async-events/map-dispatcher"
	map_router "go-async-events/map-router"
	"go-async-events/test"
	"testing"
	"time"
)

func TestRedisListener_Run(t *testing.T) {
	listener := getListener()

	var value string
	handler := func(ev async_events.Event) error {
		value = ev.GetPayload().(test.TestPayload).Text
		return nil
	}

	listener.(async_events.Listener).GetDispatcher().RegisterEventHandler(test.NewBlankTestEvent(), handler)
	chn := NewRedisChannel("event_listener_test:test_grp", getRepo(), getSerializer())
	_ = listener.AddChannel(chn)

	getRepo().KeepAlive("event_listener_test:test_grp")

	ev := test.NewTestEvent(test.TestPayload{
		Number:   1234,
		Text:     "test text",
		SubField: test.TestSubPayload{Field: "test sub field"},
	})
	ev.SetName("test")
	getEmitter().Emit(ev)
	listener.Run(context.Background())
	time.Sleep(time.Millisecond * time.Duration(200))

	if value != ev.GetPayload().(test.TestPayload).Text {
		t.Errorf("Value text dont match! Dispatcher wasn't triggered?")
	}
}

func getListener() async_events.Listener {
	repo := NewRepository(test.NewTestRdb(), RedisOptions{})
	disp := map_dispatcher.NewMapEventDispatcher()

	return NewRedisListener(repo, disp)
}

func getEmitter() async_events.Emitter {
	rtr := map_router.NewRouter()
	emtr := NewRedisEmitter(getRepo(), rtr, getSerializer())

	return emtr
}

func getRepo() RedisRepositoryInterface {
	return NewRepository(test.NewTestRdb(), RedisOptions{KeepAliveTimeout: 5})
}

func getSerializer() async_events.Serializer {
	srlr := json_serializer.NewJsonSerializer()
	srlr.RegisterEvent("test", test.NewBlankTestEvent)

	return srlr
}
