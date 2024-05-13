package gb

import "fmt"

type GenericRAM struct {
	memory [0x8000]uint8 // generic ram that covers 0x8000 - 0xFFFF
}

const (
	IF_ADDR     = 0xFF0F
	IE_ADDR     = 0xFFFF
	INTRUPT_MSK = 0x1F

	VBLANK_INTRUPT_BIT = 0
	LCD_INTRUPT_BIT    = 1
	TIMER_INTRUPT_BIT  = 2
	SERIAL_INTRUPT_BIT = 3
	JOYPAD_INTRUPT_BIT = 4

	VBLANK_INTRUPT_VEC = 0x40
	STAT_INTRUPT_VEC   = 0x48
	TIMER_INTRUPT_VEC  = 0x50
	SERIAL_INTRUPT_VEC = 0x58
	JOYPAD_INTRUPT_VEC = 0x60
)

func newGenericRAM() *GenericRAM {
	r := &GenericRAM{}
	r.init()

	return r
}

func (r *GenericRAM) init() {
	r.memory = [0x8000]uint8{}
	r.memory[0xFF10-0x8000] = 0x80
	r.memory[0xFF11-0x8000] = 0xBF
	r.memory[0xFF12-0x8000] = 0xF3
	r.memory[0xFF14-0x8000] = 0xBF
	r.memory[0xFF16-0x8000] = 0x3F
	r.memory[0xFF17-0x8000] = 0x00
	r.memory[0xFF19-0x8000] = 0xBF
	r.memory[0xFF1A-0x8000] = 0x7F
	r.memory[0xFF1B-0x8000] = 0xFF
	r.memory[0xFF1C-0x8000] = 0x9F
	r.memory[0xFF1E-0x8000] = 0xBF
	r.memory[0xFF20-0x8000] = 0xFF
	r.memory[0xFF21-0x8000] = 0x00
	r.memory[0xFF22-0x8000] = 0x00
	r.memory[0xFF23-0x8000] = 0xBF
	r.memory[0xFF24-0x8000] = 0x77
	r.memory[0xFF25-0x8000] = 0xF3
	r.memory[0xFF26-0x8000] = 0xF1
	r.memory[0xFF40-0x8000] = 0x91
	r.memory[0xFF42-0x8000] = 0x00
	r.memory[0xFF43-0x8000] = 0x00
	r.memory[0xFF45-0x8000] = 0x00
	r.memory[0xFF47-0x8000] = 0xFC
	r.memory[0xFF48-0x8000] = 0xFF
	r.memory[0xFF49-0x8000] = 0xFF
	r.memory[0xFF4A-0x8000] = 0x00
	r.memory[0xFF4B-0x8000] = 0x00
	r.memory[IE_ADDR-0x8000] = 0x00
}

func (r *GenericRAM) contains(addr uint16) bool {
	return addr >= 8000 && addr <= 0xFFFF
}

func (r *GenericRAM) read(addr uint16) uint8 {
	switch addr {
	case IF_ADDR:
		return r.memory[addr-0x8000] & INTRUPT_MSK
	case IE_ADDR:
		return r.memory[addr-0x8000] & INTRUPT_MSK
	default:
		return r.memory[addr-0x8000]
	}
}

func (r *GenericRAM) write(addr uint16, data uint8) {
	if addr == 0xFF01 {
		// fmt.Printf("Write at: [0x%04x] = 0x%02x\n", address, data)
		fmt.Printf("%c", data)
	}

	switch addr {
	case IF_ADDR:
		r.memory[addr-0x8000] = data & INTRUPT_MSK
	case IE_ADDR:
		r.memory[addr-0x8000] = data & INTRUPT_MSK
	default:
		r.memory[addr-0x8000] = data
	}
}
