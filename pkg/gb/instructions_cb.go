package gb

import "github.com/BeralaWoolies/GameboyGo/pkg/bits"

func (gb *Gameboy) instrRlc(setHandler func(result uint8), val uint8) {
	result := uint8((val << 1)) | (val >> 7)
	setHandler(result)

	gb.cpu.setZFlag(result == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag(false)
	gb.cpu.setCFlag(val >= 0x80)
}

func (gb *Gameboy) instrRrc(setHandler func(result uint8), val uint8) {
	result := (val >> 1) | ((val & 1) << 7)
	setHandler(result)

	gb.cpu.setZFlag(result == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag(false)
	gb.cpu.setCFlag(result >= 0x80)
}

func (gb *Gameboy) instrRl(setHandler func(result uint8), val uint8) {
	var carry uint8 = 0
	if gb.cpu.cFlag() {
		carry = 1
	}

	result := uint8(val<<1) + carry
	setHandler(result)

	gb.cpu.setZFlag(result == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag(false)
	gb.cpu.setCFlag(val >= 0x80)
}

func (gb *Gameboy) instrRr(setHandler func(result uint8), val uint8) {
	var carry uint8 = 0
	if gb.cpu.cFlag() {
		carry = 0x80
	}

	result := uint8(val>>1) | carry
	setHandler(result)

	gb.cpu.setZFlag(result == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag(false)
	gb.cpu.setCFlag((val & 1) == 1)
}

func (gb *Gameboy) instrSla(setHandler func(result uint8), val uint8) {
	result := uint8(val << 1)
	setHandler(result)

	gb.cpu.setZFlag(result == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag(false)
	gb.cpu.setCFlag(val >= 0x80)
}

func (gb *Gameboy) instrSra(setHandler func(result uint8), val uint8) {
	result := (val & 0x80) | (val >> 1)
	setHandler(result)

	gb.cpu.setZFlag(result == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag(false)
	gb.cpu.setCFlag((val & 1) == 1)
}

func (gb *Gameboy) instrSrl(setHandler func(result uint8), val uint8) {
	result := val >> 1
	setHandler(result)

	gb.cpu.setZFlag(result == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag(false)
	gb.cpu.setCFlag((val & 1) == 1)
}

func (gb *Gameboy) instrSwap(setHandler func(result uint8), val uint8) {
	result := ((val >> 4) & 0xF) | ((val << 4) & 0xF0)
	setHandler(result)

	gb.cpu.setZFlag(result == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag(false)
	gb.cpu.setCFlag(false)
}

func (gb *Gameboy) instrBit(val uint8, pos uint8) {
	gb.cpu.setZFlag((val>>pos)&1 == 0)
	gb.cpu.setNFlag(false)
	gb.cpu.setHFlag(true)
}

func (gb *Gameboy) initCbInstructions() [0x100]func() {
	instructions := [0x100]func(){}

	setHandlers := [8]func(result uint8){
		0: gb.cpu.setB,
		1: gb.cpu.setC,
		2: gb.cpu.setD,
		3: gb.cpu.setE,
		4: gb.cpu.setH,
		5: gb.cpu.setL,
		6: func(result uint8) { gb.mmu.write(gb.cpu.getHL(), result) },
		7: gb.cpu.setA,
	}

	getHandlers := [8]func() uint8{
		0: func() uint8 { return gb.cpu.reg.B },
		1: func() uint8 { return gb.cpu.reg.C },
		2: func() uint8 { return gb.cpu.reg.D },
		3: func() uint8 { return gb.cpu.reg.E },
		4: func() uint8 { return gb.cpu.reg.H },
		5: func() uint8 { return gb.cpu.reg.L },
		6: func() uint8 { return gb.mmu.read(gb.cpu.getHL()) },
		7: func() uint8 { return gb.cpu.reg.A },
	}

	for x := 0; x < 8; x++ {
		i := x

		instructions[0x00+i] = func() { gb.instrRlc(setHandlers[i], getHandlers[i]()) }
		instructions[0x08+i] = func() { gb.instrRrc(setHandlers[i], getHandlers[i]()) }
		instructions[0x10+i] = func() { gb.instrRl(setHandlers[i], getHandlers[i]()) }
		instructions[0x18+i] = func() { gb.instrRr(setHandlers[i], getHandlers[i]()) }
		instructions[0x20+i] = func() { gb.instrSla(setHandlers[i], getHandlers[i]()) }
		instructions[0x28+i] = func() { gb.instrSra(setHandlers[i], getHandlers[i]()) }
		instructions[0x30+i] = func() { gb.instrSwap(setHandlers[i], getHandlers[i]()) }
		instructions[0x38+i] = func() { gb.instrSrl(setHandlers[i], getHandlers[i]()) }
		instructions[0x40+i] = func() { gb.instrBit(getHandlers[i](), 0) }
		instructions[0x48+i] = func() { gb.instrBit(getHandlers[i](), 1) }
		instructions[0x50+i] = func() { gb.instrBit(getHandlers[i](), 2) }
		instructions[0x58+i] = func() { gb.instrBit(getHandlers[i](), 3) }
		instructions[0x60+i] = func() { gb.instrBit(getHandlers[i](), 4) }
		instructions[0x68+i] = func() { gb.instrBit(getHandlers[i](), 5) }
		instructions[0x70+i] = func() { gb.instrBit(getHandlers[i](), 6) }
		instructions[0x78+i] = func() { gb.instrBit(getHandlers[i](), 7) }
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
