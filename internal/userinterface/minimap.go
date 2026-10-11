package userinterface

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Point struct {
	Real float64
	Imag float64
}
type Minimap struct {
	Rect       sdl.Rect
	Color      sdl.Color
	PointColor sdl.Color
	Point      Point
	MouseState MouseState
	Min        float64
	Max        float64
}

func NewMinimap(location sdl.Rect) *Minimap {
	return &Minimap{
		Rect:       location,
		Color:      sdl.Color{R: 100, G: 100, B: 100, A: 0},
		PointColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
		Point:      Point{Real: 0, Imag: 0},
		Min:        -1,
		Max:        1,
	}
}

func (m *Minimap) HandleEvent(event sdl.Event) bool {
	switch t := event.(type) {
	case *sdl.MouseButtonEvent:
		mouseLoc := sdl.Point{X: t.X, Y: t.Y}
		if !mouseLoc.InRect(&m.Rect) {
			return false
		}
		if t.Button == sdl.BUTTON_LEFT && t.Type == sdl.MOUSEBUTTONDOWN {
			m.MouseState.IsDragging = true
			m.MouseState.StartLocation.X = t.X
			m.MouseState.StartLocation.Y = t.Y
			m.MouseState.CurrentLocation.X = t.X
			m.MouseState.CurrentLocation.Y = t.Y
		}
		if t.Button == sdl.BUTTON_LEFT && t.Type == sdl.MOUSEBUTTONUP {
			m.MouseState.IsDragging = false
			// m.Bounds = state.NextBounds
		}
		m.Point = m.LocationToPoint(m.MouseState.CurrentLocation)
		return true
	case *sdl.MouseMotionEvent:
		if m.MouseState.IsDragging {
			m.MouseState.CurrentLocation.X = t.X
			m.MouseState.CurrentLocation.Y = t.Y
			m.Point = m.LocationToPoint(m.MouseState.CurrentLocation)
			return true
		} else {
			return false
		}
	}
	return false
}

func (m *Minimap) LocationToPoint(l Location) Point {
	pointX := m.Min + float64(l.X-m.Rect.X)/float64(m.Rect.W)*(m.Max-m.Min)
	pointY := m.Min + float64(l.Y-m.Rect.Y)/float64(m.Rect.H)*(m.Max-m.Min)
	return Point{Real: pointX, Imag: pointY}
}

func (m *Minimap) PointToLocation(p Point) Location {
	locX := float64(m.Rect.X) + ((p.Real-m.Min)/(m.Max-m.Min))*float64(m.Rect.W)
	locY := float64(m.Rect.Y) + ((p.Imag-m.Min)/(m.Max-m.Min))*float64(m.Rect.H)
	return Location{X: int32(locX), Y: int32(locY)}
}

func (m *Minimap) Render(renderer *sdl.Renderer) error {
	renderer.SetDrawColor(m.Color.R, m.Color.G, m.Color.B, m.Color.A)
	err := renderer.FillRect(&m.Rect)
	if err != nil {
		return err
	}

	renderer.SetDrawColor(m.PointColor.R, m.PointColor.G, m.PointColor.B, m.PointColor.A)
	loc := m.PointToLocation(m.Point)
	return DrawCircle(renderer, loc.X, loc.Y, 3)
}
