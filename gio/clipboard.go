package gio

type ClipboardData = DataTransfer

type ClipboardEvent struct {
	BaseEvent
	Data ClipboardData // Data associated with the clipboard event
}
