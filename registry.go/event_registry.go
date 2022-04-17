package registry

import (
	"github.com/dino16m/golearn/dependencies"
	"github.com/dino16m/golearn/lib/event"
)

// RegisterEventHandlers binds events to handlers which listen to them
func RegisterEventHandlers(c dependencies.EventHandlersContainer,
	app dependencies.App) {
	dispatcher := app.EventDispatcher
	dispatcher.AddListeners(event.UserCreated, c.UserCreatedHandler)
}
