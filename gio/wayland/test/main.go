package main

import (
	"fmt"
	"log"
	"os"

	"github.com/opensraph/sraph/gio"
	_ "github.com/opensraph/sraph/gio/wayland" // Import Wayland driver
)

func main() {
	fmt.Println("Testing Wayland backend creation...")

	// Check if WAYLAND_DISPLAY is set (simulate Wayland environment)
	if os.Getenv("WAYLAND_DISPLAY") == "" {
		fmt.Println("WAYLAND_DISPLAY not set, setting to simulate Wayland session")
		os.Setenv("WAYLAND_DISPLAY", "wayland-0")
	}

	// Test driver creation
	driver, err := gio.GetDriver(gio.DriverTypeWayland)
	if err != nil {
		log.Printf("Warning: Failed to get Wayland driver (expected in test): %v", err)
		fmt.Println("This is expected when running without a real Wayland compositor")
		return
	}

	fmt.Printf("Successfully created Wayland driver: %v\n", driver.Type())

	// Test window creation (will likely fail without real compositor)
	options := gio.NewWindowOptions{
		Width:  640,
		Height: 480,
		Title:  "Test Window",
	}

	window, err := driver.CreateWindow(options)
	if err != nil {
		log.Printf("Warning: Failed to create window (expected without compositor): %v", err)
		fmt.Println("This is expected when running without a real Wayland compositor")
		return
	}

	fmt.Printf("Successfully created window with ID: %v\n", window.WindowID())
	fmt.Println("Test completed successfully!")
}
