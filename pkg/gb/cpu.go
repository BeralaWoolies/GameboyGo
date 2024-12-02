package gb

import "github.com/BeralaWoolies/GameboyGo/pkg/bits"

type CPU struct {
	reg            *Registers
	mmu            *MMU
	cbInstructions [0x100]func()
	ticks          int
	halted         bool
	IME            bool
	IMEDelay       bool
}

type Registers struct {
	A  uint8
	B  uint8
	C  uint8
	D  uint8
	E  uint8
	F  uint8
	H  uint8
	L  uint8
	SP uint16
	PC uint16
}

const (
	ZERO_FLAG_BIT       = 7
	SUB_FLAG_BIT        = 6
	HALF_CARRY_FLAG_BIT = 5
	CARRY_FLAG_BIT      = 4
)

func (cpu *CPU) init(mmu *MMU) {
	cpu.reg = &Registers{}
	cpu.mmu = mmu
	cpu.cbInstructions = cpu.initCbInstructions()
	cpu.ticks = 0
	cpu.halted = false
	cpu.IME = false
	cpu.IMEDelay = false
}

func (cpu *CPU) step() int {
	opcode := cpu.nextPC()
	return cpu.executeInstr(opcode)
}

func (cpu *CPU) executeInstr(opcode uint8) int {
	// fmt.Printf("Executing OPCODE: 0x%02x at PC: 0x%04x\n", opcode, gb.cpu.reg.PC-1)
	branchTicks := instructions[opcode](cpu)
	return instrBaseTicks[opcode] + branchTicks
}

func (cpu *CPU) pushStack(addr uint16) {
	cpu.mmu.write(cpu.reg.SP-1, bits.HiByte(addr))
	cpu.mmu.write(cpu.reg.SP-2, bits.LoByte(addr))

	cpu.setSP(cpu.reg.SP - 2)
}

func (cpu *CPU) popStack() uint16 {
	loByte := cpu.mmu.read(cpu.reg.SP)
	hiByte := cpu.mmu.read(cpu.reg.SP + 1)

	cpu.setSP(cpu.reg.SP + 2)
	return uint16(hiByte)<<8 | uint16(loByte)
}

func (cpu *CPU) nextPC() uint8 {
	data := cpu.mmu.read(cpu.reg.PC)
	cpu.reg.PC++
	return data
}

func (cpu *CPU) nextPC16() uint16 {
	loByte := cpu.nextPC()
	hiByte := cpu.nextPC()
	return uint16(hiByte)<<8 | uint16(loByte)
}

func (cpu *CPU) setA(val uint8) {
	cpu.reg.A = val
}

func (cpu *CPU) setB(val uint8) {
	cpu.reg.B = val
}

func (cpu *CPU) setC(val uint8) {
	cpu.reg.C = val
}

func (cpu *CPU) setD(val uint8) {
	cpu.reg.D = val
}

func (cpu *CPU) setE(val uint8) {
	cpu.reg.E = val
}

func (cpu *CPU) setH(val uint8) {
	cpu.reg.H = val
}

func (cpu *CPU) setL(val uint8) {
	cpu.reg.L = val
}

func (cpu *CPU) setSP(val uint16) {
	cpu.reg.SP = val
}

func (cpu *CPU) setPC(val uint16) {
	cpu.reg.PC = val
}

func (cpu *CPU) getAF() uint16 {
	return uint16(cpu.reg.A)<<8 | uint16(cpu.reg.F)
}

func (cpu *CPU) setAF(val uint16) {
	cpu.reg.A = bits.HiByte(val)
	cpu.reg.F = bits.LoByte(val) & 0xF0
}

func (cpu *CPU) getBC() uint16 {
	return uint16(cpu.reg.B)<<8 | uint16(cpu.reg.C)
}

func (cpu *CPU) setBC(val uint16) {
	cpu.reg.B = bits.HiByte(val)
	cpu.reg.C = bits.LoByte(val)
}

func (cpu *CPU) getDE() uint16 {
	return uint16(cpu.reg.D)<<8 | uint16(cpu.reg.E)
}

func (cpu *CPU) setDE(val uint16) {
	cpu.reg.D = bits.HiByte(val)
	cpu.reg.E = bits.LoByte(val)
}

func (cpu *CPU) getHL() uint16 {
	return uint16(cpu.reg.H)<<8 | uint16(cpu.reg.L)
}

func (cpu *CPU) setHL(val uint16) {
	cpu.reg.H = bits.HiByte(val)
	cpu.reg.L = bits.LoByte(val)
}

func (cpu *CPU) setFlag(pos uint8, set bool) {
	if set {
		cpu.reg.F = bits.Set(cpu.reg.F, pos)
	} else {
		cpu.reg.F = bits.Reset(cpu.reg.F, pos)
	}
}

func (cpu *CPU) setZFlag(set bool) {
	cpu.setFlag(ZERO_FLAG_BIT, set)
}

func (cpu *CPU) setNFlag(set bool) {
	cpu.setFlag(SUB_FLAG_BIT, set)
}

func (cpu *CPU) setHFlag(set bool) {
	cpu.setFlag(HALF_CARRY_FLAG_BIT, set)
}

func (cpu *CPU) setCFlag(set bool) {
	cpu.setFlag(CARRY_FLAG_BIT, set)
}

func (cpu *CPU) zFlag() bool {
	return (cpu.reg.F>>ZERO_FLAG_BIT)&1 == 1
}

func (cpu *CPU) nFlag() bool {
	return (cpu.reg.F>>SUB_FLAG_BIT)&1 == 1
}

func (cpu *CPU) hFlag() bool {
	return (cpu.reg.F>>HALF_CARRY_FLAG_BIT)&1 == 1
}

func (cpu *CPU) cFlag() bool {
	return (cpu.reg.F>>CARRY_FLAG_BIT)&1 == 1
}

func (cpu *CPU) setIME() {
	cpu.IME = true
}

func (cpu *CPU) clearIME() {
	cpu.IME = false
}

func (cpu *CPU) setIMEDelay() {
	cpu.IMEDelay = true
}

func (cpu *CPU) clearIMEDelay() {
	cpu.IMEDelay = false
}

func (cpu *CPU) enterHaltedState() {
	cpu.halted = true
}

func (cpu *CPU) exitHaltedState() {
	cpu.halted = false
}
