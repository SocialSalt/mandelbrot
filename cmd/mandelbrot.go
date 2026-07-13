package main

import (
	"fmt"
	"log"
	"time"

	userinterface "github.com/socialsalt/mandelbrot/internal"
	_ "github.com/socialsalt/mandelbrot/lib"
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
	real_min := -2.0
	real_max := 1.0
	imag_min := -0.85
	imag_max := 0.8375

	lastDraw := time.Now()

Outer:
	for {
		switch userinterface.PollInput() {
		case false:
			break Outer
		default:
			if time.Since(lastDraw).Milliseconds() > 50 {
				s := time.Now()
				// C.CpuMandelbrot(
				C.LaunchMandelbrot(
					(*C.uchar)(&buffer[0]),
					C.double(real_min),
					C.double(real_max),
					C.double(imag_min),
					C.double(imag_max),
					C.int(width),
					C.int(height),
					C.int(channels),
				)
				fmt.Printf("frame time was %v\n", time.Now().Sub(s).Microseconds())
				userinterface.UpdateDisplay(buffer, renderer, texture, width, height, 3*width)
				lastDraw = time.Now()
			}
			sdl.Delay(1)
		}
	}
}
