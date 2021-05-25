package redis

import (
	"github.com/go-redis/redis/v8"
	"go-async-events/async-events"
)

type RedisChannel struct {
	repo RedisRepositoryInterface
	name string
	slr  async_events.Serializer
}

func NewRedisChannel(name string, repo RedisRepositoryInterface, serializer async_events.Serializer) *RedisChannel {
	return &RedisChannel{
		repo: repo,
		name: name,
		slr:  serializer,
	}
}

func (channel *RedisChannel) Push(asyncEvent async_events.Event) error {
	serialized, err := channel.slr.Serialize(asyncEvent)
	if err != nil {
		return err
	}

	err = channel.repo.PushToChannel(channel.GetName(), serialized)

	return err
}

func (channel *RedisChannel) Pop() (async_events.Event, error) {
	channel.KeepAlive()
	value, err := channel.repo.PopFromChannel(channel.GetName())

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	return channel.slr.Deserialize(value)
}

func (channel *RedisChannel) GetName() string {
	return channel.name
}

func (channel *RedisChannel) KeepAlive() {
	channel.repo.KeepAlive(channel.GetName())
}

func (channel *RedisChannel) Close() {
	channel.repo.ClearChannel(channel.GetName())
}
