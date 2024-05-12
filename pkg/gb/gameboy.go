package gb

import (
	"fmt"
	"log"
	"os"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
)

type Gameboy struct {
	cpu            *CPU
	mmu            *MMU
	cbInstructions [0x100]func()
	stopped        bool
}

const (
	CLOCK_SPEED           = 4194304
	FPS                   = 60
	CLOCK_TICKS_PER_FRAME = CLOCK_SPEED / FPS

	ISR_CLOCK_TICKS = 20
)

func NewGameboy(romName string) *Gameboy {
	gb := &Gameboy{}
	gb.init()
	gb.loadRom(romName)
	return gb
}

func (gb *Gameboy) init() {
	gb.cpu = &CPU{}
	gb.cpu.init()

	gb.mmu = &MMU{}
	gb.mmu.init()

	gb.cbInstructions = gb.initCbInstructions()
	gb.stopped = false
}

func (gb *Gameboy) loadRom(romName string) {
	rom, err := os.ReadFile(romName)
	if err != nil {
		log.Fatal(err)
	}

	gb.mmu.mapRom(rom)
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

	clockTicks := 0
	for clockTicks < CLOCK_TICKS_PER_FRAME || !gb.stopped {
		ticks := 4
		if !gb.cpu.halted {
			ticks = gb.stepCPU()
		}

		clockTicks += ticks
		gb.stepTimer(ticks)
		clockTicks += gb.handleIntrupts()
	}
}

func (gb *Gameboy) stepCPU() int {
	opcode := gb.nextPC()
	return gb.executeInstr(opcode)
}

func (gb *Gameboy) executeInstr(opcode uint8) int {
	// fmt.Printf("Executing OPCODE: 0x%02x at PC: 0x%04x\n", opcode, gb.cpu.reg.PC-1)
	instructions[opcode](gb)
	gb.cpu.ticks += instrClockTicks[opcode]

	// fmt.Println("Registers After: ")
	// gb.printRegisters()
	return instrClockTicks[opcode]
}

func (gb *Gameboy) pushStack(address uint16) {
	gb.mmu.write(gb.cpu.reg.SP-1, bits.HiByte(address))
	gb.mmu.write(gb.cpu.reg.SP-2, bits.LoByte(address))

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

func (gb *Gameboy) stepTimer(clockTicks int) {
	gb.stepDIV(clockTicks)
	if gb.timerEnabled() {
		gb.stepTIMA(clockTicks)
	}
}

func (gb *Gameboy) stepDIV(clockTicks int) {
	gb.mmu.DIVTicks += clockTicks
	if gb.mmu.DIVTicks >= 256 {
		gb.mmu.DIVTicks -= 256
		gb.mmu.incDIV()
	}
}

func (gb *Gameboy) stepTIMA(clockTicks int) {
	gb.mmu.TIMATicks += clockTicks

	ovflowTicks := 0
	switch freqDiv := gb.readTAC() & TAC_FREQ_DIV_MSK; freqDiv {
	case HZ_4096:
		ovflowTicks = 1024
	case HZ_262144:
		ovflowTicks = 16
	case HZ_65536:
		ovflowTicks = 64
	case HZ_16386:
		ovflowTicks = 256
	}

	for gb.mmu.TIMATicks >= ovflowTicks {
		gb.mmu.TIMATicks -= ovflowTicks
		gb.writeTIMA(gb.readTIMA() + 1)

		if gb.readTIMA() == 0x00 {
			gb.requestIntrupt(TIMER_INTRUPT_BIT)
			gb.writeTIMA(gb.readTMA())
		}
	}
}

func (gb *Gameboy) timerEnabled() bool {
	return bits.IsSet(gb.readTAC(), TAC_TIMER_ENABLE_BIT)
}

func (gb *Gameboy) requestIntrupt(intruptBit uint8) {
	gb.mmu.write(IF_ADDR, bits.Set(gb.readIF(), intruptBit))
}

func (gb *Gameboy) writeTIMA(val uint8) {
	gb.mmu.write(TIMA_ADDR, val)
}

func (gb *Gameboy) readTIMA() uint8 {
	return gb.mmu.read(TIMA_ADDR)
}

func (gb *Gameboy) readTMA() uint8 {
	return gb.mmu.read(TMA_ADDR)
}

func (gb *Gameboy) readTAC() uint8 {
	return gb.mmu.read(TAC_ADDR)
}
