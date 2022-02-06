package events

import "github.com/dino16m/golearn/lib/event"

// Event names that are available for registration and dispatch
const (
	// expects map[string]interface{} with keys firstName, lastName, email,
	// and id which should have an int value
	UserCreated = event.EventName("usercreated")
)
