package gb

import "github.com/BeralaWoolies/GameboyGo/pkg/bits"

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

func (cpu *CPU) instrRlc(setHandler func(result uint8), val uint8) {
	result := uint8((val << 1)) | (val >> 7)
	setHandler(result)

	cpu.setZFlag(result == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(false)
	cpu.setCFlag(val >= 0x80)
}

func (cpu *CPU) instrRrc(setHandler func(result uint8), val uint8) {
	result := (val >> 1) | ((val & 1) << 7)
	setHandler(result)

	cpu.setZFlag(result == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(false)
	cpu.setCFlag(result >= 0x80)
}

func (cpu *CPU) instrRl(setHandler func(result uint8), val uint8) {
	var carry uint8 = 0
	if cpu.cFlag() {
		carry = 1
	}

	result := uint8(val<<1) + carry
	setHandler(result)

	cpu.setZFlag(result == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(false)
	cpu.setCFlag(val >= 0x80)
}

func (cpu *CPU) instrRr(setHandler func(result uint8), val uint8) {
	var carry uint8 = 0
	if cpu.cFlag() {
		carry = 0x80
	}

	result := uint8(val>>1) | carry
	setHandler(result)

	cpu.setZFlag(result == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(false)
	cpu.setCFlag((val & 1) == 1)
}

func (cpu *CPU) instrSla(setHandler func(result uint8), val uint8) {
	result := uint8(val << 1)
	setHandler(result)

	cpu.setZFlag(result == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(false)
	cpu.setCFlag(val >= 0x80)
}

func (cpu *CPU) instrSra(setHandler func(result uint8), val uint8) {
	result := (val & 0x80) | (val >> 1)
	setHandler(result)

	cpu.setZFlag(result == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(false)
	cpu.setCFlag((val & 1) == 1)
}

func (cpu *CPU) instrSrl(setHandler func(result uint8), val uint8) {
	result := val >> 1
	setHandler(result)

	cpu.setZFlag(result == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(false)
	cpu.setCFlag((val & 1) == 1)
}

func (cpu *CPU) instrSwap(setHandler func(result uint8), val uint8) {
	result := ((val >> 4) & 0xF) | ((val << 4) & 0xF0)
	setHandler(result)

	cpu.setZFlag(result == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(false)
	cpu.setCFlag(false)
}

func (cpu *CPU) instrBit(val uint8, pos uint8) {
	cpu.setZFlag((val>>pos)&1 == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(true)
}

func (cpu *CPU) initCbInstructions() [0x100]func() {
	instructions := [0x100]func(){}

	setHandlers := [8]func(result uint8){
		0: cpu.setB,
		1: cpu.setC,
		2: cpu.setD,
		3: cpu.setE,
		4: cpu.setH,
		5: cpu.setL,
		6: func(result uint8) { cpu.mmu.write(cpu.getHL(), result) },
		7: cpu.setA,
	}

	getHandlers := [8]func() uint8{
		0: func() uint8 { return cpu.reg.B },
		1: func() uint8 { return cpu.reg.C },
		2: func() uint8 { return cpu.reg.D },
		3: func() uint8 { return cpu.reg.E },
		4: func() uint8 { return cpu.reg.H },
		5: func() uint8 { return cpu.reg.L },
		6: func() uint8 { return cpu.mmu.read(cpu.getHL()) },
		7: func() uint8 { return cpu.reg.A },
	}

	for x := 0; x < 8; x++ {
		i := x

		instructions[0x00+i] = func() { cpu.instrRlc(setHandlers[i], getHandlers[i]()) }
		instructions[0x08+i] = func() { cpu.instrRrc(setHandlers[i], getHandlers[i]()) }
		instructions[0x10+i] = func() { cpu.instrRl(setHandlers[i], getHandlers[i]()) }
		instructions[0x18+i] = func() { cpu.instrRr(setHandlers[i], getHandlers[i]()) }
		instructions[0x20+i] = func() { cpu.instrSla(setHandlers[i], getHandlers[i]()) }
		instructions[0x28+i] = func() { cpu.instrSra(setHandlers[i], getHandlers[i]()) }
		instructions[0x30+i] = func() { cpu.instrSwap(setHandlers[i], getHandlers[i]()) }
		instructions[0x38+i] = func() { cpu.instrSrl(setHandlers[i], getHandlers[i]()) }
		instructions[0x40+i] = func() { cpu.instrBit(getHandlers[i](), 0) }
		instructions[0x48+i] = func() { cpu.instrBit(getHandlers[i](), 1) }
		instructions[0x50+i] = func() { cpu.instrBit(getHandlers[i](), 2) }
		instructions[0x58+i] = func() { cpu.instrBit(getHandlers[i](), 3) }
		instructions[0x60+i] = func() { cpu.instrBit(getHandlers[i](), 4) }
		instructions[0x68+i] = func() { cpu.instrBit(getHandlers[i](), 5) }
		instructions[0x70+i] = func() { cpu.instrBit(getHandlers[i](), 6) }
		instructions[0x78+i] = func() { cpu.instrBit(getHandlers[i](), 7) }
		instructions[0x80+i] = func() { setHandlers[i](bits.Reset(getHandlers[i](), 0)) }
		instructions[0x88+i] = func() { setHandlers[i](bits.Reset(getHandlers[i](), 1)) }
		instructions[0x90+i] = func() { setHandlers[i](bits.Reset(getHandlers[i](), 2)) }
		instructions[0x98+i] = func() { setHandlers[i](bits.Reset(getHandlers[i](), 3)) }
		instructions[0xA0+i] = func() { setHandlers[i](bits.Reset(getHandlers[i](), 4)) }
		instructions[0xA8+i] = func() { setHandlers[i](bits.Reset(getHandlers[i](), 5)) }
		instructions[0xB0+i] = func() { setHandlers[i](bits.Reset(getHandlers[i](), 6)) }
		instructions[0xB8+i] = func() { setHandlers[i](bits.Reset(getHandlers[i](), 7)) }
		instructions[0xC0+i] = func() { setHandlers[i](bits.Set(getHandlers[i](), 0)) }
		instructions[0xC8+i] = func() { setHandlers[i](bits.Set(getHandlers[i](), 1)) }
		instructions[0xD0+i] = func() { setHandlers[i](bits.Set(getHandlers[i](), 2)) }
		instructions[0xD8+i] = func() { setHandlers[i](bits.Set(getHandlers[i](), 3)) }
		instructions[0xE0+i] = func() { setHandlers[i](bits.Set(getHandlers[i](), 4)) }
		instructions[0xE8+i] = func() { setHandlers[i](bits.Set(getHandlers[i](), 5)) }
		instructions[0xF0+i] = func() { setHandlers[i](bits.Set(getHandlers[i](), 6)) }
		instructions[0xF8+i] = func() { setHandlers[i](bits.Set(getHandlers[i](), 7)) }
	}
	return instructions
}
