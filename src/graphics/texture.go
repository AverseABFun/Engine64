package graphics

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	logger "github.com/averseabfun/engine64/logger"
	"golang.org/x/exp/shiny/screen"
)

type Triangle struct {
	Points [3]image.Point
}

type sortedTriangle struct {
	Points  [3]image.Point
	TopY    int
	LeftX   int
	BottomY int
	RightX  int
}

func (triangle Triangle) CreateSortedTriangle() sortedTriangle {
	var sortedPoints [3]image.Point = triangle.Points
	var topY, bottomY int = sortedPoints[0].Y, sortedPoints[0].Y
	var leftX, rightX int = sortedPoints[0].X, sortedPoints[0].X
	for i := 0; i < 3; i++ {
		if triangle.Points[i].Y < topY {
			topY = triangle.Points[i].Y
		}
		if triangle.Points[i].Y > bottomY {
			bottomY = triangle.Points[i].Y
		}
		if triangle.Points[i].X < leftX {
			leftX = triangle.Points[i].X
		}
		if triangle.Points[i].X > rightX {
			rightX = triangle.Points[i].X
		}
	}
	if sortedPoints[0].Y > sortedPoints[1].Y {
		sortedPoints[0], sortedPoints[1] = sortedPoints[1], sortedPoints[0]
	}
	if sortedPoints[1].Y > sortedPoints[2].Y {
		sortedPoints[1], sortedPoints[2] = sortedPoints[2], sortedPoints[1]
	}
	if sortedPoints[0].Y > sortedPoints[1].Y {
		sortedPoints[0], sortedPoints[1] = sortedPoints[1], sortedPoints[0]
	}
	return sortedTriangle{Points: sortedPoints, TopY: topY, BottomY: bottomY, LeftX: leftX, RightX: rightX}
}

func CreateSortedTriangle(triangle Triangle) sortedTriangle {
	return triangle.CreateSortedTriangle()
}

func (t sortedTriangle) PointInTriangle(p image.Point) bool {
	b1 := sign(p, t.Points[0], t.Points[1]) < 0.0
	b2 := sign(p, t.Points[1], t.Points[2]) < 0.0
	b3 := sign(p, t.Points[2], t.Points[0]) < 0.0

	return ((b1 == b2) && (b2 == b3))
}

func sign(p1, p2, p3 image.Point) float64 {
	return float64(p1.X-p3.X)*float64(p2.Y-p3.Y) - float64(p2.X-p3.X)*float64(p1.Y-p3.Y)
}

func CreateTexture(window screen.Window, actualScreen screen.Screen, size image.Point) screen.Texture {
	var texture screen.Texture

	texture, err := actualScreen.NewTexture(size)
	if err != nil {
		logger.Log(fmt.Sprintf("Failed to create texture - %v", err), logger.LogError)
		return nil
	}
	return texture
}

func CreateRect(texture screen.Texture, x, y, width, height int, clr color.Color) {
	texture.Fill(image.Rect(x, y, x+width, y+height), clr, screen.Over)
}

func CreateTri(texture screen.Texture, triangle sortedTriangle, clr color.Color) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == i {
				continue
			}
			DrawLine(texture, triangle.Points[i].X, triangle.Points[i].Y, triangle.Points[j].X, triangle.Points[j].Y, clr)
		}
	}
	for y := triangle.TopY; y < triangle.BottomY; y++ {
		var lineStart, lineEnd int = triangle.RightX, triangle.RightX
		var inTriangle bool = false
		for x := triangle.LeftX; x < triangle.RightX; x++ {
			if triangle.PointInTriangle(image.Point{x, y}) && !inTriangle {
				inTriangle = true
				lineStart = x
			} else if inTriangle && !triangle.PointInTriangle(image.Point{x, y}) {
				inTriangle = false
				lineEnd = x
				break
			}
		}
		if lineEnd >= lineStart {
			texture.Fill(image.Rect(lineStart, y, lineEnd, y+1), clr, screen.Over)
		} else {
			logger.Log("LineEnd is less than LineStart", logger.LogError)
			logger.Log("LineStart: "+strconv.Itoa(lineStart), logger.LogError)
			logger.Log("LineEnd: "+strconv.Itoa(lineEnd), logger.LogError)
		}
	}
}

func DrawLine(w screen.Texture, x0, y0, x1, y1 int, col color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy

	for {
		w.Fill(image.Rect(x0, y0, x0+1, y0+1), col, screen.Over)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
