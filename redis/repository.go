package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	async_events "go-async-events/async-events"
	"strings"
	"time"
)

const KEEP_ALIVE_PREFIX = "keep_alive"
const EVENT_LISTENER_CHANNEL_PREFIX = "event_listener_"

type RedisRepositoryInterface interface {
	async_events.Repository
	ClearChannel(channel string)
	GetKeepAliveKey(channel string) string
}

type Repository struct {
	rdb  *redis.Client
	opts RedisOptions
	ctx  context.Context
}

type RedisOptions struct {
	Prefix           string
	KeepAliveTimeout int
}

func NewRepository(rdb *redis.Client, opts RedisOptions) RedisRepositoryInterface {
	return &Repository{
		rdb:  rdb,
		opts: opts,
		ctx:  context.Background(),
	}
}

func (repo *Repository) PushToChannel(channel string, ev string) (err error) {
	return repo.rdb.RPush(repo.ctx, repo.GetPrefixedKey(channel), ev).Err()
}

func (repo *Repository) PopFromChannel(channel string) (value string, err error) {
	return repo.rdb.LPop(repo.ctx, repo.GetPrefixedKey(channel)).Result()
}

func (repo *Repository) GetKeepAliveKey(channel string) string {
	return repo.getKeepAlivePrefix() + channel
}

func (repo *Repository) GetPrefixedKey(name string) string {
	return repo.opts.Prefix + name
}

func (repo *Repository) FindActiveChannelsByName(name string) (keys []string, err error) {
	pattern := fmt.Sprintf("%s:*", repo.GetKeepAliveKey(name))
	keys, err = repo.rdb.Keys(repo.ctx, pattern).Result()

	if err == nil {
		keys = repo.getChannelFromKeepAlive(keys)
	}
	return
}

func (repo *Repository) ClearChannel(name string) {
	repo.rdb.Expire(repo.ctx, repo.GetPrefixedKey(name), 0)
	repo.rdb.Del(repo.ctx, repo.GetKeepAliveKey(name))
}

func (repo *Repository) GetEventListenerChannelName(evName string) string {
	return EVENT_LISTENER_CHANNEL_PREFIX + evName
}

func (repo *Repository) KeepAlive(channel string) {
	repo.rdb.Set(repo.ctx, repo.GetKeepAliveKey(channel), true, time.Duration(repo.opts.KeepAliveTimeout)*time.Second)
}

func (repo *Repository) getKeepAlivePrefix() string {
	return repo.GetPrefixedKey(KEEP_ALIVE_PREFIX + ":")
}

func (repo *Repository) getChannelFromKeepAlive(chns []string) []string {
	for i, chn := range chns {
		chns[i] = strings.Replace(chn, repo.getKeepAlivePrefix(), "", 1)
	}

	return chns
}
