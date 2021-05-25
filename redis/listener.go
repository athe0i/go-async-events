package redis

import (
	"go-async-events/async-events"
)

func NewRedisListener(repo RedisRepositoryInterface, disp async_events.EventDispatcher) async_events.Listener {
	return async_events.NewEventListener(repo, disp)
}
