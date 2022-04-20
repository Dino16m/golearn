package container

import (
	"sync"

	"github.com/dino16m/golearn-core/event"
	"github.com/google/wire"
)

var eventDispatcher *event.EventDispatcher
var once sync.Once

func provideEventDispatcher() *event.EventDispatcher {
	once.Do(func() { eventDispatcher = event.NewEventDispatcher() })
	return eventDispatcher
}

var EventSet = wire.NewSet(wire.Bind(new(event.Dispatcher), new(*event.EventDispatcher)), provideEventDispatcher)
