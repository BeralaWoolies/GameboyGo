package gb

type OAMState struct {
	ppu *PPU
}

const OAM_SCAN_TICKS = 80

func newOAMState(ppu *PPU) PPUState {
	s := &OAMState{ppu: ppu}
	return s
}

func (s *OAMState) tick() {
	if s.ppu.ticks >= OAM_SCAN_TICKS {
		// end of OAM scan, move to pixel transfer
		s.ppu.xDraw = 0
		s.ppu.setState(newDrawState(s.ppu))
	}
}
