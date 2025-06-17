package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/opensraph/sraph/gio"
	_ "github.com/opensraph/sraph/gio/wayland" // Import Wayland driver
)

func main() {
	// Set environment variable to prefer Wayland
	os.Setenv("GIO_DRIVER", "wayland")

	fmt.Println("Creating Wayland window...")

	driver, err := gio.GetDriver(gio.DriverTypeWayland)
	if err != nil {
		log.Fatalf("Failed to get Wayland driver: %v", err)
	}

	fmt.Printf("Driver type: %v\n", driver.Type())

	options := gio.NewWindowOptions{
		Width:  800,
		Height: 600,
		Title:  "Wayland Test Window",
	}

	window, err := driver.CreateWindow(options)
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}

	fmt.Printf("Window created with ID: %v\n", window.WindowID())

	fmt.Println("Window created successfully. Setting up event handlers...")

	// Subscribe to keyboard events
	err = window.Subscribe(gio.EventTypeKeyboard, func(e gio.Event) error {
		if keyEvent, ok := e.(*gio.KeyboardEvent); ok {
			fmt.Printf("Key event: %v (modifiers: %v)\n",
				keyEvent.Code, keyEvent.ModifierKey)
			if keyEvent.Code == gio.KeyCodeEscape {
				fmt.Println("Escape pressed, closing window...")
				attr := window.Attr()
				attr.State |= gio.WindowStateClosed
				window.SetAttr(attr)
				return fmt.Errorf("window closed by user")
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to keyboard events: %v", err)
	}

	// Subscribe to mouse events
	err = window.Subscribe(gio.EventTypeMouse, func(e gio.Event) error {
		if mouseEvent, ok := e.(*gio.MouseEvent); ok {
			fmt.Printf("Mouse event: button=%v, pos=(%d,%d), modifiers=%v\n",
				mouseEvent.Button, mouseEvent.Position.X, mouseEvent.Position.Y, mouseEvent.ModifierKey)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to mouse events: %v", err)
	}

	// Subscribe to wheel events
	err = window.Subscribe(gio.EventTypeWheel, func(e gio.Event) error {
		if wheelEvent, ok := e.(*gio.WheelEvent); ok {
			fmt.Printf("Wheel event: delta=(%f,%f), pos=(%d,%d)\n",
				wheelEvent.DeltaX, wheelEvent.DeltaY, wheelEvent.Position.X, wheelEvent.Position.Y)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to wheel events: %v", err)
	}

	// Subscribe to window events
	err = window.Subscribe(gio.EventTypeWindow, func(e gio.Event) error {
		if windowEvent, ok := e.(*gio.WindowEvent); ok {
			fmt.Printf("Window event: %T\n", windowEvent)
			// Check if window is closed
			if window.IsClosed() {
				fmt.Println("Window was closed by user")
				return fmt.Errorf("window closed")
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to window events: %v", err)
	}

	fmt.Println("Event handlers set up. Press Ctrl+C to exit")

	// Simple event loop - just wait a bit to allow events
	for i := 0; i < 100; i++ { // Run for a limited time
		time.Sleep(100 * time.Millisecond)
		if window.IsClosed() {
			break
		}
	}

	fmt.Println("Example completed successfully")
}
