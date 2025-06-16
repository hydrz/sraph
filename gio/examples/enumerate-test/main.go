package main

import (
	"fmt"
	"log"

	"github.com/opensraph/sraph/gio"
	_ "github.com/opensraph/sraph/gio/x11" // Import X11 driver to register it
)

func main() {
	fmt.Println("=== Gio Backend Enumeration Test ===")

	// Show platform information
	platformInfo := gio.GetPlatformInfo()
	fmt.Printf("Platform: %s\n", platformInfo.String())
	fmt.Println()

	// Enumerate all available drivers
	fmt.Println("Enumerating drivers:")
	drivers := gio.EnumerateDrivers()
	if len(drivers) == 0 {
		fmt.Println("  No drivers available")
	} else {
		for i, driverType := range drivers {
			available := driverType.IsAvailable()
			fmt.Printf("  %d. %s (available: %v)\n", i+1, driverType.String(), available)
		}
	}
	fmt.Println()

	// Try to get the best driver
	fmt.Println("Getting best driver...")
	bestDriver, err := gio.GetBestDriver()
	if err != nil {
		fmt.Printf("Failed to get best driver: %v\n", err)
		return
	}

	fmt.Printf("Best driver found: %s\n", bestDriver.Type().String())
	fmt.Println()

	// Test creating a window (without showing it)
	fmt.Println("Testing window creation...")
	window, err := gio.CreateWindow(gio.NewWindowOptions{
		Title:  "Test Window",
		Width:  400,
		Height: 300,
	})
	if err != nil {
		log.Printf("Failed to create window: %v", err)
		return
	}

	fmt.Printf("Successfully created window with ID: %v\n", window.WindowID())
	fmt.Printf("Window size: %dx%d\n", window.Width(), window.Height())
	fmt.Printf("Window title: %s\n", window.Title())
	fmt.Println()

	// Test specific driver creation
	fmt.Println("Testing specific driver creation...")
	for _, driverType := range drivers {
		if driverType.IsAvailable() {
			fmt.Printf("Creating window with %s driver...\n", driverType.String())
			specificWindow, err := gio.CreateWindowWithDriver(driverType, gio.NewWindowOptions{
				Title:  fmt.Sprintf("Test Window (%s)", driverType.String()),
				Width:  300,
				Height: 200,
			})
			if err != nil {
				fmt.Printf("  Failed: %v\n", err)
			} else {
				fmt.Printf("  Success! Window ID: %v\n", specificWindow.WindowID())
				// Close the test window
				ctx := gio.NewContext(nil)
				specificWindow.Close(ctx)
			}
		}
	}

	// Close the main test window
	ctx := gio.NewContext(nil)
	window.Close(ctx)

	fmt.Println("\nAll tests completed successfully!")
}
