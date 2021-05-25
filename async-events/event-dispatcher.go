package async_events

type EventHandler func(event Event) (err error)

type EventDispatcher interface {
	RegisterEventHandler(event Event, handler EventHandler)
	HandleEvent(asyncEvent Event) (err error)
}
