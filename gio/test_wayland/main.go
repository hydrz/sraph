package main

import (
	"log"
	"time"

	"github.com/opensraph/sraph/gio"
	_ "github.com/opensraph/sraph/gio/wayland" // Import Wayland driver
)

func main() {
	log.Println("Testing Wayland driver...")

	// Try to create a window using Wayland driver
	window, err := gio.CreateWindowWithDriver(gio.DriverTypeWayland, gio.NewWindowOptions{
		Title:  "Wayland Test Window",
		Width:  800,
		Height: 600,
	})
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}

	log.Println("Window created successfully!")

	// Show the window
	if err := window.Show(); err != nil {
		log.Fatalf("Failed to show window: %v", err)
	}

	log.Println("Window shown successfully!")

	// Keep the window open for a few seconds
	time.Sleep(5 * time.Second)

	// Close the window
	if err := window.Close(); err != nil {
		log.Printf("Failed to close window: %v", err)
	}

	log.Println("Test completed!")
}
