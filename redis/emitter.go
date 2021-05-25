package redis

import (
	"go-async-events/async-events"
)

type RedisEmitter struct {
	repo RedisRepositoryInterface
	rtr  async_events.Router
	srlr async_events.Serializer
}

func NewRedisEmitter(repo RedisRepositoryInterface, rtr async_events.Router, srlr async_events.Serializer) async_events.Emitter {
	return &RedisEmitter{
		repo: repo,
		rtr:  rtr,
		srlr: srlr,
	}
}

func (re *RedisEmitter) Emit(ev async_events.Event) error {
	channels := re.findAllActiveChannels(ev)

	for _, chn := range channels {
		rchn := NewRedisChannel(chn, re.repo, re.srlr)
		err := rchn.Push(ev)

		if err != nil {
			return err
		}
	}

	return nil
}

func (re *RedisEmitter) GetRouter() async_events.Router {
	return re.rtr
}

func (re *RedisEmitter) GetSerializer() async_events.Serializer {
	return re.srlr
}

func (re *RedisEmitter) findAllActiveChannels(ev async_events.Event) []string {
	routerChannels := re.findActiveRouterChannels(ev)
	eventListeners, _ := re.findEventListenersChannels(ev)

	return append(routerChannels, eventListeners...)
}

func (re *RedisEmitter) findActiveRouterChannels(ev async_events.Event) []string {
	routerchn, _ := re.rtr.GetEventRoutes(ev)

	active := []string{}

	for _, chn := range routerchn {
		activeChannels, _ := re.findActiveChannels(chn)
		active = append(active, activeChannels...)
	}

	return active
}

func (re *RedisEmitter) findEventListenersChannels(ev async_events.Event) (chns []string, err error) {
	evName, err := re.srlr.GetName(ev)

	if err != nil {
		return
	}

	listenerName := re.repo.GetEventListenerChannelName(evName)

	chns, err = re.findActiveChannels(listenerName)

	return
}

func (re *RedisEmitter) findActiveChannels(name string) (channels []string, err error) {
	return re.repo.FindActiveChannelsByName(name)
}
