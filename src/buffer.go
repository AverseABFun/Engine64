package main

import (
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

func CreateRect(buffer screen.Buffer, x, y, width, height int, r, g, b, a uint8) {
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			buffer.RGBA().Set(i, j, color.RGBA{r, g, b, a})
		}
	}
}
