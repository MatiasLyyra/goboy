package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/MatiasLyyra/goboy/goboy"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	f, err := os.Open("./cpu_instrs/cpu_instrs/cpu_instrs.gb")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	rom, err := goboy.LoadCartridge(f)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rom)
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Goboy", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		160*4, 144*4, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	rend, err := window.GetRenderer()
	if err != nil {
		panic(err)
	}
	err = rend.SetLogicalSize(160, 144)
	surface.Set(1, 1, color.White)
	window.UpdateSurface()
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}
