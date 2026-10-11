package main

import (
	"log"
	"time"

	"github.com/socialsalt/mandelbrot/internal/userinterface"
	"github.com/veandco/go-sdl2/sdl"
)

/*
#cgo LDFLAGS: -L/usr/local/cuda/lib64 -lcudart
#cgo LDFLAGS: -L${SRCDIR}/../lib -Wl,-rpath,${SRCDIR}/../lib -lmandelbrot
#include "../lib/include/mandelbrot.h"
#include "../lib/include/julia.h"
*/
import "C"

func main() {
	width := 1920
	height := 1080
	window, renderer, texture, err := userinterface.CreateSDL(width, height)
	if err != nil {
		log.Fatalf("failed to create sdl with error %v", err)
	}
	defer sdl.Quit()
	defer window.Destroy()

	font, err := userinterface.CreateFont(renderer)
	if err != nil {
		log.Fatalf("failed to create sdl with error %v", err)
	}
	defer font.Close()
	button, surface, textTexture, err := userinterface.NewButton(
		font,
		renderer,
		sdl.Rect{X: 10, Y: 10, W: 100, H: 50},
		sdl.Color{R: 50, G: 50, B: 50, A: 255},
		sdl.Color{R: 100, G: 100, B: 100, A: 255},
		"Julia",
	)
	if err != nil {
		log.Fatalf("failed to create sdl with error %v", err)
	}
	defer surface.Free()
	defer textTexture.Destroy()

	minimap := userinterface.NewMinimap(sdl.Rect{X: 10, Y: 60, W: 200, H: 200})

	channels := 3
	buffer := make([]byte, width*height*channels)

	w, h := window.GetSize()
	state := userinterface.State{
		MouseState: userinterface.MouseState{},
		Center: userinterface.ImageCenter{
			Real: 0.5,
			Imag: 0.0,
		},
		NextCenter: userinterface.ImageCenter{
			Real: 0.5,
			Imag: 0.0,
		},
		Bounds: userinterface.ViewBoundaries{
			RealWidth:  3.0,
			ImagHeight: 1.6875,
		},
		Window: userinterface.WindowShape{W: w, H: h},
	}
	state.NextBounds = state.Bounds

	lastDraw := time.Now()

Outer:
	for {
		switch userinterface.PollInput(&state, button, minimap) {
		case false:
			break Outer
		default:
			if time.Since(lastDraw).Milliseconds() > 7 {
				w, h := window.GetSize()
				state.Window.W = w
				state.Window.H = h
				// s := time.Now()

				if button.IsPressed {
					C.ComputeJulia(
						(*C.uchar)(&buffer[0]),
						C.int(width),
						C.int(height),
						C.double(state.NextCenter.Real),
						C.double(state.NextCenter.Imag),
						C.double(state.Bounds.RealWidth),
						C.double(state.Bounds.ImagHeight),
						C.int(channels),
						C.double(minimap.Point.Real),
						C.double(-1*minimap.Point.Imag),
					)
				} else {
					C.ComputeMandelbrot(
						(*C.uchar)(&buffer[0]),
						C.int(width),
						C.int(height),
						C.double(state.NextCenter.Real),
						C.double(state.NextCenter.Imag),
						C.double(state.Bounds.RealWidth),
						C.double(state.Bounds.ImagHeight),
						C.int(channels),
					)

				}

				// fmt.Printf("time to compute data was %v\n", time.Since(s).Milliseconds())
				userinterface.UpdateDisplay(buffer, renderer, texture, textTexture, button, minimap, width, height, 3*width)
				lastDraw = time.Now()
			}
			sdl.Delay(1)
		}
	}
}
