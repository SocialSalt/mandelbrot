package main

import (
	"fmt"
	"log"
	"time"

	userinterface "github.com/socialsalt/mandelbrot/internal"
	"github.com/veandco/go-sdl2/sdl"
)

// cpu stuff

/*
#include "../lib/include/cpumandelbrot.h"
*/

// cuda stuff

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
		Bounds: userinterface.ViewBoundaries{
			MinReal: -2.0,
			MaxReal: 1.0,
			MinImag: -0.85,
			MaxImag: 0.8375,
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
			if time.Since(lastDraw).Milliseconds() > 50 {
				w, h := window.GetSize()
				state.Window.W = w
				state.Window.H = h
				if state.MouseState.IsDragging {
					fmt.Printf("width: %#v\n", state.NextBounds.MaxReal-state.NextBounds.MinReal)
					fmt.Printf("height: %#v\n", state.NextBounds.MaxImag-state.NextBounds.MinImag)
				}
				// s := time.Now()
				C.LaunchMandelbrot(
					(*C.uchar)(&buffer[0]),
					C.double(state.NextBounds.MinReal),
					C.double(state.NextBounds.MaxReal),
					C.double(state.NextBounds.MinImag),
					C.double(state.NextBounds.MaxImag),
					C.int(width),
					C.int(height),
					C.int(channels),
				)
				// fmt.Printf("frame time was %v\n", time.Now().Sub(s).Milliseconds())
				userinterface.UpdateDisplay(buffer, renderer, texture, width, height, 3*width)
				lastDraw = time.Now()
			}
			sdl.Delay(1)
		}
	}
}
