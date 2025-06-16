package gio

type DragEvent struct {
	BaseEvent
	Data DataTransfer // Data associated with the drag event
}
