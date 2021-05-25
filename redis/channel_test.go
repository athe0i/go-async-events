package redis

import (
	"context"
	"go-async-events/async-events"
	"go-async-events/json-serializer"
	"go-async-events/test"
	"reflect"
	"testing"
)

func TestRedisChannel_Push(t *testing.T) {
	srlzr := json_serializer.NewJsonSerializer()
	chn := GetTestChannel(srlzr)
	srlzr.RegisterEvent("test", test.NewBlankTestEvent)
	err := chn.Push(GetTestEvent())

	if err != nil {
		t.Errorf("Failed to push event to channel: " + err.Error())
		return
	}

	rdb := test.NewTestRdb()

	lastEvent, err := rdb.LRange(context.Background(), "test", -1, -1).Result()

	if err != nil {
		t.Errorf("Failed to fetch events: " + err.Error())
	}

	testPayload := "{\"name\":\"test\",\"payload\":{\"Number\":124,\"Text\":\"Test\",\"SubField\":{\"Field\":\"SubField\"}}}"
	if lastEvent[0] != testPayload {
		t.Errorf("Latest event has incorrect payload!\nActual:" + lastEvent[0] + "\nExpected" + testPayload)
		return
	}
}

func TestRedisChannel_Pop(t *testing.T) {
	slr := json_serializer.NewJsonSerializer()
	chn := GetTestChannel(slr)

	ev := GetTestEvent()

	slr.(*json_serializer.JsonSerializer).EventRegistry.RegisterEvent("test", test.NewBlankTestEvent)
	err := chn.Push(ev)

	if err != nil {
		t.Errorf("Failed to push event to channel: " + err.Error())
		return
	}

	popEv, err := chn.Pop()

	if err != nil {
		t.Errorf("Failed to pop event from channel: " + err.Error())
		return
	}

	if reflect.DeepEqual(ev, popEv) == false {
		t.Errorf("Events are not equal!")
	}
}

func GetTestChannel(slr async_events.Serializer) *RedisChannel {
	return NewRedisChannel("test", NewRepository(test.NewTestRdb(), RedisOptions{}), slr)
}

func GetTestEvent() async_events.Event {
	e := test.TestEvent{
		BaseAsyncEvent: &async_events.BaseAsyncEvent{},
		Payload: test.TestPayload{
			Number: 124,
			Text:   "Test",
			SubField: test.TestSubPayload{
				Field: "SubField",
			},
		},
	}
	e.SetName("test")

	return &e
}
