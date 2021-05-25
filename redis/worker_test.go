package redis

import (
	"context"
	"fmt"
	"go-async-events/async-events"
	//json_serializer "go-async-events/json-serializer"
	//map_dispatcher "go-async-events/map-dispatcher"
	//map_router "go-async-events/map-router"
	"go-async-events/test"
	"testing"
	"time"
)

func TestEventWorker_Run(t *testing.T) {
	lstnr := getListener()
	chn := NewRedisChannel("event_listener_test:test_grp", getRepo(), getSerializer())
	_ = lstnr.AddChannel(chn)

	wrkr := async_events.NewEventWorker(time.Duration(2)*time.Second, time.Duration(200)*time.Millisecond, lstnr)
	emtr := getEmitter()

	var value string
	handler := func(ev async_events.Event) error {
		value = ev.GetPayload().(test.TestPayload).Text
		fmt.Println("in handler")
		return nil
	}
	lstnr.GetDispatcher().RegisterEventHandler(test.NewBlankTestEvent(), handler)

	go wrkr.Run(context.Background())
	ev := test.NewTestEvent(test.TestPayload{
		Number:   1234,
		Text:     "test text",
		SubField: test.TestSubPayload{Field: "test sub field"},
	})
	time.Sleep(time.Duration(100) * time.Millisecond)
	emtr.Emit(ev)

	time.Sleep(time.Second)

	if value != ev.GetPayload().(test.TestPayload).Text {
		t.Errorf("Value text dont match! Dispatcher wasn't triggered?")
	}
}

//func getListener() async_events.Listener {
//	repo := NewRepository(test.NewTestRdb(), RedisOptions{})
//	disp := map_dispatcher.NewMapEventDispatcher()
//
//	return NewRedisListener(repo, disp)
//}
//
//func getEmitter() async_events.Emitter {
//	rtr := map_router.NewRouter()
//	emtr := NewRedisEmitter(getRepo(), rtr, getSerializer())
//
//	return emtr
//}
//
//func getRepo() *Repository {
//	return NewRepository(test.NewTestRdb(), RedisOptions{KeepAliveTimeout: 5})
//}
//
//func getSerializer() async_events.Serializer {
//	srlr := json_serializer.NewJsonSerializer()
//	srlr.RegisterEvent("test", test.NewBlankTestEvent)
//
//	return srlr
//}
