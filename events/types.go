package events

import "github.com/dino16m/golearn/lib/event"

// Event names that are available for registration and dispatch
const (
	// expects map[string]interface{} with keys firstName, lastName, email,
	// and id which should have an int value
	UserCreated = event.EventName("usercreated")
)

type User struct {
	name string
}

type Topic struct {
	id string
}

func (topic *Topic) ID() event.EventName {
	return event.EventName(topic.id)
}

type newUserCreated struct {
	Topic
}

func (topic *newUserCreated) With(data User) event.Event {
	return event.NewEvent(topic.ID(), data)
}

var NewUserCreated = newUserCreated{Topic{"new-user-created"}}
