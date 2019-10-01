package margui

import (
	"reflect"
)

type EventSubscriber struct {
	id    int
	event *Event
}

func (e *EventSubscriber) Unlisten() {
	if e.event != nil {
		e.event.unlisten(e.id)
		e.event = nil
	}
}

type Event struct {
	listeners []eventListener
	nextId    int
}

type eventListener struct {
	Id       int
	Function reflect.Value
}

func (e *Event) unlisten(id int) {
	for i, l := range e.listeners {
		if l.Id == id {
			copy(e.listeners[i:], e.listeners[i+1:])
			e.listeners = e.listeners[:len(e.listeners)-1]
			return
		}
	}
}

func (e *Event) Fire(args ...interface{}) {
	argVals := make([]reflect.Value, len(args))
	for i, arg := range args {
		argVals[i] = reflect.ValueOf(arg)
	}
	for _, l := range e.listeners {
		l.Function.Call(argVals)
	}
}
func (e *Event) Listen(listener interface{}) EventSubscriber {
	var function reflect.Value

	reflectTy := reflect.TypeOf(listener)
	if reflectTy.Kind() == reflect.Func {
		function = reflect.ValueOf(listener)

		id := e.nextId
		e.nextId++

		e.listeners = append(e.listeners, eventListener{
			Id:       id,
			Function: function,
		})
		return EventSubscriber{id: id, event: e}
	} else {
		panic(reflectTy.Kind())
	}
}

type MouseEventArgs struct {
	//Button (left, mid, right)
	//State ()
	//Point
	//Window
	//ScrollX, ScrollY int
	//Modifier KeyboardModifier
}

type KeyboardEventArgs struct {
	//Key
	//Modifier KeyboardModifier
}

var mouseDownEvent Event

func OnMouseDown(f func(MouseEventArgs)) EventSubscriber {
	return mouseDownEvent.Listen(f)
}

func xx() {
	mouseDownEvent.Fire(MouseEventArgs{})
}
