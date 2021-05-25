package redis

import (
	"context"
	"go-async-events/json-serializer"
	map_router "go-async-events/map-router"
	"go-async-events/test"
	"testing"
	"time"
)

func TestRedisEmitter_Emit(t *testing.T) {
	repo := NewRepository(test.NewTestRdb(), RedisOptions{})
	rtr := map_router.NewRouter()
	srlr := json_serializer.NewJsonSerializer()
	emtr := NewRedisEmitter(repo, rtr, srlr)

	ev := test.NewTestEvent(test.TestPayload{
		Number:   1234,
		Text:     "test text",
		SubField: test.TestSubPayload{Field: "test sub field"},
	})
	ev.SetName("test")

	emtr.GetSerializer().RegisterEvent("test", test.NewBlankTestEvent)
	emtr.GetRouter().RegisterEventRoute(ev, "router_channel")

	rdb := test.NewTestRdb()
	ctx := context.Background()
	re := emtr.(*RedisEmitter)
	_, err := rdb.Set(ctx, re.repo.GetKeepAliveKey("router_channel")+":test", true, time.Second*5).Result()
	if err != nil {
		t.Errorf("Cannot set keep alive key: %s", err.Error())
	}

	evName, err := re.GetSerializer().GetName(ev)
	if err != nil {
		t.Errorf("Failed to get name of event")
		return
	}
	listenerChn := re.repo.GetEventListenerChannelName(evName)

	_, err = rdb.Set(ctx, re.repo.GetKeepAliveKey(listenerChn)+":test", true, time.Second*5).Result()
	if err != nil {
		t.Errorf("Cannot set listener key: %s", err.Error())
		return
	}

	err = emtr.Emit(ev)
	if err != nil {
		t.Errorf("Failed to emit event: %s", err.Error())
		return
	}

	serialized, err := emtr.GetSerializer().Serialize(ev)
	if err != nil {
		t.Errorf("Failed to serialize event: %s", err.Error())
	}

	lastChannelEvent, err := rdb.LPop(context.Background(), "router_channel:test").Result()

	if err != nil {
		t.Errorf("Failed to fetch events: " + err.Error())
		return
	}

	if lastChannelEvent != serialized {
		t.Errorf("Event is not the original one!")
		return
	}

	lastListenerEvent, err := rdb.LPop(context.Background(), "event_listener_test:test").Result()

	if err != nil {
		t.Errorf("Failed to fetch events: " + err.Error())
		return
	}

	if lastListenerEvent != serialized {
		t.Errorf("Event is not the original one!")
		return
	}
}
