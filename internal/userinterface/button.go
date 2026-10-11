package userinterface

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Button struct {
	Rect       sdl.Rect
	Color      sdl.Color
	BaseColor  sdl.Color
	HoverColor sdl.Color
	ClickColor sdl.Color
	IsPressed  bool
	IsHovered  bool
}

func NewButton(font *ttf.Font, renderer *sdl.Renderer, rect sdl.Rect, color1 sdl.Color, color2 sdl.Color, message string) (*Button, *sdl.Surface, *sdl.Texture, error) {

	textColor := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	textSurface, err := font.RenderUTF8Blended(message, textColor)
	if err != nil {
		return nil, nil, nil, err
	}

	textTexture, err := renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		return nil, nil, nil, err
	}

	return &Button{Rect: rect, Color: color1, BaseColor: color1, ClickColor: color2}, textSurface, textTexture, nil
}

func (b *Button) HandleEvent(event sdl.Event) bool {
	e, ok := event.(*sdl.MouseButtonEvent)
	if !ok {
		return false
	}
	if e.Button != sdl.BUTTON_LEFT {
		return false
	}
	mouseLoc := sdl.Point{X: e.X, Y: e.Y}
	if !mouseLoc.InRect(&b.Rect) {
		return false
	}

	if e.Type != sdl.MOUSEBUTTONUP {
		return true
	}
	b.IsPressed = !b.IsPressed
	if b.IsPressed {
		b.Color = b.ClickColor
	} else {
		b.Color = b.BaseColor
	}
	return true
}

func (b *Button) Render(renderer *sdl.Renderer, texture *sdl.Texture) error {
	renderer.SetDrawColor(b.Color.R, b.Color.G, b.Color.B, b.Color.A)
	err := renderer.FillRect(&b.Rect)
	if err != nil {
		return err
	}
	padding := int32(5)
	txtRect := sdl.Rect{
		X: b.Rect.X + padding,
		Y: b.Rect.Y + padding,
		W: b.Rect.W - (padding * 2),
		H: b.Rect.H - (padding * 2),
	}
	return renderer.Copy(texture, nil, &txtRect)
}
