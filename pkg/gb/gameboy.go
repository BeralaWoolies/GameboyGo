package gb

import (
	"fmt"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
)

type Gameboy struct {
	cpu            *CPU
	mmu            *MMU
	timer          *Timer
	cbInstructions [0x100]func()
	stopped        bool
}

const (
	CLOCK_SPEED       = 4194304
	FPS               = 60
	C_TICKS_PER_FRAME = CLOCK_SPEED / FPS

	ISR_CLOCK_TICKS = 20
)

func NewGameboy(filename string) *Gameboy {
	gb := &Gameboy{}
	gb.init(filename)
	return gb
}

func (gb *Gameboy) init(filename string) {
	gb.cpu = &CPU{}
	gb.cpu.init()

	gb.mmu = &MMU{}

	gb.timer = &Timer{}
	gb.timer.init(gb.mmu)

	gb.cbInstructions = gb.initCbInstructions()
	gb.stopped = false

	gb.initMemoryMap(filename)
}

func (gb *Gameboy) initMemoryMap(filename string) {
	gb.mmu.mapAddrSpace(newBootROM("boot_rom.bin", gb.mmu))
	gb.mmu.mapAddrSpace(newROM(filename))
	gb.mmu.mapAddrSpace(gb.timer)

	// for now have our generic RAM be last in precedence to "catch" unimplemented addresses
	gb.mmu.mapAddrSpace(newGenericRAM())

	// manually disable boot rom for now
	gb.mmu.write(BOOT_ROM_ENABLE_ADDR, 1)
}

func (gb *Gameboy) printRegisters() {
	fmt.Printf("A: 0x%02x | %d\n", gb.cpu.reg.A, gb.cpu.reg.A)
	fmt.Printf("B: 0x%02x | %d\n", gb.cpu.reg.B, gb.cpu.reg.B)
	fmt.Printf("C: 0x%02x | %d\n", gb.cpu.reg.C, gb.cpu.reg.C)
	fmt.Printf("D: 0x%02x | %d\n", gb.cpu.reg.D, gb.cpu.reg.D)
	fmt.Printf("E: 0x%02x | %d\n", gb.cpu.reg.E, gb.cpu.reg.E)
	fmt.Printf("F: 0x%02x | %d\n", gb.cpu.reg.F, gb.cpu.reg.F)
	fmt.Printf("H: 0x%02x | %d\n", gb.cpu.reg.H, gb.cpu.reg.H)
	fmt.Printf("L: 0x%02x | %d\n", gb.cpu.reg.L, gb.cpu.reg.L)
	fmt.Printf("SP: 0x%04x | %d\n", gb.cpu.reg.SP, gb.cpu.reg.SP)
	fmt.Printf("PC: 0x%04x | %d\n", gb.cpu.reg.PC, gb.cpu.reg.PC)
	fmt.Println()
}

func (gb *Gameboy) Start() {
	fmt.Printf("Starting gameboy...\n\n")

	for !gb.stopped {
		gb.update()
	}
}

func (gb *Gameboy) update() {
	cTicksInUpdate := 0
	for cTicksInUpdate < C_TICKS_PER_FRAME {
		cTicks := 4
		if !gb.cpu.halted {
			cTicks = gb.stepCPU()
		}

		cTicksInUpdate += cTicks
		gb.timer.step(cTicks)
		cTicksInUpdate += gb.handleIntrupts()
	}
}

func (gb *Gameboy) stepCPU() int {
	opcode := gb.nextPC()
	return gb.executeInstr(opcode)
}

func (gb *Gameboy) executeInstr(opcode uint8) int {
	// fmt.Printf("Executing OPCODE: 0x%02x at PC: 0x%04x\n", opcode, gb.cpu.reg.PC-1)
	instructions[opcode](gb)
	cTicks := instrClockTicks[opcode]
	gb.cpu.ticks += cTicks

	// fmt.Println("Registers After: ")
	// gb.printRegisters()
	return cTicks
}

func (gb *Gameboy) pushStack(addr uint16) {
	gb.mmu.write(gb.cpu.reg.SP-1, bits.HiByte(addr))
	gb.mmu.write(gb.cpu.reg.SP-2, bits.LoByte(addr))

	gb.cpu.setSP(gb.cpu.reg.SP - 2)
}

func (gb *Gameboy) popStack() uint16 {
	loByte := gb.mmu.read(gb.cpu.reg.SP)
	hiByte := gb.mmu.read(gb.cpu.reg.SP + 1)

	gb.cpu.setSP(gb.cpu.reg.SP + 2)
	return uint16(hiByte)<<8 | uint16(loByte)
}

func (gb *Gameboy) handleIntrupts() int {
	if gb.cpu.IMEDelay {
		gb.cpu.clearIMEDelay()
		gb.cpu.setIME()
		return 0
	}

	if !gb.cpu.IME && !gb.cpu.halted {
		return 0
	}

	IE := gb.readIE()
	IF := gb.readIF()

	if bits.IsSetInBoth(IE, IF, VBLANK_INTRUPT_BIT) {
		gb.serviceIntrupt(VBLANK_INTRUPT_BIT)
	} else if bits.IsSetInBoth(IE, IF, LCD_INTRUPT_BIT) {
		gb.serviceIntrupt(LCD_INTRUPT_BIT)
	} else if bits.IsSetInBoth(IE, IF, TIMER_INTRUPT_BIT) {
		gb.serviceIntrupt(TIMER_INTRUPT_BIT)
	} else if bits.IsSetInBoth(IE, IF, SERIAL_INTRUPT_BIT) {
		gb.serviceIntrupt(SERIAL_INTRUPT_BIT)
	} else if bits.IsSetInBoth(IE, IF, JOYPAD_INTRUPT_BIT) {
		gb.serviceIntrupt(JOYPAD_INTRUPT_BIT)
	} else {
		return 0
	}

	return ISR_CLOCK_TICKS
}

func (gb *Gameboy) serviceIntrupt(intruptBit uint8) {
	if !gb.cpu.IME && gb.cpu.halted {
		gb.cpu.exitHaltedState()
		return
	}

	gb.cpu.exitHaltedState()
	gb.cpu.clearIME()
	gb.clearIFBit(intruptBit)
	gb.pushStack(gb.cpu.reg.PC)

	switch intruptBit {
	case VBLANK_INTRUPT_BIT:
		gb.cpu.setPC(VBLANK_INTRUPT_VEC)
	case LCD_INTRUPT_BIT:
		gb.cpu.setPC(STAT_INTRUPT_VEC)
	case TIMER_INTRUPT_BIT:
		gb.cpu.setPC(TIMER_INTRUPT_VEC)
	case SERIAL_INTRUPT_BIT:
		gb.cpu.setPC(SERIAL_INTRUPT_VEC)
	case JOYPAD_INTRUPT_BIT:
		gb.cpu.setPC(JOYPAD_INTRUPT_VEC)
	}
}

func (gb *Gameboy) clearIFBit(pos uint8) {
	gb.mmu.write(IF_ADDR, bits.Reset(gb.readIF(), pos))
}

func (gb *Gameboy) readIE() uint8 {
	return gb.mmu.read(IE_ADDR)
}

func (gb *Gameboy) readIF() uint8 {
	return gb.mmu.read(IF_ADDR)
}
