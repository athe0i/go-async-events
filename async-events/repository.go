package async_events

type Repository interface {
	PushToChannel(channel string, ev string) error
	PopFromChannel(channel string) (string, error)
	KeepAlive(channel string)
	GetEventListenerChannelName(evName string) string
	FindActiveChannelsByName(name string) (keys []string, err error)
}
