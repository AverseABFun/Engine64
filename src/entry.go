package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	formatColor "github.com/fatih/color"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/size"
)

const (
	// LogInfo is used to log information messages
	LogInfo = iota
	// LogError is used to log error messages
	LogError
	// LogWarning is used to log warning messages
	LogWarning
	// LogDebug is used to log debug messages
	LogDebug
)

var logger = log.New(os.Stderr, "", log.Lmsgprefix|log.Ltime)

func Log(msg string, logType int) {
	msg = formatColor.MagentaString(msg)
	switch logType {
	case LogInfo:
		logger.SetPrefix(formatColor.GreenString("[INFO] "))
	case LogError:
		logger.SetPrefix(formatColor.RedString("[ERROR] "))
	case LogWarning:
		logger.SetPrefix(formatColor.YellowString("[WARNING] "))
	case LogDebug:
		logger.SetPrefix(formatColor.BlueString("[DEBUG] "))
	default:
		logger.SetPrefix(formatColor.WhiteString("[INFO] "))
	}
	logger.Println(msg)
}

func LogEmptyNewline() {
	logger.SetPrefix("")
	logger.SetFlags(0)
	logger.Print("\n")
	logger.SetFlags(log.Lmsgprefix | log.Ltime)
}

var tri = CreateSortedTriangle(Triangle{[3]image.Point{{0, 0}, {100, 0}, {0, 100}}})

func Draw(window screen.Window, actualScreen screen.Screen, texture screen.Texture, size image.Point) {
	window.Fill(image.Rect(0, 0, size.X, size.Y), color.NRGBA{0, 0, 0, 255}, screen.Over)

	CreateTri(texture, tri, color.RGBA{255, 0, 0, 255})

	window.Draw(f64.Aff3{1, 0, 0, 0, 1, 0}, texture, texture.Bounds(), draw.Over, nil)
	window.Publish()
}

func CreateTexture(window screen.Window, actualScreen screen.Screen, size image.Point) screen.Texture {
	var texture screen.Texture

	texture, err := actualScreen.NewTexture(size)
	if err != nil {
		Log(fmt.Sprintf("Failed to create texture - %v", err), LogError)
		return nil
	}
	return texture
}

func main() {
	Log("Initalizing Engine64", LogInfo)
	OpenWindow("Engine64", 800, 650, func(window screen.Window, actualScreen screen.Screen) {
		Log("Window created", LogDebug)

		var windowSize = image.Point{800, 650}

		var texture = CreateTexture(window, actualScreen, windowSize)
		defer texture.Release()

		Draw(window, actualScreen, texture, windowSize)

		defer window.Release()

		var cnt int

		for {
			Draw(window, actualScreen, texture, windowSize)
			switch e := window.NextEvent().(type) {

			case lifecycle.Event:
				cnt++
				Log(fmt.Sprintf("Window Event %d: From %s To %s", cnt, e.From, e.To), LogDebug)
				if e.To == lifecycle.StageDead {
					Log("Window closed by user", LogInfo)
					return // quit the application.
				}

				if e.To == lifecycle.StageFocused {
					Log("Window now has the focus", LogDebug)
				}
				if e.From == lifecycle.StageFocused {
					Log("Window has lost the focus", LogDebug)
				}
				LogEmptyNewline()
			case size.Event:
				cnt++
				Log(fmt.Sprintf("Size Event %d: Width %d Height %d", cnt, e.WidthPx, e.HeightPx), LogDebug)
				windowSize = image.Point{e.WidthPx, e.HeightPx}
				texture.Release()
				texture = CreateTexture(window, actualScreen, windowSize)
				LogEmptyNewline()
			}
		}
	})
}
