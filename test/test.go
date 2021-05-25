package test

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go-async-events/async-events"
)

type TestPayload struct {
	Number   int
	Text     string
	SubField TestSubPayload
}

type TestSubPayload struct {
	Field string
}

type TestEvent struct {
	*async_events.BaseAsyncEvent
	Payload TestPayload
}

func (tEvent *TestEvent) GetPayload() interface{} {
	return tEvent.Payload
}

func (tEvent *TestEvent) ParsePayload(payload map[string]interface{}) (err error) {
	marshalled, err := json.Marshal(payload)

	if err != nil {
		return
	}

	err = json.Unmarshal(marshalled, &tEvent.Payload)

	return
}

func NewBlankTestEvent() async_events.Event {
	return &TestEvent{
		BaseAsyncEvent: &async_events.BaseAsyncEvent{},
		Payload:        TestPayload{},
	}
}

func NewTestEvent(payload TestPayload) async_events.Event {
	return &TestEvent{
		BaseAsyncEvent: &async_events.BaseAsyncEvent{},
		Payload:        payload,
	}
}

func NewTestRdb() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "redisdb:6379",
	})
}
