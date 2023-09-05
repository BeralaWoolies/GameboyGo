package gb

import (
	"fmt"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
)

var instrClockTicks = [0x100]int{
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

var cbInstrClockTicks = [0x100]int{
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x04, 0x02,
}

func (gb *Gameboy) nextPC() uint8 {
	data := gb.mmu.read(gb.cpu.reg.PC)
	gb.cpu.reg.PC++
	return data
}

func (gb *Gameboy) nextPC16() uint16 {
	loByte := gb.nextPC()
	hiByte := gb.nextPC()
	return uint16(hiByte)<<8 | uint16(loByte)
}

func (gb *Gameboy) instrInc(setHandler func(result uint8), val uint8) {
	result := val + 1
	setHandler(result)

	gb.cpu.setZFlag(result == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag((val&0xF)+(1&0xF) > 0xF)
}

func (gb *Gameboy) instrDec(setHandler func(result uint8), val uint8) {
	result := val - 1
	setHandler(result)

	gb.cpu.setZFlag(result == 0)
	gb.cpu.setNFlag(true)
	gb.cpu.setHFlag(val&0xF == 0)
}

func (gb *Gameboy) instrAddA(rhs uint8, addCarry bool) {
	var carry int16 = 0
	if gb.cpu.cFlag() && addCarry {
		carry = 1
	}

	lhs := gb.cpu.reg.A
	result := int16(lhs) + int16(rhs) + carry
	gb.cpu.setA(uint8(result))

	gb.cpu.setZFlag(gb.cpu.reg.A == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag((lhs&0xF)+(rhs&0xF)+uint8(carry) > 0xF)
	gb.cpu.setCFlag(result > 0xFF)
}

func (gb *Gameboy) instrSubA(rhs uint8, subCarry bool) {
	var carry int16 = 0
	if gb.cpu.cFlag() && subCarry {
		carry = 1
	}

	lhs := gb.cpu.reg.A
	result := int16(lhs) - int16(rhs) - carry
	gb.cpu.setA(uint8(result))

	gb.cpu.setZFlag(gb.cpu.reg.A == 0)
	gb.cpu.setNFlag(true)
	gb.cpu.setHFlag(int16(lhs&0xF)-int16(rhs&0xF)-carry < 0)
	gb.cpu.setCFlag(result < 0)
}

func (gb *Gameboy) instrAdd16(setHandler func(result uint16), val1 uint16, val2 uint16) {
	result := int32(val1) + int32(val2)
	setHandler(uint16(result))

	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag(int32(val1&0xFFF) > (result & 0xFFF))
	gb.cpu.setCFlag(result > 0xFFFF)
}

func (gb *Gameboy) instrAdd16Signed(setHandler func(result uint16), val1 uint16, val2 int8) {
	result := uint16(int32(val1) + int32(val2))
	setHandler(result)

	carryBits := val1 ^ uint16(val2) ^ result

	gb.cpu.setZFlag(false)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag((carryBits & 0x10) == 0x10)
	gb.cpu.setCFlag((carryBits & 0x100) == 0x100)
}

func (gb *Gameboy) instrAndA(rhs uint8) {
	gb.cpu.setA(gb.cpu.reg.A & rhs)

	gb.cpu.setZFlag(gb.cpu.reg.A == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag(true)
	gb.cpu.setCFlag(false)
}

func (gb *Gameboy) instrXorA(rhs uint8) {
	gb.cpu.setA(gb.cpu.reg.A ^ rhs)

	gb.cpu.setZFlag(gb.cpu.reg.A == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag(false)
	gb.cpu.setCFlag(false)
}

func (gb *Gameboy) instrOrA(rhs uint8) {
	gb.cpu.setA(gb.cpu.reg.A | rhs)

	gb.cpu.setZFlag(gb.cpu.reg.A == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag(false)
	gb.cpu.setCFlag(false)
}

func (gb *Gameboy) instrCpA(rhs uint8) {
	lhs := gb.cpu.reg.A
	cmpResult := lhs - rhs

	gb.cpu.setZFlag(cmpResult == 0)
	gb.cpu.setNFlag(true)
	gb.cpu.setHFlag((rhs & 0xF) > (lhs & 0xF))
	gb.cpu.setCFlag(rhs > lhs)
}

func (gb *Gameboy) instrJump(jumpAddress uint16) {
	gb.cpu.setPC(jumpAddress)
}

func (gb *Gameboy) instrRet() {
	gb.cpu.setPC(gb.popStack())
}

func (gb *Gameboy) instrCall(jumpAddress uint16) {
	gb.pushStack(gb.cpu.reg.PC)
	gb.cpu.setPC(jumpAddress)
}

var instructions = [0x100]func(gb *Gameboy){
	0x00: func(gb *Gameboy) {
		// NOP
		fmt.Println("Decoded OPCODE: NOP")
		return
	},
	0x01: func(gb *Gameboy) {
		// LD BC, u16
		fmt.Println("Decoded OPCODE: LD BC, u16")
		gb.cpu.setBC(gb.nextPC16())
	},
	0x02: func(gb *Gameboy) {
		// LD [BC], A
		fmt.Println("Decoded OPCODE: LD [BC], A")
		gb.mmu.write(gb.cpu.getBC(), gb.cpu.reg.A)
	},
	0x03: func(gb *Gameboy) {
		// INC BC
		fmt.Println("Decoded OPCODE: INC BC")
		gb.cpu.setBC(gb.cpu.getBC() + 1)
	},
	0x04: func(gb *Gameboy) {
		// INC B
		fmt.Println("Decoded OPCODE: INC B")
		gb.instrInc(gb.cpu.setB, gb.cpu.reg.B)
	},
	0x05: func(gb *Gameboy) {
		// DEC B
		fmt.Println("Decoded OPCODE: DEC B")
		gb.instrDec(gb.cpu.setB, gb.cpu.reg.B)
	},
	0x06: func(gb *Gameboy) {
		// LD B, u8
		fmt.Println("Decoded OPCODE: LD B, u8")
		gb.cpu.setB(gb.nextPC())
	},
	0x07: func(gb *Gameboy) {
		// RLCA
		fmt.Println("Decoded OPCODE: RLCA")
		val := gb.cpu.reg.A
		gb.cpu.setA(uint8((val << 1)) | (val >> 7))

		gb.cpu.setZFlag(false)
		gb.cpu.setNFlag(false)
		gb.cpu.setHFlag(false)
		gb.cpu.setCFlag(val >= 0x80)
	},
	0x08: func(gb *Gameboy) {
		// LD [u16], SP
		fmt.Println("Decoded OPCODE: LD [u16], SP")
		address := gb.nextPC16()
		gb.mmu.write(address, bits.LoByte(gb.cpu.reg.SP))
		gb.mmu.write(address+1, bits.HiByte(gb.cpu.reg.SP))
	},
	0x09: func(gb *Gameboy) {
		// ADD HL, BC
		fmt.Println("Decoded OPCODE: ADD HL, BC")
		gb.instrAdd16(gb.cpu.setHL, gb.cpu.getHL(), gb.cpu.getBC())
	},
	0x0A: func(gb *Gameboy) {
		// LD A, [BC]
		fmt.Println("Decoded OPCODE: LD A, [BC]")
		gb.cpu.setA(gb.mmu.read(gb.cpu.getBC()))
	},
	0x0B: func(gb *Gameboy) {
		// DEC BC
		fmt.Println("Decoded OPCODE: DEC BC")
		gb.cpu.setBC(gb.cpu.getBC() - 1)
	},
	0x0C: func(gb *Gameboy) {
		// INC C
		fmt.Println("Decoded OPCODE: INC C")
		gb.instrInc(gb.cpu.setC, gb.cpu.reg.C)
	},
	0x0D: func(gb *Gameboy) {
		// DEC C
		fmt.Println("Decoded OPCODE: DEC C")
		gb.instrDec(gb.cpu.setC, gb.cpu.reg.C)
	},
	0x0E: func(gb *Gameboy) {
		// LD C, u8
		fmt.Println("Decoded OPCODE: LD C, u8")
		gb.cpu.setC(gb.nextPC())
	},
	0x0F: func(gb *Gameboy) {
		// RRCA
		fmt.Println("Decoded OPCODE: RRCA")
		val := gb.cpu.reg.A
		gb.cpu.setA((val >> 1) | ((val & 1) << 7))

		gb.cpu.setZFlag(false)
		gb.cpu.setNFlag(false)
		gb.cpu.setHFlag(false)
		gb.cpu.setCFlag(gb.cpu.reg.A >= 0x80)
	},
	0x10: func(gb *Gameboy) {
		// STOP
		fmt.Println("Decoded OPCODE: STOP")
		gb.cpu.halted = true
		gb.nextPC()
	},
	0x11: func(gb *Gameboy) {
		// LD DE, u16
		fmt.Println("Decoded OPCODE: LD DE, u16")
		gb.cpu.setDE(gb.nextPC16())
	},
	0x12: func(gb *Gameboy) {
		// LD [DE], A
		fmt.Println("Decoded OPCODE: LD [DE], A")
		gb.mmu.write(gb.cpu.getDE(), gb.cpu.reg.A)
	},
	0x13: func(gb *Gameboy) {
		// INC DE
		fmt.Println("Decoded OPCODE: INC DE")
		gb.cpu.setDE(gb.cpu.getDE() + 1)
	},
	0x14: func(gb *Gameboy) {
		// INC D
		fmt.Println("Decoded OPCODE: INC D")
		gb.instrInc(gb.cpu.setD, gb.cpu.reg.D)
	},
	0x15: func(gb *Gameboy) {
		// DEC D
		fmt.Println("Decoded OPCODE: DEC D")
		gb.instrDec(gb.cpu.setD, gb.cpu.reg.D)
	},
	0x16: func(gb *Gameboy) {
		// LD D, u8
		fmt.Println("Decoded OPCODE: LD D, u8")
		gb.cpu.setD(gb.nextPC())
	},
	0x17: func(gb *Gameboy) {
		// RLA
		fmt.Println("Decoded OPCODE: RLA")
		var carry uint8 = 0
		if gb.cpu.cFlag() {
			carry = 1
		}

		val := gb.cpu.reg.A
		gb.cpu.setA(uint8(val<<1) | carry)

		gb.cpu.setZFlag(false)
		gb.cpu.setNFlag(false)
		gb.cpu.setHFlag(false)
		gb.cpu.setCFlag(val >= 0x80)
	},
	0x18: func(gb *Gameboy) {
		// JR i8
		fmt.Println("Decoded OPCODE: JR i8")
		jumpAddress := int32(gb.cpu.reg.PC) + int32(int8(gb.nextPC()))
		gb.instrJump(uint16(jumpAddress))
	},
	0x19: func(gb *Gameboy) {
		// ADD HL, DE
		fmt.Println("Decoded OPCODE: ADD HL, DE")
		gb.instrAdd16(gb.cpu.setHL, gb.cpu.getHL(), gb.cpu.getDE())
	},
	0x1A: func(gb *Gameboy) {
		// LD A, [DE]
		fmt.Println("Decoded OPCODE: LD A, [DE]")
		gb.cpu.setA(gb.mmu.read(gb.cpu.getDE()))
	},
	0x1B: func(gb *Gameboy) {
		// DEC DE
		fmt.Println("Decoded OPCODE: DEC DE")
		gb.cpu.setDE(gb.cpu.getDE() - 1)
	},
	0x1C: func(gb *Gameboy) {
		// INC E
		fmt.Println("Decoded OPCODE: INC E")
		gb.instrInc(gb.cpu.setE, gb.cpu.reg.E)
	},
	0x1D: func(gb *Gameboy) {
		// DEC E
		fmt.Println("Decoded OPCODE: DEC E")
		gb.instrDec(gb.cpu.setE, gb.cpu.reg.E)
	},
	0x1E: func(gb *Gameboy) {
		// LD E, u8
		fmt.Println("Decoded OPCODE: LD E, u8")
		gb.cpu.setE(gb.nextPC())
	},
	0x1F: func(gb *Gameboy) {
		// RRA
		fmt.Println("Decoded OPCODE: RRA")
		var carry uint8 = 0
		if gb.cpu.cFlag() {
			carry = 0x80
		}

		val := gb.cpu.reg.A
		gb.cpu.setA(uint8(val>>1) | carry)

		gb.cpu.setZFlag(false)
		gb.cpu.setNFlag(false)
		gb.cpu.setHFlag(false)
		gb.cpu.setCFlag((val & 1) == 1)
	},
	0x20: func(gb *Gameboy) {
		// JR NZ, i8
		fmt.Println("Decoded OPCODE: JR NZ, i8")
		offset := int8(gb.nextPC())
		if !gb.cpu.zFlag() {
			jumpAddress := int32(gb.cpu.reg.PC) + int32(offset)
			gb.instrJump(uint16(jumpAddress))
			gb.cpu.ticks += 4
		}
	},
	0x21: func(gb *Gameboy) {
		// LD HL, u16
		fmt.Println("Decoded OPCODE: LD HL, u16")
		gb.cpu.setHL(gb.nextPC16())
	},
	0x22: func(gb *Gameboy) {
		// LD [HL+], A
		fmt.Println("Decoded OPCODE: LD [HL+], A")
		gb.mmu.write(gb.cpu.getHL(), gb.cpu.reg.A)
		gb.cpu.setHL(gb.cpu.getHL() + 1)
	},
	0x23: func(gb *Gameboy) {
		// INC HL
		fmt.Println("Decoded OPCODE: INC HL")
		gb.cpu.setHL(gb.cpu.getHL() + 1)
	},
	0x24: func(gb *Gameboy) {
		// INC H
		fmt.Println("Decoded OPCODE: INC H")
		gb.instrInc(gb.cpu.setH, gb.cpu.reg.H)
	},
	0x25: func(gb *Gameboy) {
		// DEC H
		fmt.Println("Decoded OPCODE: DEC H")
		gb.instrDec(gb.cpu.setH, gb.cpu.reg.H)
	},
	0x26: func(gb *Gameboy) {
		// LD H, u8
		fmt.Println("Decoded OPCODE: LD H, u8")
		gb.cpu.setH(gb.nextPC())
	},
	0x27: func(gb *Gameboy) {
		// DAA
		fmt.Println("Decoded OPCODE: DAA")
		if !gb.cpu.nFlag() {
			if gb.cpu.cFlag() || gb.cpu.reg.A > 0x99 {
				gb.cpu.reg.A += 0x60
				gb.cpu.setCFlag(true)
			}
			if gb.cpu.hFlag() || (gb.cpu.reg.A&0xF) > 0x09 {
				gb.cpu.reg.A += 0x06
			}
		} else {
			if gb.cpu.cFlag() {
				gb.cpu.reg.A -= 0x60
			}
			if gb.cpu.hFlag() {
				gb.cpu.reg.A -= 0x06
			}
		}

		gb.cpu.setZFlag(gb.cpu.reg.A == 0)
		gb.cpu.setHFlag(false)
	},
	0x28: func(gb *Gameboy) {
		// JR Z, i8
		fmt.Println("Decoded OPCODE: JR Z, i8")
		offset := int8(gb.nextPC())
		if gb.cpu.zFlag() {
			jumpAddress := int32(gb.cpu.reg.PC) + int32(offset)
			gb.instrJump(uint16(jumpAddress))
			gb.cpu.ticks += 4
		}
	},
	0x29: func(gb *Gameboy) {
		// ADD HL, HL
		fmt.Println("Decoded OPCODE: ADD HL, HL")
		gb.instrAdd16(gb.cpu.setHL, gb.cpu.getHL(), gb.cpu.getHL())
	},
	0x2A: func(gb *Gameboy) {
		// LD A, [HL+]
		fmt.Println("Decoded OPCODE: LD A, [HL+]")
		gb.cpu.setA(gb.mmu.read(gb.cpu.getHL()))
		gb.cpu.setHL(gb.cpu.getHL() + 1)
	},
	0x2B: func(gb *Gameboy) {
		// DEC HL
		fmt.Println("Decoded OPCODE: DEC HL")
		gb.cpu.setHL(gb.cpu.getHL() - 1)
	},
	0x2C: func(gb *Gameboy) {
		// INC L
		fmt.Println("Decoded OPCODE: INC L")
		gb.instrInc(gb.cpu.setL, gb.cpu.reg.L)
	},
	0x2D: func(gb *Gameboy) {
		// DEC L
		fmt.Println("Decoded OPCODE: DEC L")
		gb.instrDec(gb.cpu.setL, gb.cpu.reg.L)
	},
	0x2E: func(gb *Gameboy) {
		// LD L, u8
		fmt.Println("Decoded OPCODE: LD L, u8")
		gb.cpu.setL(gb.nextPC())
	},
	0x2F: func(gb *Gameboy) {
		// CPL
		fmt.Println("Decoded OPCODE: CPL")
		gb.cpu.reg.A = ^(gb.cpu.reg.A)
		gb.cpu.setNFlag(true)
		gb.cpu.setHFlag(true)
	},
	0x30: func(gb *Gameboy) {
		// JR NC, i8
		fmt.Println("Decoded OPCODE: JR NC, i8")
		offset := int8(gb.nextPC())
		if !gb.cpu.cFlag() {
			jumpAddress := int32(gb.cpu.reg.PC) + int32(offset)
			gb.instrJump(uint16(jumpAddress))
			gb.cpu.ticks += 4
		}
	},
	0x31: func(gb *Gameboy) {
		// LD SP, u16
		fmt.Println("Decoded OPCODE: LD SP, u16")
		gb.cpu.setSP(gb.nextPC16())
	},
	0x32: func(gb *Gameboy) {
		// LD [HL-], A
		fmt.Println("Decoded OPCODE: LD [HL-], A")
		gb.mmu.write(gb.cpu.getHL(), gb.cpu.reg.A)
		gb.cpu.setHL(gb.cpu.getHL() - 1)
	},
	0x33: func(gb *Gameboy) {
		// INC SP
		fmt.Println("Decoded OPCODE: INC SP")
		gb.cpu.setSP(gb.cpu.reg.SP + 1)
	},
	0x34: func(gb *Gameboy) {
		// INC [HL]
		fmt.Println("Decoded OPCODE: INC [HL]")
		val := gb.mmu.read(gb.cpu.getHL())
		gb.instrInc(func(result uint8) { gb.mmu.write(gb.cpu.getHL(), result) }, val)
	},
	0x35: func(gb *Gameboy) {
		// DEC [HL]
		fmt.Println("Decoded OPCODE: DEC [HL]")
		val := gb.mmu.read(gb.cpu.getHL())
		gb.instrDec(func(result uint8) { gb.mmu.write(gb.cpu.getHL(), result) }, val)
	},
	0x36: func(gb *Gameboy) {
		// LD [HL], u8
		fmt.Println("Decoded OPCODE: LD [HL], u8")
		gb.mmu.write(gb.cpu.getHL(), gb.nextPC())
	},
	0x37: func(gb *Gameboy) {
		// SCF
		fmt.Println("Decoded OPCODE: SCF")
		gb.cpu.setNFlag(false)
		gb.cpu.setHFlag(false)
		gb.cpu.setCFlag(true)
	},
	0x38: func(gb *Gameboy) {
		// JR C, i8
		fmt.Println("Decoded OPCODE: JR C, i8")
		offset := int8(gb.nextPC())
		if gb.cpu.cFlag() {
			jumpAddress := int32(gb.cpu.reg.PC) + int32(offset)
			gb.instrJump(uint16(jumpAddress))
			gb.cpu.ticks += 4
		}
	},
	0x39: func(gb *Gameboy) {
		// ADD HL, SP
		fmt.Println("Decoded OPCODE: ADD HL, SP")
		gb.instrAdd16(gb.cpu.setHL, gb.cpu.getHL(), gb.cpu.reg.SP)
	},
	0x3A: func(gb *Gameboy) {
		// LD A, [HL-]
		fmt.Println("Decoded OPCODE: LD A, [HL-]")
		gb.cpu.setA(gb.mmu.read(gb.cpu.getHL()))
		gb.cpu.setHL(gb.cpu.getHL() - 1)
	},
	0x3B: func(gb *Gameboy) {
		// DEC SP
		fmt.Println("Decoded OPCODE: DEC SP")
		gb.cpu.setSP(gb.cpu.reg.SP - 1)
	},
	0x3C: func(gb *Gameboy) {
		// INC A
		fmt.Println("Decoded OPCODE: INC A")
		gb.instrInc(gb.cpu.setA, gb.cpu.reg.A)
	},
	0x3D: func(gb *Gameboy) {
		// DEC A
		fmt.Println("Decoded OPCODE: DEC A")
		gb.instrDec(gb.cpu.setA, gb.cpu.reg.A)
	},
	0x3E: func(gb *Gameboy) {
		// LD A, u8
		fmt.Println("Decoded OPCODE: LD A, u8")
		gb.cpu.setA(gb.nextPC())
	},
	0x3F: func(gb *Gameboy) {
		// CCF
		fmt.Println("Decoded OPCODE: CCF")
		gb.cpu.setNFlag(false)
		gb.cpu.setHFlag(false)
		gb.cpu.setCFlag(!gb.cpu.cFlag())
	},
	0x40: func(gb *Gameboy) {
		// LD B, B
		fmt.Println("Decoded OPCODE: LD B, B")
		gb.cpu.setB(gb.cpu.reg.B)
	},
	0x41: func(gb *Gameboy) {
		// LD B, C
		fmt.Println("Decoded OPCODE: LD B, C")
		gb.cpu.setB(gb.cpu.reg.C)
	},
	0x42: func(gb *Gameboy) {
		// LD B, D
		fmt.Println("Decoded OPCODE: LD B, D")
		gb.cpu.setB(gb.cpu.reg.D)
	},
	0x43: func(gb *Gameboy) {
		// LD B, E
		fmt.Println("Decoded OPCODE: LD B, E")
		gb.cpu.setB(gb.cpu.reg.E)
	},
	0x44: func(gb *Gameboy) {
		// LD B, H
		fmt.Println("Decoded OPCODE: LD B, H")
		gb.cpu.setB(gb.cpu.reg.H)
	},
	0x45: func(gb *Gameboy) {
		// LD B, L
		fmt.Println("Decoded OPCODE: LD B, L")
		gb.cpu.setB(gb.cpu.reg.L)
	},
	0x46: func(gb *Gameboy) {
		// LD B, [HL]
		fmt.Println("Decoded OPCODE: LD B, [HL]")
		gb.cpu.setB(gb.mmu.read(gb.cpu.getHL()))
	},
	0x47: func(gb *Gameboy) {
		// LD B, A
		fmt.Println("Decoded OPCODE: LD B, A")
		gb.cpu.setB(gb.cpu.reg.A)
	},
	0x48: func(gb *Gameboy) {
		// LD C, B
		fmt.Println("Decoded OPCODE: LD C, B")
		gb.cpu.setC(gb.cpu.reg.B)
	},
	0x49: func(gb *Gameboy) {
		// LD C, C
		fmt.Println("Decoded OPCODE: LD C, C")
		gb.cpu.setC(gb.cpu.reg.C)
	},
	0x4A: func(gb *Gameboy) {
		// LD C, D
		fmt.Println("Decoded OPCODE: LD C, D")
		gb.cpu.setC(gb.cpu.reg.D)
	},
	0x4B: func(gb *Gameboy) {
		// LD C, E
		fmt.Println("Decoded OPCODE: LD C, E")
		gb.cpu.setC(gb.cpu.reg.E)
	},
	0x4C: func(gb *Gameboy) {
		// LD C, H
		fmt.Println("Decoded OPCODE: LD C, H")
		gb.cpu.setC(gb.cpu.reg.H)
	},
	0x4D: func(gb *Gameboy) {
		// LD C, L
		fmt.Println("Decoded OPCODE: LD C, L")
		gb.cpu.setC(gb.cpu.reg.L)
	},
	0x4E: func(gb *Gameboy) {
		// LD C, [HL]
		fmt.Println("Decoded OPCODE: LD C, [HL]")
		gb.cpu.setC(gb.mmu.read(gb.cpu.getHL()))
	},
	0x4F: func(gb *Gameboy) {
		// LD C, A
		fmt.Println("Decoded OPCODE: LD C, A")
		gb.cpu.setC(gb.cpu.reg.A)
	},
	0x50: func(gb *Gameboy) {
		// LD D, B
		fmt.Println("Decoded OPCODE: LD D, B")
		gb.cpu.setD(gb.cpu.reg.B)
	},
	0x51: func(gb *Gameboy) {
		// LD D, C
		fmt.Println("Decoded OPCODE: LD D, C")
		gb.cpu.setD(gb.cpu.reg.C)
	},
	0x52: func(gb *Gameboy) {
		// LD D, D
		fmt.Println("Decoded OPCODE: LD D, D")
		gb.cpu.setD(gb.cpu.reg.D)
	},
	0x53: func(gb *Gameboy) {
		// LD D, E
		fmt.Println("Decoded OPCODE: LD D, E")
		gb.cpu.setD(gb.cpu.reg.E)
	},
	0x54: func(gb *Gameboy) {
		// LD D, H
		fmt.Println("Decoded OPCODE: LD D, H")
		gb.cpu.setD(gb.cpu.reg.H)
	},
	0x55: func(gb *Gameboy) {
		// LD D, L
		fmt.Println("Decoded OPCODE: LD D, L")
		gb.cpu.setD(gb.cpu.reg.L)
	},
	0x56: func(gb *Gameboy) {
		// LD D, [HL]
		fmt.Println("Decoded OPCODE: LD D, [HL]")
		gb.cpu.setD(gb.mmu.read(gb.cpu.getHL()))
	},
	0x57: func(gb *Gameboy) {
		// LD D, A
		fmt.Println("Decoded OPCODE: LD D, A")
		gb.cpu.setD(gb.cpu.reg.A)
	},
	0x58: func(gb *Gameboy) {
		// LD E, B
		fmt.Println("Decoded OPCODE: LD E, B")
		gb.cpu.setE(gb.cpu.reg.B)
	},
	0x59: func(gb *Gameboy) {
		// LD E, C
		fmt.Println("Decoded OPCODE: LD E, C")
		gb.cpu.setE(gb.cpu.reg.C)
	},
	0x5A: func(gb *Gameboy) {
		// LD E, D
		fmt.Println("Decoded OPCODE: LD E, D")
		gb.cpu.setE(gb.cpu.reg.D)
	},
	0x5B: func(gb *Gameboy) {
		// LD E, E
		fmt.Println("Decoded OPCODE: LD E, E")
		gb.cpu.setE(gb.cpu.reg.E)
	},
	0x5C: func(gb *Gameboy) {
		// LD E, H
		fmt.Println("Decoded OPCODE: LD E, H")
		gb.cpu.setE(gb.cpu.reg.H)
	},
	0x5D: func(gb *Gameboy) {
		// LD E, L
		fmt.Println("Decoded OPCODE: LD E, L")
		gb.cpu.setE(gb.cpu.reg.L)
	},
	0x5E: func(gb *Gameboy) {
		// LD E, [HL]
		fmt.Println("Decoded OPCODE: LD E, [HL]")
		gb.cpu.setE(gb.mmu.read(gb.cpu.getHL()))
	},
	0x5F: func(gb *Gameboy) {
		// LD E, A
		fmt.Println("Decoded OPCODE: LD E, A")
		gb.cpu.setE(gb.cpu.reg.A)
	},
	0x60: func(gb *Gameboy) {
		// LD H, B
		fmt.Println("Decoded OPCODE: LD H, B")
		gb.cpu.setH(gb.cpu.reg.B)
	},
	0x61: func(gb *Gameboy) {
		// LD H, C
		fmt.Println("Decoded OPCODE: LD H, C")
		gb.cpu.setH(gb.cpu.reg.C)
	},
	0x62: func(gb *Gameboy) {
		// LD H, D
		fmt.Println("Decoded OPCODE: LD H, D")
		gb.cpu.setH(gb.cpu.reg.D)
	},
	0x63: func(gb *Gameboy) {
		// LD H, E
		fmt.Println("Decoded OPCODE: LD H, E")
		gb.cpu.setH(gb.cpu.reg.E)
	},
	0x64: func(gb *Gameboy) {
		// LD H, H
		fmt.Println("Decoded OPCODE: LD H, H")
		gb.cpu.setH(gb.cpu.reg.H)
	},
	0x65: func(gb *Gameboy) {
		// LD H, L
		fmt.Println("Decoded OPCODE: LD H, L")
		gb.cpu.setH(gb.cpu.reg.L)
	},
	0x66: func(gb *Gameboy) {
		// LD H, [HL]
		fmt.Println("Decoded OPCODE: LD H, [HL]")
		gb.cpu.setH(gb.mmu.read(gb.cpu.getHL()))
	},
	0x67: func(gb *Gameboy) {
		// LD H, A
		fmt.Println("Decoded OPCODE: LD H, A")
		gb.cpu.setH(gb.cpu.reg.A)
	},
	0x68: func(gb *Gameboy) {
		// LD L, B
		fmt.Println("Decoded OPCODE: LD L, B")
		gb.cpu.setL(gb.cpu.reg.B)
	},
	0x69: func(gb *Gameboy) {
		// LD L, C
		fmt.Println("Decoded OPCODE: LD L, C")
		gb.cpu.setL(gb.cpu.reg.C)
	},
	0x6A: func(gb *Gameboy) {
		// LD L, D
		fmt.Println("Decoded OPCODE: LD L, D")
		gb.cpu.setL(gb.cpu.reg.D)
	},
	0x6B: func(gb *Gameboy) {
		// LD L, E
		fmt.Println("Decoded OPCODE: LD L, E")
		gb.cpu.setL(gb.cpu.reg.E)
	},
	0x6C: func(gb *Gameboy) {
		// LD L, H
		fmt.Println("Decoded OPCODE: LD L, H")
		gb.cpu.setL(gb.cpu.reg.H)
	},
	0x6D: func(gb *Gameboy) {
		// LD L, L
		fmt.Println("Decoded OPCODE: LD L, L")
		gb.cpu.setL(gb.cpu.reg.L)
	},
	0x6E: func(gb *Gameboy) {
		// LD L, [HL]
		fmt.Println("Decoded OPCODE: LD L, [HL]")
		gb.cpu.setL(gb.mmu.read(gb.cpu.getHL()))
	},
	0x6F: func(gb *Gameboy) {
		// LD L, A
		fmt.Println("Decoded OPCODE: LD L, A")
		gb.cpu.setL(gb.cpu.reg.A)
	},
	0x70: func(gb *Gameboy) {
		// LD [HL], B
		fmt.Println("Decoded OPCODE: LD [HL], B")
		gb.mmu.write(gb.cpu.getHL(), gb.cpu.reg.B)
	},
	0x71: func(gb *Gameboy) {
		// LD [HL], C
		fmt.Println("Decoded OPCODE: LD [HL], C")
		gb.mmu.write(gb.cpu.getHL(), gb.cpu.reg.C)
	},
	0x72: func(gb *Gameboy) {
		// LD [HL], D
		fmt.Println("Decoded OPCODE: LD [HL], D")
		gb.mmu.write(gb.cpu.getHL(), gb.cpu.reg.D)
	},
	0x73: func(gb *Gameboy) {
		// LD [HL], E
		fmt.Println("Decoded OPCODE: LD [HL], E")
		gb.mmu.write(gb.cpu.getHL(), gb.cpu.reg.E)
	},
	0x74: func(gb *Gameboy) {
		// LD [HL], H
		fmt.Println("Decoded OPCODE: LD [HL], H")
		gb.mmu.write(gb.cpu.getHL(), gb.cpu.reg.H)
	},
	0x75: func(gb *Gameboy) {
		// LD [HL], L
		fmt.Println("Decoded OPCODE: LD [HL], L")
		gb.mmu.write(gb.cpu.getHL(), gb.cpu.reg.L)
	},
	0x76: func(gb *Gameboy) {
		// HALT
		fmt.Println("Decoded OPCODE: HALT")
		gb.cpu.halted = true
	},
	0x77: func(gb *Gameboy) {
		// LD [HL], A
		fmt.Println("Decoded OPCODE: LD [HL], A")
		gb.mmu.write(gb.cpu.getHL(), gb.cpu.reg.A)
	},
	0x78: func(gb *Gameboy) {
		// LD A, B
		fmt.Println("Decoded OPCODE: LD A, B")
		gb.cpu.setA(gb.cpu.reg.B)
	},
	0x79: func(gb *Gameboy) {
		// LD A, C
		fmt.Println("Decoded OPCODE: LD A, C")
		gb.cpu.setA(gb.cpu.reg.C)
	},
	0x7A: func(gb *Gameboy) {
		// LD A, D
		fmt.Println("Decoded OPCODE: LD A, D")
		gb.cpu.setA(gb.cpu.reg.D)
	},
	0x7B: func(gb *Gameboy) {
		// LD A, E
		fmt.Println("Decoded OPCODE: LD A, E")
		gb.cpu.setA(gb.cpu.reg.E)
	},
	0x7C: func(gb *Gameboy) {
		// LD A, H
		fmt.Println("Decoded OPCODE: LD A, H")
		gb.cpu.setA(gb.cpu.reg.H)
	},
	0x7D: func(gb *Gameboy) {
		// LD A, L
		fmt.Println("Decoded OPCODE: LD A, L")
		gb.cpu.setA(gb.cpu.reg.L)
	},
	0x7E: func(gb *Gameboy) {
		// LD A, [HL]
		fmt.Println("Decoded OPCODE: LD A, [HL]")
		gb.cpu.setA(gb.mmu.read(gb.cpu.getHL()))
	},
	0x7F: func(gb *Gameboy) {
		// LD A, A
		fmt.Println("Decoded OPCODE: LD A, A")
		gb.cpu.setA(gb.cpu.reg.A)
	},
	0x80: func(gb *Gameboy) {
		// ADD A, B
		fmt.Println("Decoded OPCODE: ADD A, B")
		gb.instrAddA(gb.cpu.reg.B, false)
	},
	0x81: func(gb *Gameboy) {
		// ADD A, C
		fmt.Println("Decoded OPCODE: ADD A, C")
		gb.instrAddA(gb.cpu.reg.C, false)
	},
	0x82: func(gb *Gameboy) {
		// ADD A, D
		fmt.Println("Decoded OPCODE: ADD A, D")
		gb.instrAddA(gb.cpu.reg.D, false)
	},
	0x83: func(gb *Gameboy) {
		// ADD A, E
		fmt.Println("Decoded OPCODE: ADD A, E")
		gb.instrAddA(gb.cpu.reg.E, false)
	},
	0x84: func(gb *Gameboy) {
		// ADD A, H
		fmt.Println("Decoded OPCODE: ADD A, H")
		gb.instrAddA(gb.cpu.reg.H, false)
	},
	0x85: func(gb *Gameboy) {
		// ADD A, L
		fmt.Println("Decoded OPCODE: ADD A, L")
		gb.instrAddA(gb.cpu.reg.L, false)
	},
	0x86: func(gb *Gameboy) {
		// ADD A, [HL]
		fmt.Println("Decoded OPCODE: ADD A, [HL]")
		gb.instrAddA(gb.mmu.read(gb.cpu.getHL()), false)
	},
	0x87: func(gb *Gameboy) {
		// ADD A, A
		fmt.Println("Decoded OPCODE: ADD A, A")
		gb.instrAddA(gb.cpu.reg.A, false)
	},
	0x88: func(gb *Gameboy) {
		// ADC A, B
		fmt.Println("Decoded OPCODE: ADC A, B")
		gb.instrAddA(gb.cpu.reg.B, true)
	},
	0x89: func(gb *Gameboy) {
		// ADC A, C
		fmt.Println("Decoded OPCODE: ADC A, C")
		gb.instrAddA(gb.cpu.reg.C, true)
	},
	0x8A: func(gb *Gameboy) {
		// ADC A, D
		fmt.Println("Decoded OPCODE: ADC A, D")
		gb.instrAddA(gb.cpu.reg.D, true)
	},
	0x8B: func(gb *Gameboy) {
		// ADC A, E
		fmt.Println("Decoded OPCODE: ADC A, E")
		gb.instrAddA(gb.cpu.reg.E, true)
	},
	0x8C: func(gb *Gameboy) {
		// ADC A, H
		fmt.Println("Decoded OPCODE: ADC A, H")
		gb.instrAddA(gb.cpu.reg.H, true)
	},
	0x8D: func(gb *Gameboy) {
		// ADC A, L
		fmt.Println("Decoded OPCODE: ADC A, L")
		gb.instrAddA(gb.cpu.reg.L, true)
	},
	0x8E: func(gb *Gameboy) {
		// ADC A, [HL]
		fmt.Println("Decoded OPCODE: ADC A, [HL]")
		gb.instrAddA(gb.mmu.read(gb.cpu.getHL()), true)
	},
	0x8F: func(gb *Gameboy) {
		// ADC A, A
		fmt.Println("Decoded OPCODE: ADC A, A")
		gb.instrAddA(gb.cpu.reg.A, true)
	},
	0x90: func(gb *Gameboy) {
		// SUB A, B
		fmt.Println("Decoded OPCODE: SUB A, B")
		gb.instrSubA(gb.cpu.reg.B, false)
	},
	0x91: func(gb *Gameboy) {
		// SUB A, C
		fmt.Println("Decoded OPCODE: SUB A, C")
		gb.instrSubA(gb.cpu.reg.C, false)
	},
	0x92: func(gb *Gameboy) {
		// SUB A, D
		fmt.Println("Decoded OPCODE: SUB A, D")
		gb.instrSubA(gb.cpu.reg.D, false)
	},
	0x93: func(gb *Gameboy) {
		// SUB A, E
		fmt.Println("Decoded OPCODE: SUB A, E")
		gb.instrSubA(gb.cpu.reg.E, false)
	},
	0x94: func(gb *Gameboy) {
		// SUB A, H
		fmt.Println("Decoded OPCODE: SUB A, H")
		gb.instrSubA(gb.cpu.reg.H, false)
	},
	0x95: func(gb *Gameboy) {
		// SUB A, L
		fmt.Println("Decoded OPCODE: SUB A, L")
		gb.instrSubA(gb.cpu.reg.L, false)
	},
	0x96: func(gb *Gameboy) {
		// SUB A, [HL]
		fmt.Println("Decoded OPCODE: SUB A, [HL]")
		gb.instrSubA(gb.mmu.read(gb.cpu.getHL()), false)
	},
	0x97: func(gb *Gameboy) {
		// SUB A, A
		fmt.Println("Decoded OPCODE: SUB A, A")
		gb.instrSubA(gb.cpu.reg.A, false)
	},
	0x98: func(gb *Gameboy) {
		// SBC A, B
		fmt.Println("Decoded OPCODE: SBC A, B")
		gb.instrSubA(gb.cpu.reg.B, true)
	},
	0x99: func(gb *Gameboy) {
		// SBC A, C
		fmt.Println("Decoded OPCODE: SBC A, C")
		gb.instrSubA(gb.cpu.reg.C, true)
	},
	0x9A: func(gb *Gameboy) {
		// SBC A, D
		fmt.Println("Decoded OPCODE: SBC A, D")
		gb.instrSubA(gb.cpu.reg.D, true)
	},
	0x9B: func(gb *Gameboy) {
		// SBC A, E
		fmt.Println("Decoded OPCODE: SBC A, E")
		gb.instrSubA(gb.cpu.reg.E, true)
	},
	0x9C: func(gb *Gameboy) {
		// SBC A, H
		fmt.Println("Decoded OPCODE: SBC A, H")
		gb.instrSubA(gb.cpu.reg.H, true)
	},
	0x9D: func(gb *Gameboy) {
		// SBC A, L
		fmt.Println("Decoded OPCODE: SBC A, L")
		gb.instrSubA(gb.cpu.reg.L, true)
	},
	0x9E: func(gb *Gameboy) {
		// SBC A, [HL]
		fmt.Println("Decoded OPCODE: SBC A, [HL]")
		gb.instrSubA(gb.mmu.read(gb.cpu.getHL()), true)
	},
	0x9F: func(gb *Gameboy) {
		// SBC A, A
		fmt.Println("Decoded OPCODE: SBC A, A")
		gb.instrSubA(gb.cpu.reg.A, true)
	},
	0xA0: func(gb *Gameboy) {
		// AND A, B
		fmt.Println("Decoded OPCODE: AND A, B")
		gb.instrAndA(gb.cpu.reg.B)
	},
	0xA1: func(gb *Gameboy) {
		// AND A, C
		fmt.Println("Decoded OPCODE: AND A, C")
		gb.instrAndA(gb.cpu.reg.C)
	},
	0xA2: func(gb *Gameboy) {
		// AND A, D
		fmt.Println("Decoded OPCODE: AND A, D")
		gb.instrAndA(gb.cpu.reg.D)
	},
	0xA3: func(gb *Gameboy) {
		// AND A, E
		fmt.Println("Decoded OPCODE: AND A, E")
		gb.instrAndA(gb.cpu.reg.E)
	},
	0xA4: func(gb *Gameboy) {
		// AND A, H
		fmt.Println("Decoded OPCODE: AND A, H")
		gb.instrAndA(gb.cpu.reg.H)
	},
	0xA5: func(gb *Gameboy) {
		// AND A, L
		fmt.Println("Decoded OPCODE: AND A, L")
		gb.instrAndA(gb.cpu.reg.L)
	},
	0xA6: func(gb *Gameboy) {
		// AND A, [HL]
		fmt.Println("Decoded OPCODE: AND A, [HL]")
		gb.instrAndA(gb.mmu.read(gb.cpu.getHL()))
	},
	0xA7: func(gb *Gameboy) {
		// AND A, A
		fmt.Println("Decoded OPCODE: AND A, A")
		gb.instrAndA(gb.cpu.reg.A)
	},
	0xA8: func(gb *Gameboy) {
		// XOR A, B
		fmt.Println("Decoded OPCODE: XOR A, B")
		gb.instrXorA(gb.cpu.reg.B)
	},
	0xA9: func(gb *Gameboy) {
		// XOR A, C
		fmt.Println("Decoded OPCODE: XOR A, C")
		gb.instrXorA(gb.cpu.reg.C)
	},
	0xAA: func(gb *Gameboy) {
		// XOR A, D
		fmt.Println("Decoded OPCODE: XOR A, D")
		gb.instrXorA(gb.cpu.reg.D)
	},
	0xAB: func(gb *Gameboy) {
		// XOR A, E
		fmt.Println("Decoded OPCODE: XOR A, E")
		gb.instrXorA(gb.cpu.reg.E)
	},
	0xAC: func(gb *Gameboy) {
		// XOR A, H
		fmt.Println("Decoded OPCODE: XOR A, H")
		gb.instrXorA(gb.cpu.reg.H)
	},
	0xAD: func(gb *Gameboy) {
		// XOR A, L
		fmt.Println("Decoded OPCODE: XOR A, L")
		gb.instrXorA(gb.cpu.reg.L)
	},
	0xAE: func(gb *Gameboy) {
		// XOR A, [HL]
		fmt.Println("Decoded OPCODE: XOR A, [HL]")
		gb.instrXorA(gb.mmu.read(gb.cpu.getHL()))
	},
	0xAF: func(gb *Gameboy) {
		// XOR A, A
		fmt.Println("Decoded OPCODE: XOR A, A")
		gb.instrXorA(gb.cpu.reg.A)
	},
	0xB0: func(gb *Gameboy) {
		// OR A, B
		fmt.Println("Decoded OPCODE: OR A, B")
		gb.instrOrA(gb.cpu.reg.B)
	},
	0xB1: func(gb *Gameboy) {
		// OR A, C
		fmt.Println("Decoded OPCODE: OR A, C")
		gb.instrOrA(gb.cpu.reg.C)
	},
	0xB2: func(gb *Gameboy) {
		// OR A, D
		fmt.Println("Decoded OPCODE: OR A, D")
		gb.instrOrA(gb.cpu.reg.D)
	},
	0xB3: func(gb *Gameboy) {
		// OR A, E
		fmt.Println("Decoded OPCODE: OR A, E")
		gb.instrOrA(gb.cpu.reg.E)
	},
	0xB4: func(gb *Gameboy) {
		// OR A, H
		fmt.Println("Decoded OPCODE: OR A, H")
		gb.instrOrA(gb.cpu.reg.H)
	},
	0xB5: func(gb *Gameboy) {
		// OR A, L
		fmt.Println("Decoded OPCODE: OR A, L")
		gb.instrOrA(gb.cpu.reg.L)
	},
	0xB6: func(gb *Gameboy) {
		// OR A, [HL]
		fmt.Println("Decoded OPCODE: OR A, [HL]")
		gb.instrOrA(gb.mmu.read(gb.cpu.getHL()))
	},
	0xB7: func(gb *Gameboy) {
		// OR A, A
		fmt.Println("Decoded OPCODE: OR A, A")
		gb.instrOrA(gb.cpu.reg.A)
	},
	0xB8: func(gb *Gameboy) {
		// CP A, B
		fmt.Println("Decoded OPCODE: CP A, B")
		gb.instrCpA(gb.cpu.reg.B)
	},
	0xB9: func(gb *Gameboy) {
		// CP A, C
		fmt.Println("Decoded OPCODE: CP A, C")
		gb.instrCpA(gb.cpu.reg.C)
	},
	0xBA: func(gb *Gameboy) {
		// CP A, D
		fmt.Println("Decoded OPCODE: CP A, D")
		gb.instrCpA(gb.cpu.reg.D)
	},
	0xBB: func(gb *Gameboy) {
		// CP A, E
		fmt.Println("Decoded OPCODE: CP A, E")
		gb.instrCpA(gb.cpu.reg.E)
	},
	0xBC: func(gb *Gameboy) {
		// CP A, H
		fmt.Println("Decoded OPCODE: CP A, H")
		gb.instrCpA(gb.cpu.reg.H)
	},
	0xBD: func(gb *Gameboy) {
		// CP A, L
		fmt.Println("Decoded OPCODE: CP A, L")
		gb.instrCpA(gb.cpu.reg.L)
	},
	0xBE: func(gb *Gameboy) {
		// CP A, [HL]
		fmt.Println("Decoded OPCODE: CP A, [HL]")
		gb.instrCpA(gb.mmu.read(gb.cpu.getHL()))
	},
	0xBF: func(gb *Gameboy) {
		// CP A, A
		fmt.Println("Decoded OPCODE: CP A, A")
		gb.instrCpA(gb.cpu.reg.A)
	},
	0xC0: func(gb *Gameboy) {
		// RET NZ
		fmt.Println("Decoded OPCODE: RET NZ")
		if !gb.cpu.zFlag() {
			gb.instrRet()
			gb.cpu.ticks += 12
		}
	},
	0xC1: func(gb *Gameboy) {
		// POP BC
		fmt.Println("Decoded OPCODE: POP BC")
		gb.cpu.setBC(gb.popStack())
	},
	0xC2: func(gb *Gameboy) {
		// JP NZ, u16
		fmt.Println("Decoded OPCODE: JP NZ, u16")
		jumpAddress := gb.nextPC16()
		if !gb.cpu.zFlag() {
			gb.instrJump(jumpAddress)
			gb.cpu.ticks += 4
		}
	},
	0xC3: func(gb *Gameboy) {
		// JP u16
		fmt.Println("Decoded OPCODE: JP u16")
		gb.instrJump(gb.nextPC16())
	},
	0xC4: func(gb *Gameboy) {
		// CALL NZ, u16
		fmt.Println("Decoded OPCODE: CALL NZ, u16")
		jumpAddress := gb.nextPC16()
		if !gb.cpu.zFlag() {
			gb.instrCall(jumpAddress)
			gb.cpu.ticks += 12
		}
	},
	0xC5: func(gb *Gameboy) {
		// PUSH BC
		fmt.Println("Decoded OPCODE: PUSH BC")
		gb.pushStack(gb.cpu.getBC())
	},
	0xC6: func(gb *Gameboy) {
		// ADD A, u8
		fmt.Println("Decoded OPCODE: ADD A, u8")
		gb.instrAddA(gb.nextPC(), false)
	},
	0xC7: func(gb *Gameboy) {
		// RST $00
		fmt.Println("Decoded OPCODE: RST $00")
		gb.instrCall(0x0000)
	},
	0xC8: func(gb *Gameboy) {
		// RET Z
		fmt.Println("Decoded OPCODE: RET Z")
		if gb.cpu.zFlag() {
			gb.instrRet()
			gb.cpu.ticks += 12
		}
	},
	0xC9: func(gb *Gameboy) {
		// RET
		fmt.Println("Decoded OPCODE: RET")
		gb.instrRet()
	},
	0xCA: func(gb *Gameboy) {
		// JP Z, u16
		fmt.Println("Decoded OPCODE: JP Z, u16")
		jumpAddress := gb.nextPC16()
		if gb.cpu.zFlag() {
			gb.instrJump(jumpAddress)
			gb.cpu.ticks += 4
		}
	},
	0xCB: func(gb *Gameboy) {
		// PREFIX
		fmt.Println("Decoded OPCODE: PREFIX")

		cbOpcode := gb.nextPC()
		fmt.Printf("CB OPCODE: 0x%02x\n", cbOpcode)

		gb.cpu.ticks += cbInstrClockTicks[cbOpcode]
		gb.cbInstructions[cbOpcode]()
	},
	0xCC: func(gb *Gameboy) {
		// CALL Z, u16
		fmt.Println("Decoded OPCODE: CALL Z, u16")
		jumpAddress := gb.nextPC16()
		if gb.cpu.zFlag() {
			gb.instrCall(jumpAddress)
			gb.cpu.ticks += 12
		}
	},
	0xCD: func(gb *Gameboy) {
		// CALL u16
		fmt.Println("Decoded OPCODE: CALL u16")
		gb.instrCall(gb.nextPC16())
	},
	0xCE: func(gb *Gameboy) {
		// ADC A, u8
		fmt.Println("Decoded OPCODE: ADC A, u8")
		gb.instrAddA(gb.nextPC(), true)
	},
	0xCF: func(gb *Gameboy) {
		// RST $08
		fmt.Println("Decoded OPCODE: RST $08")
		gb.instrCall(0x0008)
	},
	0xD0: func(gb *Gameboy) {
		// RET NC
		fmt.Println("Decoded OPCODE: RET NC")
		if !gb.cpu.cFlag() {
			gb.instrRet()
			gb.cpu.ticks += 12
		}
	},
	0xD1: func(gb *Gameboy) {
		// POP DE
		fmt.Println("Decoded OPCODE: POP DE")
		gb.cpu.setDE(gb.popStack())
	},
	0xD2: func(gb *Gameboy) {
		// JP NC, u16
		fmt.Println("Decoded OPCODE: JP NC, u16")
		jumpAddress := gb.nextPC16()
		if !gb.cpu.cFlag() {
			gb.instrJump(jumpAddress)
			gb.cpu.ticks += 4
		}
	},
	0xD3: func(gb *Gameboy) {
		// ILLEGAL OPCODE 0xD3
		fmt.Println("ILLEGAL OPCODE: 0xD3")
	},
	0xD4: func(gb *Gameboy) {
		// CALL NC, u16
		fmt.Println("Decoded OPCODE: CALL NC, u16")
		jumpAddress := gb.nextPC16()
		if !gb.cpu.cFlag() {
			gb.instrCall(jumpAddress)
			gb.cpu.ticks += 12
		}
	},
	0xD5: func(gb *Gameboy) {
		// PUSH DE
		fmt.Println("Decoded OPCODE: PUSH DE")
		gb.pushStack(gb.cpu.getDE())
	},
	0xD6: func(gb *Gameboy) {
		// SUB A, u8
		fmt.Println("Decoded OPCODE: SUB A, u8")
		gb.instrSubA(gb.nextPC(), false)
	},
	0xD7: func(gb *Gameboy) {
		// RST $10
		fmt.Println("Decoded OPCODE: RST $10")
		gb.instrCall(0x0010)
	},
	0xD8: func(gb *Gameboy) {
		// RET C
		fmt.Println("Decoded OPCODE: RET C")
		if gb.cpu.cFlag() {
			gb.instrRet()
			gb.cpu.ticks += 12
		}
	},
	0xD9: func(gb *Gameboy) {
		// RETI
		fmt.Println("Decoded OPCODE: RETI")
		gb.instrRet()
		gb.interruptsPendingEnabled = true
	},
	0xDA: func(gb *Gameboy) {
		// JP C, u16
		fmt.Println("Decoded OPCODE: JP C, u16")
		jumpAddress := gb.nextPC16()
		if gb.cpu.cFlag() {
			gb.instrJump(jumpAddress)
			gb.cpu.ticks += 4
		}
	},
	0xDB: func(gb *Gameboy) {
		// ILLEGAL OPCODE 0xDB
		fmt.Println("ILLEGAL OPCODE: 0xDB")
	},
	0xDC: func(gb *Gameboy) {
		// CALL C, u16
		fmt.Println("Decoded OPCODE: CALL C, u16")
		jumpAddress := gb.nextPC16()
		if gb.cpu.cFlag() {
			gb.instrCall(jumpAddress)
			gb.cpu.ticks += 12
		}
	},
	0xDD: func(gb *Gameboy) {
		// ILLEGAL OPCODE 0xDD
		fmt.Println("ILLEGAL OPCODE: 0xDD")
	},
	0xDE: func(gb *Gameboy) {
		// SBC A, u8
		fmt.Println("Decoded OPCODE: SBC A, u8")
		gb.instrSubA(gb.nextPC(), true)
	},
	0xDF: func(gb *Gameboy) {
		// RST $18
		fmt.Println("Decoded OPCODE: RST $18")
		gb.instrCall(0x0018)
	},
	0xE0: func(gb *Gameboy) {
		// LD (FF00+u8), A
		fmt.Println("Decoded OPCODE: LD (FF00+u8), A")
		address := 0xFF00 + uint16(gb.nextPC())
		gb.mmu.write(address, gb.cpu.reg.A)
	},
	0xE1: func(gb *Gameboy) {
		// POP HL
		fmt.Println("Decoded OPCODE: POP HL")
		gb.cpu.setHL(gb.popStack())
	},
	0xE2: func(gb *Gameboy) {
		// LD (FF00+C), A
		fmt.Println("Decoded OPCODE: LD (FF00+C), A")
		address := 0xFF00 + uint16(gb.cpu.reg.C)
		gb.mmu.write(address, gb.cpu.reg.A)
	},
	0xE3: func(gb *Gameboy) {
		// ILLEGAL OPCODE 0xE3
		fmt.Println("ILLEGAL OPCODE: 0xE3")
	},
	0xE4: func(gb *Gameboy) {
		// ILLEGAL OPCODE 0xE4
		fmt.Println("ILLEGAL OPCODE: 0xE4")
	},
	0xE5: func(gb *Gameboy) {
		// PUSH HL
		fmt.Println("Decoded OPCODE: PUSH HL")
		gb.pushStack(gb.cpu.getHL())
	},
	0xE6: func(gb *Gameboy) {
		// AND A, u8
		fmt.Println("Decoded OPCODE: AND A, u8")
		gb.instrAndA(gb.nextPC())
	},
	0xE7: func(gb *Gameboy) {
		// RST $20
		fmt.Println("Decoded OPCODE: RST $20")
		gb.instrCall(0x0020)
	},
	0xE8: func(gb *Gameboy) {
		// ADD SP, i8
		fmt.Println("Decoded OPCODE: ADD SP, i8")
		gb.instrAdd16Signed(gb.cpu.setSP, gb.cpu.reg.SP, int8(gb.nextPC()))
	},
	0xE9: func(gb *Gameboy) {
		// JP HL
		fmt.Println("Decoded OPCODE: JP HL")
		gb.instrJump(gb.cpu.getHL())
	},
	0xEA: func(gb *Gameboy) {
		// LD [u16], A
		fmt.Println("Decoded OPCODE: LD [u16], A")
		gb.mmu.write(gb.nextPC16(), gb.cpu.reg.A)
	},
	0xEB: func(gb *Gameboy) {
		// ILLEGAL OPCODE 0xEB
		fmt.Println("ILLEGAL OPCODE: 0xEB")
	},
	0xEC: func(gb *Gameboy) {
		// ILLEGAL OPCODE 0xEC
		fmt.Println("ILLEGAL OPCODE: 0xEC")
	},
	0xED: func(gb *Gameboy) {
		// ILLEGAL OPCODE 0xED
		fmt.Println("ILLEGAL OPCODE: 0xED")
	},
	0xEE: func(gb *Gameboy) {
		// XOR A, u8
		fmt.Println("Decoded OPCODE: XOR A, u8")
		gb.instrXorA(gb.nextPC())
	},
	0xEF: func(gb *Gameboy) {
		// RST $28
		fmt.Println("Decoded OPCODE: RST $28")
		gb.instrCall(0x0028)
	},
	0xF0: func(gb *Gameboy) {
		// LD A, (FF00+u8)
		fmt.Println("Decoded OPCODE: LD A, (FF00+u8)")
		address := 0xFF00 + uint16(gb.nextPC())
		gb.cpu.setA(gb.mmu.read(address))
	},
	0xF1: func(gb *Gameboy) {
		// POP AF
		fmt.Println("Decoded OPCODE: POP AF")
		gb.cpu.setAF(gb.popStack())
	},
	0xF2: func(gb *Gameboy) {
		// LD A, (FF00+C)
		fmt.Println("Decoded OPCODE: LD A, (FF00+C)")
		address := 0xFF00 + uint16(gb.cpu.reg.C)
		gb.cpu.setA(gb.mmu.read(address))
	},
	0xF3: func(gb *Gameboy) {
		// DI
		fmt.Println("Decoded OPCODE: DI")
		gb.interruptsOn = false
	},
	0xF4: func(gb *Gameboy) {
		// ILLEGAL OPCODE 0xF4
		fmt.Println("ILLEGAL OPCODE: 0xF4")
	},
	0xF5: func(gb *Gameboy) {
		// PUSH AF
		fmt.Println("Decoded OPCODE: PUSH AF")
		gb.pushStack(gb.cpu.getAF())
	},
	0xF6: func(gb *Gameboy) {
		// OR A, u8
		fmt.Println("Decoded OPCODE: OR A, u8")
		gb.instrOrA(gb.nextPC())
	},
	0xF7: func(gb *Gameboy) {
		// RST $30
		fmt.Println("Decoded OPCODE: RST $30")
		gb.instrCall(0x0030)
	},
	0xF8: func(gb *Gameboy) {
		// LD HL, SP + e8
		fmt.Println("Decoded OPCODE: LD HL, SP + e8")
		gb.instrAdd16Signed(gb.cpu.setHL, gb.cpu.reg.SP, int8(gb.nextPC()))
	},
	0xF9: func(gb *Gameboy) {
		// LD SP, HL
		fmt.Println("Decoded OPCODE: LD SP, HL")
		gb.cpu.setSP(gb.cpu.getHL())
	},
	0xFA: func(gb *Gameboy) {
		// LD A, [u16]
		fmt.Println("Decoded OPCODE: LD A, [u16]")
		gb.cpu.setA(gb.mmu.read(gb.nextPC16()))
	},
	0xFB: func(gb *Gameboy) {
		// EI
		fmt.Println("Decoded OPCODE: EI")
		gb.interruptsPendingEnabled = true
	},
	0xFC: func(gb *Gameboy) {
		// ILLEGAL OPCODE 0xFC
		fmt.Println("ILLEGAL OPCODE: 0xFC")
	},
	0xFD: func(gb *Gameboy) {
		// ILLEGAL OPCODE 0xFD
		fmt.Println("ILLEGAL OPCODE: 0xFD")
	},
	0xFE: func(gb *Gameboy) {
		// CP A, u8
		fmt.Println("Decoded OPCODE: CP A, u8")
		gb.instrCpA(gb.nextPC())
	},
	0xFF: func(gb *Gameboy) {
		// RST $38
		fmt.Println("Decoded OPCODE: RST $38")
		gb.instrCall(0x0038)
	},
}
