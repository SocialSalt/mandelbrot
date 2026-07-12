package userinterface

import "github.com/veandco/go-sdl2/sdl"

func PollInput() bool {
	event := sdl.PollEvent()
	switch event.(type) {
	case *sdl.QuitEvent:
		return false
	default:
		return true
	}
}
