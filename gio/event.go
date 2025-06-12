package gio

import (
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
)

type KeyEvent = key.Event

type LifecycleEvent = lifecycle.Event

type MouseEvent = mouse.Event

type PaintEvent = paint.Event

type SizeEvent = size.Event

type TouchEvent = touch.Event

// EventDeque is an infinitely buffered double-ended queue of events.
type EventDeque interface{}
