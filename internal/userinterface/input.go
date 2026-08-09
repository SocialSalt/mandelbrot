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
	RealWidth  float64
	ImagHeight float64
}

type ImageCenter struct {
	Real float64
	Imag float64
}

type WindowShape struct {
	W int32
	H int32
}

type State struct {
	MouseState MouseState
	Center     ImageCenter
	NextCenter ImageCenter
	Bounds     ViewBoundaries
	NextBounds ViewBoundaries
	Window     WindowShape
}

func ComputeDrag(start Location, end Location, width int32, height int32, bounds ViewBoundaries, center ImageCenter) ImageCenter {
	xPx := start.X - end.X
	yPx := start.Y - end.Y

	xMovement := (float64(xPx) / float64(width)) * bounds.RealWidth
	yMovement := (float64(yPx) / float64(height)) * bounds.ImagHeight

	return ImageCenter{
		Real: center.Real + xMovement,
		Imag: center.Imag + yMovement,
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
			// state.Bounds = state.NextBounds
			state.Center = state.NextCenter
		}
	case *sdl.MouseMotionEvent:
		if state.MouseState.IsDragging {
			state.MouseState.CurrentLocation.X = t.X
			state.MouseState.CurrentLocation.Y = t.Y
			state.NextCenter = ComputeDrag(
				state.MouseState.StartLocation,
				state.MouseState.CurrentLocation,
				state.Window.W,
				state.Window.H,
				state.Bounds,
				state.Center,
			)
		}
	case *sdl.MouseWheelEvent:
		switch t.PreciseY {
		case 1:
			state.Bounds.RealWidth *= 0.9
			state.Bounds.ImagHeight *= 0.9
		case -1:
			state.Bounds.RealWidth *= 1.1
			state.Bounds.ImagHeight *= 1.1
		}
		// state.MouseState.Scroll += t.PreciseY
		// fmt.Printf("PreciseY: %v\n", t.PreciseY)
		// fmt.Printf("Scroll: %v\n", state.MouseState.Scroll)

	}
	return true
}
