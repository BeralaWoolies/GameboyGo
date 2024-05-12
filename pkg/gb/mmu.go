package gb

import (
	"fmt"
)

type MMU struct {
	memory    [0x10000]uint8
	DIVTicks  int
	TIMATicks int
}

const (
	DIV_ADDR  = 0xFF04
	TIMA_ADDR = 0xFF05
	TMA_ADDR  = 0xFF06
	TAC_ADDR  = 0xFF07

	TAC_TIMER_ENABLE_BIT = 2
	TAC_FREQ_DIV_MSK     = 0x3
	TAC_MSK              = 0x7

	HZ_4096   = 0
	HZ_262144 = 1
	HZ_65536  = 2
	HZ_16386  = 3

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

func (mmu *MMU) init() {
	mmu.memory = [0x10000]uint8{}
	mmu.memory[DIV_ADDR] = 0xAB
	mmu.memory[TIMA_ADDR] = 0x0
	mmu.memory[TMA_ADDR] = 0x0
	mmu.memory[TAC_ADDR] = 0x0
	mmu.memory[0xFF10] = 0x80
	mmu.memory[0xFF11] = 0xBF
	mmu.memory[0xFF12] = 0xF3
	mmu.memory[0xFF14] = 0xBF
	mmu.memory[0xFF16] = 0x3F
	mmu.memory[0xFF17] = 0x00
	mmu.memory[0xFF19] = 0xBF
	mmu.memory[0xFF1A] = 0x7F
	mmu.memory[0xFF1B] = 0xFF
	mmu.memory[0xFF1C] = 0x9F
	mmu.memory[0xFF1E] = 0xBF
	mmu.memory[0xFF20] = 0xFF
	mmu.memory[0xFF21] = 0x00
	mmu.memory[0xFF22] = 0x00
	mmu.memory[0xFF23] = 0xBF
	mmu.memory[0xFF24] = 0x77
	mmu.memory[0xFF25] = 0xF3
	mmu.memory[0xFF26] = 0xF1
	mmu.memory[0xFF40] = 0x91
	mmu.memory[0xFF42] = 0x00
	mmu.memory[0xFF43] = 0x00
	mmu.memory[0xFF45] = 0x00
	mmu.memory[0xFF47] = 0xFC
	mmu.memory[0xFF48] = 0xFF
	mmu.memory[0xFF49] = 0xFF
	mmu.memory[0xFF4A] = 0x00
	mmu.memory[0xFF4B] = 0x00
	mmu.memory[IE_ADDR] = 0x00

	mmu.DIVTicks = 0
	mmu.TIMATicks = 0
}

func (mmu *MMU) mapRom(rom []byte) {
	copy(mmu.memory[:], rom)
}

func (mmu *MMU) read(address uint16) uint8 {
	switch address {
	case TAC_ADDR:
		return mmu.memory[address] & TAC_MSK
	case IF_ADDR:
		return mmu.memory[address] & INTRUPT_MSK
	case IE_ADDR:
		return mmu.memory[address] & INTRUPT_MSK
	default:
		return mmu.memory[address]
	}
}

func (mmu *MMU) write(address uint16, data uint8) {
	if address == 0xFF01 {
		// fmt.Printf("Write at: [0x%04x] = 0x%02x\n", address, data)
		fmt.Printf("%c", data)
	}

	switch address {
	case DIV_ADDR:
		// not allowed to write to Divider Register
		mmu.memory[address] = 0
	case TAC_ADDR:
		mmu.memory[address] = data & TAC_MSK
	case IF_ADDR:
		mmu.memory[address] = data & INTRUPT_MSK
	case IE_ADDR:
		mmu.memory[address] = data & INTRUPT_MSK
	default:
		mmu.memory[address] = data
	}
}

func (mmu *MMU) incDIV() {
	// increment DIV without resetting it to 0 in write()
	mmu.memory[DIV_ADDR]++
}
