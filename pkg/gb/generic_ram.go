package gb

import "fmt"

type GenericRAM struct {
	memory [RAM_SIZE]uint8 // generic ram that covers 0xA000 - 0xFFFF
	ly     uint8
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

	RAM_SIZE = 0x6000
	RAM_BASE = 0xA000
	RAM_TOP  = 0xFFFF
)

func newGenericRAM() *GenericRAM {
	r := &GenericRAM{}
	r.init()

	return r
}

func (r *GenericRAM) init() {
	r.memory = [RAM_SIZE]uint8{}
	r.memory[0xFF10-RAM_BASE] = 0x80
	r.memory[0xFF11-RAM_BASE] = 0xBF
	r.memory[0xFF12-RAM_BASE] = 0xF3
	r.memory[0xFF14-RAM_BASE] = 0xBF
	r.memory[0xFF16-RAM_BASE] = 0x3F
	r.memory[0xFF17-RAM_BASE] = 0x00
	r.memory[0xFF19-RAM_BASE] = 0xBF
	r.memory[0xFF1A-RAM_BASE] = 0x7F
	r.memory[0xFF1B-RAM_BASE] = 0xFF
	r.memory[0xFF1C-RAM_BASE] = 0x9F
	r.memory[0xFF1E-RAM_BASE] = 0xBF
	r.memory[0xFF20-RAM_BASE] = 0xFF
	r.memory[0xFF21-RAM_BASE] = 0x00
	r.memory[0xFF22-RAM_BASE] = 0x00
	r.memory[0xFF23-RAM_BASE] = 0xBF
	r.memory[0xFF24-RAM_BASE] = 0x77
	r.memory[0xFF25-RAM_BASE] = 0xF3
	r.memory[0xFF26-RAM_BASE] = 0xF1
	r.memory[0xFF40-RAM_BASE] = 0x91
	r.memory[0xFF42-RAM_BASE] = 0x00
	r.memory[0xFF43-RAM_BASE] = 0x00
	r.memory[0xFF45-RAM_BASE] = 0x00
	r.memory[0xFF47-RAM_BASE] = 0xFC
	r.memory[0xFF48-RAM_BASE] = 0xFF
	r.memory[0xFF49-RAM_BASE] = 0xFF
	r.memory[0xFF4A-RAM_BASE] = 0x00
	r.memory[0xFF4B-RAM_BASE] = 0x00
	r.memory[IE_ADDR-RAM_BASE] = 0x00
}

func (r *GenericRAM) contains(addr uint16) bool {
	return inRange(addr, RAM_BASE, RAM_TOP)
}

func (r *GenericRAM) read(addr uint16) uint8 {
	switch addr {
	case 0xFF44:
		res := r.ly
		r.ly++
		return res
	case IF_ADDR:
		return r.memory[addr-RAM_BASE] & INTRUPT_MSK
	case IE_ADDR:
		return r.memory[addr-RAM_BASE] & INTRUPT_MSK
	default:
		return r.memory[addr-RAM_BASE]
	}
}

func (r *GenericRAM) write(addr uint16, data uint8) {
	if addr == 0xFF01 {
		// fmt.Printf("Write at: [0x%04x] = 0x%02x\n", address, data)
		fmt.Printf("%c", data)
	}

	switch addr {
	case IF_ADDR:
		r.memory[addr-RAM_BASE] = data & INTRUPT_MSK
	case IE_ADDR:
		r.memory[addr-RAM_BASE] = data & INTRUPT_MSK
	default:
		r.memory[addr-RAM_BASE] = data
	}
}
