package goboy

import (
	"fmt"
)

type Keystate struct {
	Up     bool
	Down   bool
	Left   bool
	Right  bool
	A      bool
	B      bool
	Start  bool
	Select bool
}

type Joypad struct {
	state         Keystate
	mmu           *MMU
	buttonKeys    bool
	directionKeys bool
}

// Bit 5 - P15 Select Button Keys      (0=Select)
// Bit 4 - P14 Select Direction Keys   (0=Select)
// Bit 3 - P13 Input Down  or Start    (0=Pressed) (Read Only)
// Bit 2 - P12 Input Up    or Select   (0=Pressed) (Read Only)
// Bit 1 - P11 Input Left  or Button B (0=Pressed) (Read Only)
// Bit 0 - P10 Input Right or Button A (0=Pressed) (Read Only)

func (j *Joypad) Update(state Keystate) {
	ifReg := j.mmu.registers[AddrIF]
	var requestInt bool
	// if j.buttonKeys {
	requestInt = requestInt || (!j.state.Select && state.Select)
	requestInt = requestInt || (!j.state.Start && state.Start)
	requestInt = requestInt || (!j.state.B && state.B)
	requestInt = requestInt || (!j.state.A && state.A)
	// }
	// if j.directionKeys {
	requestInt = requestInt || (!j.state.Up && state.Up)
	requestInt = requestInt || (!j.state.Down && state.Down)
	requestInt = requestInt || (!j.state.Left && state.Left)
	requestInt = requestInt || (!j.state.Right && state.Right)
	// }
	j.state = state
	if requestInt {
		fmt.Println("Request int")
		ifReg.RawSet(setBit(ifReg.Get(), JoypadInt))
	}
}

func (j *Joypad) Set(data uint8) {
	j.buttonKeys = data&(1<<5) == 0
	j.directionKeys = data&(1<<4) == 0
}

func (j *Joypad) Get() uint8 {
	var (
		buttonKey    uint8 = 1
		directionKey uint8 = 1
	)
	if j.buttonKeys {
		buttonKey = 0
	}
	if j.directionKeys {
		directionKey = 0
	}
	var data uint8 = (buttonKey << 5) | (directionKey << 4) | 0xF
	if j.buttonKeys {
		if j.state.Start {
			mask := ^uint8(1 << 3)
			data &= mask
		}
		if j.state.Select {
			mask := ^uint8(1 << 2)
			data &= mask
		}
		if j.state.B {
			mask := ^uint8(1 << 1)
			data &= mask
		}
		if j.state.A {
			mask := ^uint8(1 << 0)
			data &= mask
		}
	}
	if j.directionKeys {
		if j.state.Down {
			mask := ^uint8(1 << 3)
			data &= mask
		}
		if j.state.Up {
			mask := ^uint8(1 << 2)
			data &= mask
		}
		if j.state.Left {
			mask := ^uint8(1 << 1)
			data &= mask
		}
		if j.state.Right {
			mask := ^uint8(1 << 0)
			data &= mask
		}
	}
	return data
}

func (j *Joypad) RawSet(data uint8) {}
