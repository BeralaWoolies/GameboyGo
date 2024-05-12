package gb

import "github.com/BeralaWoolies/GameboyGo/pkg/bits"

type CPU struct {
	reg      *Registers
	ticks    int
	halted   bool
	IME      bool
	IMEDelay bool
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
	zeroFlagBitPos      = 7
	subFlagBitPos       = 6
	halfCarryFlagBitPos = 5
	carryFlagBitPos     = 4
)

func (cpu *CPU) init() {
	cpu.reg = &Registers{}

	cpu.reg.A = 0x01
	cpu.reg.B = 0x00
	cpu.reg.C = 0x13
	cpu.reg.D = 0x00
	cpu.reg.E = 0xD8
	cpu.reg.F = 0xB0
	cpu.reg.H = 0x01
	cpu.reg.L = 0x4D
	cpu.reg.SP = 0xFFFE
	cpu.reg.PC = 0x100

	cpu.ticks = 0
	cpu.halted = false
	cpu.IME = false
	cpu.IMEDelay = false
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
	cpu.setFlag(zeroFlagBitPos, set)
}

func (cpu *CPU) setNFlag(set bool) {
	cpu.setFlag(subFlagBitPos, set)
}

func (cpu *CPU) setHFlag(set bool) {
	cpu.setFlag(halfCarryFlagBitPos, set)
}

func (cpu *CPU) setCFlag(set bool) {
	cpu.setFlag(carryFlagBitPos, set)
}

func (cpu *CPU) zFlag() bool {
	return (cpu.reg.F>>zeroFlagBitPos)&1 == 1
}

func (cpu *CPU) nFlag() bool {
	return (cpu.reg.F>>subFlagBitPos)&1 == 1
}

func (cpu *CPU) hFlag() bool {
	return (cpu.reg.F>>halfCarryFlagBitPos)&1 == 1
}

func (cpu *CPU) cFlag() bool {
	return (cpu.reg.F>>carryFlagBitPos)&1 == 1
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
