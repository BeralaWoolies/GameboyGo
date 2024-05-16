package gb

type DrawState struct {
	ppu *PPU
}

func newDrawState(ppu *PPU) PPUState {
	s := &DrawState{ppu: ppu}
	return s
}

func (s *DrawState) tick() {
	// TODO: tick the pixel FIFO and only inc xDraw if it has outputted a pixel
	s.ppu.xDraw++
	if s.ppu.xDraw >= GB_SCREEN_WIDTH {
		// end of pixel transfer, move to HBLANK
		s.ppu.setState(newHBlankState(s.ppu))
	}
}
