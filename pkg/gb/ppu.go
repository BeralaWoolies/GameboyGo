package gb

import (
	"image/color"
	"log"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
	"github.com/hajimehoshi/ebiten/v2"
)

type PPU struct {
	vram  [VRAM_SIZE]uint8
	oam   [OAM_SIZE]uint8
	mmu   *MMU
	ticks int
}

const (
	VRAM_SIZE = 0x2000
	VRAM_BASE = 0x8000
	VRAM_TOP  = 0x9FFF

	OAM_SIZE = 0xA0
	OAM_BASE = 0xFE00
	OAM_TOP  = 0xFE9F

	TILE_SIZE  = 16
	TILE_WIDTH = 8
)

var pallete = [4]color.RGBA{
	0: color.RGBA{
		R: 64,
		G: 80,
		B: 16,
		A: 255,
	},
	1: color.RGBA{
		R: 112,
		G: 128,
		B: 40,
		A: 255,
	},
	2: color.RGBA{
		R: 160,
		G: 168,
		B: 64,
		A: 255,
	},
	3: color.RGBA{
		R: 208,
		G: 208,
		B: 88,
		A: 255,
	},
}

func (ppu *PPU) init(mmu *MMU) {
	ppu.vram = [VRAM_SIZE]uint8{}
	ppu.oam = [OAM_SIZE]uint8{}
	ppu.mmu = mmu
	ppu.ticks = 0
}

func (ppu *PPU) contains(addr uint16) bool {
	return inRange(addr, VRAM_BASE, VRAM_TOP) || inRange(addr, OAM_BASE, OAM_TOP)
}

func (ppu *PPU) write(addr uint16, data uint8) {
	if inRange(addr, VRAM_BASE, VRAM_TOP) {
		ppu.vram[addr-VRAM_BASE] = data
	} else if inRange(addr, OAM_BASE, OAM_TOP) {
		ppu.oam[addr-OAM_BASE] = data
	} else {
		log.Fatal("MMU mapped an illegal write address: 0x%02x to PPU", addr)
	}
}

func (ppu *PPU) read(addr uint16) uint8 {
	if inRange(addr, VRAM_BASE, VRAM_TOP) {
		return ppu.vram[addr-VRAM_BASE]
	} else if inRange(addr, OAM_BASE, OAM_TOP) {
		return ppu.oam[addr-OAM_BASE]
	} else {
		log.Fatal("MMU mapped an illegal read address: 0x%02x to PPU", addr)
		return 0xFF
	}
}

func (ppu *PPU) updateDebugScreen(screen *ebiten.Image, xOff int, yOff int) {
	var tileId uint16 = 0

	for y := 0; y < DEBUG_SCREEN_HEIGHT/TILE_WIDTH; y++ {
		for x := 0; x < DEBUG_SCREEN_WIDTH/TILE_WIDTH; x++ {
			ppu.drawTile(screen, tileId, xOff+(x*TILE_WIDTH), yOff+(y*TILE_WIDTH))
			tileId++
		}
	}
}

func (ppu *PPU) drawTile(screen *ebiten.Image, tileId uint16, x int, y int) {
	tileOff := tileId * TILE_SIZE

	for tileRow := 0; tileRow < 16; tileRow += 2 {
		loByte := ppu.vram[tileOff+uint16(tileRow)]
		hiByte := ppu.vram[tileOff+uint16(tileRow)+1]

		for bit := 7; bit >= 0; bit-- {
			pixLoBit := bits.GetBit(loByte, uint8(bit))
			pixHiBit := bits.GetBit(hiByte, uint8(bit)) << 1

			color := pallete[pixHiBit|pixLoBit]
			screen.Set(x+(7-bit), y+(tileRow/2), color)
		}
	}
}
