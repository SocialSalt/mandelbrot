package main

import (
	"log"

	userinterface "github.com/socialsalt/mandelbrot/internal"
	_ "github.com/socialsalt/mandelbrot/lib"
	"github.com/veandco/go-sdl2/sdl"
)

// cuda stuff

/*
#cgo LDFLAGS: -L${SRCDIR}/.. -lmandelbrot
#cgo LDFLAGS: -L/usr/local/cuda/lib64 -lcudart
#include "../lib/include/mandelbrot.h"
*/

// cpu stuff

/*
// #cgo CFLAGS: -I${SRCDIR}/../lib/include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lcpumandelbrot
#include "../lib/include/cpumandelbrot.h"
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

Outer:
	for {
		switch userinterface.PollInput() {
		case false:
			break Outer
		default:
			C.CpuMandelbrot(
				(*C.uchar)(&buffer[0]),
				C.double(real_min),
				C.double(real_max),
				C.double(imag_min),
				C.double(imag_max),
				C.int(width),
				C.int(height),
				C.int(channels),
			)
			userinterface.UpdateDisplay(buffer, renderer, texture, width, height, 3*width)
		}
	}
}
