package gb

import (
	"cmp"
	"image/color"
	"log"
	"slices"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
	"github.com/hajimehoshi/ebiten/v2"
)

type PPU struct {
	mmu         *MMU
	dmac        *DMAController
	ic          *IntruptController
	pxF         *PixelFIFO
	frameBuffer [GB_SCREEN_WIDTH][GB_SCREEN_HEIGHT]color.RGBA

	vram          [VRAM_SIZE]uint8
	oam           [OAM_SIZE]uint8
	oamScan       uint16
	spriteBuffer  []Sprite
	spritesOnLine uint8
	ticks         int

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
	wly        uint8

	currState  PPUState
	lx         uint8
	inWindow   bool
	scxDropped uint8
}

type PPUState uint8

type Sprite struct {
	x      uint8
	y      uint8
	tileId uint8
	flags  uint8
}

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
	OBP0_ADDR             = 0xFF48
	OBP1_ADDR             = 0xFF49
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

	LCDC_BGWIN_ENABLE   = 0
	LCDC_OBJ_ENABLE     = 1
	LCDC_OBJ_SIZE       = 2
	LCDC_BG_TILE_MAP    = 3
	LCDC_TILE_DATA_AREA = 4
	LCDC_WIN_ENABLE     = 5
	LCDC_WIN_TILE_MAP   = 6

	SPRITES_PER_SCANLINE = 10

	OAM_SCAN       PPUState = 2
	PIXEL_TRANSFER PPUState = 3
	HBLANK         PPUState = 0
	VBLANK         PPUState = 1

	OAM_SCAN_TICKS = 80
)

var pallete = [4]color.RGBA{
	0: color.RGBA{
		R: 208,
		G: 208,
		B: 88,
		A: 255,
	},
	1: color.RGBA{
		R: 160,
		G: 168,
		B: 64,
		A: 255,
	},
	2: color.RGBA{
		R: 112,
		G: 128,
		B: 40,
		A: 255,
	},
	3: color.RGBA{
		R: 64,
		G: 80,
		B: 16,
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
	ppu.spriteBuffer = make([]Sprite, 0, SPRITES_PER_SCANLINE)
	ppu.setState(OAM_SCAN)
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
		ppu.scanOAM()

		if ppu.ticks >= OAM_SCAN_TICKS {
			// end of OAM scan, move to pixel transfer
			slices.SortStableFunc(ppu.spriteBuffer, func(a, b Sprite) int {
				return cmp.Compare(a.x, b.x)
			})
			ppu.setState(PIXEL_TRANSFER)
		}
	case PIXEL_TRANSFER:
		if ppu.winEnabled() && ppu.winEncountered() {
			if !ppu.pxF.windowFetch {
				ppu.pxF.startWinFetch()
			}
		} else {
			if ppu.pxF.windowFetch {
				ppu.pxF.startBGFetch()
			}
		}

		for ppu.spritesEnabled() && ppu.spriteEncountered() {
			ppu.pxF.spriteFetchFIFO.Add(ppu.spriteBuffer[0])
			ppu.spriteBuffer = ppu.spriteBuffer[1:]
		}

		ppu.pxF.tick()

		if pxFItem, err := ppu.pxF.pop(); err == nil && ppu.handleSCXDrop() {
			var color color.RGBA

			switch pxFItem.palette {
			case BGP:
				color = getLCDColor(ppu.bgPalette, pxFItem.color)
			case OBP0:
				color = getLCDColor(ppu.spPalettes[0], pxFItem.color)
			case OBP1:
				color = getLCDColor(ppu.spPalettes[1], pxFItem.color)
			}

			ppu.frameBuffer[ppu.lx][ppu.ly] = color
			ppu.lx++
		}

		if ppu.lx >= GB_SCREEN_WIDTH {
			// end of pixel transfer, move to HBLANK
			ppu.setState(HBLANK)

			if ppu.winEnabled() && ppu.winEncountered() {
				ppu.wly++
			}
		}
	case HBLANK:
		if ppu.ticks >= TICKS_PER_SCANLINE {
			ppu.resetTicks()
			ppu.incLY()

			if ppu.winEnabled() && ppu.wy == ppu.ly {
				ppu.inWindow = true
			}

			if ppu.ly >= GB_SCREEN_HEIGHT {
				ppu.inWindow = false
				ppu.wly = 0
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

func (ppu *PPU) handleSCXDrop() bool {
	if ppu.scxDropped < ppu.scx%8 {
		ppu.scxDropped++
		return false
	}

	return true
}

func (ppu *PPU) scanOAM() {
	if ppu.ticks&1 == 1 {
		return
	}

	if ppu.spritesOnLine >= SPRITES_PER_SCANLINE {
		return
	}

	spriteY := ppu.mmu.read(ppu.oamScan)
	spriteX := ppu.mmu.read(ppu.oamScan + 1)
	spriteTileId := ppu.mmu.read(ppu.oamScan + 2)
	spriteFlags := ppu.mmu.read(ppu.oamScan + 3)

	if spriteX > 0 && ppu.ly >= spriteY-16 && ppu.ly < spriteY+ppu.getSpriteHeight()-16 {
		ppu.spriteBuffer = append(ppu.spriteBuffer, Sprite{
			x:      spriteX,
			y:      spriteY,
			tileId: spriteTileId,
			flags:  spriteFlags,
		})
		ppu.spritesOnLine++
	}

	ppu.oamScan += 4
}

func (ppu *PPU) getSpriteHeight() uint8 {
	if bits.IsSet(ppu.lcdc, LCDC_OBJ_SIZE) {
		return 16
	}

	return 8
}

func (sp *Sprite) getPalette() Palette {
	if bits.IsSet(sp.flags, 4) {
		return OBP1
	}

	return OBP0
}

func (sp *Sprite) flippedX() bool {
	return bits.IsSet(sp.flags, 5)
}

func (sp *Sprite) flippedY() bool {
	return bits.IsSet(sp.flags, 6)
}

func (sp *Sprite) getBGPriority() bool {
	return bits.IsSet(sp.flags, 7)
}

func (ppu *PPU) spritesEnabled() bool {
	return bits.IsSet(ppu.lcdc, LCDC_OBJ_ENABLE)
}

func (ppu *PPU) spriteEncountered() bool {
	for _, sprite := range ppu.spriteBuffer {
		if sprite.x-8 <= ppu.lx {
			return true
		}
	}

	return false
}

func (ppu *PPU) winEncountered() bool {
	return ppu.inWindow && ppu.lx >= ppu.wx-7
}

func (ppu *PPU) bgWinEnabled() bool {
	return bits.IsSet(ppu.lcdc, LCDC_BGWIN_ENABLE)
}

func (ppu *PPU) winEnabled() bool {
	return bits.IsSet(ppu.lcdc, LCDC_WIN_ENABLE)
}

func (ppu *PPU) setState(state PPUState) {
	ppu.currState = state
	ppu.updateStat(state)
}

func (ppu *PPU) updateStat(state PPUState) {
	ppu.stat = bits.Reset(ppu.stat, 0)
	ppu.stat = bits.Reset(ppu.stat, 1)
	ppu.stat |= uint8(state)

	switch state {
	case OAM_SCAN:
		ppu.oamScan = OAM_BASE
		ppu.spritesOnLine = 0
		ppu.spriteBuffer = ppu.spriteBuffer[:0]

		if bits.IsSet(ppu.stat, STAT_SELECT_OAM) {
			ppu.ic.requestIntrupt(LCD_INTRUPT_BIT)
		}
	case PIXEL_TRANSFER:
		ppu.lx = 0
		ppu.scxDropped = 0
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

func (ppu *PPU) getWinTileMap() uint16 {
	if bits.IsSet(ppu.lcdc, LCDC_WIN_TILE_MAP) {
		return 0x9C00
	}

	return 0x9800
}

func (ppu *PPU) getTileDataArea() (tileDataArea uint16, unsig bool) {
	if bits.IsSet(ppu.lcdc, LCDC_TILE_DATA_AREA) || ppu.pxF.spriteFetch {
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

func (ppu *PPU) setStatDirect(val uint8) {
	ppu.stat = val
}

func (ppu *PPU) setDMADirect(val uint8) {
	ppu.dma = val
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
		ppu.stat = (data & STAT_RW_MSK) | 0x80
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
	case OBP0_ADDR:
		ppu.spPalettes[0] = data
	case OBP1_ADDR:
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
		return ppu.stat | 0x80
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
	case OBP0_ADDR:
		return ppu.spPalettes[0]
	case OBP1_ADDR:
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
			color := getLCDColor(ppu.bgPalette, getColor(loByte, hiByte, uint8(bit)))
			screen.Set(x+(7-bit), y+(tileRow/2), color)
		}
	}
}

func getColor(loByte uint8, hiByte uint8, pos uint8) uint8 {
	pixLoBit := bits.GetBit(loByte, pos)
	pixHiBit := bits.GetBit(hiByte, pos) << 1

	return pixHiBit | pixLoBit
}

func getLCDColor(pal uint8, color uint8) color.RGBA {
	return pallete[(pal>>(2*color))&3]
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
