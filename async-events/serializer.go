package async_events

type Serializer interface {
	Serialize(event Event) (string, error)
	Deserialize(str string) (Event, error)
	EventRegistry
}
