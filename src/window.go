package main

import (
	"fmt"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
)

// OpenWindow creates a new window with the specified title, width, and height.
// If the width or height is less than or equal to 0, the default values are used.
// The default width is 800 and the default height is 650.
// The callback function is called with the created window and the current screen object. Please note that the window is *not* closed automatically. Also note that no operations are performed on the window by this function EVER.
func OpenWindow(title string, width int, height int, callback func(window screen.Window, screen screen.Screen)) {
	if width <= 0 {
		width = 800
	}
	if height <= 0 {
		height = 650
	}
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title:  title,
			Width:  width,
			Height: height,
		})
		if err != nil {
			fmt.Printf("Failed to create Window - %v", err)
			return
		}
		callback(w, s)
	})
}
