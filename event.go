package eventbus

import (
	"fmt"
)

type EventName string

type Payload []interface{}

// Event
type Event struct {
	Name    EventName
	Payload Payload
}

func (e Event) String() string {

	return fmt.Sprintf("Event `%s` \n Payload: \n%v\n", e.Name, e.Payload)
}
