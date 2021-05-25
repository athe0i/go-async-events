package json_serializer

import (
	"errors"
	"go-async-events/async-events"
	"go-async-events/test"
	"reflect"
	"testing"
)

func TestJsonSerializer_Serialize(t *testing.T) {
	e := test.TestEvent{
		BaseAsyncEvent: &async_events.BaseAsyncEvent{},
		Payload: test.TestPayload{
			Number: 124,
			Text:   "Test",
			SubField: test.TestSubPayload{
				Field: "SubField",
			},
		},
	}
	e.SetName("test")

	serializer := NewJsonSerializer()
	serializer.(async_events.EventRegistry).RegisterEvent("test", test.NewBlankTestEvent)

	serialized, err := serializer.Serialize(&e)
	if err != nil {
		t.Errorf("Cannot Serialize Event: %s", err.Error())
		return
	}

	testJson := "{\"name\":\"test\",\"payload\":{\"Number\":124,\"Text\":\"Test\",\"SubField\":{\"Field\":\"SubField\"}}}"
	if serialized != testJson {
		t.Errorf("Serialization is not working correctly!")
		return
	}

	parsed, err := serializer.Deserialize(serialized)
	if err != nil {
		t.Errorf("Cannot Serialize Event: %s", err.Error())
		return
	}

	if reflect.DeepEqual(&e, parsed) == false {
		t.Errorf("Deserialization failed! Two objects are not equal")
		return
	}

	testErrorJson := "{\"name\":\"non-existing\",\"payload\":{\"Number\":124,\"Text\":\"Test\",\"SubField\":{\"Field\":\"SubField\"}}}"
	_, err = serializer.Deserialize(testErrorJson)

	if errors.As(err, &async_events.NoEventRegistered{}) == false {
		t.Errorf("It should return error if there is no such event registered!")
		return
	}
}
