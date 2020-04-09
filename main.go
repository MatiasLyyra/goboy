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
	sink := make(chan [160 * 144]uint8, 2)
	go func() {
		mmu := goboy.NewMMU(rom)
		cpu := goboy.CPU{
			Memory:   mmu,
			PC:       0x0100,
			AutoExec: true,
		}
		const cyclesInSecond = 4213440
		// var cycleCount uint64
		for {
			cycles := cpu.RunSingleOpcode()
			mmu.GPU.Run(cycles, sink)
			// time.Sleep(time.Duration(cycles) * 237 * time.Nanosecond)
			// cycleCount += uint64(cycles)
			// if cycleCount > cyclesInSecond*20 {
			// 	fmt.Println("Written")
			// 	ioutil.WriteFile("./vram.bin", mmu.GPU.VRAM[:], 0755)
			// 	os.Exit(0)
			// }
		}
	}()
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
		w.Draw(<-sink)
	}
}
