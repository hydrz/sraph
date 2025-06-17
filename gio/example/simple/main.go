package main

import (
	"fmt"
	"log"
	"time"

	"github.com/opensraph/sraph/gio"
	_ "github.com/opensraph/sraph/gio/wayland" // Import Wayland driver to register it
	_ "github.com/opensraph/sraph/gio/x11"     // Import X11 driver to register it
)

func main() {
	fmt.Println("=== Gio Graphics Backend Demo ===")

	// Show platform information
	platformInfo := gio.GetPlatformInfo()
	fmt.Printf("Platform Info: %s\n", platformInfo.String())
	fmt.Println()

	// Enumerate all available drivers
	fmt.Println("Available drivers:")
	drivers := gio.EnumerateDrivers()
	for _, driverType := range drivers {
		fmt.Printf("  - %s (available: %v)\n", driverType.String(), driverType.IsAvailable())
	}
	fmt.Println()

	// Get the best driver
	fmt.Println("Getting best driver...")
	bestDriver, err := gio.GetBestDriver()
	if err != nil {
		log.Fatalf("Failed to get best driver: %v", err)
	}
	fmt.Printf("Best driver: %s\n", bestDriver.Type().String())
	fmt.Println()

	// Create a window using the automatic driver selection
	fmt.Println("Creating window with automatic driver selection...")
	window, err := gio.CreateWindow(gio.NewWindowOptions{
		Title:  "Gio Auto-Driver Test",
		Width:  800,
		Height: 600,
		State:  gio.WindowStateVisible | gio.WindowStateResizable,
	})
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}

	fmt.Printf("Successfully created window with ID: %v\n", window.WindowID())
	fmt.Printf("Window properties: %dx%d at %v\n",
		window.Width(), window.Height(), window.Position())
	fmt.Println()

	// Subscribe to keyboard events
	window.Subscribe(gio.EventTypeKeyboard, func(e gio.Event) error {
		if keyEvent, ok := e.(*gio.KeyboardEvent); ok {
			fmt.Printf("[KEYBOARD] %s\n", keyEvent)
		}
		return nil
	})

	// Subscribe to pointer events
	window.Subscribe(gio.EventTypePointer, func(e gio.Event) error {
		if pointerEvent, ok := e.(*gio.PointerEvent); ok {
			fmt.Printf("[POINTER]  %s\n", pointerEvent)
		}
		return nil
	})

	// Subscribe to wheel events
	window.Subscribe(gio.EventTypeWheel, func(e gio.Event) error {
		if wheelEvent, ok := e.(*gio.WheelEvent); ok {
			fmt.Printf("[WHEEL] %s\n", wheelEvent)
		}
		return nil
	})

	// Subscribe to clipboard events
	window.Subscribe(gio.EventTypeClipboard, func(e gio.Event) error {
		if clipboardEvent, ok := e.(*gio.ClipboardEvent); ok {
			fmt.Printf("[CLIPBOARD] %s\n", clipboardEvent)
		}
		return nil
	})

	// Subscribe to drag and drop events
	window.Subscribe(gio.EventTypeDrag, func(e gio.Event) error {
		if dragEvent, ok := e.(*gio.DragEvent); ok {
			fmt.Printf("[DRAG] %s\n", dragEvent)
		}
		return nil
	})

	// Subscribe to window events
	window.Subscribe(gio.EventTypeWindow, func(e gio.Event) error {
		if windowEvent, ok := e.(*gio.WindowEvent); ok {
			fmt.Printf("[WINDOW] %s\n", windowEvent)

			// Check if window is closed
			if window.IsClosed() {
				fmt.Println("Window was closed by user")
				return fmt.Errorf("window closed")
			}
		}
		return nil
	})

	// Show the window
	if err := window.Show(); err != nil {
		log.Printf("Warning: Failed to show window: %v", err)
	}

	fmt.Println("Window is now visible!")
	fmt.Println("Try interacting with the window:")
	fmt.Println("  - Press keys to see keyboard events")
	fmt.Println("  - Click and move mouse to see mouse events")
	fmt.Println("  - Use mouse wheel to see scroll events")
	fmt.Println("  - Close the window to exit")
	fmt.Println()
	fmt.Println("The demo will run for 60 seconds or until you close the window...")

	// Keep the program running and monitor for window closure
	timeout := time.After(60 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("\nTimeout reached. Closing window...")
			goto cleanup

		case <-ticker.C:
			// Check if window is still open
			if window.IsClosed() {
				fmt.Println("\nWindow was closed. Exiting...")
				goto cleanup
			}

		default:
			// Brief pause to prevent busy waiting
			time.Sleep(10 * time.Millisecond)
		}
	}

cleanup:
	// Close the window if it's still open
	if !window.IsClosed() {
		if err := window.Close(); err != nil {
			log.Printf("Failed to close window: %v", err)
		}
	}

	fmt.Println("Demo completed successfully!")

	// Show final platform info
	fmt.Println("\nFinal platform info:")
	fmt.Printf("  Platform: %s\n", platformInfo.String())
	fmt.Printf("  Used driver: %s\n", bestDriver.Type().String())
}
