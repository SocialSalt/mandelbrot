package userinterface

import (
	"fmt"
	"unsafe"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

func CreateSDL(width int, height int) (*sdl.Window, *sdl.Renderer, *sdl.Texture, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to init sdl")
	}
	window, err := sdl.CreateWindow(
		"Mandelbrot",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(width),
		int32(height),
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE,
	)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to create window")
	}

	renderer, err := sdl.CreateRenderer(
		window,
		-1,
		sdl.RENDERER_ACCELERATED|sdl.RENDERER_TARGETTEXTURE,
	)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to get renderer")
	}
	if renderer == nil {
		return nil, nil, nil, fmt.Errorf("renderer was nil")
	}

	texture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_RGB24,
		sdl.TEXTUREACCESS_STREAMING,
		int32(width),
		int32(height),
	)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to create texture")
	}
	return window, renderer, texture, nil
}

func UpdateDisplay(buffer []byte, renderer *sdl.Renderer, texture *sdl.Texture, width int, height int, pitch int) error {
	rect := sdl.Rect{
		X: 0,
		Y: 0,
		W: int32(width),
		H: int32(height),
	}
	ptr := unsafe.SliceData(buffer)
	unsafeptr := unsafe.Pointer(ptr)
	if err := texture.Update(&rect, unsafeptr, pitch); err != nil {
		return errors.Wrap(err, "failed to update texture")
	}
	if err := renderer.Clear(); err != nil {
		return errors.Wrap(err, "failed to clear the renderer")
	}

	if err := renderer.Copy(texture, nil, nil); err != nil {
		return errors.Wrap(err, "failed to copy the renderer")
	}
	renderer.Present()
	return nil
}
