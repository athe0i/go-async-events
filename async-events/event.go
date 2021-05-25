package async_events

type Event interface {
	//GetName() string
	SetName(name string)
	GetPayload() interface{}
	ParsePayload(payload map[string]interface{}) (err error)
}

type BaseAsyncEvent struct {
	name            string
}

func (b *BaseAsyncEvent) SetName(name string) {
	b.name = name
}
