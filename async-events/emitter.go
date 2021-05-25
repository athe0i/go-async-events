package async_events

type Emitter interface {
	Emit(event Event) error
	GetRouter() Router
	GetSerializer() Serializer
}
