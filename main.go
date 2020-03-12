package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MatiasLyyra/goboy/goboy"
)

func main() {
	f, err := os.Open("./cpu_instrs/cpu_instrs/cpu_instrs.gb")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	rom, err := goboy.LoadROM(f)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rom)
}
