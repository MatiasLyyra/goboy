package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MatiasLyyra/goboy/goboy"
	"github.com/MatiasLyyra/goboy/gui"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	// f, err := os.Open("/home/malyy/src/gb-test-roms/cpu_instrs/individual/02-interrupts.gb")
	f, err := os.Open("/home/malyy/src/gb-test-roms/cpu_instrs/cpu_instrs.gb")
	// f, err := os.Open("/home/malyy/roms/tetris.gb")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	rom, err := goboy.LoadCartridge(f)
	if err != nil {
		log.Fatalln(err)
	}
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	fmt.Println(rom)
	defer sdl.Quit()
	w, err := gui.NewWindow("Goyboy", 4)
	defer w.Close()
	if err != nil {
		panic(err)
	}
	running := true
	// sink := make(chan [160 * 144]uint8, 0)
	// go func() {
		mmu := goboy.NewMMU(rom)
		cpu := goboy.CPU{
			Memory: mmu,
			PC:     0x0100,
			SP:     0xFFFE,
		}
		cpu.SetF(0x80)
	// debugger := debug.Debugger{
	// 	CPU: &cpu,
	// 	Breakpoints: map[uint16]struct{}{
	// 		0x0100: struct{}{},
	// 	},
	// }
	// debug.StartDebugger(debugger, sink)
		const cyclesInSecond = 4213440
		// var cycleCount uint64
	// for {

	// time.Sleep(time.Duration(cycles) * time.Nanosecond)

			// cycleCount += uint64(cycles)
			// if cycleCount > cyclesInSecond*20 {
			// 	fmt.Println("Written")
			// 	ioutil.WriteFile("./vram.bin", mmu.GPU.VRAM[:], 0755)
			// 	os.Exit(0)
			// }
	// }
	// }()
	var keys goboy.Keystate
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:
				switch e.Type {
				case sdl.KEYDOWN, sdl.KEYUP:
					updateKeystate(&keys, e)
				}
			}
		}
		var drawCount int
		for drawCount < 10 {
			mmu.Pad.Update(keys)
			cycles := cpu.RunSingleOpcode()
			hasDrawn := mmu.GPU.Run(cycles)
			if hasDrawn {
				drawCount++
			}
		}
		w.Draw(mmu.GPU.ScreenBuffer())
		// <-sink
	}
}

func updateKeystate(keyState *goboy.Keystate, keyEvent *sdl.KeyboardEvent) {
	var state bool
	if keyEvent.Type == sdl.KEYDOWN {
		state = true
	}
	switch keyEvent.Keysym.Sym {
	case sdl.K_UP:
		keyState.Up = state
	case sdl.K_DOWN:
		keyState.Down = state
	case sdl.K_LEFT:
		keyState.Left = state
	case sdl.K_RIGHT:
		keyState.Right = state
	case 'z':
		keyState.B = state
	case 'x':
		keyState.A = state
	case sdl.K_RETURN:
		keyState.Start = state
	case 'a':
		keyState.Select = state
	default:
		fmt.Println(keyEvent.Keysym.Sym)
	}
}
