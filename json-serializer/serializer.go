package json_serializer

import (
	"encoding/json"
	"go-async-events/async-events"
)

const NameField = "name"
const PayloadField = "payload"

type JsonSerializer struct {
	async_events.EventRegistry
}

func NewJsonSerializer() (s async_events.Serializer) {
	return &JsonSerializer{EventRegistry: async_events.NewEventConstructorRegistry()}
}

func (serializer *JsonSerializer) Serialize(ev async_events.Event) (string, error) {
	evName, err := serializer.GetName(ev)
	if err != nil {
		return "", err
	}

	preserializedEvent := map[string]interface{}{
		NameField:            evName,
		PayloadField:         ev.GetPayload(),
	}
	valByte, err := json.Marshal(preserializedEvent)

	return string(valByte), err
}

func (serializer *JsonSerializer) Deserialize(jsonEvent string) (event async_events.Event, err error) {
	var unmarshalled map[string]interface{}
	err = json.Unmarshal([]byte(jsonEvent), &unmarshalled)

	if err != nil {
		return
	}

	e, err := serializer.CreateEvent(
		unmarshalled[NameField].(string),
		unmarshalled[PayloadField].(map[string]interface{}))

	return e, err
}
