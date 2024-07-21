package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	graphics "github.com/averseabfun/engine64/graphics"
	logger "github.com/averseabfun/engine64/logger"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var tri = graphics.CreateSortedTriangle(graphics.Triangle{[3]image.Point{{0, 0}, {100, 0}, {0, 100}}})
var tri2 = graphics.CreateSortedTriangle(graphics.Triangle{[3]image.Point{{100, 0}, {0, 100}, {100, 100}}})

func Draw(window screen.Window, actualScreen screen.Screen, texture *screen.Texture, size image.Point) {
	if !(*texture).Bounds().Size().Eq(size) {
		(*texture).Release()
		(*texture) = graphics.CreateTexture(window, actualScreen, size)
	}
	(*texture).Fill(image.Rect(0, 0, size.X, size.Y), color.RGBA{0, 0, 0, 255}, screen.Over)

	graphics.CreateTri((*texture), tri, color.RGBA{255, 0, 0, 255})
	graphics.CreateTri((*texture), tri2, color.RGBA{0, 255, 0, 255})

	window.Draw(f64.Aff3{1, 0, 0, 0, 1, 0}, (*texture), (*texture).Bounds(), draw.Over, nil)
	window.Publish()
}

func HandleEvents(window screen.Window, actualScreen screen.Screen, windowSize *image.Point, texture *screen.Texture) {
	switch e := window.NextEvent().(type) {

	case lifecycle.Event:
		logger.Log(fmt.Sprintf("Window Event: From %s To %s", e.From, e.To), logger.LogDebug)
		if e.To == lifecycle.StageDead {
			logger.Log("Window closed by user", logger.LogInfo)
			os.Exit(0)
		}

		if e.To == lifecycle.StageFocused {
			logger.Log("Window now has the focus", logger.LogDebug)
		}
		if e.From == lifecycle.StageFocused {
			logger.Log("Window has lost the focus", logger.LogDebug)
		}
		logger.LogEmptyNewline()
	case size.Event:
		logger.Log(fmt.Sprintf("Size Event: Width %d Height %d", e.WidthPx, e.HeightPx), logger.LogDebug)
		*windowSize = e.Size()
		(*texture).Release()
		*texture = logger.CreateTexture(window, actualScreen, *windowSize)
		logger.LogEmptyNewline()
		Draw(window, actualScreen, texture, *windowSize)
	case paint.Event:
		if !e.External {
			Draw(window, actualScreen, texture, *windowSize)
		}
	}
}

func main() {
	logger.Log("Initalizing Engine64", logger.LogInfo)
	graphics.OpenWindow("Engine64", 800, 650, func(window screen.Window, actualScreen screen.Screen) {
		logger.Log("Window created", logger.LogDebug)

		var windowSize = image.Point{800, 650}

		var texture = graphics.CreateTexture(window, actualScreen, windowSize)
		defer texture.Release()
		defer window.Release()

		Draw(window, actualScreen, &texture, windowSize)

		for {
			HandleEvents(window, actualScreen, &windowSize, &texture)
		}
	})
}
