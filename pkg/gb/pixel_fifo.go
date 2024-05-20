package gb

import (
	"github.com/BeralaWoolies/GameboyGo/pkg/queue"
)

type PixelFIFO struct {
	ppu *PPU

	ticks int

	currState  PixelFIFOState
	FIFO       *queue.Queue[uint8]
	tileLoByte uint8
	tileHiByte uint8

	tileId   uint8
	tileMapY uint16
	tileLine uint16
	fetchX   uint8
}

type PixelFIFOState uint8

const (
	ReadTileID     PixelFIFOState = 0
	ReadTileDataLo PixelFIFOState = 1
	ReadTileDataHi PixelFIFOState = 2
	IDLE           PixelFIFOState = 3
	PushFIFO       PixelFIFOState = 4
)

func (pxF *PixelFIFO) init(ppu *PPU) {
	pxF.ppu = ppu
	pxF.FIFO = queue.NewQueue[uint8](16)
}

func (pxF *PixelFIFO) setState(state PixelFIFOState) {
	pxF.currState = state
}

func (pxF *PixelFIFO) tick() {
	pxF.ticks++

	if pxF.ticks < 2 {
		return
	}

	pxF.ticks = 0

	switch pxF.currState {
	case ReadTileID:
		tileMapYOff := (pxF.tileMapY / TILE_WIDTH) * TILE_MAP_WIDTH
		tileMapXOff := (uint16(pxF.fetchX) + (uint16(pxF.ppu.scx) / TILE_WIDTH)) % TILE_MAP_WIDTH
		tileMapAddr := pxF.ppu.getBGTileMap() + tileMapYOff + tileMapXOff

		pxF.tileId = pxF.ppu.vram[tileMapAddr-VRAM_BASE]

		pxF.setState(ReadTileDataLo)
	case ReadTileDataLo:
		pxF.tileLoByte, _ = pxF.getTileLine(pxF.tileId, pxF.tileLine)

		pxF.setState(ReadTileDataHi)
	case ReadTileDataHi:
		_, pxF.tileHiByte = pxF.getTileLine(pxF.tileId, pxF.tileLine)

		pxF.setState(IDLE)
	case IDLE:
		pxF.setState(PushFIFO)
	case PushFIFO:
		if pxF.FIFO.Size() <= 8 {
			for bit := 7; bit >= 0; bit-- {
				pxF.FIFO.Enqueue(getPaletteId(pxF.tileLoByte, pxF.tileHiByte, uint8(bit)))
			}

			pxF.fetchX++
			pxF.setState(ReadTileID)
		}
	}
}

func (pxF *PixelFIFO) start() {
	pxF.currState = ReadTileID
	pxF.FIFO.Clear()
	pxF.tileLoByte = 0
	pxF.tileHiByte = 0
	pxF.tileMapY = uint16(pxF.ppu.ly + pxF.ppu.scy)
	pxF.tileLine = pxF.tileMapY % TILE_WIDTH
	pxF.fetchX = 0
}

func (pxF *PixelFIFO) getTileLine(tileId uint8, tileLine uint16) (loByte uint8, hiByte uint8) {
	addr, unsig := pxF.ppu.getTileDataArea()
	addr -= VRAM_BASE

	if unsig {
		addr += (uint16(tileId) * TILE_SIZE) + (tileLine * 2)
	} else {
		addr += (uint16(int(int8(tileId)) * TILE_SIZE)) + (tileLine * 2)
	}

	return pxF.ppu.vram[addr], pxF.ppu.vram[addr+1]
}

func (pxF *PixelFIFO) push() (uint8, error) {
	return pxF.FIFO.Dequeue()
}
