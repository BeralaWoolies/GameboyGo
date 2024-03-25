package gb

import (
	"fmt"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
)

type Gameboy struct {
	cpu                      *CPU
	mmu                      *MMU
	interruptsOn             bool
	interruptsPendingEnabled bool
	cbInstructions           [0x100]func()
}

func NewGameboy() *Gameboy {
	gb := &Gameboy{}
	gb.init()
	return gb
}

func (gb *Gameboy) init() {
	gb.cpu = &CPU{}
	gb.cpu.init()

	gb.mmu = &MMU{}
	gb.mmu.init()

	gb.cbInstructions = gb.initCbInstructions()
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
	fmt.Println("Initial Registers:")
	gb.printRegisters()
	gb.step()
}

func (gb *Gameboy) step() {
	opcode := gb.nextPC()
	opcode = 0x04
	cycles := gb.executeInstr(opcode)
	fmt.Printf("Took %d Clock Ticks\n", cycles)
}

func (gb *Gameboy) executeInstr(opcode uint8) int {
	fmt.Printf("Executing OPCODE: 0x%02x at PC: 0x%04x\n", opcode, gb.cpu.reg.PC-1)

	gb.cpu.ticks = instrClockTicks[opcode]
	instructions[opcode](gb)

	fmt.Println("Registers After: ")
	gb.printRegisters()
	return gb.cpu.ticks
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
