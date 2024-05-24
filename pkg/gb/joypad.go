package gb

import (
	"log"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
)

type Joypad struct {
	ic    *IntruptController
	a     *Button
	b     *Button
	up    *Button
	down  *Button
	right *Button
	left  *Button
	start *Button
	sel   *Button
	reg   uint8
}

type Button struct {
	pressed uint8
}

const (
	JOYP_ADDR        = 0xFF00
	JOYP_A_RIGHT     = 0
	JOYP_B_LEFT      = 1
	JOYP_SELECT_UP   = 2
	JOYP_START_DOWN  = 3
	JOYP_DPAD_SELECT = 4
	JOYP_BTN_SELECT  = 5
)

func (joyp *Joypad) init(ic *IntruptController) {
	joyp.ic = ic
	joyp.a = &Button{}
	joyp.b = &Button{}
	joyp.up = &Button{}
	joyp.down = &Button{}
	joyp.right = &Button{}
	joyp.left = &Button{}
	joyp.start = &Button{}
	joyp.sel = &Button{}
	joyp.reg = 0xFF
}

func (joyp *Joypad) contains(addr uint16) bool {
	return addr == JOYP_ADDR
}

func (joyp *Joypad) read(addr uint16) uint8 {
	switch addr {
	case JOYP_ADDR:
		return joyp.output()
	default:
		log.Fatalf("MMU mapped an illegal read address: 0x%02x to Joypad", addr)
		return 0xFF
	}
}

func (joyp *Joypad) write(addr uint16, data uint8) {
	switch addr {
	case JOYP_ADDR:
		joyp.reg = ((data | 0xC0) & 0x30) | (joyp.reg & 0xF)
	default:
		log.Fatalf("MMU mapped an illegal write address: 0x%02x to Joypad", addr)
	}
}

func (joyp *Joypad) output() uint8 {
	if bits.IsSet(joyp.reg, JOYP_DPAD_SELECT) && bits.IsSet(joyp.reg, JOYP_BTN_SELECT) {
		return 0xFF
	}

	joyp.reg |= 0xCF

	if !bits.IsSet(joyp.reg, JOYP_DPAD_SELECT) && !bits.IsSet(joyp.reg, JOYP_BTN_SELECT) {
		joyp.reg = bits.SetCond(joyp.reg, JOYP_A_RIGHT, joyp.a.eitherPressed(joyp.right))
		joyp.reg = bits.SetCond(joyp.reg, JOYP_B_LEFT, joyp.b.eitherPressed(joyp.left))
		joyp.reg = bits.SetCond(joyp.reg, JOYP_SELECT_UP, joyp.sel.eitherPressed(joyp.up))
		joyp.reg = bits.SetCond(joyp.reg, JOYP_START_DOWN, joyp.start.eitherPressed(joyp.down))
	} else if bits.IsSet(joyp.reg, JOYP_BTN_SELECT) {
		joyp.reg = bits.SetCond(joyp.reg, JOYP_A_RIGHT, joyp.right.pressed)
		joyp.reg = bits.SetCond(joyp.reg, JOYP_B_LEFT, joyp.left.pressed)
		joyp.reg = bits.SetCond(joyp.reg, JOYP_SELECT_UP, joyp.up.pressed)
		joyp.reg = bits.SetCond(joyp.reg, JOYP_START_DOWN, joyp.down.pressed)
	} else {
		joyp.reg = bits.SetCond(joyp.reg, JOYP_A_RIGHT, joyp.a.pressed)
		joyp.reg = bits.SetCond(joyp.reg, JOYP_B_LEFT, joyp.b.pressed)
		joyp.reg = bits.SetCond(joyp.reg, JOYP_SELECT_UP, joyp.sel.pressed)
		joyp.reg = bits.SetCond(joyp.reg, JOYP_START_DOWN, joyp.start.pressed)
	}

	return joyp.reg
}

func (btn *Button) eitherPressed(other *Button) uint8 {
	return btn.pressed & other.pressed
}

func (btn *Button) press(pressed bool) {
	// pressed if 0
	btn.pressed = bits.BoolToUint8(!pressed)
}
