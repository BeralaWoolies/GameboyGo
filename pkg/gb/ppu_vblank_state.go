package gb

type VBlankState struct {
	ppu *PPU
}

func newVBlankState(ppu *PPU) PPUState {
	s := &VBlankState{ppu: ppu}
	return s
}

func (s *VBlankState) tick() {
	// end of scanline, check if we are moving into next pseudo-scanline or back to the top
	if s.ppu.ticks >= TICKS_PER_SCANLINE {
		s.ppu.ticks = 0
		s.ppu.incLY()

		if s.ppu.ly >= SCANLINES_PER_FRAME {
			s.ppu.ly = 0
			s.ppu.setState(newOAMState(s.ppu))
		}
	}
}
