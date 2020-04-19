package debug

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/MatiasLyyra/goboy/goboy"
)

func StartDebugger(d Debugger, sink chan<- [160 * 144]uint8) {
	scan := bufio.NewReader(os.Stdin)
	for {
		decoded, lookup := d.DecodeROM()
		var options []string
		fmt.Print("> ")
		option, _ := scan.ReadString('\n')
		option = strings.TrimSpace(option)
		options = strings.Split(option, " ")
		if len(options) > 0 {
			option = options[0]
			options = options[1:]
		}
		switch option {
		case "v", "view":
			printSnippet(d, decoded, lookup)
		case "s", "step":
			runSingleStep(d, sink)
			printSnippet(d, decoded, lookup)
		case "r", "run":
			d.StepToNextBreakpoint(sink)
			decoded, lookup := d.DecodeROM()
			printSnippet(d, decoded, lookup)
		case "b", "break":
			if len(options) > 0 {
				val, err := strconv.ParseInt(options[0], 16, 64)
				if err != nil && val < (1<<16) {
					fmt.Println("invalid value")
				} else {
					d.ToggleBreakpoint(uint16(val))
				}
			} else {
				d.ToggleBreakpoint(d.CPU.PC)
			}
			printSnippet(d, decoded, lookup)
		case "read":
			if len(options) == 0 {
				continue
			}
			val, err := strconv.ParseInt(options[0], 16, 64)
			if err != nil || val >= (1<<16) {
				fmt.Println("invalid value")
			} else {
				fmt.Printf("%02X\n", d.CPU.Memory.Read(uint16(val)))
			}
		case "write":
			if len(options) < 2 {
				continue
			}
			addr, err := strconv.ParseInt(options[0], 16, 64)
			if err != nil || addr >= (1<<16) {
				fmt.Println("invalid addr")
				continue
			}
			val, err := strconv.ParseInt(options[1], 16, 64)
			if err != nil || val >= (1<<8) {
				fmt.Println("invalid value")
				continue
			}
			d.CPU.Memory.Write(uint16(addr), uint8(val))
		case "mbc1":
			mbc1, ok := d.CPU.Memory.Cartridge.MBC.(*goboy.MBC1)
			if len(options) == 0 || !ok {
				if !ok {
					fmt.Println("Not MBC1")
				}
				continue
			}
			switch options[0] {
			case "rom":
				fmt.Printf("Selected rom bank: %d\n", mbc1.SelectedROM())
			case "ram":
				fmt.Printf("Selected ram bank: %d\n", mbc1.SelectedRAM())
			}
		case "quit":
			os.Exit(0)
		}
	}
}

func runSingleStep(d Debugger, sink chan<- [160 * 144]uint8) {
	cycles := d.CPU.RunSingleOpcode()
	d.CPU.Memory.GPU.Run(cycles, sink)
}

func printSnippet(d Debugger, decoded []DecodedInsturction, lookup map[uint16]int) {
	startAddr := d.CPU.PC
	ops := []DecodedInsturction{
		decoded[lookup[d.CPU.PC]],
	}
outer:
	for i := 0; i < 5; i++ {
		for {
			startAddr--
			if startAddr <= 0 {
				break outer
			}
			if ind, found := lookup[startAddr]; found {

				ops = append([]DecodedInsturction{decoded[ind]}, ops...)
				break
			}
		}
	}
	currentOP := decoded[lookup[d.CPU.PC]]
	for i := 0; i < 5; i++ {
		addr := currentOP.Addr + uint16(currentOP.Len)
		op := decoded[lookup[addr]]
		currentOP = op
		ops = append(ops, op)
	}
	var reg int
	for _, op := range ops {
		if op.Addr == d.CPU.PC {
			fmt.Print("> ")
		} else if _, found := d.Breakpoints[op.Addr]; found {
			fmt.Print("* ")
		} else {
			fmt.Print("  ")
		}
		fmt.Printf("%04X: %v", op.Addr, op)
		switch reg {
		case 0:
			fmt.Printf("\t\tAF: $%04X", d.CPU.AF())
		case 1:
			fmt.Printf("\t\tBC: $%04X", d.CPU.BC())
		case 2:
			fmt.Printf("\t\tDE: $%04X", d.CPU.DE())
		case 3:
			fmt.Printf("\t\tHL: $%04X", d.CPU.HL())
		case 4:
			fmt.Printf("\t\tSP: $%04X", d.CPU.SP)
		case 5:
			fmt.Printf("\t\tPC: $%04X", d.CPU.PC)
		}
		reg++
		fmt.Println()
	}
}
