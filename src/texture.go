package main

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

type Triangle struct {
	Points [3]image.Point
}

type sortedTriangle struct {
	Points [3]image.Point
}

func CreateSortedTriangle(triangle Triangle) sortedTriangle {
	var sortedPoints [3]image.Point = triangle.Points
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if sortedPoints[i].Y < sortedPoints[j].Y {
				temp := sortedPoints[i]
				sortedPoints[i] = sortedPoints[j]
				sortedPoints[j] = temp
			}
		}
	}
	return sortedTriangle{sortedPoints}
}

func (t sortedTriangle) PointInTriangle(p image.Point) bool {
	dx0 := t.Points[1].X - t.Points[0].X
	dy0 := t.Points[1].Y - t.Points[0].Y
	dx1 := t.Points[2].X - t.Points[1].X
	dy1 := t.Points[2].Y - t.Points[1].Y
	dx2 := t.Points[0].X - t.Points[2].X
	dy2 := t.Points[0].Y - t.Points[2].Y

	s0 := float64(dx0*(p.Y-t.Points[0].Y) - dy0*(p.X-t.Points[0].X))
	s1 := float64(dx1*(p.Y-t.Points[1].Y) - dy1*(p.X-t.Points[1].X))
	s2 := float64(dx2*(p.Y-t.Points[2].Y) - dy2*(p.X-t.Points[2].X))

	return (s0 > 0 && s1 > 0 && s2 > 0) || (s0 < 0 && s1 < 0 && s2 < 0)
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
	FloodTriFill(texture, triangle.Points[0].X-5, triangle.Points[0].Y+1, clr, triangle)
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

func FloodTriFill(texture screen.Texture, x, y int, clr color.Color, triangle sortedTriangle) {

	queue := [][2]int{{x, y}}
	for len(queue) > 0 {
		point := queue[0]
		queue = queue[1:]
		x, y := point[0], point[1]

		if x < 0 || x >= texture.Size().X || y < 0 || y >= texture.Size().Y {
			continue
		}
		if !triangle.PointInTriangle(image.Point{x, y}) {
			continue
		}

		texture.Fill(image.Rect(x, y, x+1, y+1), clr, screen.Over)

		queue = append(queue, [2]int{x + 1, y})
		queue = append(queue, [2]int{x - 1, y})
		queue = append(queue, [2]int{x, y + 1})
		queue = append(queue, [2]int{x, y - 1})
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
