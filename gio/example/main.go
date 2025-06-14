package main

import (
	"fmt"

	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/gio"
)

func main() {
	driver, err := gio.GetDriver(gio.DriverTypeX11)
	if err != nil {
		panic(err)
	}

	window, err := driver.CreateWindow(gio.NewWindowOptions{
		Title: "Hello World",
		Size:  geom.NewSize[geom.F32](800, 600),
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Window created with ID:", window.Receive())
}
