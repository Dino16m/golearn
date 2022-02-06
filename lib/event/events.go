package event

import (
	"math/rand"
	"sync"
)

// AuthEventDispatcher ...
type AuthEventDispatcher struct {
	listeners map[EventName][]Listener
	ids       map[int]bool
	mutex     *sync.Mutex
}

// NewAuthEventDispatcher construct the event dispatcher
func NewAuthEventDispatcher() *AuthEventDispatcher {
	return &AuthEventDispatcher{
		mutex:     &sync.Mutex{},
		ids:       make(map[int]bool),
		listeners: make(map[EventName][]Listener),
	}
}

// AddListeners ...
func (a *AuthEventDispatcher) AddListeners(name EventName, listeners ...Listener) {
	a.listeners[name] = append(a.listeners[name], listeners...)
	for _, listener := range listeners {
		id := a.getID()
		listener.SetId(id)
	}
}

func (a *AuthEventDispatcher) getID() int {

	var id int
	a.mutex.Lock()
	for {
		id = rand.Int()
		exists := a.ids[id]
		if exists == false {
			a.ids[id] = true
			break
		}
	}
	a.mutex.Unlock()
	return id
}

// Dispatch emits an event
func (a *AuthEventDispatcher) Dispatch(name EventName, payload interface{}) {
	listeners := a.listeners[name]
	for _, listener := range listeners {
		listener.Handle(payload)
	}
}

// RemoveListener unsubscribes a listener from a particular events
func (a *AuthEventDispatcher) RemoveListener(
	name EventName, listenerToRemove Listener,
) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	listeners := a.listeners[name]
	if len(listeners) < 1 {
		return
	}
	listenersCopy := listeners
	for i, listener := range listeners {
		if listener.GetId() == listenerToRemove.GetId() {
			listenersCopy = spliceSlice(listenersCopy, i)
		}
	}
	a.listeners[name] = listenersCopy
	delete(a.ids, listenerToRemove.GetId())
}

func spliceSlice(slice []Listener, index int) []Listener {
	if len(slice) == 1 {
		return slice[0:0]
	}
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
