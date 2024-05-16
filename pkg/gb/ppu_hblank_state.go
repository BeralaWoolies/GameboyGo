package gb

import (
	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
)

type HBlankState struct {
	ppu *PPU
}

func newHBlankState(ppu *PPU) PPUState {
	s := &HBlankState{ppu: ppu}
	return s
}

func (s *HBlankState) tick() {
	// end of scanline, check if we are moving into VBLANK or the next scanline
	if s.ppu.ticks >= TICKS_PER_SCANLINE {
		s.ppu.ticks = 0
		s.ppu.incLY()

		if s.ppu.ly >= GB_SCREEN_HEIGHT {
			s.ppu.setState(newVBlankState(s.ppu))
			s.ppu.ic.requestIntrupt(VBLANK_INTRUPT_BIT)

			if bits.IsSet(s.ppu.stat, STAT_VBLANK_ENABLE_BIT) {
				s.ppu.ic.requestIntrupt(LCD_INTRUPT_BIT)
			}
		} else {
			s.ppu.setState(newOAMState(s.ppu))
		}
	}
}
