package gb

import (
	"image/color"
	"log"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
	"github.com/hajimehoshi/ebiten/v2"
)

type PPU struct {
	mmu         *MMU
	dmac        *DMAController
	ic          *IntruptController
	pxF         *PixelFIFO
	frameBuffer [GB_SCREEN_WIDTH][GB_SCREEN_HEIGHT]color.RGBA

	vram  [VRAM_SIZE]uint8
	oam   [OAM_SIZE]uint8
	ticks int

	lcdc       uint8
	stat       uint8
	scy        uint8
	scx        uint8
	ly         uint8
	lyc        uint8
	dma        uint8
	bgPalette  uint8
	spPalettes [NUM_SP_PALETTES]uint8
	wy         uint8
	wx         uint8

	currState PPUState
	xDraw     uint8
}

type PPUState uint8

const (
	VRAM_SIZE = 0x2000
	VRAM_BASE = 0x8000
	VRAM_TOP  = 0x9FFF

	OAM_SIZE = 0xA0
	OAM_BASE = 0xFE00
	OAM_TOP  = 0xFE9F

	TILE_SIZE       = 16
	TILE_WIDTH      = 8
	TILE_MAP_WIDTH  = 32
	NUM_SP_PALETTES = 2

	LCDC_ADDR             = 0xFF40
	STAT_ADDR             = 0xFF41
	SCY_ADDR              = 0xFF42
	SCX_ADDR              = 0xFF43
	LY_ADDR               = 0xFF44
	LYC_ADDR              = 0xFF45
	OAM_DMA_TRANSFER_ADDR = 0xFF46
	BG_PALETTE_ADDR       = 0xFF47
	SP_PALETTE_BASE       = 0xFF48
	WY_ADDR               = 0xFF4A
	WX_ADDR               = 0xFF4B

	TICKS_PER_SCANLINE  = 456
	SCANLINES_PER_FRAME = GB_SCREEN_HEIGHT + 10

	STAT_LYC           = 2
	STAT_SELECT_HBLANK = 3
	STAT_SELECT_VBLANK = 4
	STAT_SELECT_OAM    = 5
	STAT_SELECT_LYC    = 6
	STAT_MSK           = 0x7F
	STAT_RW_MSK        = 0x78

	LCDC_BG_TILE_MAP    = 3
	LCDC_TILE_DATA_AREA = 4

	OAM_SCAN       PPUState = 2
	PIXEL_TRANSFER PPUState = 3
	HBLANK         PPUState = 0
	VBLANK         PPUState = 1

	OAM_SCAN_TICKS = 80
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

func (ppu *PPU) init(mmu *MMU, dmac *DMAController, ic *IntruptController) {
	ppu.mmu = mmu
	ppu.dmac = dmac
	ppu.ic = ic
	ppu.pxF = &PixelFIFO{}
	ppu.pxF.init(ppu)
	ppu.frameBuffer = [GB_SCREEN_WIDTH][GB_SCREEN_HEIGHT]color.RGBA{}

	ppu.vram = [VRAM_SIZE]uint8{}
	ppu.oam = [OAM_SIZE]uint8{}
	ppu.currState = OAM_SCAN
	ppu.lcdc = 0x91
	ppu.ticks = 0
	ppu.ly = 0
	ppu.xDraw = 0
}

func (ppu *PPU) step(cTicks int) {
	for i := 0; i < cTicks; i++ {
		ppu.tick()
	}
}

func (ppu *PPU) tick() {
	ppu.ticks++

	switch ppu.currState {
	case OAM_SCAN:
		if ppu.ticks >= OAM_SCAN_TICKS {
			// end of OAM scan, move to pixel transfer
			ppu.setState(PIXEL_TRANSFER)
		}
	case PIXEL_TRANSFER:
		ppu.pxF.tick()

		if pId, err := ppu.pxF.push(); err == nil {
			ppu.frameBuffer[ppu.xDraw][ppu.ly] = pallete[pId]
			ppu.xDraw++
		}

		if ppu.xDraw >= GB_SCREEN_WIDTH {
			// end of pixel transfer, move to HBLANK
			ppu.setState(HBLANK)
		}
	case HBLANK:
		if ppu.ticks >= TICKS_PER_SCANLINE {
			ppu.resetTicks()
			ppu.incLY()

			if ppu.ly >= GB_SCREEN_HEIGHT {
				ppu.setState(VBLANK)
			} else {
				ppu.setState(OAM_SCAN)
			}
		}
	case VBLANK:
		if ppu.ticks >= TICKS_PER_SCANLINE {
			ppu.resetTicks()
			ppu.incLY()

			if ppu.ly >= SCANLINES_PER_FRAME {
				ppu.resetLY()
				ppu.setState(OAM_SCAN)
			}
		}
	default:
		log.Fatalf("PPU is in an unimplemented state")
	}
}

func (ppu *PPU) setState(state PPUState) {
	ppu.currState = state
	ppu.stat = bits.Reset(ppu.stat, 0)
	ppu.stat = bits.Reset(ppu.stat, 1)
	ppu.stat |= uint8(state)

	ppu.updateStat(state)
}

func (ppu *PPU) updateStat(state PPUState) {
	switch state {
	case OAM_SCAN:
		if bits.IsSet(ppu.stat, STAT_SELECT_OAM) {
			ppu.ic.requestIntrupt(LCD_INTRUPT_BIT)
		}
	case PIXEL_TRANSFER:
		ppu.xDraw = 0
		ppu.pxF.start()
	case HBLANK:
		if bits.IsSet(ppu.stat, STAT_SELECT_HBLANK) {
			ppu.ic.requestIntrupt(LCD_INTRUPT_BIT)
		}
	case VBLANK:
		ppu.ic.requestIntrupt(VBLANK_INTRUPT_BIT)

		if bits.IsSet(ppu.stat, STAT_SELECT_VBLANK) {
			ppu.ic.requestIntrupt(LCD_INTRUPT_BIT)
		}
	}
}

func (ppu *PPU) getBGTileMap() uint16 {
	if bits.IsSet(ppu.lcdc, LCDC_BG_TILE_MAP) {
		return 0x9C00
	}

	return 0x9800
}

func (ppu *PPU) getTileDataArea() (tileDataArea uint16, unsig bool) {
	if bits.IsSet(ppu.lcdc, LCDC_TILE_DATA_AREA) {
		return VRAM_BASE, true
	}

	return 0x9000, false
}

func (ppu *PPU) resetTicks() {
	ppu.ticks = 0
}

func (ppu *PPU) resetLY() {
	ppu.ly = 0
}

func (ppu *PPU) incLY() {
	ppu.ly++

	if ppu.ly == ppu.lyc {
		ppu.stat = bits.Set(ppu.stat, STAT_LYC)

		if bits.IsSet(ppu.stat, STAT_SELECT_LYC) {
			ppu.ic.requestIntrupt(LCD_INTRUPT_BIT)
		}
	} else {
		ppu.stat = bits.Reset(ppu.stat, STAT_LYC)
	}
}

func (ppu *PPU) contains(addr uint16) bool {
	return (inRange(addr, VRAM_BASE, VRAM_TOP) ||
		inRange(addr, OAM_BASE, OAM_TOP) ||
		inRange(addr, LCDC_ADDR, WX_ADDR))
}

func (ppu *PPU) write(addr uint16, data uint8) {
	if inRange(addr, VRAM_BASE, VRAM_TOP) {
		ppu.vram[addr-VRAM_BASE] = data
		return
	} else if inRange(addr, OAM_BASE, OAM_TOP) {
		ppu.oam[addr-OAM_BASE] = data
		return
	}

	switch addr {
	case LCDC_ADDR:
		ppu.lcdc = data
	case STAT_ADDR:
		ppu.stat = data & STAT_RW_MSK
	case SCY_ADDR:
		ppu.scy = data
	case SCX_ADDR:
		ppu.scx = data
	case LY_ADDR:
		// LY is read only
		return
	case LYC_ADDR:
		ppu.lyc = data
	case OAM_DMA_TRANSFER_ADDR:
		ppu.dma = data
		ppu.dmac.initOAMTransfer(data)
	case BG_PALETTE_ADDR:
		ppu.bgPalette = data
	case SP_PALETTE_BASE:
		ppu.spPalettes[0] = data
	case SP_PALETTE_BASE + 1:
		ppu.spPalettes[1] = data
	case WY_ADDR:
		ppu.wy = data
	case WX_ADDR:
		ppu.wx = data
	default:
		log.Fatalf("MMU mapped an illegal write address: 0x%02x to PPU", addr)
	}
}

func (ppu *PPU) read(addr uint16) uint8 {
	if inRange(addr, VRAM_BASE, VRAM_TOP) {
		return ppu.vram[addr-VRAM_BASE]
	} else if inRange(addr, OAM_BASE, OAM_TOP) {
		return ppu.oam[addr-OAM_BASE]
	}

	switch addr {
	case LCDC_ADDR:
		return ppu.lcdc
	case STAT_ADDR:
		return ppu.stat & STAT_MSK
	case SCY_ADDR:
		return ppu.scy
	case SCX_ADDR:
		return ppu.scx
	case LY_ADDR:
		return ppu.ly
	case LYC_ADDR:
		return ppu.lyc
	case OAM_DMA_TRANSFER_ADDR:
		return ppu.dma
	case BG_PALETTE_ADDR:
		return ppu.bgPalette
	case SP_PALETTE_BASE:
		return ppu.spPalettes[0]
	case SP_PALETTE_BASE + 1:
		return ppu.spPalettes[1]
	case WY_ADDR:
		return ppu.wy
	case WX_ADDR:
		return ppu.wx
	default:
		log.Fatalf("MMU mapped an illegal read address: 0x%02x to PPU", addr)
		return 0xFF
	}
}

func (ppu *PPU) updateGBScreen(screen *ebiten.Image, xOff int, yOff int) {
	for y := 0; y < GB_SCREEN_HEIGHT; y++ {
		for x := 0; x < GB_SCREEN_WIDTH; x++ {
			screen.Set(x+xOff, y+yOff, ppu.frameBuffer[x][y])
		}
	}
}

func (ppu *PPU) writeTile(screen *ebiten.Image, tileId uint16, x int, y int) {
	addr, unsig := ppu.getTileDataArea()
	addr -= VRAM_BASE

	if unsig {
		addr += (uint16(tileId) * TILE_SIZE)
	} else {
		addr += (uint16(int(int8(tileId)) * TILE_SIZE))
	}

	for tileRow := 0; tileRow < 16; tileRow += 2 {
		loByte := ppu.vram[addr+uint16(tileRow)]
		hiByte := ppu.vram[addr+uint16(tileRow)+1]

		for bit := 7; bit >= 0; bit-- {
			color := pallete[getPaletteId(loByte, hiByte, uint8(bit))]
			screen.Set(x+(7-bit), y+(tileRow/2), color)
		}
	}
}

func getPaletteId(loByte uint8, hiByte uint8, pos uint8) uint8 {
	pixLoBit := bits.GetBit(loByte, pos)
	pixHiBit := bits.GetBit(hiByte, pos) << 1

	return pixHiBit | pixLoBit
}

// ============================= Debug Functions ===============================
func (ppu *PPU) updateTileDataScreen(screen *ebiten.Image, xOff int, yOff int) {
	var tileId uint16 = 0

	for y := 0; y < TILE_DATA_SCREEN_HEIGHT/TILE_WIDTH; y++ {
		for x := 0; x < TILE_DATA_SCREEN_WIDTH/TILE_WIDTH; x++ {
			ppu.writeTile(screen, tileId, xOff+(x*TILE_WIDTH), yOff+(y*TILE_WIDTH))
			tileId++
		}
	}
}

func (ppu *PPU) updateTileMaps(screen *ebiten.Image, xOff int, yOff int) {
	var tileMap1 uint16 = 0x9800
	var tileMap2 uint16 = 0x9C00

	for y := 0; y < TILE_MAP_SCREEN_WIDTH/TILE_WIDTH; y++ {
		for x := 0; x < TILE_MAP_SCREEN_WIDTH/TILE_WIDTH; x++ {
			tileMap1Addr := tileMap1 + uint16(y)*TILE_MAP_WIDTH + uint16(x)
			tileMap2Addr := tileMap2 + uint16(y)*TILE_MAP_WIDTH + uint16(x)

			ppu.writeTile(screen, uint16(ppu.vram[tileMap1Addr-VRAM_BASE]), xOff+(x*TILE_WIDTH), yOff+(y*TILE_WIDTH))
			ppu.writeTile(screen, uint16(ppu.vram[tileMap2Addr-VRAM_BASE]), xOff+TILE_MAP_SCREEN_WIDTH+(x*TILE_WIDTH), yOff+(y*TILE_WIDTH))
		}
	}
}
