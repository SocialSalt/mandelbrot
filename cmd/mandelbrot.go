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
		switch userinterface.PollInput(&state) {
		case false:
			break Outer
		default:
			if time.Since(lastDraw).Milliseconds() > 7 {
				w, h := window.GetSize()
				state.Window.W = w
				state.Window.H = h
				// s := time.Now()

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

				// fmt.Printf("time to compute data was %v\n", time.Since(s).Milliseconds())
				userinterface.UpdateDisplay(buffer, renderer, texture, width, height, 3*width)
				lastDraw = time.Now()
			}
			sdl.Delay(1)
		}
	}
}
