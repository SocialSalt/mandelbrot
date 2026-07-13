package userinterface

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Location struct {
	X int32
	Y int32
}

type MouseState struct {
	IsDragging      bool
	StartLocation   Location
	CurrentLocation Location
	Scroll          float32
}

type ViewBoundaries struct {
	MinReal float64
	MaxReal float64
	MinImag float64
	MaxImag float64
}
type WindowShape struct {
	W int32
	H int32
}

type State struct {
	MouseState MouseState
	Bounds     ViewBoundaries
	NextBounds ViewBoundaries
	Window     WindowShape
}

func ComputeDrag(start Location, end Location, width int32, height int32, bounds ViewBoundaries) ViewBoundaries {
	xPx := start.X - end.X
	yPx := start.Y - end.Y

	xMovement := float64(xPx) / float64(width)
	yMovement := float64(yPx) / float64(height)

	xDist := (bounds.MaxReal - bounds.MinReal) * xMovement
	yDist := (bounds.MaxImag - bounds.MinImag) * yMovement
	// fmt.Printf("xPx: %v, yPx %v\n", xPx, yPx)
	// fmt.Printf("xMovement: %v, yMovement %v\n", xMovement, yMovement)
	// fmt.Printf("xDist: %v, yDist %v\n", xDist, yDist)

	return ViewBoundaries{
		MinReal: bounds.MinReal + xDist,
		MaxReal: bounds.MaxImag + xDist,
		MinImag: bounds.MinImag + yDist,
		MaxImag: bounds.MaxImag + yDist,
	}
}

func PollInput(state *State) bool {
	event := sdl.PollEvent()
	switch t := event.(type) {
	case *sdl.QuitEvent:
		return false
	case *sdl.MouseButtonEvent:
		if t.Button == sdl.BUTTON_LEFT && t.Type == sdl.MOUSEBUTTONDOWN {
			state.MouseState.IsDragging = true
			state.MouseState.StartLocation.X = t.X
			state.MouseState.StartLocation.Y = t.Y
			state.MouseState.CurrentLocation.X = t.X
			state.MouseState.CurrentLocation.Y = t.Y
		}
		if t.Button == sdl.BUTTON_LEFT && t.Type == sdl.MOUSEBUTTONUP {
			state.MouseState.IsDragging = false
			state.Bounds = state.NextBounds
		}
	case *sdl.MouseMotionEvent:
		if state.MouseState.IsDragging {
			state.MouseState.CurrentLocation.X = t.X
			state.MouseState.CurrentLocation.Y = t.Y
			state.NextBounds = ComputeDrag(
				state.MouseState.StartLocation,
				state.MouseState.CurrentLocation,
				state.Window.W,
				state.Window.H,
				state.Bounds,
			)
		}
	case *sdl.MouseWheelEvent:
		state.MouseState.Scroll += t.PreciseY
	}
	return true
}
