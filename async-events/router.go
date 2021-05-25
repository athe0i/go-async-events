package async_events

type Router interface {
	RegisterEventRoute(ev Event, channelName string) error
	GetEventRoutes(ev Event) ([]string, error)
}
