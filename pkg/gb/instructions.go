package gb

import (
	"log"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
)

var instrBaseTicks = [0x100]int{
	0x04, 0x0C, 0x08, 0x08, 0x04, 0x04, 0x08, 0x04, 0x14, 0x08, 0x08, 0x08, 0x04, 0x04, 0x08, 0x04,
	0x04, 0x0C, 0x08, 0x08, 0x04, 0x04, 0x08, 0x04, 0x0C, 0x08, 0x08, 0x08, 0x04, 0x04, 0x08, 0x04,
	0x08, 0x0C, 0x08, 0x08, 0x04, 0x04, 0x08, 0x04, 0x08, 0x08, 0x08, 0x08, 0x04, 0x04, 0x08, 0x04,
	0x08, 0x0C, 0x08, 0x08, 0x0C, 0x0C, 0x0C, 0x04, 0x08, 0x08, 0x08, 0x08, 0x04, 0x04, 0x08, 0x04,
	0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04,
	0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04,
	0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x04, 0x08, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04,
	0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04,
	0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04,
	0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04,
	0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x04,
	0x08, 0x0C, 0x0C, 0x10, 0x0C, 0x10, 0x08, 0x10, 0x08, 0x10, 0x0C, 0x04, 0x0C, 0x18, 0x08, 0x10,
	0x08, 0x0C, 0x0C, 0x00, 0x0C, 0x10, 0x08, 0x10, 0x08, 0x10, 0x0C, 0x00, 0x0C, 0x00, 0x08, 0x10,
	0x0C, 0x0C, 0x08, 0x00, 0x00, 0x10, 0x08, 0x10, 0x10, 0x04, 0x10, 0x00, 0x00, 0x00, 0x08, 0x10,
	0x0C, 0x0C, 0x08, 0x04, 0x00, 0x10, 0x08, 0x10, 0x0C, 0x08, 0x10, 0x04, 0x00, 0x00, 0x08, 0x10,
}

func (cpu *CPU) instrInc(setHandler func(result uint8), val uint8) {
	result := val + 1
	setHandler(result)

	cpu.setZFlag(result == 0)
	cpu.setNFlag(false)
	cpu.setHFlag((val&0xF)+(1&0xF) > 0xF)
}

func (cpu *CPU) instrDec(setHandler func(result uint8), val uint8) {
	result := val - 1
	setHandler(result)

	cpu.setZFlag(result == 0)
	cpu.setNFlag(true)
	cpu.setHFlag(val&0xF == 0)
}

func (cpu *CPU) instrAddA(rhs uint8, addCarry bool) {
	var carry int16 = 0
	if cpu.cFlag() && addCarry {
		carry = 1
	}

	lhs := cpu.reg.A
	result := int16(lhs) + int16(rhs) + carry
	cpu.setA(uint8(result))

	cpu.setZFlag(cpu.reg.A == 0)
	cpu.setNFlag(false)
	cpu.setHFlag((lhs&0xF)+(rhs&0xF)+uint8(carry) > 0xF)
	cpu.setCFlag(result > 0xFF)
}

func (cpu *CPU) instrSubA(rhs uint8, subCarry bool) {
	var carry int16 = 0
	if cpu.cFlag() && subCarry {
		carry = 1
	}

	lhs := cpu.reg.A
	result := int16(lhs) - int16(rhs) - carry
	cpu.setA(uint8(result))

	cpu.setZFlag(cpu.reg.A == 0)
	cpu.setNFlag(true)
	cpu.setHFlag(int16(lhs&0xF)-int16(rhs&0xF)-carry < 0)
	cpu.setCFlag(result < 0)
}

func (cpu *CPU) instrAdd16(setHandler func(result uint16), val1 uint16, val2 uint16) {
	result := int32(val1) + int32(val2)
	setHandler(uint16(result))

	cpu.setNFlag(false)
	cpu.setHFlag(int32(val1&0xFFF) > (result & 0xFFF))
	cpu.setCFlag(result > 0xFFFF)
}

func (cpu *CPU) instrAdd16Signed(setHandler func(result uint16), val1 uint16, val2 int8) {
	result := uint16(int32(val1) + int32(val2))
	setHandler(result)

	carryBits := val1 ^ uint16(val2) ^ result

	cpu.setZFlag(false)
	cpu.setNFlag(false)
	cpu.setHFlag((carryBits & 0x10) == 0x10)
	cpu.setCFlag((carryBits & 0x100) == 0x100)
}

func (cpu *CPU) instrAndA(rhs uint8) {
	cpu.setA(cpu.reg.A & rhs)

	cpu.setZFlag(cpu.reg.A == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(true)
	cpu.setCFlag(false)
}

func (cpu *CPU) instrXorA(rhs uint8) {
	cpu.setA(cpu.reg.A ^ rhs)

	cpu.setZFlag(cpu.reg.A == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(false)
	cpu.setCFlag(false)
}

func (cpu *CPU) instrOrA(rhs uint8) {
	cpu.setA(cpu.reg.A | rhs)

	cpu.setZFlag(cpu.reg.A == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(false)
	cpu.setCFlag(false)
}

func (cpu *CPU) instrCpA(rhs uint8) {
	lhs := cpu.reg.A
	cmpResult := lhs - rhs

	cpu.setZFlag(cmpResult == 0)
	cpu.setNFlag(true)
	cpu.setHFlag((rhs & 0xF) > (lhs & 0xF))
	cpu.setCFlag(rhs > lhs)
}

func (cpu *CPU) instrJump(jumpAddress uint16) {
	cpu.setPC(jumpAddress)
}

func (cpu *CPU) instrRet() {
	cpu.setPC(cpu.popStack())
}

func (cpu *CPU) instrCall(jumpAddress uint16) {
	cpu.pushStack(cpu.reg.PC)
	cpu.setPC(jumpAddress)
}

var instructions = [0x100]func(cpu *CPU) int{
	0x00: func(cpu *CPU) int {
		// NOP
		// fmt.Println("Decoded OPCODE: NOP")
		return 0
	},
	0x01: func(cpu *CPU) int {
		// LD BC, u16
		// fmt.Println("Decoded OPCODE: LD BC, u16")
		cpu.setBC(cpu.nextPC16())
		return 0
	},
	0x02: func(cpu *CPU) int {
		// LD [BC], A
		// fmt.Println("Decoded OPCODE: LD [BC], A")
		cpu.mmu.write(cpu.getBC(), cpu.reg.A)
		return 0
	},
	0x03: func(cpu *CPU) int {
		// INC BC
		// fmt.Println("Decoded OPCODE: INC BC")
		cpu.setBC(cpu.getBC() + 1)
		return 0
	},
	0x04: func(cpu *CPU) int {
		// INC B
		// fmt.Println("Decoded OPCODE: INC B")
		cpu.instrInc(cpu.setB, cpu.reg.B)
		return 0
	},
	0x05: func(cpu *CPU) int {
		// DEC B
		// fmt.Println("Decoded OPCODE: DEC B")
		cpu.instrDec(cpu.setB, cpu.reg.B)
		return 0
	},
	0x06: func(cpu *CPU) int {
		// LD B, u8
		// fmt.Println("Decoded OPCODE: LD B, u8")
		cpu.setB(cpu.nextPC())
		return 0
	},
	0x07: func(cpu *CPU) int {
		// RLCA
		// fmt.Println("Decoded OPCODE: RLCA")
		val := cpu.reg.A
		cpu.setA(uint8((val << 1)) | (val >> 7))

		cpu.setZFlag(false)
		cpu.setNFlag(false)
		cpu.setHFlag(false)
		cpu.setCFlag(val >= 0x80)
		return 0
	},
	0x08: func(cpu *CPU) int {
		// LD [u16], SP
		// fmt.Println("Decoded OPCODE: LD [u16], SP")
		address := cpu.nextPC16()
		cpu.mmu.write(address, bits.LoByte(cpu.reg.SP))
		cpu.mmu.write(address+1, bits.HiByte(cpu.reg.SP))
		return 0
	},
	0x09: func(cpu *CPU) int {
		// ADD HL, BC
		// fmt.Println("Decoded OPCODE: ADD HL, BC")
		cpu.instrAdd16(cpu.setHL, cpu.getHL(), cpu.getBC())
		return 0
	},
	0x0A: func(cpu *CPU) int {
		// LD A, [BC]
		// fmt.Println("Decoded OPCODE: LD A, [BC]")
		cpu.setA(cpu.mmu.read(cpu.getBC()))
		return 0
	},
	0x0B: func(cpu *CPU) int {
		// DEC BC
		// fmt.Println("Decoded OPCODE: DEC BC")
		cpu.setBC(cpu.getBC() - 1)
		return 0
	},
	0x0C: func(cpu *CPU) int {
		// INC C
		// fmt.Println("Decoded OPCODE: INC C")
		cpu.instrInc(cpu.setC, cpu.reg.C)
		return 0
	},
	0x0D: func(cpu *CPU) int {
		// DEC C
		// fmt.Println("Decoded OPCODE: DEC C")
		cpu.instrDec(cpu.setC, cpu.reg.C)
		return 0
	},
	0x0E: func(cpu *CPU) int {
		// LD C, u8
		// fmt.Println("Decoded OPCODE: LD C, u8")
		cpu.setC(cpu.nextPC())
		return 0
	},
	0x0F: func(cpu *CPU) int {
		// RRCA
		// fmt.Println("Decoded OPCODE: RRCA")
		val := cpu.reg.A
		cpu.setA((val >> 1) | ((val & 1) << 7))

		cpu.setZFlag(false)
		cpu.setNFlag(false)
		cpu.setHFlag(false)
		cpu.setCFlag(cpu.reg.A >= 0x80)
		return 0
	},
	0x10: func(cpu *CPU) int {
		// STOP
		// fmt.Println("Decoded OPCODE: STOP")
		cpu.enterHaltedState()
		cpu.nextPC()
		return 0
	},
	0x11: func(cpu *CPU) int {
		// LD DE, u16
		// fmt.Println("Decoded OPCODE: LD DE, u16")
		cpu.setDE(cpu.nextPC16())
		return 0
	},
	0x12: func(cpu *CPU) int {
		// LD [DE], A
		// fmt.Println("Decoded OPCODE: LD [DE], A")
		cpu.mmu.write(cpu.getDE(), cpu.reg.A)
		return 0
	},
	0x13: func(cpu *CPU) int {
		// INC DE
		// fmt.Println("Decoded OPCODE: INC DE")
		cpu.setDE(cpu.getDE() + 1)
		return 0
	},
	0x14: func(cpu *CPU) int {
		// INC D
		// fmt.Println("Decoded OPCODE: INC D")
		cpu.instrInc(cpu.setD, cpu.reg.D)
		return 0
	},
	0x15: func(cpu *CPU) int {
		// DEC D
		// fmt.Println("Decoded OPCODE: DEC D")
		cpu.instrDec(cpu.setD, cpu.reg.D)
		return 0
	},
	0x16: func(cpu *CPU) int {
		// LD D, u8
		// fmt.Println("Decoded OPCODE: LD D, u8")
		cpu.setD(cpu.nextPC())
		return 0
	},
	0x17: func(cpu *CPU) int {
		// RLA
		// fmt.Println("Decoded OPCODE: RLA")
		var carry uint8 = 0
		if cpu.cFlag() {
			carry = 1
		}

		val := cpu.reg.A
		cpu.setA(uint8(val<<1) | carry)

		cpu.setZFlag(false)
		cpu.setNFlag(false)
		cpu.setHFlag(false)
		cpu.setCFlag(val >= 0x80)
		return 0
	},
	0x18: func(cpu *CPU) int {
		// JR i8
		// fmt.Println("Decoded OPCODE: JR i8")
		jumpAddress := int32(cpu.reg.PC) + int32(int8(cpu.nextPC()))
		cpu.instrJump(uint16(jumpAddress))
		return 0
	},
	0x19: func(cpu *CPU) int {
		// ADD HL, DE
		// fmt.Println("Decoded OPCODE: ADD HL, DE")
		cpu.instrAdd16(cpu.setHL, cpu.getHL(), cpu.getDE())
		return 0
	},
	0x1A: func(cpu *CPU) int {
		// LD A, [DE]
		// fmt.Println("Decoded OPCODE: LD A, [DE]")
		cpu.setA(cpu.mmu.read(cpu.getDE()))
		return 0
	},
	0x1B: func(cpu *CPU) int {
		// DEC DE
		// fmt.Println("Decoded OPCODE: DEC DE")
		cpu.setDE(cpu.getDE() - 1)
		return 0
	},
	0x1C: func(cpu *CPU) int {
		// INC E
		// fmt.Println("Decoded OPCODE: INC E")
		cpu.instrInc(cpu.setE, cpu.reg.E)
		return 0
	},
	0x1D: func(cpu *CPU) int {
		// DEC E
		// fmt.Println("Decoded OPCODE: DEC E")
		cpu.instrDec(cpu.setE, cpu.reg.E)
		return 0
	},
	0x1E: func(cpu *CPU) int {
		// LD E, u8
		// fmt.Println("Decoded OPCODE: LD E, u8")
		cpu.setE(cpu.nextPC())
		return 0
	},
	0x1F: func(cpu *CPU) int {
		// RRA
		// fmt.Println("Decoded OPCODE: RRA")
		var carry uint8 = 0
		if cpu.cFlag() {
			carry = 0x80
		}

		val := cpu.reg.A
		cpu.setA(uint8(val>>1) | carry)

		cpu.setZFlag(false)
		cpu.setNFlag(false)
		cpu.setHFlag(false)
		cpu.setCFlag((val & 1) == 1)
		return 0
	},
	0x20: func(cpu *CPU) int {
		// JR NZ, i8
		// fmt.Println("Decoded OPCODE: JR NZ, i8")
		offset := int8(cpu.nextPC())
		if !cpu.zFlag() {
			jumpAddress := int32(cpu.reg.PC) + int32(offset)
			cpu.instrJump(uint16(jumpAddress))
			return 4
		}
		return 0
	},
	0x21: func(cpu *CPU) int {
		// LD HL, u16
		// fmt.Println("Decoded OPCODE: LD HL, u16")
		cpu.setHL(cpu.nextPC16())
		return 0
	},
	0x22: func(cpu *CPU) int {
		// LD [HL+], A
		// fmt.Println("Decoded OPCODE: LD [HL+], A")
		cpu.mmu.write(cpu.getHL(), cpu.reg.A)
		cpu.setHL(cpu.getHL() + 1)
		return 0
	},
	0x23: func(cpu *CPU) int {
		// INC HL
		// fmt.Println("Decoded OPCODE: INC HL")
		cpu.setHL(cpu.getHL() + 1)
		return 0
	},
	0x24: func(cpu *CPU) int {
		// INC H
		// fmt.Println("Decoded OPCODE: INC H")
		cpu.instrInc(cpu.setH, cpu.reg.H)
		return 0
	},
	0x25: func(cpu *CPU) int {
		// DEC H
		// fmt.Println("Decoded OPCODE: DEC H")
		cpu.instrDec(cpu.setH, cpu.reg.H)
		return 0
	},
	0x26: func(cpu *CPU) int {
		// LD H, u8
		// fmt.Println("Decoded OPCODE: LD H, u8")
		cpu.setH(cpu.nextPC())
		return 0
	},
	0x27: func(cpu *CPU) int {
		// DAA
		// fmt.Println("Decoded OPCODE: DAA")
		if !cpu.nFlag() {
			if cpu.cFlag() || cpu.reg.A > 0x99 {
				cpu.reg.A += 0x60
				cpu.setCFlag(true)
			}
			if cpu.hFlag() || (cpu.reg.A&0xF) > 0x09 {
				cpu.reg.A += 0x06
			}
		} else {
			if cpu.cFlag() {
				cpu.reg.A -= 0x60
			}
			if cpu.hFlag() {
				cpu.reg.A -= 0x06
			}
		}

		cpu.setZFlag(cpu.reg.A == 0)
		cpu.setHFlag(false)
		return 0
	},
	0x28: func(cpu *CPU) int {
		// JR Z, i8
		// fmt.Println("Decoded OPCODE: JR Z, i8")
		offset := int8(cpu.nextPC())
		if cpu.zFlag() {
			jumpAddress := int32(cpu.reg.PC) + int32(offset)
			cpu.instrJump(uint16(jumpAddress))
			return 4
		}
		return 0
	},
	0x29: func(cpu *CPU) int {
		// ADD HL, HL
		// fmt.Println("Decoded OPCODE: ADD HL, HL")
		cpu.instrAdd16(cpu.setHL, cpu.getHL(), cpu.getHL())
		return 0
	},
	0x2A: func(cpu *CPU) int {
		// LD A, [HL+]
		// fmt.Println("Decoded OPCODE: LD A, [HL+]")
		cpu.setA(cpu.mmu.read(cpu.getHL()))
		cpu.setHL(cpu.getHL() + 1)
		return 0
	},
	0x2B: func(cpu *CPU) int {
		// DEC HL
		// fmt.Println("Decoded OPCODE: DEC HL")
		cpu.setHL(cpu.getHL() - 1)
		return 0
	},
	0x2C: func(cpu *CPU) int {
		// INC L
		// fmt.Println("Decoded OPCODE: INC L")
		cpu.instrInc(cpu.setL, cpu.reg.L)
		return 0
	},
	0x2D: func(cpu *CPU) int {
		// DEC L
		// fmt.Println("Decoded OPCODE: DEC L")
		cpu.instrDec(cpu.setL, cpu.reg.L)
		return 0
	},
	0x2E: func(cpu *CPU) int {
		// LD L, u8
		// fmt.Println("Decoded OPCODE: LD L, u8")
		cpu.setL(cpu.nextPC())
		return 0
	},
	0x2F: func(cpu *CPU) int {
		// CPL
		// fmt.Println("Decoded OPCODE: CPL")
		cpu.reg.A = ^(cpu.reg.A)
		cpu.setNFlag(true)
		cpu.setHFlag(true)
		return 0
	},
	0x30: func(cpu *CPU) int {
		// JR NC, i8
		// fmt.Println("Decoded OPCODE: JR NC, i8")
		offset := int8(cpu.nextPC())
		if !cpu.cFlag() {
			jumpAddress := int32(cpu.reg.PC) + int32(offset)
			cpu.instrJump(uint16(jumpAddress))
			return 4
		}
		return 0
	},
	0x31: func(cpu *CPU) int {
		// LD SP, u16
		// fmt.Println("Decoded OPCODE: LD SP, u16")
		cpu.setSP(cpu.nextPC16())
		return 0
	},
	0x32: func(cpu *CPU) int {
		// LD [HL-], A
		// fmt.Println("Decoded OPCODE: LD [HL-], A")
		cpu.mmu.write(cpu.getHL(), cpu.reg.A)
		cpu.setHL(cpu.getHL() - 1)
		return 0
	},
	0x33: func(cpu *CPU) int {
		// INC SP
		// fmt.Println("Decoded OPCODE: INC SP")
		cpu.setSP(cpu.reg.SP + 1)
		return 0
	},
	0x34: func(cpu *CPU) int {
		// INC [HL]
		// fmt.Println("Decoded OPCODE: INC [HL]")
		val := cpu.mmu.read(cpu.getHL())
		cpu.instrInc(func(result uint8) { cpu.mmu.write(cpu.getHL(), result) }, val)
		return 0
	},
	0x35: func(cpu *CPU) int {
		// DEC [HL]
		// fmt.Println("Decoded OPCODE: DEC [HL]")
		val := cpu.mmu.read(cpu.getHL())
		cpu.instrDec(func(result uint8) { cpu.mmu.write(cpu.getHL(), result) }, val)
		return 0
	},
	0x36: func(cpu *CPU) int {
		// LD [HL], u8
		// fmt.Println("Decoded OPCODE: LD [HL], u8")
		cpu.mmu.write(cpu.getHL(), cpu.nextPC())
		return 0
	},
	0x37: func(cpu *CPU) int {
		// SCF
		// fmt.Println("Decoded OPCODE: SCF")
		cpu.setNFlag(false)
		cpu.setHFlag(false)
		cpu.setCFlag(true)
		return 0
	},
	0x38: func(cpu *CPU) int {
		// JR C, i8
		// fmt.Println("Decoded OPCODE: JR C, i8")
		offset := int8(cpu.nextPC())
		if cpu.cFlag() {
			jumpAddress := int32(cpu.reg.PC) + int32(offset)
			cpu.instrJump(uint16(jumpAddress))
			return 4
		}
		return 0
	},
	0x39: func(cpu *CPU) int {
		// ADD HL, SP
		// fmt.Println("Decoded OPCODE: ADD HL, SP")
		cpu.instrAdd16(cpu.setHL, cpu.getHL(), cpu.reg.SP)
		return 0
	},
	0x3A: func(cpu *CPU) int {
		// LD A, [HL-]
		// fmt.Println("Decoded OPCODE: LD A, [HL-]")
		cpu.setA(cpu.mmu.read(cpu.getHL()))
		cpu.setHL(cpu.getHL() - 1)
		return 0
	},
	0x3B: func(cpu *CPU) int {
		// DEC SP
		// fmt.Println("Decoded OPCODE: DEC SP")
		cpu.setSP(cpu.reg.SP - 1)
		return 0
	},
	0x3C: func(cpu *CPU) int {
		// INC A
		// fmt.Println("Decoded OPCODE: INC A")
		cpu.instrInc(cpu.setA, cpu.reg.A)
		return 0
	},
	0x3D: func(cpu *CPU) int {
		// DEC A
		// fmt.Println("Decoded OPCODE: DEC A")
		cpu.instrDec(cpu.setA, cpu.reg.A)
		return 0
	},
	0x3E: func(cpu *CPU) int {
		// LD A, u8
		// fmt.Println("Decoded OPCODE: LD A, u8")
		cpu.setA(cpu.nextPC())
		return 0
	},
	0x3F: func(cpu *CPU) int {
		// CCF
		// fmt.Println("Decoded OPCODE: CCF")
		cpu.setNFlag(false)
		cpu.setHFlag(false)
		cpu.setCFlag(!cpu.cFlag())
		return 0
	},
	0x40: func(cpu *CPU) int {
		// LD B, B
		// fmt.Println("Decoded OPCODE: LD B, B")
		cpu.setB(cpu.reg.B)
		return 0
	},
	0x41: func(cpu *CPU) int {
		// LD B, C
		// fmt.Println("Decoded OPCODE: LD B, C")
		cpu.setB(cpu.reg.C)
		return 0
	},
	0x42: func(cpu *CPU) int {
		// LD B, D
		// fmt.Println("Decoded OPCODE: LD B, D")
		cpu.setB(cpu.reg.D)
		return 0
	},
	0x43: func(cpu *CPU) int {
		// LD B, E
		// fmt.Println("Decoded OPCODE: LD B, E")
		cpu.setB(cpu.reg.E)
		return 0
	},
	0x44: func(cpu *CPU) int {
		// LD B, H
		// fmt.Println("Decoded OPCODE: LD B, H")
		cpu.setB(cpu.reg.H)
		return 0
	},
	0x45: func(cpu *CPU) int {
		// LD B, L
		// fmt.Println("Decoded OPCODE: LD B, L")
		cpu.setB(cpu.reg.L)
		return 0
	},
	0x46: func(cpu *CPU) int {
		// LD B, [HL]
		// fmt.Println("Decoded OPCODE: LD B, [HL]")
		cpu.setB(cpu.mmu.read(cpu.getHL()))
		return 0
	},
	0x47: func(cpu *CPU) int {
		// LD B, A
		// fmt.Println("Decoded OPCODE: LD B, A")
		cpu.setB(cpu.reg.A)
		return 0
	},
	0x48: func(cpu *CPU) int {
		// LD C, B
		// fmt.Println("Decoded OPCODE: LD C, B")
		cpu.setC(cpu.reg.B)
		return 0
	},
	0x49: func(cpu *CPU) int {
		// LD C, C
		// fmt.Println("Decoded OPCODE: LD C, C")
		cpu.setC(cpu.reg.C)
		return 0
	},
	0x4A: func(cpu *CPU) int {
		// LD C, D
		// fmt.Println("Decoded OPCODE: LD C, D")
		cpu.setC(cpu.reg.D)
		return 0
	},
	0x4B: func(cpu *CPU) int {
		// LD C, E
		// fmt.Println("Decoded OPCODE: LD C, E")
		cpu.setC(cpu.reg.E)
		return 0
	},
	0x4C: func(cpu *CPU) int {
		// LD C, H
		// fmt.Println("Decoded OPCODE: LD C, H")
		cpu.setC(cpu.reg.H)
		return 0
	},
	0x4D: func(cpu *CPU) int {
		// LD C, L
		// fmt.Println("Decoded OPCODE: LD C, L")
		cpu.setC(cpu.reg.L)
		return 0
	},
	0x4E: func(cpu *CPU) int {
		// LD C, [HL]
		// fmt.Println("Decoded OPCODE: LD C, [HL]")
		cpu.setC(cpu.mmu.read(cpu.getHL()))
		return 0
	},
	0x4F: func(cpu *CPU) int {
		// LD C, A
		// fmt.Println("Decoded OPCODE: LD C, A")
		cpu.setC(cpu.reg.A)
		return 0
	},
	0x50: func(cpu *CPU) int {
		// LD D, B
		// fmt.Println("Decoded OPCODE: LD D, B")
		cpu.setD(cpu.reg.B)
		return 0
	},
	0x51: func(cpu *CPU) int {
		// LD D, C
		// fmt.Println("Decoded OPCODE: LD D, C")
		cpu.setD(cpu.reg.C)
		return 0
	},
	0x52: func(cpu *CPU) int {
		// LD D, D
		// fmt.Println("Decoded OPCODE: LD D, D")
		cpu.setD(cpu.reg.D)
		return 0
	},
	0x53: func(cpu *CPU) int {
		// LD D, E
		// fmt.Println("Decoded OPCODE: LD D, E")
		cpu.setD(cpu.reg.E)
		return 0
	},
	0x54: func(cpu *CPU) int {
		// LD D, H
		// fmt.Println("Decoded OPCODE: LD D, H")
		cpu.setD(cpu.reg.H)
		return 0
	},
	0x55: func(cpu *CPU) int {
		// LD D, L
		// fmt.Println("Decoded OPCODE: LD D, L")
		cpu.setD(cpu.reg.L)
		return 0
	},
	0x56: func(cpu *CPU) int {
		// LD D, [HL]
		// fmt.Println("Decoded OPCODE: LD D, [HL]")
		cpu.setD(cpu.mmu.read(cpu.getHL()))
		return 0
	},
	0x57: func(cpu *CPU) int {
		// LD D, A
		// fmt.Println("Decoded OPCODE: LD D, A")
		cpu.setD(cpu.reg.A)
		return 0
	},
	0x58: func(cpu *CPU) int {
		// LD E, B
		// fmt.Println("Decoded OPCODE: LD E, B")
		cpu.setE(cpu.reg.B)
		return 0
	},
	0x59: func(cpu *CPU) int {
		// LD E, C
		// fmt.Println("Decoded OPCODE: LD E, C")
		cpu.setE(cpu.reg.C)
		return 0
	},
	0x5A: func(cpu *CPU) int {
		// LD E, D
		// fmt.Println("Decoded OPCODE: LD E, D")
		cpu.setE(cpu.reg.D)
		return 0
	},
	0x5B: func(cpu *CPU) int {
		// LD E, E
		// fmt.Println("Decoded OPCODE: LD E, E")
		cpu.setE(cpu.reg.E)
		return 0
	},
	0x5C: func(cpu *CPU) int {
		// LD E, H
		// fmt.Println("Decoded OPCODE: LD E, H")
		cpu.setE(cpu.reg.H)
		return 0
	},
	0x5D: func(cpu *CPU) int {
		// LD E, L
		// fmt.Println("Decoded OPCODE: LD E, L")
		cpu.setE(cpu.reg.L)
		return 0
	},
	0x5E: func(cpu *CPU) int {
		// LD E, [HL]
		// fmt.Println("Decoded OPCODE: LD E, [HL]")
		cpu.setE(cpu.mmu.read(cpu.getHL()))
		return 0
	},
	0x5F: func(cpu *CPU) int {
		// LD E, A
		// fmt.Println("Decoded OPCODE: LD E, A")
		cpu.setE(cpu.reg.A)
		return 0
	},
	0x60: func(cpu *CPU) int {
		// LD H, B
		// fmt.Println("Decoded OPCODE: LD H, B")
		cpu.setH(cpu.reg.B)
		return 0
	},
	0x61: func(cpu *CPU) int {
		// LD H, C
		// fmt.Println("Decoded OPCODE: LD H, C")
		cpu.setH(cpu.reg.C)
		return 0
	},
	0x62: func(cpu *CPU) int {
		// LD H, D
		// fmt.Println("Decoded OPCODE: LD H, D")
		cpu.setH(cpu.reg.D)
		return 0
	},
	0x63: func(cpu *CPU) int {
		// LD H, E
		// fmt.Println("Decoded OPCODE: LD H, E")
		cpu.setH(cpu.reg.E)
		return 0
	},
	0x64: func(cpu *CPU) int {
		// LD H, H
		// fmt.Println("Decoded OPCODE: LD H, H")
		cpu.setH(cpu.reg.H)
		return 0
	},
	0x65: func(cpu *CPU) int {
		// LD H, L
		// fmt.Println("Decoded OPCODE: LD H, L")
		cpu.setH(cpu.reg.L)
		return 0
	},
	0x66: func(cpu *CPU) int {
		// LD H, [HL]
		// fmt.Println("Decoded OPCODE: LD H, [HL]")
		cpu.setH(cpu.mmu.read(cpu.getHL()))
		return 0
	},
	0x67: func(cpu *CPU) int {
		// LD H, A
		// fmt.Println("Decoded OPCODE: LD H, A")
		cpu.setH(cpu.reg.A)
		return 0
	},
	0x68: func(cpu *CPU) int {
		// LD L, B
		// fmt.Println("Decoded OPCODE: LD L, B")
		cpu.setL(cpu.reg.B)
		return 0
	},
	0x69: func(cpu *CPU) int {
		// LD L, C
		// fmt.Println("Decoded OPCODE: LD L, C")
		cpu.setL(cpu.reg.C)
		return 0
	},
	0x6A: func(cpu *CPU) int {
		// LD L, D
		// fmt.Println("Decoded OPCODE: LD L, D")
		cpu.setL(cpu.reg.D)
		return 0
	},
	0x6B: func(cpu *CPU) int {
		// LD L, E
		// fmt.Println("Decoded OPCODE: LD L, E")
		cpu.setL(cpu.reg.E)
		return 0
	},
	0x6C: func(cpu *CPU) int {
		// LD L, H
		// fmt.Println("Decoded OPCODE: LD L, H")
		cpu.setL(cpu.reg.H)
		return 0
	},
	0x6D: func(cpu *CPU) int {
		// LD L, L
		// fmt.Println("Decoded OPCODE: LD L, L")
		cpu.setL(cpu.reg.L)
		return 0
	},
	0x6E: func(cpu *CPU) int {
		// LD L, [HL]
		// fmt.Println("Decoded OPCODE: LD L, [HL]")
		cpu.setL(cpu.mmu.read(cpu.getHL()))
		return 0
	},
	0x6F: func(cpu *CPU) int {
		// LD L, A
		// fmt.Println("Decoded OPCODE: LD L, A")
		cpu.setL(cpu.reg.A)
		return 0
	},
	0x70: func(cpu *CPU) int {
		// LD [HL], B
		// fmt.Println("Decoded OPCODE: LD [HL], B")
		cpu.mmu.write(cpu.getHL(), cpu.reg.B)
		return 0
	},
	0x71: func(cpu *CPU) int {
		// LD [HL], C
		// fmt.Println("Decoded OPCODE: LD [HL], C")
		cpu.mmu.write(cpu.getHL(), cpu.reg.C)
		return 0
	},
	0x72: func(cpu *CPU) int {
		// LD [HL], D
		// fmt.Println("Decoded OPCODE: LD [HL], D")
		cpu.mmu.write(cpu.getHL(), cpu.reg.D)
		return 0
	},
	0x73: func(cpu *CPU) int {
		// LD [HL], E
		// fmt.Println("Decoded OPCODE: LD [HL], E")
		cpu.mmu.write(cpu.getHL(), cpu.reg.E)
		return 0
	},
	0x74: func(cpu *CPU) int {
		// LD [HL], H
		// fmt.Println("Decoded OPCODE: LD [HL], H")
		cpu.mmu.write(cpu.getHL(), cpu.reg.H)
		return 0
	},
	0x75: func(cpu *CPU) int {
		// LD [HL], L
		// fmt.Println("Decoded OPCODE: LD [HL], L")
		cpu.mmu.write(cpu.getHL(), cpu.reg.L)
		return 0
	},
	0x76: func(cpu *CPU) int {
		// HALT
		// fmt.Println("Decoded OPCODE: HALT")
		cpu.enterHaltedState()
		return 0
	},
	0x77: func(cpu *CPU) int {
		// LD [HL], A
		// fmt.Println("Decoded OPCODE: LD [HL], A")
		cpu.mmu.write(cpu.getHL(), cpu.reg.A)
		return 0
	},
	0x78: func(cpu *CPU) int {
		// LD A, B
		// fmt.Println("Decoded OPCODE: LD A, B")
		cpu.setA(cpu.reg.B)
		return 0
	},
	0x79: func(cpu *CPU) int {
		// LD A, C
		// fmt.Println("Decoded OPCODE: LD A, C")
		cpu.setA(cpu.reg.C)
		return 0
	},
	0x7A: func(cpu *CPU) int {
		// LD A, D
		// fmt.Println("Decoded OPCODE: LD A, D")
		cpu.setA(cpu.reg.D)
		return 0
	},
	0x7B: func(cpu *CPU) int {
		// LD A, E
		// fmt.Println("Decoded OPCODE: LD A, E")
		cpu.setA(cpu.reg.E)
		return 0
	},
	0x7C: func(cpu *CPU) int {
		// LD A, H
		// fmt.Println("Decoded OPCODE: LD A, H")
		cpu.setA(cpu.reg.H)
		return 0
	},
	0x7D: func(cpu *CPU) int {
		// LD A, L
		// fmt.Println("Decoded OPCODE: LD A, L")
		cpu.setA(cpu.reg.L)
		return 0
	},
	0x7E: func(cpu *CPU) int {
		// LD A, [HL]
		// fmt.Println("Decoded OPCODE: LD A, [HL]")
		cpu.setA(cpu.mmu.read(cpu.getHL()))
		return 0
	},
	0x7F: func(cpu *CPU) int {
		// LD A, A
		// fmt.Println("Decoded OPCODE: LD A, A")
		cpu.setA(cpu.reg.A)
		return 0
	},
	0x80: func(cpu *CPU) int {
		// ADD A, B
		// fmt.Println("Decoded OPCODE: ADD A, B")
		cpu.instrAddA(cpu.reg.B, false)
		return 0
	},
	0x81: func(cpu *CPU) int {
		// ADD A, C
		// fmt.Println("Decoded OPCODE: ADD A, C")
		cpu.instrAddA(cpu.reg.C, false)
		return 0
	},
	0x82: func(cpu *CPU) int {
		// ADD A, D
		// fmt.Println("Decoded OPCODE: ADD A, D")
		cpu.instrAddA(cpu.reg.D, false)
		return 0
	},
	0x83: func(cpu *CPU) int {
		// ADD A, E
		// fmt.Println("Decoded OPCODE: ADD A, E")
		cpu.instrAddA(cpu.reg.E, false)
		return 0
	},
	0x84: func(cpu *CPU) int {
		// ADD A, H
		// fmt.Println("Decoded OPCODE: ADD A, H")
		cpu.instrAddA(cpu.reg.H, false)
		return 0
	},
	0x85: func(cpu *CPU) int {
		// ADD A, L
		// fmt.Println("Decoded OPCODE: ADD A, L")
		cpu.instrAddA(cpu.reg.L, false)
		return 0
	},
	0x86: func(cpu *CPU) int {
		// ADD A, [HL]
		// fmt.Println("Decoded OPCODE: ADD A, [HL]")
		cpu.instrAddA(cpu.mmu.read(cpu.getHL()), false)
		return 0
	},
	0x87: func(cpu *CPU) int {
		// ADD A, A
		// fmt.Println("Decoded OPCODE: ADD A, A")
		cpu.instrAddA(cpu.reg.A, false)
		return 0
	},
	0x88: func(cpu *CPU) int {
		// ADC A, B
		// fmt.Println("Decoded OPCODE: ADC A, B")
		cpu.instrAddA(cpu.reg.B, true)
		return 0
	},
	0x89: func(cpu *CPU) int {
		// ADC A, C
		// fmt.Println("Decoded OPCODE: ADC A, C")
		cpu.instrAddA(cpu.reg.C, true)
		return 0
	},
	0x8A: func(cpu *CPU) int {
		// ADC A, D
		// fmt.Println("Decoded OPCODE: ADC A, D")
		cpu.instrAddA(cpu.reg.D, true)
		return 0
	},
	0x8B: func(cpu *CPU) int {
		// ADC A, E
		// fmt.Println("Decoded OPCODE: ADC A, E")
		cpu.instrAddA(cpu.reg.E, true)
		return 0
	},
	0x8C: func(cpu *CPU) int {
		// ADC A, H
		// fmt.Println("Decoded OPCODE: ADC A, H")
		cpu.instrAddA(cpu.reg.H, true)
		return 0
	},
	0x8D: func(cpu *CPU) int {
		// ADC A, L
		// fmt.Println("Decoded OPCODE: ADC A, L")
		cpu.instrAddA(cpu.reg.L, true)
		return 0
	},
	0x8E: func(cpu *CPU) int {
		// ADC A, [HL]
		// fmt.Println("Decoded OPCODE: ADC A, [HL]")
		cpu.instrAddA(cpu.mmu.read(cpu.getHL()), true)
		return 0
	},
	0x8F: func(cpu *CPU) int {
		// ADC A, A
		// fmt.Println("Decoded OPCODE: ADC A, A")
		cpu.instrAddA(cpu.reg.A, true)
		return 0
	},
	0x90: func(cpu *CPU) int {
		// SUB A, B
		// fmt.Println("Decoded OPCODE: SUB A, B")
		cpu.instrSubA(cpu.reg.B, false)
		return 0
	},
	0x91: func(cpu *CPU) int {
		// SUB A, C
		// fmt.Println("Decoded OPCODE: SUB A, C")
		cpu.instrSubA(cpu.reg.C, false)
		return 0
	},
	0x92: func(cpu *CPU) int {
		// SUB A, D
		// fmt.Println("Decoded OPCODE: SUB A, D")
		cpu.instrSubA(cpu.reg.D, false)
		return 0
	},
	0x93: func(cpu *CPU) int {
		// SUB A, E
		// fmt.Println("Decoded OPCODE: SUB A, E")
		cpu.instrSubA(cpu.reg.E, false)
		return 0
	},
	0x94: func(cpu *CPU) int {
		// SUB A, H
		// fmt.Println("Decoded OPCODE: SUB A, H")
		cpu.instrSubA(cpu.reg.H, false)
		return 0
	},
	0x95: func(cpu *CPU) int {
		// SUB A, L
		// fmt.Println("Decoded OPCODE: SUB A, L")
		cpu.instrSubA(cpu.reg.L, false)
		return 0
	},
	0x96: func(cpu *CPU) int {
		// SUB A, [HL]
		// fmt.Println("Decoded OPCODE: SUB A, [HL]")
		cpu.instrSubA(cpu.mmu.read(cpu.getHL()), false)
		return 0
	},
	0x97: func(cpu *CPU) int {
		// SUB A, A
		// fmt.Println("Decoded OPCODE: SUB A, A")
		cpu.instrSubA(cpu.reg.A, false)
		return 0
	},
	0x98: func(cpu *CPU) int {
		// SBC A, B
		// fmt.Println("Decoded OPCODE: SBC A, B")
		cpu.instrSubA(cpu.reg.B, true)
		return 0
	},
	0x99: func(cpu *CPU) int {
		// SBC A, C
		// fmt.Println("Decoded OPCODE: SBC A, C")
		cpu.instrSubA(cpu.reg.C, true)
		return 0
	},
	0x9A: func(cpu *CPU) int {
		// SBC A, D
		// fmt.Println("Decoded OPCODE: SBC A, D")
		cpu.instrSubA(cpu.reg.D, true)
		return 0
	},
	0x9B: func(cpu *CPU) int {
		// SBC A, E
		// fmt.Println("Decoded OPCODE: SBC A, E")
		cpu.instrSubA(cpu.reg.E, true)
		return 0
	},
	0x9C: func(cpu *CPU) int {
		// SBC A, H
		// fmt.Println("Decoded OPCODE: SBC A, H")
		cpu.instrSubA(cpu.reg.H, true)
		return 0
	},
	0x9D: func(cpu *CPU) int {
		// SBC A, L
		// fmt.Println("Decoded OPCODE: SBC A, L")
		cpu.instrSubA(cpu.reg.L, true)
		return 0
	},
	0x9E: func(cpu *CPU) int {
		// SBC A, [HL]
		// fmt.Println("Decoded OPCODE: SBC A, [HL]")
		cpu.instrSubA(cpu.mmu.read(cpu.getHL()), true)
		return 0
	},
	0x9F: func(cpu *CPU) int {
		// SBC A, A
		// fmt.Println("Decoded OPCODE: SBC A, A")
		cpu.instrSubA(cpu.reg.A, true)
		return 0
	},
	0xA0: func(cpu *CPU) int {
		// AND A, B
		// fmt.Println("Decoded OPCODE: AND A, B")
		cpu.instrAndA(cpu.reg.B)
		return 0
	},
	0xA1: func(cpu *CPU) int {
		// AND A, C
		// fmt.Println("Decoded OPCODE: AND A, C")
		cpu.instrAndA(cpu.reg.C)
		return 0
	},
	0xA2: func(cpu *CPU) int {
		// AND A, D
		// fmt.Println("Decoded OPCODE: AND A, D")
		cpu.instrAndA(cpu.reg.D)
		return 0
	},
	0xA3: func(cpu *CPU) int {
		// AND A, E
		// fmt.Println("Decoded OPCODE: AND A, E")
		cpu.instrAndA(cpu.reg.E)
		return 0
	},
	0xA4: func(cpu *CPU) int {
		// AND A, H
		// fmt.Println("Decoded OPCODE: AND A, H")
		cpu.instrAndA(cpu.reg.H)
		return 0
	},
	0xA5: func(cpu *CPU) int {
		// AND A, L
		// fmt.Println("Decoded OPCODE: AND A, L")
		cpu.instrAndA(cpu.reg.L)
		return 0
	},
	0xA6: func(cpu *CPU) int {
		// AND A, [HL]
		// fmt.Println("Decoded OPCODE: AND A, [HL]")
		cpu.instrAndA(cpu.mmu.read(cpu.getHL()))
		return 0
	},
	0xA7: func(cpu *CPU) int {
		// AND A, A
		// fmt.Println("Decoded OPCODE: AND A, A")
		cpu.instrAndA(cpu.reg.A)
		return 0
	},
	0xA8: func(cpu *CPU) int {
		// XOR A, B
		// fmt.Println("Decoded OPCODE: XOR A, B")
		cpu.instrXorA(cpu.reg.B)
		return 0
	},
	0xA9: func(cpu *CPU) int {
		// XOR A, C
		// fmt.Println("Decoded OPCODE: XOR A, C")
		cpu.instrXorA(cpu.reg.C)
		return 0
	},
	0xAA: func(cpu *CPU) int {
		// XOR A, D
		// fmt.Println("Decoded OPCODE: XOR A, D")
		cpu.instrXorA(cpu.reg.D)
		return 0
	},
	0xAB: func(cpu *CPU) int {
		// XOR A, E
		// fmt.Println("Decoded OPCODE: XOR A, E")
		cpu.instrXorA(cpu.reg.E)
		return 0
	},
	0xAC: func(cpu *CPU) int {
		// XOR A, H
		// fmt.Println("Decoded OPCODE: XOR A, H")
		cpu.instrXorA(cpu.reg.H)
		return 0
	},
	0xAD: func(cpu *CPU) int {
		// XOR A, L
		// fmt.Println("Decoded OPCODE: XOR A, L")
		cpu.instrXorA(cpu.reg.L)
		return 0
	},
	0xAE: func(cpu *CPU) int {
		// XOR A, [HL]
		// fmt.Println("Decoded OPCODE: XOR A, [HL]")
		cpu.instrXorA(cpu.mmu.read(cpu.getHL()))
		return 0
	},
	0xAF: func(cpu *CPU) int {
		// XOR A, A
		// fmt.Println("Decoded OPCODE: XOR A, A")
		cpu.instrXorA(cpu.reg.A)
		return 0
	},
	0xB0: func(cpu *CPU) int {
		// OR A, B
		// fmt.Println("Decoded OPCODE: OR A, B")
		cpu.instrOrA(cpu.reg.B)
		return 0
	},
	0xB1: func(cpu *CPU) int {
		// OR A, C
		// fmt.Println("Decoded OPCODE: OR A, C")
		cpu.instrOrA(cpu.reg.C)
		return 0
	},
	0xB2: func(cpu *CPU) int {
		// OR A, D
		// fmt.Println("Decoded OPCODE: OR A, D")
		cpu.instrOrA(cpu.reg.D)
		return 0
	},
	0xB3: func(cpu *CPU) int {
		// OR A, E
		// fmt.Println("Decoded OPCODE: OR A, E")
		cpu.instrOrA(cpu.reg.E)
		return 0
	},
	0xB4: func(cpu *CPU) int {
		// OR A, H
		// fmt.Println("Decoded OPCODE: OR A, H")
		cpu.instrOrA(cpu.reg.H)
		return 0
	},
	0xB5: func(cpu *CPU) int {
		// OR A, L
		// fmt.Println("Decoded OPCODE: OR A, L")
		cpu.instrOrA(cpu.reg.L)
		return 0
	},
	0xB6: func(cpu *CPU) int {
		// OR A, [HL]
		// fmt.Println("Decoded OPCODE: OR A, [HL]")
		cpu.instrOrA(cpu.mmu.read(cpu.getHL()))
		return 0
	},
	0xB7: func(cpu *CPU) int {
		// OR A, A
		// fmt.Println("Decoded OPCODE: OR A, A")
		cpu.instrOrA(cpu.reg.A)
		return 0
	},
	0xB8: func(cpu *CPU) int {
		// CP A, B
		// fmt.Println("Decoded OPCODE: CP A, B")
		cpu.instrCpA(cpu.reg.B)
		return 0
	},
	0xB9: func(cpu *CPU) int {
		// CP A, C
		// fmt.Println("Decoded OPCODE: CP A, C")
		cpu.instrCpA(cpu.reg.C)
		return 0
	},
	0xBA: func(cpu *CPU) int {
		// CP A, D
		// fmt.Println("Decoded OPCODE: CP A, D")
		cpu.instrCpA(cpu.reg.D)
		return 0
	},
	0xBB: func(cpu *CPU) int {
		// CP A, E
		// fmt.Println("Decoded OPCODE: CP A, E")
		cpu.instrCpA(cpu.reg.E)
		return 0
	},
	0xBC: func(cpu *CPU) int {
		// CP A, H
		// fmt.Println("Decoded OPCODE: CP A, H")
		cpu.instrCpA(cpu.reg.H)
		return 0
	},
	0xBD: func(cpu *CPU) int {
		// CP A, L
		// fmt.Println("Decoded OPCODE: CP A, L")
		cpu.instrCpA(cpu.reg.L)
		return 0
	},
	0xBE: func(cpu *CPU) int {
		// CP A, [HL]
		// fmt.Println("Decoded OPCODE: CP A, [HL]")
		cpu.instrCpA(cpu.mmu.read(cpu.getHL()))
		return 0
	},
	0xBF: func(cpu *CPU) int {
		// CP A, A
		// fmt.Println("Decoded OPCODE: CP A, A")
		cpu.instrCpA(cpu.reg.A)
		return 0
	},
	0xC0: func(cpu *CPU) int {
		// RET NZ
		// fmt.Println("Decoded OPCODE: RET NZ")
		if !cpu.zFlag() {
			cpu.instrRet()
			return 12
		}
		return 0
	},
	0xC1: func(cpu *CPU) int {
		// POP BC
		// fmt.Println("Decoded OPCODE: POP BC")
		cpu.setBC(cpu.popStack())
		return 0
	},
	0xC2: func(cpu *CPU) int {
		// JP NZ, u16
		// fmt.Println("Decoded OPCODE: JP NZ, u16")
		jumpAddress := cpu.nextPC16()
		if !cpu.zFlag() {
			cpu.instrJump(jumpAddress)
			return 4
		}
		return 0
	},
	0xC3: func(cpu *CPU) int {
		// JP u16
		// fmt.Println("Decoded OPCODE: JP u16")
		cpu.instrJump(cpu.nextPC16())
		return 0
	},
	0xC4: func(cpu *CPU) int {
		// CALL NZ, u16
		// fmt.Println("Decoded OPCODE: CALL NZ, u16")
		jumpAddress := cpu.nextPC16()
		if !cpu.zFlag() {
			cpu.instrCall(jumpAddress)
			return 12
		}
		return 0
	},
	0xC5: func(cpu *CPU) int {
		// PUSH BC
		// fmt.Println("Decoded OPCODE: PUSH BC")
		cpu.pushStack(cpu.getBC())
		return 0
	},
	0xC6: func(cpu *CPU) int {
		// ADD A, u8
		// fmt.Println("Decoded OPCODE: ADD A, u8")
		cpu.instrAddA(cpu.nextPC(), false)
		return 0
	},
	0xC7: func(cpu *CPU) int {
		// RST $00
		// fmt.Println("Decoded OPCODE: RST $00")
		cpu.instrCall(0x0000)
		return 0
	},
	0xC8: func(cpu *CPU) int {
		// RET Z
		// fmt.Println("Decoded OPCODE: RET Z")
		if cpu.zFlag() {
			cpu.instrRet()
			return 12
		}
		return 0
	},
	0xC9: func(cpu *CPU) int {
		// RET
		// fmt.Println("Decoded OPCODE: RET")
		cpu.instrRet()
		return 0
	},
	0xCA: func(cpu *CPU) int {
		// JP Z, u16
		// fmt.Println("Decoded OPCODE: JP Z, u16")
		jumpAddress := cpu.nextPC16()
		if cpu.zFlag() {
			cpu.instrJump(jumpAddress)
			return 4
		}
		return 0
	},
	0xCB: func(cpu *CPU) int {
		// PREFIX
		// fmt.Println("Decoded OPCODE: PREFIX")

		cbOpcode := cpu.nextPC()
		// fmt.Printf("CB OPCODE: 0x%02x\n", cbOpcode)

		cpu.cbInstructions[cbOpcode]()
		return cbInstrBaseTicks[cbOpcode]
	},
	0xCC: func(cpu *CPU) int {
		// CALL Z, u16
		// fmt.Println("Decoded OPCODE: CALL Z, u16")
		jumpAddress := cpu.nextPC16()
		if cpu.zFlag() {
			cpu.instrCall(jumpAddress)
			return 12
		}
		return 0
	},
	0xCD: func(cpu *CPU) int {
		// CALL u16
		// fmt.Println("Decoded OPCODE: CALL u16")
		cpu.instrCall(cpu.nextPC16())
		return 0
	},
	0xCE: func(cpu *CPU) int {
		// ADC A, u8
		// fmt.Println("Decoded OPCODE: ADC A, u8")
		cpu.instrAddA(cpu.nextPC(), true)
		return 0
	},
	0xCF: func(cpu *CPU) int {
		// RST $08
		// fmt.Println("Decoded OPCODE: RST $08")
		cpu.instrCall(0x0008)
		return 0
	},
	0xD0: func(cpu *CPU) int {
		// RET NC
		// fmt.Println("Decoded OPCODE: RET NC")
		if !cpu.cFlag() {
			cpu.instrRet()
			return 12
		}
		return 0
	},
	0xD1: func(cpu *CPU) int {
		// POP DE
		// fmt.Println("Decoded OPCODE: POP DE")
		cpu.setDE(cpu.popStack())
		return 0
	},
	0xD2: func(cpu *CPU) int {
		// JP NC, u16
		// fmt.Println("Decoded OPCODE: JP NC, u16")
		jumpAddress := cpu.nextPC16()
		if !cpu.cFlag() {
			cpu.instrJump(jumpAddress)
			return 4
		}
		return 0
	},
	0xD3: func(cpu *CPU) int {
		// ILLEGAL OPCODE 0xD3
		// fmt.Println("ILLEGAL OPCODE: 0xD3")
		return 0
	},
	0xD4: func(cpu *CPU) int {
		// CALL NC, u16
		// fmt.Println("Decoded OPCODE: CALL NC, u16")
		jumpAddress := cpu.nextPC16()
		if !cpu.cFlag() {
			cpu.instrCall(jumpAddress)
			return 12
		}
		return 0
	},
	0xD5: func(cpu *CPU) int {
		// PUSH DE
		// fmt.Println("Decoded OPCODE: PUSH DE")
		cpu.pushStack(cpu.getDE())
		return 0
	},
	0xD6: func(cpu *CPU) int {
		// SUB A, u8
		// fmt.Println("Decoded OPCODE: SUB A, u8")
		cpu.instrSubA(cpu.nextPC(), false)
		return 0
	},
	0xD7: func(cpu *CPU) int {
		// RST $10
		// fmt.Println("Decoded OPCODE: RST $10")
		cpu.instrCall(0x0010)
		return 0
	},
	0xD8: func(cpu *CPU) int {
		// RET C
		// fmt.Println("Decoded OPCODE: RET C")
		if cpu.cFlag() {
			cpu.instrRet()
			return 12
		}
		return 0
	},
	0xD9: func(cpu *CPU) int {
		// RETI
		// fmt.Println("Decoded OPCODE: RETI")
		cpu.instrRet()
		cpu.setIMEDelay()
		return 0
	},
	0xDA: func(cpu *CPU) int {
		// JP C, u16
		// fmt.Println("Decoded OPCODE: JP C, u16")
		jumpAddress := cpu.nextPC16()
		if cpu.cFlag() {
			cpu.instrJump(jumpAddress)
			return 4
		}
		return 0
	},
	0xDB: func(cpu *CPU) int {
		// ILLEGAL OPCODE 0xDB
		// fmt.Println("ILLEGAL OPCODE: 0xDB")
		return 0
	},
	0xDC: func(cpu *CPU) int {
		// CALL C, u16
		// fmt.Println("Decoded OPCODE: CALL C, u16")
		jumpAddress := cpu.nextPC16()
		if cpu.cFlag() {
			cpu.instrCall(jumpAddress)
			return 12
		}
		return 0
	},
	0xDD: func(cpu *CPU) int {
		// ILLEGAL OPCODE 0xDD
		// fmt.Println("ILLEGAL OPCODE: 0xDD")
		return 0
	},
	0xDE: func(cpu *CPU) int {
		// SBC A, u8
		// fmt.Println("Decoded OPCODE: SBC A, u8")
		cpu.instrSubA(cpu.nextPC(), true)
		return 0
	},
	0xDF: func(cpu *CPU) int {
		// RST $18
		// fmt.Println("Decoded OPCODE: RST $18")
		cpu.instrCall(0x0018)
		return 0
	},
	0xE0: func(cpu *CPU) int {
		// LD (FF00+u8), A
		// fmt.Println("Decoded OPCODE: LD (FF00+u8), A")
		address := 0xFF00 + uint16(cpu.nextPC())
		cpu.mmu.write(address, cpu.reg.A)
		return 0
	},
	0xE1: func(cpu *CPU) int {
		// POP HL
		// fmt.Println("Decoded OPCODE: POP HL")
		cpu.setHL(cpu.popStack())
		return 0
	},
	0xE2: func(cpu *CPU) int {
		// LD (FF00+C), A
		// fmt.Println("Decoded OPCODE: LD (FF00+C), A")
		address := 0xFF00 + uint16(cpu.reg.C)
		cpu.mmu.write(address, cpu.reg.A)
		return 0
	},
	0xE3: func(cpu *CPU) int {
		// ILLEGAL OPCODE 0xE3
		// fmt.Println("ILLEGAL OPCODE: 0xE3")
		return 0
	},
	0xE4: func(cpu *CPU) int {
		// ILLEGAL OPCODE 0xE4
		// fmt.Println("ILLEGAL OPCODE: 0xE4")
		return 0
	},
	0xE5: func(cpu *CPU) int {
		// PUSH HL
		// fmt.Println("Decoded OPCODE: PUSH HL")
		cpu.pushStack(cpu.getHL())
		return 0
	},
	0xE6: func(cpu *CPU) int {
		// AND A, u8
		// fmt.Println("Decoded OPCODE: AND A, u8")
		cpu.instrAndA(cpu.nextPC())
		return 0
	},
	0xE7: func(cpu *CPU) int {
		// RST $20
		// fmt.Println("Decoded OPCODE: RST $20")
		cpu.instrCall(0x0020)
		return 0
	},
	0xE8: func(cpu *CPU) int {
		// ADD SP, i8
		// fmt.Println("Decoded OPCODE: ADD SP, i8")
		cpu.instrAdd16Signed(cpu.setSP, cpu.reg.SP, int8(cpu.nextPC()))
		return 0
	},
	0xE9: func(cpu *CPU) int {
		// JP HL
		// fmt.Println("Decoded OPCODE: JP HL")
		cpu.instrJump(cpu.getHL())
		return 0
	},
	0xEA: func(cpu *CPU) int {
		// LD [u16], A
		// fmt.Println("Decoded OPCODE: LD [u16], A")
		cpu.mmu.write(cpu.nextPC16(), cpu.reg.A)
		return 0
	},
	0xEB: func(cpu *CPU) int {
		// ILLEGAL OPCODE 0xEB
		// fmt.Println("ILLEGAL OPCODE: 0xEB")
		return 0
	},
	0xEC: func(cpu *CPU) int {
		// ILLEGAL OPCODE 0xEC
		// fmt.Println("ILLEGAL OPCODE: 0xEC")
		return 0
	},
	0xED: func(cpu *CPU) int {
		// ILLEGAL OPCODE 0xED
		// fmt.Println("ILLEGAL OPCODE: 0xED")
		return 0
	},
	0xEE: func(cpu *CPU) int {
		// XOR A, u8
		// fmt.Println("Decoded OPCODE: XOR A, u8")
		cpu.instrXorA(cpu.nextPC())
		return 0
	},
	0xEF: func(cpu *CPU) int {
		// RST $28
		// fmt.Println("Decoded OPCODE: RST $28")
		cpu.instrCall(0x0028)
		return 0
	},
	0xF0: func(cpu *CPU) int {
		// LD A, (FF00+u8)
		// fmt.Println("Decoded OPCODE: LD A, (FF00+u8)")
		address := 0xFF00 + uint16(cpu.nextPC())
		cpu.setA(cpu.mmu.read(address))
		return 0
	},
	0xF1: func(cpu *CPU) int {
		// POP AF
		// fmt.Println("Decoded OPCODE: POP AF")
		cpu.setAF(cpu.popStack())
		return 0
	},
	0xF2: func(cpu *CPU) int {
		// LD A, (FF00+C)
		// fmt.Println("Decoded OPCODE: LD A, (FF00+C)")
		address := 0xFF00 + uint16(cpu.reg.C)
		cpu.setA(cpu.mmu.read(address))
		return 0
	},
	0xF3: func(cpu *CPU) int {
		// DI
		// fmt.Println("Decoded OPCODE: DI")
		cpu.clearIME()
		return 0
	},
	0xF4: func(cpu *CPU) int {
		// ILLEGAL OPCODE 0xF4
		// fmt.Println("ILLEGAL OPCODE: 0xF4")
		return 0
	},
	0xF5: func(cpu *CPU) int {
		// PUSH AF
		// fmt.Println("Decoded OPCODE: PUSH AF")
		cpu.pushStack(cpu.getAF())
		return 0
	},
	0xF6: func(cpu *CPU) int {
		// OR A, u8
		// fmt.Println("Decoded OPCODE: OR A, u8")
		cpu.instrOrA(cpu.nextPC())
		return 0
	},
	0xF7: func(cpu *CPU) int {
		// RST $30
		// fmt.Println("Decoded OPCODE: RST $30")
		cpu.instrCall(0x0030)
		return 0
	},
	0xF8: func(cpu *CPU) int {
		// LD HL, SP + e8
		// fmt.Println("Decoded OPCODE: LD HL, SP + e8")
		cpu.instrAdd16Signed(cpu.setHL, cpu.reg.SP, int8(cpu.nextPC()))
		return 0
	},
	0xF9: func(cpu *CPU) int {
		// LD SP, HL
		// fmt.Println("Decoded OPCODE: LD SP, HL")
		cpu.setSP(cpu.getHL())
		return 0
	},
	0xFA: func(cpu *CPU) int {
		// LD A, [u16]
		// fmt.Println("Decoded OPCODE: LD A, [u16]")
		cpu.setA(cpu.mmu.read(cpu.nextPC16()))
		return 0
	},
	0xFB: func(cpu *CPU) int {
		// EI
		// fmt.Println("Decoded OPCODE: EI")
		cpu.setIMEDelay()
		return 0
	},
	0xFC: func(cpu *CPU) int {
		// ILLEGAL OPCODE 0xFC
		// fmt.Println("ILLEGAL OPCODE: 0xFC")
		return 0
	},
	0xFD: func(cpu *CPU) int {
		// ILLEGAL OPCODE 0xFD
		// fmt.Println("ILLEGAL OPCODE: 0xFD")
		return 0
	},
	0xFE: func(cpu *CPU) int {
		// CP A, u8
		// fmt.Println("Decoded OPCODE: CP A, u8")
		cpu.instrCpA(cpu.nextPC())
		return 0
	},
	0xFF: func(cpu *CPU) int {
		// RST $38
		// fmt.Println("Decoded OPCODE: RST $38")
		if cpu.mmu.read(0x0038) == 0xFF {
			log.Fatal("RST $38 loop")
		}
		cpu.instrCall(0x0038)
		return 0
	},
}
