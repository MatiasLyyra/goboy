package goboy

type MemoryRegister interface {
	RawSet(uint8)
	Set(uint8)
	Get() uint8
}

type CallbackRegister struct {
	fn func(data uint8)
}

func (r CallbackRegister) RawSet(uint8)   {}
func (r CallbackRegister) Set(data uint8) { r.fn(data) }
func (r CallbackRegister) Get() uint8     { return 0 }

const (
	Bit0 = 1 << 0
	Bit1 = 1 << 1
	Bit2 = 1 << 2
	Bit3 = 1 << 3
	Bit4 = 1 << 4
	Bit5 = 1 << 5
	Bit6 = 1 << 6
	Bit7 = 1 << 7
)

type RWRegister struct {
	value     uint8
	writeMask uint8
}

func (reg *RWRegister) RawSet(value uint8) {
	reg.value = value
}

func (reg *RWRegister) Set(value uint8) {
	reg.value = value & reg.writeMask
}

func (reg *RWRegister) Get() uint8 {
	return reg.value
}

func NewRWRegister(initialValue uint8, readOnlyBits uint8) *RWRegister {
	reg := &RWRegister{
		value:     initialValue,
		writeMask: ^readOnlyBits,
	}
	return reg
}
