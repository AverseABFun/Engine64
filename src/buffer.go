package main

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

func CreateRect(texture screen.Texture, x, y, width, height int, clr color.RGBA) {
	texture.Fill(image.Rect(x, y, x+width, y+height), clr, screen.Over)
}
