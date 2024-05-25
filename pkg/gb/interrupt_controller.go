package gb

import (
	"log"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
)

type IntruptController struct {
	mmu              *MMU
	cpu              *CPU
	intruptFlagReg   uint8
	intruptEnableReg uint8
}

const (
	IF_ADDR     = 0xFF0F
	IE_ADDR     = 0xFFFF
	INTRUPT_MSK = 0xE0

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

	ISR_CLOCK_TICKS = 20
)

func (ic *IntruptController) init(mmu *MMU, cpu *CPU) {
	ic.mmu = mmu
	ic.cpu = cpu
}

func (ic *IntruptController) contains(addr uint16) bool {
	return addr == IF_ADDR || addr == IE_ADDR
}

func (ic *IntruptController) read(addr uint16) uint8 {
	switch addr {
	case IF_ADDR:
		return ic.intruptFlagReg | INTRUPT_MSK
	case IE_ADDR:
		return ic.intruptEnableReg | INTRUPT_MSK
	default:
		log.Fatalf("MMU mapped an illegal read address: 0x%02x to interrupt controller", addr)
		return 0xFF
	}
}

func (ic *IntruptController) write(addr uint16, data uint8) {
	switch addr {
	case IF_ADDR:
		ic.intruptFlagReg = data | INTRUPT_MSK
	case IE_ADDR:
		ic.intruptEnableReg = data | INTRUPT_MSK
	default:
		log.Fatalf("MMU mapped an illegal write address: 0x%02x to interrupt controller", addr)
	}
}

func (ic *IntruptController) handleIntrupts() int {
	if ic.cpu.IMEDelay {
		ic.cpu.clearIMEDelay()
		ic.cpu.setIME()
		return 0
	}

	if !ic.cpu.IME && !ic.cpu.halted {
		return 0
	}

	IE := ic.intruptEnableReg
	IF := ic.intruptFlagReg

	if bits.IsSetInBoth(IE, IF, VBLANK_INTRUPT_BIT) {
		ic.serviceIntrupt(VBLANK_INTRUPT_BIT)
	} else if bits.IsSetInBoth(IE, IF, LCD_INTRUPT_BIT) {
		ic.serviceIntrupt(LCD_INTRUPT_BIT)
	} else if bits.IsSetInBoth(IE, IF, TIMER_INTRUPT_BIT) {
		ic.serviceIntrupt(TIMER_INTRUPT_BIT)
	} else if bits.IsSetInBoth(IE, IF, SERIAL_INTRUPT_BIT) {
		ic.serviceIntrupt(SERIAL_INTRUPT_BIT)
	} else if bits.IsSetInBoth(IE, IF, JOYPAD_INTRUPT_BIT) {
		ic.serviceIntrupt(JOYPAD_INTRUPT_BIT)
	} else {
		return 0
	}

	return ISR_CLOCK_TICKS
}

func (ic *IntruptController) serviceIntrupt(intruptBit uint8) {
	if !ic.cpu.IME && ic.cpu.halted {
		ic.cpu.exitHaltedState()
		return
	}

	ic.cpu.exitHaltedState()
	ic.cpu.clearIME()
	ic.clearIFBit(intruptBit)
	ic.cpu.pushStack(ic.cpu.reg.PC)

	switch intruptBit {
	case VBLANK_INTRUPT_BIT:
		ic.cpu.setPC(VBLANK_INTRUPT_VEC)
	case LCD_INTRUPT_BIT:
		ic.cpu.setPC(STAT_INTRUPT_VEC)
	case TIMER_INTRUPT_BIT:
		ic.cpu.setPC(TIMER_INTRUPT_VEC)
	case SERIAL_INTRUPT_BIT:
		ic.cpu.setPC(SERIAL_INTRUPT_VEC)
	case JOYPAD_INTRUPT_BIT:
		ic.cpu.setPC(JOYPAD_INTRUPT_VEC)
	}
}

func (ic *IntruptController) clearIFBit(pos uint8) {
	ic.intruptFlagReg = bits.Reset(ic.intruptFlagReg, pos)
}

func (ic *IntruptController) requestIntrupt(intruptBit uint8) {
	if !inRange(uint16(intruptBit), VBLANK_INTRUPT_BIT, JOYPAD_INTRUPT_BIT) {
		log.Fatalf("Illegal interrupt requested of bit: %d", intruptBit)
	}

	ic.intruptFlagReg = bits.Set(ic.intruptFlagReg, intruptBit)
}
