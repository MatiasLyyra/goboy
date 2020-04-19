package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Available colors
const (
	// 0xFFFFFF
	ColorWhite = iota
	// 0xAAAAAA
	ColorLightGray
	// 0x555555
	ColorDarkGray
	// 0x000000
	ColorBlack
)

type Window struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

func NewWindow(title string, scale int) (*Window, error) {
	guiWindow := &Window{}
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		160*int32(scale), 144*int32(scale), sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}
	guiWindow.window = window
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}
	guiWindow.renderer = renderer
	renderer.SetScale(float32(scale), float32(scale))
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.Present()

	return guiWindow, nil
}

func (w *Window) Close() {
	if w.renderer != nil {
		w.renderer.Destroy()
	}
	if w.window != nil {
		w.window.Destroy()
	}
}

func (w *Window) Draw(buffer []uint8) {
	for i, val := range buffer {
		y := i / 160
		x := i % 160
		switch val {
		case ColorBlack:
			w.renderer.SetDrawColor(0x0f, 0x38, 0x0f, 255)
		case ColorDarkGray:
			w.renderer.SetDrawColor(0x30, 0x62, 0x30, 255)
		case ColorLightGray:
			w.renderer.SetDrawColor(0x8b, 0xac, 0x0f, 255)
		case ColorWhite:
			w.renderer.SetDrawColor(0x9b, 0xbc, 0x0f, 255)
		}
		w.renderer.DrawPoint(int32(x), int32(y))
	}
	w.renderer.Present()
}
