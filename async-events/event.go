package async_events

type Event interface {
	//GetName() string
	SetName(name string)
	GetPayload() interface{}
	ParsePayload(payload map[string]interface{}) (err error)
	GetCallbackChannel() string
	SetCallbackChannel(callbackChannel string)
}

type BaseAsyncEvent struct {
	name            string
	callbackChannel string
}

func (b *BaseAsyncEvent) SetName(name string) {
	b.name = name
}

func (b *BaseAsyncEvent) GetCallbackChannel() string {
	return b.callbackChannel
}

func (b *BaseAsyncEvent) SetCallbackChannel(callbackChannel string) {
	b.callbackChannel = callbackChannel
}
