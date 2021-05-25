package async_events

type WriterChannel interface {
	Push(asyncEvent Event) error
	NamedChannel
}

type ReaderChannel interface {
	Pop() (Event, error)
	NamedChannel
}

type NamedChannel interface {
	GetName() string
}